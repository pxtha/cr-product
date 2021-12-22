package worker

import (
	"bytes"
	"cr-product/internal/app/model"
	"cr-product/internal/pkg/rabbitmq"
	"cr-product/internal/utils"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/streadway/amqp"
)

func (w *Worker) GetProductHP(job *model.MessageReceive, ch *amqp.Channel) error {
	// Load the HTML document
	html, err := w.GetHttpHtmlContent(job.Link)
	if err != nil {
		return err
	}
	if html == "" {
		return errors.New("can't get html document")
	}

	dom, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		log.Fatal(err)
	}

	product_detail := &model.RawProduct{}
	related_product := &model.Variant{}

	//get product image
	dom.Find("img.fotorama__img").Each(func(i int, s *goquery.Selection) {
		image, _ := s.Attr("src")
		related_product.Images = append(related_product.Images, image)
	})
	dom.Find("div.value.description-item > div.block-size > p > img").Each(func(i int, s *goquery.Selection) {
		image, _ := s.Attr("src")
		related_product.Images = append(related_product.Images, image)
	})
	product_detail.CateID = job.CateID
	product_detail.VendorID = job.VendorID
	product_detail.MadeIn = ""
	product_detail.EcProductID = dom.Find("table.data.table.additional-attributes> tbody > tr > td.col.data[data-th='Model']").Text()
	product_detail.Description = dom.Find("div.value.description-item > p").Text()

	related_product.DiscountPrice = utils.FMPrice(dom.Find("span.normal-price > span.price-container.price-final_price.tax.weee > span.price-wrapper > span.price").First().Text())
	related_product.Price = utils.FMPrice(dom.Find("span.old-price.sly-old-price > span.price-container.price-final_price.tax.weee > span.price-wrapper > span.price").First().Text())
	related_product.SKU = strings.Replace(dom.Find("div.product.attribute.sku > div.value").First().Text(), " ", "", -1)
	related_product.Name = dom.Find("h1.page-title>span.base").Text()
	related_product.IsAvailable = true
	related_product.Link = job.Link

	//product_detail.CATEGORY = dom.Find("table.data.table.additional-attributes> tbody > tr > td.col.data[data-th='Giới tính']").Text()
	//product_detail.SEASON = dom.Find("table.data.table.additional-attributes> tbody > tr >td.col.data[data-th='Season']").Text()
	//product_detail.VENDOR = dom.Find("table.data.table.additional-attributes> tbody > tr > td.col.data[data-th='Thương hiệu']").Text()
	//product_detail.ID, _ = dom.Find("div.price-box.price-final_price").Attr("data-product-id")

	item, err := utils.Split(related_product.Name, product_detail.EcProductID)
	if err != nil {
		return err
	}

	related_product.Color = strings.TrimSpace(item[1])
	product_detail.Title = strings.Trim(related_product.Name, item[1])

	if dom.Find("div").HasClass("stock unavailable") {
		related_product.IsAvailable = false
		product_detail.Variant = append(product_detail.Variant, *related_product)
	}

	dom.Find("div.swatch-option.text").Each(func(i int, s *goquery.Selection) {
		related_product.Size = s.Text()
		fmt.Println(related_product.Size)
		product_detail.Variant = append(product_detail.Variant, *related_product)
	})

	err = rabbitmq.Produce(product_detail, utils.Default_redelivered, utils.Exchange, utils.RouteKey_dataload, ch)
	if err != nil {
		utils.FailOnError(err, "Failed to publish a message to the queue", "")
	}

	/* 	productJson, err := json.MarshalIndent(product_detail, "", "   ")
	   	utils.CheckError(err)
	   	err = ioutil.WriteFile("data.json", productJson, 0644)
	   	utils.CheckError(err) */
	return nil
}

func (w *Worker) GetHttpHtmlContent(link string) (string, error) {

	url := "http://188.166.220.131:1003/ferret/"
	method := "POST"

	msg := fmt.Sprintf(`{
		"text": "LET doc = DOCUMENT(@url, {driver: 'cdp'}) WAIT_ELEMENT(doc, '.modals-wrapper', 10000) RETURN doc",
		"params": {
			"url": "%v"
		}
	}`, link)
	payload := strings.NewReader(msg)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		return "", err
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(res.Body)
	newStr := buf.String()

	newHTML := strings.Trim(newStr, "\"")
	newHTML = strings.Replace(newHTML, `\n`, "", -1)
	newHTML = strings.Replace(newHTML, "\\", "", -1)

	return newHTML, nil
}
