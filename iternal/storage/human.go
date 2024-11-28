package storage

import (
	"context"
	"log/slog"

	utils "github.com/didikizi/RedSoft/packege"
	"github.com/jackc/pgx/v5"
)

type Human struct {
	Id         int    `db:"id"`
	Age        int    `db:"age"`
	Status     string `db:"status"`
	Surname    string `db:"surname"`
	Name       string `db:"name"`
	Fatherland string `db:"fatherland"`
	National   string `db:"national"`
	Sex        string `db:"sex"`
}

type PutHuman interface {
	GetNational() string
	GetSex() string
	GetSurname() string
	GetName() string
	GetFatherland() string
	GetAge() int
}

func (s *Storage) GetHumanFromId(ctx context.Context, id int) (*Human, error) {
	args := pgx.NamedArgs{"id": id}

	rows, err := s.pool.Query(ctx, s.selectHumanFromId, args)
	if err != nil {
		slog.Info(utils.GetCallerInfo(), slog.String("err:", err.Error()))

		return nil, err
	}

	defer rows.Close()

	rows.Next()

	human, err := pgx.RowToAddrOfStructByNameLax[Human](rows)
	if err != nil {
		slog.Info(utils.GetCallerInfo(), slog.String("err:", err.Error()))

		return nil, err
	}

	return human, nil
}

func (s *Storage) GetHumanFromSurname(ctx context.Context, surname string) (result []*Human, err error) {
	args := pgx.NamedArgs{"surname": surname}

	rows, err := s.pool.Query(ctx, s.selectHumanFromSurname, args)
	if err != nil {
		slog.Info(utils.GetCallerInfo(), slog.String("err:", err.Error()))

		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		human, err := pgx.RowToAddrOfStructByNameLax[Human](rows)
		if err != nil {
			slog.Info(utils.GetCallerInfo(), slog.String("err:", err.Error()))

			return nil, err
		}

		result = append(result, human)
	}

	return
}

func (s *Storage) GetHumanList(ctx context.Context) (result []*Human, err error) {
	rows, err := s.pool.Query(ctx, s.selectAllHuman)
	if err != nil {
		slog.Info(utils.GetCallerInfo(), slog.String("err:", err.Error()))

		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		human, err := pgx.RowToAddrOfStructByNameLax[Human](rows)
		if err != nil {
			slog.Info(utils.GetCallerInfo(), slog.String("err:", err.Error()))

			return nil, err
		}

		result = append(result, human)
	}

	return
}

func (s *Storage) PutHuman(ctx context.Context, human PutHuman) error {
	args := pgx.NamedArgs{
		"status":     1,
		"age":        human.GetAge(),
		"surname":    human.GetSurname(),
		"name":       human.GetName(),
		"fatherland": human.GetFatherland(),
		"national":   human.GetNational(),
		"sex":        human.GetSex(),
	}

	rows, err := s.pool.Query(ctx, s.insertHuman, args)
	if err != nil {
		slog.Info(utils.GetCallerInfo(), slog.String("err:", err.Error()))
	}

	defer rows.Close()

	return err
}

func (s *Storage) DeleteHumanFromId(ctx context.Context, id int) (bool, error) {
	args := pgx.NamedArgs{"id": id}

	rows, err := s.pool.Query(ctx, s.deleteHuman, args)
	if err != nil {
		slog.Info(utils.GetCallerInfo(), slog.String("err:", err.Error()))

		return false, err
	}

	defer rows.Close()

	return true, nil
}

func (s *Storage) PostHuman(ctx context.Context, human PutHuman, id int) error {
	args := pgx.NamedArgs{
		"id":         id,
		"status":     2,
		"age":        human.GetAge(),
		"surname":    human.GetSurname(),
		"name":       human.GetName(),
		"fatherland": human.GetFatherland(),
		"national":   human.GetNational(),
		"sex":        human.GetSex(),
	}

	rows, err := s.pool.Query(ctx, s.updateHumanFromId, args)
	if err != nil {
		slog.Info(utils.GetCallerInfo(), slog.String("err:", err.Error()))
	}

	defer rows.Close()

	return err
}
