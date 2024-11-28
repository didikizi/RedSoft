package service

import (
	"context"

	"github.com/didikizi/RedSoft/iternal/storage"
)

type Service struct {
	Storage
}

type Storage interface {
	GetHumanFromSurname(context.Context, string) ([]*storage.Human, error)
	GetHumanFromId(context.Context, int) (*storage.Human, error)
	GetHumanList(context.Context) ([]*storage.Human, error)
	PostHuman(context.Context, storage.PutHuman, int) error
	PutHuman(context.Context, storage.PutHuman) error
	DeleteHumanFromId(context.Context, int) (bool, error)

	PutMailForHuman(context.Context, storage.PutMail) error
	GetMailForHuman(context.Context, int) ([]*storage.Mail, error)
	DeleteMailForHuman(context.Context, int) (bool, error)
}

func New(storage Storage) *Service {
	return &Service{Storage: storage}
}
