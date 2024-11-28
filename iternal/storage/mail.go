package storage

import (
	"context"
	"log/slog"

	utils "github.com/didikizi/RedSoft/packege"
	"github.com/jackc/pgx/v5"
)

type PutMail interface {
	GetMail() string
	GetHumanId() int
	GetDescription() string
}

func (s *Storage) PutMailForHuman(ctx context.Context, mail PutMail) error {
	args := pgx.NamedArgs{
		"human_id":    mail.GetHumanId(),
		"description": mail.GetDescription(),
		"mail":        mail.GetMail(),
	}

	rows, err := s.pool.Query(ctx, s.insertMail, args)
	if err != nil {
		slog.Info(utils.GetCallerInfo(), slog.String("err:", err.Error()))
	}

	defer rows.Close()

	return err
}

type Mail struct {
	Id          int    `db:"id"`
	HumanId     int    `db:"human_id"`
	Mail        string `db:"mail"`
	Description string `db:"description"`
}

func (s *Storage) GetMailForHuman(ctx context.Context, id int) (result []*Mail, err error) {
	args := pgx.NamedArgs{
		"human_id": id,
	}

	rows, err := s.pool.Query(ctx, s.selectMail, args)
	if err != nil {
		slog.Info(utils.GetCallerInfo(), slog.String("err:", err.Error()))

		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		mail, err := pgx.RowToAddrOfStructByNameLax[Mail](rows)
		if err != nil {
			slog.Info(utils.GetCallerInfo(), slog.String("err:", err.Error()))

			return nil, err
		}

		result = append(result, mail)
	}

	return
}

func (s *Storage) DeleteMailForHuman(ctx context.Context, id int) (bool, error) {
	args := pgx.NamedArgs{"id": id}

	rows, err := s.pool.Query(ctx, s.deleteMail, args)
	if err != nil {
		slog.Info(utils.GetCallerInfo(), slog.String("err:", err.Error()))

		return false, err
	}

	defer rows.Close()

	return true, nil
}
