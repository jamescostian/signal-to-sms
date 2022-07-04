package conversion

import (
	"github.com/RyanCarrier/dijkstra"
	"github.com/jamescostian/signal-to-sms/formats"
)

type Graph struct {
	// Used to find conversion paths
	graph *dijkstra.Graph
	// The graph is built off of these slices, and vertex IDs in the graph are actually indices in the graphFormats.
	// Using NewGraph will update them.
	graphFormats    []formats.MsgFormat
	graphConverters []MsgConverter
	// Given a formats.MsgFormat, look up its index in graphFormats
	lookupGraphFormatIndex map[string]int
	// Given an input formats.MsgFormat and an output formats.MsgFormat, what's the converter to get between them?
	msgConverterLookup map[string]map[string]MsgConverter
}

// NewGraph builds the graph used for calculating the shortest path between formats
func NewGraph(allConverters []MsgConverter, allFormats []formats.MsgFormat) (*Graph, error) {
	cg := &Graph{
		graph:                  dijkstra.NewGraph(),
		graphFormats:           allFormats,
		graphConverters:        allConverters,
		lookupGraphFormatIndex: make(map[string]int, len(allFormats)),
		msgConverterLookup:     make(map[string]map[string]MsgConverter, len(allFormats)),
	}
	for i, format := range allFormats {
		cg.msgConverterLookup[format.Name] = make(map[string]MsgConverter, len(allFormats))
		cg.lookupGraphFormatIndex[format.Name] = i
		cg.graph.AddVertex(i)
	}
	for _, converter := range allConverters {
		cg.msgConverterLookup[converter.InputFormat.Name][converter.OutputFormat.Name] = converter
		distance := int64(10)
		if converter.InputFormat.IncludesAttachments && converter.OutputFormat.IncludesAttachments {
			// Attachments may also need to be converted during this process.
			// Since there's another conversion happening, increase the weight of taking this route.
			// Assume there's a 50% chance of a conversion being needed (since there is no heuristic)
			distance += 5
		} else if converter.InputFormat.IncludesAttachments != converter.OutputFormat.IncludesAttachments {
			// Converting between these formats definitely requires some sort of attachment conversion, which is a time consuming procedure.
			// While the previous condition *may* affect things, this one *definitely* affects things, which is why it's so heavily weighted
			distance += 10
		}
		if err := cg.graph.AddArc(cg.lookupGraphFormatIndex[converter.InputFormat.Name], cg.lookupGraphFormatIndex[converter.OutputFormat.Name], distance); err != nil {
			return nil, err
		}
	}
	return cg, nil
}

func (cg *Graph) FindPath(start formats.MsgFormat, end formats.MsgFormat) (path []MsgConverter, err error) {
	formatIndexPath, err := cg.graph.Shortest(cg.lookupGraphFormatIndex[start.Name], cg.lookupGraphFormatIndex[end.Name])
	if err != nil {
		return
	}
	for i := 1; i < len(formatIndexPath.Path); i++ {
		inputFormat := cg.graphFormats[formatIndexPath.Path[i-1]]
		outputFormat := cg.graphFormats[formatIndexPath.Path[i]]
		path = append(path, cg.msgConverterLookup[inputFormat.Name][outputFormat.Name])
	}
	return
}
