package main

import (
	"github.com/RozmiDan/gameReviewHub/internal/app"
	"github.com/RozmiDan/gameReviewHub/internal/config"
)

func main() {
	cnfg := config.MustLoad()

	app.Run(cnfg)
}
