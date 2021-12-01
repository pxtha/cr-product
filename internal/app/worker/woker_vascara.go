package worker

import (
	"cr-product/internal/app/model"
	"cr-product/internal/utils"
	"encoding/json"
	"io/ioutil"
	"strings"

	"github.com/gocolly/colly"
)

func GetProductVascara(URL string, cate_id string, vendorid string) error {
	c := colly.NewCollector(
		colly.AllowedDomains("www.vascara.com"),
	)
	pr := model.RawProduct{}
	var str string
	c.OnHTML("div.product-info", func(h *colly.HTMLElement) {
		pr.Title = h.ChildText("h1.title-product")
		pr.CateId = cate_id
		pr.Link = URL
		pr.VendorId = vendorid
		h.ForEach("ul.list-oppr > li > span", func(_ int, h *colly.HTMLElement) {
			str += h.Text + "|"
		})
		tmp := strings.Split(str, "|")
		if len(tmp) > 0 {
			pr.SKU = tmp[1]
			str = strings.Replace(str, tmp[0]+"|", "", -1)
			str = strings.Replace(str, tmp[1]+"|", "", -1)
			str = strings.Replace(str, tmp[2]+"|", "", -1)
			str = strings.Replace(str, tmp[3]+"|", "", -1)
			pr.Detail = str
		}
		categoryJson1, err := json.MarshalIndent(pr, "", "   ")
		utils.FailOnError(err, "", "")
		err = ioutil.WriteFile("pr.json", categoryJson1, 0644)
		utils.FailOnError(err, "", "")
	})

	c.Visit(URL)
	return nil
}
