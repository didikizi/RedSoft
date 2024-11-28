package service

import (
	"context"
	"log/slog"

	"github.com/SteelPangolin/go-genderize"
	"github.com/didikizi/RedSoft/iternal/storage"
	utils "github.com/didikizi/RedSoft/packege"
	"github.com/masonkmeyer/nationalize"
)

type Human struct {
	Id         int    `json:"id"`
	Age        int    `json:"age"`
	Surname    string `json:"surname"`
	Name       string `json:"name"`
	Fatherland string `json:"fatherland"`
	National   string `json:"national"`
	Sex        string `json:"sex"`
}

func (h Human) GetSurname() string    { return h.Surname }
func (h Human) GetName() string       { return h.Name }
func (h Human) GetFatherland() string { return h.Fatherland }
func (h Human) GetNational() string   { return h.National }
func (h Human) GetSex() string        { return h.Sex }
func (h Human) GetAge() int           { return h.Age }

func (s *Service) GetHumanList(ctx context.Context) ([]*Human, error) {
	humans, err := s.Storage.GetHumanList(ctx)
	if err != nil {
		slog.Info(utils.GetCallerInfo(), slog.String("err:", err.Error()))

		return nil, err
	}

	return convertHuman(humans...), nil
}

func (s *Service) GetHumanFromSurname(ctx context.Context, surname string) ([]*Human, error) {
	human, err := s.Storage.GetHumanFromSurname(ctx, surname)
	if err != nil {
		slog.Info(utils.GetCallerInfo(), slog.String("err:", err.Error()))

		return nil, err
	}

	return convertHuman(human...), nil
}

func (s *Service) GetHumanFromId(ctx context.Context, id int) ([]*Human, error) {
	human, err := s.Storage.GetHumanFromId(ctx, id)
	if err != nil {
		slog.Info(utils.GetCallerInfo(), slog.String("err:", err.Error()))

		return nil, err
	}

	return convertHuman(human), nil
}

func (s *Service) DeleteHuman(ctx context.Context, id int) (bool, error) {
	return s.Storage.DeleteHumanFromId(ctx, id)
}

type PutHuman interface {
	GetSurname() string
	GetName() string
	GetFatherland() string
	GetAge() int
}

func (s *Service) PutHuman(ctx context.Context, human PutHuman) error {
	var sex string

	response, err := genderize.Get([]string{human.GetName()})
	if err != nil {
		slog.Info(utils.GetCallerInfo(), slog.String("err:", err.Error()))
		sex = "nil"
	}

	if len(sex) != 0 {
		slog.Info(utils.GetCallerInfo(), slog.String("genderize response:", "nil"))
		sex = "nil"
	} else {
		sex = response[0].Gender
	}

	client := nationalize.NewClient()
	tmp, _, err := client.Predict(human.GetSurname())
	national := tmp.Country[0].CountryId

	if err != nil {
		national = "none"
	}

	putHuman := Human{
		Name:       human.GetName(),
		Surname:    human.GetSurname(),
		Fatherland: human.GetFatherland(),
		Age:        human.GetAge(),
		National:   national,
		Sex:        sex,
	}

	err = s.Storage.PutHuman(ctx, putHuman)
	if err != nil {
		slog.Info(utils.GetCallerInfo(), slog.String("err:", err.Error()))

		return err
	}

	return nil
}

type PostHuman interface {
	GetSurname() string
	GetName() string
	GetFatherland() string
	GetAge() int
	GetNational() string
	GetSex() string
}

func (s *Service) PostHuman(ctx context.Context, human PostHuman, id int) error {
	putHuman := Human{
		Name:       human.GetName(),
		Surname:    human.GetSurname(),
		Fatherland: human.GetFatherland(),
		Age:        human.GetAge(),
		National:   human.GetNational(),
		Sex:        human.GetSex(),
	}

	err := s.Storage.PostHuman(ctx, putHuman, id)
	if err != nil {
		slog.Info(utils.GetCallerInfo(), slog.String("err:", err.Error()))

		return err
	}

	return nil
}

func convertHuman(human ...*storage.Human) (result []*Human) {
	if human == nil || human[0] == nil {
		return
	}

	for _, value := range human {
		result = append(result, &Human{
			Id:         value.Id,
			Age:        value.Age,
			Surname:    value.Surname,
			Name:       value.Name,
			Fatherland: value.Fatherland,
			National:   value.National,
			Sex:        value.Sex,
		})
	}

	return
}
