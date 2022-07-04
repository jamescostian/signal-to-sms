package conversion

import "context"

func ConvertToNextFormat(ctx context.Context, c *Conversion) error {
	return c.convertToNextFormat(ctx)
}
