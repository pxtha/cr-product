package main

import (
	"cr-product/internal/app/route"
	"cr-product/internal/app/worker"
	"cr-product/internal/utils"
	"sync"
)

func main() {
	utils.InitLogger()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		worker.Run()
	}()
	go func() {
		route.NewService()
	}()
	wg.Wait()
}
