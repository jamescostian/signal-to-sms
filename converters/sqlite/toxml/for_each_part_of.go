package toxml

import (
	"context"
	"strings"

	"github.com/cristalhq/base64"
	"github.com/jamescostian/signal-to-sms/internal/fromsignalback/mimetoext"
)

// ForEachPartOf will loop through each mmsPart (an attachment) of a MMS in the DB, and let you run a function on each part.
// Like go's for loops, the variable you get here will be reused, so make a copy of it if you need it to last after your function returns!
// If there's an error getting something from the DB, or if your function returns an error, then iteration will halt, and the error will be returned by ForEachPartOf.
func (w *writer) ForEachPartOf(ctx context.Context, mid int64, processOne func(*mmsPart) error) error {
	prepped, err := w.dbOrTx.PrepareContext(ctx, `
		SELECT
			_id,
			seq,
			ct,
			COALESCE(file_name, cast(_id as text)),
			COALESCE(chset, `+utf8Str+`),
			cd,
			fn,
			cid,
			cl,
			IIF(
				ctt_s IS NOT NULL,
				cast(ctt_s as text),
				ctt_s
			),
			ctt_t,
			caption
		FROM
			part
		WHERE
			mid = ?
		ORDER BY
			part._id`)
	if err != nil {
		return err
	}
	rows, err := prepped.QueryContext(ctx, mid)
	if err != nil {
		return err
	}
	defer rows.Close()
	var (
		part mmsPart
		id   int64
	)
	for rows.Next() {
		if err = rows.Scan(
			&id,
			&part.Order,
			&part.ContentType,
			&part.Name,
			&part.Charset,
			&part.ContentDisposition,
			&part.FileName,
			&part.ContentID,
			&part.ContentLocation,
			&part.ContentTypeStart,
			&part.ContentTypeType,
			&part.Text,
		); err != nil {
			return err
		}

		rawData, err := w.attachments.GetAttachment(context.Background(), id)
		if err != nil {
			return err
		}
		part.Data = base64.StdEncoding.EncodeToString(rawData)

		if !strings.ContainsRune(part.Name, '.') {
			ext, err := mimetoext.Convert(part.ContentType)
			if err != nil {
				return err
			}
			part.Name += "." + ext
		}

		if err := processOne(&part); err != nil {
			return err
		}
	}
	return nil
}
