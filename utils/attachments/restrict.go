package attachments

import (
	"context"
	"errors"
)

type noWrites struct {
	Store Store
}

func NoWrites(store Store) Store {
	return &noWrites{store}
}
func (r *noWrites) SetAttachment(ctx context.Context, id int64, blob []byte) error {
	return errors.New("cannot SetAttachment on a NoWrites attachment store")
}
func (r *noWrites) GetAttachment(ctx context.Context, id int64) (blob []byte, err error) {
	return r.Store.GetAttachment(ctx, id)
}
func (r *noWrites) ListAttachments(ctx context.Context) (ids []int64, err error) {
	return r.Store.ListAttachments(ctx)
}
func (r *noWrites) DeleteAttachment(ctx context.Context, id int64) (err error) {
	return errors.New("cannot DeleteAttachment on a NoWrites attachment store")
}
func (r *noWrites) Close() error {
	return r.Store.Close()
}

type noReads struct {
	Store Store
}

func NoReads(store Store) Store {
	return &noReads{store}
}
func (r *noReads) SetAttachment(ctx context.Context, id int64, blob []byte) error {
	return r.Store.SetAttachment(ctx, id, blob)
}
func (r *noReads) GetAttachment(ctx context.Context, id int64) (blob []byte, err error) {
	return nil, errors.New("cannot GetAttachment on a NoReads attachment store")
}
func (r *noReads) ListAttachments(ctx context.Context) (ids []int64, err error) {
	return nil, errors.New("cannot ListAttachments on a NoReads attachment store")
}
func (r *noReads) DeleteAttachment(ctx context.Context, id int64) (err error) {
	return r.Store.DeleteAttachment(ctx, id)
}
func (r *noReads) Close() error {
	return r.Store.Close()
}

type unclosable struct {
	Store Store
}

func Unclosable(store Store) Store {
	return &unclosable{store}
}
func (r *unclosable) SetAttachment(ctx context.Context, id int64, blob []byte) error {
	return r.Store.SetAttachment(ctx, id, blob)
}
func (r *unclosable) GetAttachment(ctx context.Context, id int64) (blob []byte, err error) {
	return r.Store.GetAttachment(ctx, id)
}
func (r *unclosable) ListAttachments(ctx context.Context) (ids []int64, err error) {
	return r.Store.ListAttachments(ctx)
}
func (r *unclosable) DeleteAttachment(ctx context.Context, id int64) (err error) {
	return r.Store.DeleteAttachment(ctx, id)
}
func (r *unclosable) Close() error {
	return errors.New("cannot close Unclosable attachment store")
}
