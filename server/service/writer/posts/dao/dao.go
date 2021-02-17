package dao

import (
	"context"
	"cxfw/db"
	"cxfw/model/writer"
	"cxfw/orm"
	"cxfw/types"
	"net/http"

	"github.com/jackc/pgx/v4"
)

func Add(m *PostWithTagIDs) (int, *PostWithTags, error) {
	tags := make([]writer.Tag, 0)

	err := orm.NewTx(db.S(), func(tx pgx.Tx) error {
		if err := tx.QueryRow(context.Background(),
			`insert into posts values (default, $1, $2, now() at time zone 'utc', now() at time zone 'utc') returning id`,
			m.Content, types.StatusActive).Scan(&m.ID); err != nil {
			return err
		}

		if len(m.TagIDs) > 0 {
			bt := pgx.Batch{}
			for _, tagID := range m.TagIDs {
				bt.Queue(`insert into post_tags values ($1, $2);`, m.ID, tagID)
			}
			br := tx.SendBatch(context.Background(), &bt)
			if err := br.Close(); err != nil {
				return err
			}

			rows, err := tx.Query(context.Background(), `select * from tags where id in (select tag_id from post_tags where post_id = $1)`, m.ID)
			if err != nil {
				return err
			}

			for rows.Next() {
				var tag writer.Tag
				err := rows.Scan(&tag.ID, &tag.Title, nil, nil)
				if err != nil {
					rows.Close()
					return err
				}
				tags = append(tags, tag)
			}
			rows.Close()
		}

		return nil
	})

	rm := PostWithTags{
		Post:  m.Post,
		Title: m.Title,
		Tags:  tags,
	}

	return http.StatusOK, &rm, err
}
