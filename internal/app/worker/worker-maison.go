package worker

import (
	"cr-product/conf"
	"cr-product/internal/app/model"
	"cr-product/internal/pkg/rabbitmq"
	"cr-product/internal/utils"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/google/uuid"
)

func (w *Worker) GetProductMaison(vendorId uuid.UUID, categoryId uuid.UUID, url string) error {
	var prodJson model.RawJunoMaiSonJson
	var prodRaw model.RawProduct
	var prodVariant model.Variant

	conf.SetEnv()
	var centerQueueConfig = rabbitmq.QueueConfig{
		Host:     conf.LoadEnv().RBHost,
		Port:     conf.LoadEnv().RBPort,
		Username: conf.LoadEnv().RBUser,
		Password: conf.LoadEnv().RBPass,
	}
	ch, _ := rabbitmq.GetRabbitmqConnChannel(centerQueueConfig)

	res, err := http.Get(url + ".js")
	if err != nil {
		panic(err.Error())
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err.Error())
	}
	json.Unmarshal(body, &prodJson)

	prodRaw.Title = prodJson.Title
	prodRaw.Description = prodJson.Description
	prodRaw.CateID = categoryId
	prodRaw.VendorID = vendorId
	prodRaw.MadeIn = prodJson.Vendor
	prodRaw.Shop = utils.MAISON

	for _, v := range prodJson.Variants {
		prodRaw.EcProductID = v.Option3
		prodVariant.SKU = v.Sku
		prodVariant.Link = "https://www.maisononline.vn" + prodJson.URL
		prodVariant.Price = strconv.Itoa(v.CompareAtPrice / 100)
		if prodVariant.Price == "0" {
			prodVariant.Price = strconv.Itoa(prodJson.Price / 100)
		}
		prodVariant.DiscountPrice = strconv.Itoa(v.Price / 100)
		prodVariant.Name = v.Title
		prodVariant.Color = v.Option1
		prodVariant.Size = v.Option2
		prodVariant.Images = prodJson.Images
		prodVariant.Videos = nil
		prodVariant.Stock = v.InventoryQuantity
		if prodVariant.Stock > 0 {
			prodVariant.IsAvailable = true
		}

		prodRaw.Variant = append(prodRaw.Variant, prodVariant)
	}

	mgs, err := json.Marshal(prodRaw)
	if err != nil {
		return err
	}
	message := model.MessageSendDataload{
		Type: "product",
		Shop: "maison",
		Body: string(mgs),
	}
	err = rabbitmq.Produce(message, utils.DefaultRedelivered, utils.Exchange, utils.RoutekeyDataload, ch)
	if err != nil {
		return err
	}

	return nil
}
