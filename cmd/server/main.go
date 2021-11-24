package main

import (
	"context"
	"os"

	"cr-product/conf"
	"cr-product/pkg/route"

	"gitlab.com/goxp/cloud0/logger"
)

const (
	APPNAME = "Crawl data"
)

func main() {
	conf.SetEnv()

	_ = os.Setenv("PORT", conf.LoadEnv().Port)

	logger.Init(APPNAME)

	app := route.NewService()
	ctx := context.Background()
	err := app.Start(ctx)
	if err != nil {
		logger.Tag("main").Error(err)
	}
	os.Clearenv()
}
