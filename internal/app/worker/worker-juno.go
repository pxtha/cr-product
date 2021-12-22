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
	"strings"

	"github.com/google/uuid"
)

func (w *Worker) GetProductJuno(vendorId uuid.UUID, categoryId uuid.UUID, URL string) error {
	var prodJson model.RawJunoJson
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

	res, err := http.Get(URL)
	if err != nil {
		panic(err.Error())
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err.Error())
	}
	json.Unmarshal(body, &prodJson)

	prodRaw.EcProductID = strconv.Itoa(prodJson.ID)
	prodRaw.Title = prodJson.Title
	prodRaw.Description = prodJson.Description
	prodRaw.CateID = categoryId
	prodRaw.VendorID = vendorId
	prodRaw.MadeIn = prodJson.Vendor
	prodRaw.Description = ""

	for _, v := range prodJson.Variants {
		prodVariant.SKU = v.Sku
		prodVariant.Link = "https://juno.vn" + prodJson.URL
		prodVariant.Price = strconv.Itoa(v.CompareAtPrice / 100)
		if prodVariant.Price == "0" {
			prodVariant.Price = strconv.Itoa(prodJson.Price / 100)
		}
		prodVariant.DiscountPrice = strconv.Itoa(v.Price / 100)
		prodVariant.Name = v.Title
		prodVariant.Color = v.Option2
		prodVariant.Size = v.Option1
		prodVariant.Images = nil
		for _, c := range prodJson.Images {
			variantColor := strings.ToLower(utils.NormalizeString(prodVariant.Color))
			variantColor = strings.ReplaceAll(variantColor, " ", "-")
			if strings.Contains(c, "/"+variantColor+"_") {
				prodVariant.Images = append(prodVariant.Images, c)
			}
		}
		prodVariant.Videos = nil
		prodVariant.Stock = v.InventoryQuantity

		prodRaw.Variant = append(prodRaw.Variant, prodVariant)
	}

	mgs, err := json.Marshal(prodRaw)
	if err != nil {
		return err
	}
	message := model.MessageSendDataload{
		Type: "product",
		Shop: "juno",
		Body: string(mgs),
	}
	err = rabbitmq.Produce(message, utils.Default_redelivered, utils.Exchange, utils.RouteKey_dataload, ch)
	if err != nil {
		return err
	}

	return nil
}
