package worker

import (
	"cr-product/internal/app/model"
	"cr-product/internal/utils"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
	"github.com/tidwall/gjson"
)

func (w *Worker) GetProductVascara(URL string, cate_id string, vendorid string) error {
	c := colly.NewCollector(
		colly.AllowedDomains("www.vascara.com"),
	)
	pr := model.RawProduct{}
	prVariant := model.Variant{}
	var str string
	var objmap map[string]json.RawMessage
	var err error

	c.OnHTML("body > div > div.page-content > div:nth-child(2) > div.container", func(h *colly.HTMLElement) {
		prId := h.ChildAttr("#productId", "value")
		pr.EcProductId = h.ChildAttr("#productCode", "value")
		prVariant.Link = URL
		err = json.Unmarshal([]byte(gjson.Get(w.GetStock(prId, pr.EcProductId, URL), "size").String()), &objmap)
	})

	c.OnHTML("div.product-info", func(h *colly.HTMLElement) {
		pr.Title = h.ChildText("h1.title-product")
		tmp := strings.Split(pr.Title, "-")
		pr.Title = tmp[0]
		pr.CateId = cate_id
		pr.VendorId = vendorid
		h.ForEach("ul.list-oppr > li > span", func(_ int, h *colly.HTMLElement) {
			str += h.Text + "|"
		})
		tmp = strings.Split(str, "|")
		if len(tmp) > 0 {
			str = strings.Replace(str, tmp[0]+"|", "", -1)
			str = strings.Replace(str, tmp[1]+"|", "", -1)
			str = strings.Replace(str, tmp[2]+"|", "", -1)
			str = strings.Replace(str, tmp[3]+"|", "", -1)
			str = strings.ReplaceAll(str, "|", ":")
			str = strings.ReplaceAll(str, " :", ";")
			pr.Detail = str
		}
	})

	c.OnHTML("body > div > div.page-content > div:nth-child(2) > div.container > div", func(e *colly.HTMLElement) {
		e.ForEach("div.group-images > div > a", func(_ int, h *colly.HTMLElement) {
			prVariant.Images = append(prVariant.Images, h.Attr("href"))
		})

		e.ForEach("li.lisize", func(_ int, h *colly.HTMLElement) {
			tmp := strings.Split(URL, "-")
			prVariant.Color = tmp[len(tmp)-2] + "-" + tmp[len(tmp)-1]
			prVariant.SKU = pr.EcProductId + "-" + prVariant.Color
			prVariant.Size = h.Text
			prVariant.Name = e.ChildText("h1.title-product")
			price, _ := strconv.Atoi(strings.Replace(e.ChildText("del > span.amount"), ".", "", -1))
			prVariant.Price = float64(price)
			price, _ = strconv.Atoi(strings.Replace(e.ChildText("ins > span.amount"), ".", "", -1))
			prVariant.DiscountPrice = float64(price)
			if prVariant.Price == 0 {
				prVariant.Price = prVariant.DiscountPrice
				prVariant.DiscountPrice = 0
			}
			prVariant.Stock, err = strconv.Atoi(string(objmap[prVariant.Size]))

			pr.Variant = append(pr.Variant, prVariant)
		})

	})

	c.Visit(URL)
	categoryJson1, err := json.MarshalIndent(pr, "", "   ")
	utils.FailOnError(err, "", "")
	err = ioutil.WriteFile("pr.json", categoryJson1, 0644)
	utils.FailOnError(err, "", "")
	if err != nil {
		return err
	}
	return nil
}

func (w *Worker) GetStock(productId string, productCode string, link string) string {
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
