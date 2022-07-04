package attachments

import "context"

func Copy(ctx context.Context, dest Store, src Store) error {
	ids, err := src.ListAttachments(ctx)
	if err != nil {
		return err
	}
	for _, id := range ids {
		blob, err := src.GetAttachment(ctx, id)
		if err != nil {
			return err
		}
		if err = dest.SetAttachment(ctx, id, blob); err != nil {
			return err
		}
	}
	return nil
}
