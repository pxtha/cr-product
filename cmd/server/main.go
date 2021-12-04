package main

import (
	"context"
	"cr-product/conf"
	"cr-product/internal/app/route"
	"cr-product/internal/app/worker"
	"cr-product/internal/utils"
	"os"
	"sync"

	"gitlab.com/goxp/cloud0/logger"
)

func main() {
	utils.InitLogger()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		conf.SetEnv()

		_ = os.Setenv("PORT", conf.LoadEnv().Port)
		logger.Init(utils.APPNAME)

		app := route.NewService()
		ctx := context.Background()
		err := app.Start(ctx)
		if err != nil {
			logger.Tag("main").Error(err)
		}
		os.Clearenv()
	}()
	go func() {
		w := worker.NewWorker()
		w.Run()
	}()
	wg.Wait()
}
