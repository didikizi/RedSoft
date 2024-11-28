package service

import (
	"context"

	"github.com/didikizi/RedSoft/iternal/storage"
)

type PutMail interface {
	GetHumanId() int
	GetMail() string
	GetDescription() string
}

func (s *Service) PutMailForHuman(ctx context.Context, mail PutMail) error {
	return s.Storage.PutMailForHuman(ctx, mail)
}

type Mail struct {
	Id          int    `json:"id"`
	HumanId     int    `json:"Human_id"`
	Mail        string `json:"mail"`
	Description string `json:"description"`
}

func (s *Service) GetMailListForHuman(ctx context.Context, id int) ([]*Mail, error) {
	mails, err := s.Storage.GetMailForHuman(ctx, id)

	return convertMails(mails...), err
}

func (s *Service) DeleteMailForHuman(ctx context.Context, id int) (bool, error) {
	return s.Storage.DeleteMailForHuman(ctx, id)
}

func convertMails(human ...*storage.Mail) (result []*Mail) {
	if human == nil || human[0] == nil {
		return
	}

	for _, value := range human {
		result = append(result, &Mail{
			Id:          value.Id,
			HumanId:     value.HumanId,
			Description: value.Description,
			Mail:        value.Mail,
		})
	}

	return
}
