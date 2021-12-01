package main

import (
	"cr-product/internal/app/worker"
	"cr-product/internal/utils"
	"sync"
)

func main() {
	utils.InitLogger()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		//worker.Run()
		err := worker.GetProductVascara("https://www.vascara.com/giay-cao-got/giay-du-tiec-quai-anh-bac-sdn-0683-mau-trang", "job.Cate_ID", "vendorId")
		if err != nil {
			utils.Log(utils.ERROR_LOG, "Error: ", err, "")
		}
	}()
	go func() {
		//route.NewService()
	}()
	wg.Wait()
}
