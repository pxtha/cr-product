package worker

import (
	"cr-product/internal/app/model"
	"cr-product/internal/pkg/rabbitmq"
	"cr-product/internal/utils"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/danilopolani/fua"
	"github.com/gocolly/colly"
	"github.com/google/uuid"
	"github.com/streadway/amqp"
)

func (w *Worker) GetProductVascara(URL string, cate_id uuid.UUID, vendorid uuid.UUID, shop string, ch *amqp.Channel) error {
	c := colly.NewCollector(
		colly.AllowedDomains("www.vascara.com"),
	)
	pr := model.RawProduct{}
	prVariant := model.Variant{}
	var str string
	// var stock string
	var err error

	c.OnHTML("body > div > div.page-content > div:nth-child(2) > div.container", func(h *colly.HTMLElement) {
		//prId := h.ChildAttr("#productId", "value")
		pr.EcProductID = h.ChildAttr("#productCode", "value")
		prVariant.Link = URL
		//stock = w.GetStockVascara(prId, pr.EcProductID, URL)
	})

	c.OnHTML("div.product-info", func(h *colly.HTMLElement) {
		pr.Title = h.ChildText("h1.title-product")
		tmp := strings.Split(pr.Title, "-")
		pr.Title = tmp[0]
		pr.CateID = cate_id
		pr.VendorID = vendorid
		h.ForEach("ul.list-oppr > li > span", func(_ int, h *colly.HTMLElement) {
			str += h.Text + "|"
		})
		str = strings.ReplaceAll(str, "|", ":")
		str = strings.ReplaceAll(str, " :", ";")
		pr.Description = str
		pr.Shop = shop
	})

	c.OnHTML("body > div > div.page-content > div:nth-child(2) > div.container > div", func(e *colly.HTMLElement) {
		e.ForEach("div.group-images > div > a", func(_ int, h *colly.HTMLElement) {
			prVariant.Images = append(prVariant.Images, h.Attr("href"))
		})

		e.ForEach("li.licolor", func(_ int, h *colly.HTMLElement) {
			prVariant.Color = h.Attr("data-key")
			e.ForEach("li.lisize", func(_ int, h *colly.HTMLElement) {
				prVariant.Size = h.Text
				prVariant.SKU = pr.EcProductID + "-" + prVariant.Color + "-" + prVariant.Size
				prVariant.Name = e.ChildText("h1.title-product")
				prVariant.Price = e.ChildText("del > span.amount")
				prVariant.DiscountPrice = e.ChildText("ins > span.amount")
				if prVariant.Price == "" {
					prVariant.Price = prVariant.DiscountPrice
					prVariant.DiscountPrice = ""
				}
				//prVariant.Stock, err = strconv.Atoi(gjson.Get(stock, "size."+prVariant.Color+"."+strings.ToLower(prVariant.Size)).String())
				// if err != nil {
				// 	prVariant.Stock, err = strconv.Atoi(gjson.Get(stock, "size."+strings.ToLower(prVariant.Size)).String())
				// }

				if prVariant.Stock != 0 {
					prVariant.IsAvailable = true
				}
				pr.Variant = append(pr.Variant, prVariant)
			})
		})

		if pr.Variant == nil {
			tmp := strings.Split(URL, "-")
			prVariant.Color = tmp[len(tmp)-2] + "-" + tmp[len(tmp)-1]
			if prVariant.Size == "" {
				prVariant.SKU = pr.EcProductID + "-" + prVariant.Color
			} else if prVariant.Color == "" {
				prVariant.SKU = pr.EcProductID + "-" + prVariant.Size
			}
			prVariant.SKU = pr.EcProductID + "-" + prVariant.Color + "-" + prVariant.Size
			prVariant.Name = e.ChildText("h1.title-product")
			prVariant.Price = e.ChildText("del > span.amount")
			prVariant.DiscountPrice = e.ChildText("ins > span.amount")
			if prVariant.Price == "" {
				prVariant.Price = prVariant.DiscountPrice
				prVariant.DiscountPrice = ""
			}
			///prVariant.Stock, err = strconv.Atoi(gjson.Get(stock, "size."+strings.ToLower(prVariant.Size)).String())
			if prVariant.Stock != 0 {
				prVariant.IsAvailable = true
			}
			pr.Variant = append(pr.Variant, prVariant)
		}
	})

	c.Visit(URL)
	if err != nil {
		return err
	}
	message, err := json.Marshal(pr)
	if err != nil {
		return err
	}
	mgs := model.MessageSendDataload{
		Type: "product",
		Shop: "vascara",
		Body: string(message),
	}
	err = rabbitmq.Produce(mgs, utils.Default_redelivered, utils.Exchange, utils.RouteKey_dataload, ch)
	if err != nil {
		return err
	}
	return nil
}

func (w *Worker) GetStockVascara(productId string, productCode string, link string) string {
	url := "https://www.vascara.com/generalrealtime/getinventory"
	method := "POST"

	payload := strings.NewReader(fmt.Sprintf("fpid=%v&fpcode=%v", productId, productCode))

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		return ""
	}
	req.Header.Add("X-Requested-With", "XMLHttpRequest")
	req.Header.Add("referer", link)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Cookie", "SHASH=flbfrkc56ja78qr8bn3q5ia8p3; _t=l0o9ad03c6lop7q0hfi5jj6gmf")
	req.Header.Add("User-Agent", fua.Random("desktop"))

	res, err := client.Do(req)
	if err != nil {
		return ""
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return ""
	}
	return string(body)
}
