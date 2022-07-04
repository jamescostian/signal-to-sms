package tosqlite

import (
	"context"
	"fmt"

	"github.com/jamescostian/signal-to-sms/utils/proto/frameio"
	"github.com/jamescostian/signal-to-sms/utils/proto/signal"
	"github.com/jamescostian/signal-to-sms/utils/sqlite"
	"github.com/pkg/errors"
)

type importer struct {
	DBOrTx             sqlite.DBOrTx
	ShouldIgnoreImport func(string) bool
	ctx                context.Context
	// Temporary storage to avoid allocating things
	importParams []interface{}
	err          error
}

// NewFrameWriter provides an efficient way to import SQL statements from Signal's backup protobufs, while filtering out some statements that cause issues
func NewFrameWriter(ctx context.Context, dbOrTx sqlite.DBOrTx, fastSlightlyInaccurateImport bool) frameio.FrameWriter {
	importer := importer{ctx: ctx, DBOrTx: dbOrTx}

	importer.ShouldIgnoreImport = shouldIgnoreImport
	if fastSlightlyInaccurateImport {
		importer.ShouldIgnoreImport = fastShouldIgnoreImport
	}

	return &importer
}

func (i *importer) Write(frame *signal.BackupFrame, blob []byte) error {
	if frame.Statement != nil {
		if i.err = i.Import(frame.Statement); i.err != nil {
			return i.err
		}
	}
	return nil
}

// Import runs an SQL statement, unless it should be ignored to avoid errors
func (i *importer) Import(sqlStatement *signal.SqlStatement) error {
	statement := sqlStatement.GetStatement()
	if i.ShouldIgnoreImport(statement) {
		return nil
	}
	s, err := i.DBOrTx.PrepareContext(i.ctx, statement)
	if err != nil {
		return errors.Wrap(err, "invalid SQL received during import")
	}

	// Reset the array of import parameters to be empty (but maintain the allocated space in memory), and then fill it up with data
	i.importParams = i.importParams[:0]
	if err = decodeImportParams(&i.importParams, sqlStatement.Parameters); err != nil {
		return err
	}

	_, err = s.ExecContext(i.ctx, i.importParams...)
	return errors.Wrap(err, "failed executing SQL during import")
}

// decodeImportParams converts SQL parameters from protobuf form to a form that can be used by database/sql
func decodeImportParams(dest *[]interface{}, source []*signal.SqlStatement_SqlParameter) error {
	for _, param := range source {
		decodedParam, err := decodeImportParam(param)
		if err != nil {
			return err
		}
		*dest = append(*dest, decodedParam)
	}
	return nil
}

// decodeImportParam extracts the true value from a signal protobuf form of an SQL parameter
func decodeImportParam(param *signal.SqlStatement_SqlParameter) (interface{}, error) {
	if param.Nullparameter != nil {
		return nil, nil
	} else if param.StringParamter != nil {
		return *param.StringParamter, nil
	} else if param.BlobParameter != nil {
		return param.BlobParameter, nil
	} else if param.DoubleParameter != nil {
		return *param.DoubleParameter, nil
	} else if param.IntegerParameter != nil {
		// Java doesn't have unsigned types, and Signal-Android is written in Java.
		// Unfortunately, someone who wrote the protobuf called the IntegerParameter a uint64, instead of a int64.
		// In Java, it doesn't matter if you specify uint64 or int64 in your protobuf, you'll get a long either way.
		// See: https://developers.google.com/protocol-buffers/docs/proto#scalar
		// However, golang has proper unsigned types, so a Java long set to -1 encoded as a uint64 in protobufs will be treated as the highest uint64 number possible.
		// In order to interpret integers the way the Android app does, uint64->int64 conversion is necessary.
		// Also, uint64s can cause issues on SQLite: https://github.com/golang/go/issues/9373#issuecomment-91424470
		return int64(*param.IntegerParameter), nil
	}
	return nil, fmt.Errorf("unknown import param - perhaps the protobuf definition has changed?")
}
