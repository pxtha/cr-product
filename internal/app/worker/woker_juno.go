package worker

import (
	"cr-product/internal/app/model"
	"cr-product/internal/utils"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

func GetProductJuno(vendorId uuid.UUID, categoryId uuid.UUID, URL string) error {
	var prodJson model.RawJunoJson
	var prodRaw model.RawProduct
	var prodVariant model.Variant

	res, err := http.Get(URL)
	if err != nil {
		panic(err.Error())
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err.Error())
	}
	json.Unmarshal(body,&prodJson)

	prodRaw.EcProductId = strconv.Itoa(prodJson.ID)
	prodRaw.Title = prodJson.Title
	prodRaw.Description = prodJson.Description
	prodRaw.CateId = categoryId.String()
	prodRaw.VendorId = vendorId.String()
	prodRaw.MadeIn = prodJson.Vendor
	prodRaw.Detail = ""

	for _,v := range prodJson.Variants {
		prodVariant.SKU = v.Sku
		prodVariant.Link = "https://juno.vn" + prodJson.URL
		prodVariant.Price = float64(v.CompareAtPrice / 100)
		if (prodVariant.Price == 0) {
			prodVariant.Price = float64(prodJson.Price / 100)
		}
		prodVariant.DiscountPrice = float64(v.Price / 100)
		prodVariant.Name = v.Title
		prodVariant.Color = v.Option2
		prodVariant.Size = v.Option1
		prodVariant.Images = nil
		for _,c := range prodJson.Images {
			variantColor := strings.ToLower(utils.NormalizeString(prodVariant.Color))
			variantColor = strings.ReplaceAll(variantColor, " ", "-")
			if strings.Contains(c,"/" + variantColor + "_") {
				prodVariant.Images = append(prodVariant.Images,c)
			}
		}
		prodVariant.Videos = nil
		prodVariant.Stock = v.InventoryQuantity
		
		prodRaw.Variant = append(prodRaw.Variant,prodVariant)
	}
	return nil
}