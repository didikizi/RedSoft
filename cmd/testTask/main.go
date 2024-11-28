package main

import (
	"log/slog"

	"github.com/didikizi/RedSoft/iternal/app"
	utils "github.com/didikizi/RedSoft/packege"
)

func main() {
	err := app.Start()
	if err != nil {
		slog.Info(utils.GetCallerInfo(), slog.String("init err:", err.Error()))
	}
}
