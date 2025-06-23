package storage

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
)

const getAdoptionStatisticsPG = `
SELECT
    u.user_uuid,
    u.firstname,
    u.lastname,
    u.email,
    u.status,
    u.role,
    u.updated_at,
    COUNT(DISTINCT c.id)     			AS conversation_count,
    COUNT(m.id)              			AS message_count
FROM   public.users          			AS u
     LEFT JOIN public.conversations 	AS c
            ON c.user_uuid = u.user_uuid
     LEFT JOIN public.messages      	AS m
            ON m.conv_id = c.id
GROUP  BY
    u.user_uuid, u.firstname, u.lastname,
    u.email, u.status, u.role, u.updated_at
ORDER  BY u.updated_at DESC;`

type AdoptionStat struct {
	ID                  string
	FirstName, LastName string
	Email               string
	Status              string
	Role                string
	UpdatedAt           time.Time
	ConversationCount   int
	MessageCount        int
}

func (u *usersRepoPgx) GetAdoptionStatistics(ctx context.Context) ([]*AdoptionStat, error) {
	out := make([]*AdoptionStat, 0)
	rows, err := u.Query(ctx, getAdoptionStatisticsPG)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}
	for rows.Next() {
		var t AdoptionStat
		err = rows.Scan(&t.ID, &t.FirstName, &t.LastName, &t.Email, &t.Status, &t.Role, &t.UpdatedAt, &t.ConversationCount, &t.MessageCount)
		if err != nil {
			return nil, err
		}
		out = append(out, &t)
	}

	return out, nil
}
