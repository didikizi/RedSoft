package router

import (
	"log/slog"
	"strconv"

	"github.com/didikizi/RedSoft/iternal/service"
	utils "github.com/didikizi/RedSoft/packege"

	"github.com/labstack/echo/v4"
)

func (r *Router) GetHumanList(ectx echo.Context) error {
	result, err := r.Service.GetHumanList(ectx.Request().Context())
	if err != nil {
		slog.Info(utils.GetCallerInfo(), slog.String("err:", err.Error()))

		return ectx.JSON(400, err)
	}

	if len(result) == 0 {
		return ectx.NoContent(204)
	}

	return ectx.JSON(200, result)
}

func (r *Router) GetHuman(ectx echo.Context) error {
	cursor := ectx.Param("cursor")
	if cursor == "" {
		return ectx.NoContent(400)
	}

	var human []*service.Human

	humanId, err := strconv.Atoi(cursor)
	if err != nil {
		human, err = r.Service.GetHumanFromSurname(ectx.Request().Context(), cursor)
	} else {
		human, err = r.Service.GetHumanFromId(ectx.Request().Context(), humanId)
	}

	if err != nil {
		slog.Info(utils.GetCallerInfo(), slog.String("err:", err.Error()))

		return ectx.JSON(500, err)
	}

	if human == nil {
		slog.Info(utils.GetCallerInfo(), slog.String("err:", "human is nil"))

		return ectx.JSON(404, err)
	}

	return ectx.JSON(200, human)
}

func (r *Router) DeleteHuman(ectx echo.Context) error {
	humanId, err := strconv.Atoi(ectx.Param("id"))
	if err != nil {
		slog.Info(utils.GetCallerInfo(), slog.String("err:", err.Error()))

		return ectx.JSON(400, err)
	}

	tmp, err := r.Service.DeleteHuman(ectx.Request().Context(), humanId)
	if err != nil {
		slog.Info(utils.GetCallerInfo(), slog.String("err:", err.Error()))

		return ectx.JSON(500, err)
	}

	if tmp {
		slog.Info(utils.GetCallerInfo(), slog.String("err:", "human is nil"))

		return ectx.JSON(404, err)
	}

	return ectx.NoContent(204)
}

type PutHuman struct {
	Surname    string `json:"surname"`
	Name       string `json:"name"`
	Fatherland string `json:"fatherland"`
	Age        int    `json:"age"`
}

func (h *PutHuman) GetSurname() string    { return h.Surname }
func (h *PutHuman) GetName() string       { return h.Name }
func (h *PutHuman) GetFatherland() string { return h.Fatherland }
func (h *PutHuman) GetAge() int           { return h.Age }

func (r *Router) PutHuman(ectx echo.Context) error {
	human := new(PutHuman)

	err := ectx.Bind(human)
	if err != nil {
		slog.Info(utils.GetCallerInfo(), slog.String("err:", err.Error()))

		return ectx.JSON(400, err)
	}

	err = r.Service.PutHuman(ectx.Request().Context(), human)
	if err != nil {
		slog.Info(utils.GetCallerInfo(), slog.String("err:", err.Error()))

		return ectx.JSON(500, err)
	}

	return ectx.NoContent(204)
}

type PostHuman struct {
	Surname    string `json:"surname"`
	Name       string `json:"name"`
	Fatherland string `json:"fatherland"`
	National   string `json:"national"`
	Age        int    `json:"age"`
	Sex        string `json:"sex"`
}

func (h *PostHuman) GetSurname() string    { return h.Surname }
func (h *PostHuman) GetName() string       { return h.Name }
func (h *PostHuman) GetFatherland() string { return h.Fatherland }
func (h *PostHuman) GetNational() string   { return h.National }
func (h *PostHuman) GetAge() int           { return h.Age }
func (h *PostHuman) GetSex() string        { return h.Sex }

// TODO: Переделать на merge path
func (r *Router) PostHuman(ectx echo.Context) error {
	human := new(PostHuman)

	err := ectx.Bind(human)
	if err != nil {
		slog.Info(utils.GetCallerInfo(), slog.String("err:", err.Error()))

		return ectx.JSON(400, err)
	}

	humanId, err := strconv.Atoi(ectx.Param("id"))
	if err != nil {
		slog.Info(utils.GetCallerInfo(), slog.String("err:", err.Error()))

		return ectx.JSON(400, err)
	}

	err = r.Service.PostHuman(ectx.Request().Context(), human, humanId)
	if err != nil {
		slog.Info(utils.GetCallerInfo(), slog.String("err:", err.Error()))

		return ectx.JSON(500, err)
	}

	return ectx.NoContent(204)
}
