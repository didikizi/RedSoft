package router

import (
	"log/slog"
	"regexp"
	"strconv"

	utils "github.com/didikizi/RedSoft/packege"
	"github.com/labstack/echo/v4"
)

func (r *Router) GetMailListForHuman(ectx echo.Context) error {
	id := ectx.Param("id")
	if id == "" {
		return ectx.NoContent(400)
	}

	humanId, err := strconv.Atoi(id)
	if err != nil {
		slog.Info(utils.GetCallerInfo(), slog.String("err:", err.Error()))

		return ectx.JSON(400, err)
	}

	human, err := r.Service.GetHumanFromId(ectx.Request().Context(), humanId)
	if err != nil {
		slog.Info(utils.GetCallerInfo(), slog.String("err:", err.Error()))

		return ectx.JSON(500, err)
	}

	if (human) == (nil) || len(human) != 1 {
		slog.Info(utils.GetCallerInfo(), slog.String("err:", "human not found"))

		return ectx.JSON(404, err)
	}

	mails, err := r.Service.GetMailListForHuman(ectx.Request().Context(), humanId)
	if err != nil {
		slog.Info(utils.GetCallerInfo(), slog.String("err:", err.Error()))

		return ectx.JSON(500, err)
	}

	if mails == nil || mails[0] == nil {
		slog.Info(utils.GetCallerInfo(), slog.String("err:", "mail not found"))

		return ectx.JSON(404, err)
	}

	return ectx.JSON(200, mails)
}

type PutMail struct {
	HumanId     int
	Mail        string `json:"mail"`
	Description string `json:"description"`
}

func (m *PutMail) GetHumanId() int        { return m.HumanId }
func (m *PutMail) GetMail() string        { return m.Mail }
func (m *PutMail) GetDescription() string { return m.Description }

func (r *Router) PutMailForHuman(ectx echo.Context) error {
	mail := new(PutMail)

	cursor := ectx.Param("id")
	if cursor == "" {
		return ectx.NoContent(400)
	}

	humanId, err := strconv.Atoi(cursor)
	if err != nil {
		slog.Info(utils.GetCallerInfo(), slog.String("err:", err.Error()))

		return ectx.JSON(400, err)
	}

	human, err := r.Service.GetHumanFromId(ectx.Request().Context(), humanId)
	if err != nil {
		slog.Info(utils.GetCallerInfo(), slog.String("err:", err.Error()))

		return ectx.JSON(500, err)
	}

	if (human) == (nil) || len(human) != 1 {
		slog.Info(utils.GetCallerInfo(), slog.String("err:", "human not found"))

		return ectx.JSON(404, err)
	}

	err = ectx.Bind(mail)
	if err != nil {
		slog.Info(utils.GetCallerInfo(), slog.String("err:", err.Error()))

		return ectx.JSON(400, err)
	}

	if !isValidEmail(mail.Mail) {
		slog.Info(utils.GetCallerInfo(), slog.String("err:", "mail not valid"))

		return ectx.JSON(400, err)
	}

	mail.HumanId = humanId

	err = r.Service.PutMailForHuman(ectx.Request().Context(), mail)
	if err != nil {
		slog.Info(utils.GetCallerInfo(), slog.String("err:", err.Error()))

		return ectx.JSON(500, err)
	}

	return ectx.NoContent(204)
}

func (r *Router) DeleteMailForHuman(ectx echo.Context) error {
	idHuman := ectx.Param("id_human")
	if idHuman == "" {
		return ectx.NoContent(400)
	}

	humanId, err := strconv.Atoi(idHuman)
	if err != nil {
		slog.Info(utils.GetCallerInfo(), slog.String("err:", err.Error()))

		return ectx.JSON(400, err)
	}

	human, err := r.Service.GetHumanFromId(ectx.Request().Context(), humanId)
	if err != nil {
		slog.Info(utils.GetCallerInfo(), slog.String("err:", err.Error()))

		return ectx.JSON(500, err)
	}

	if (human) == (nil) || len(human) != 1 {
		slog.Info(utils.GetCallerInfo(), slog.String("err:", "human not found"))

		return ectx.JSON(404, err)
	}

	idMail := ectx.Param("id_mail")
	if idMail == "" {
		return ectx.NoContent(400)
	}

	mailId, err := strconv.Atoi(idMail)
	if err != nil {
		slog.Info(utils.GetCallerInfo(), slog.String("err:", err.Error()))

		return ectx.JSON(400, err)
	}

	enable, err := r.Service.DeleteMailForHuman(ectx.Request().Context(), mailId)
	if err != nil {
		slog.Info(utils.GetCallerInfo(), slog.String("err:", err.Error()))

		return ectx.JSON(500, err)
	}

	if !enable {
		slog.Info(utils.GetCallerInfo(), slog.String("err:", "mail not found"))

		return ectx.JSON(404, err)
	}

	return ectx.NoContent(204)
}

func isValidEmail(email string) bool {
	regex := `^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`
	re := regexp.MustCompile(regex)
	return re.MatchString(email)
}
