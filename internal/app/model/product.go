package model

type (
	MessageReceive struct {
		ID   string `json:"id_trigger"`
		Shop string `json:"shop_name"`
		URL  string `json:"url"`
	}

	RawProduct struct {
		SKU         string `json:"sku"`
		Description string `json:"description"`
		Slugname    string `json:"slug_name"`
		Link        string `json:"link"`
		EcProductId string `json:"ec_product_id"`
		CateId      string `json:"category_id"`
		Title       string `json:"title"`
		VendorId    string `json:"vendor_id"`
		MadeIn      string `json:"made_in"`
		Detail      string `json:"detail"`
		Variant     Variant
	}

	Variant struct {
		Id            string   `json:"id"`
		Price         float64  `json:"Price"`
		DiscountPrice float64  `json:"discount_price"`
		Name          string   `json:"name"`
		Color         string   `json:"color"`
		Size          string   `json:"size"`
		Images        []string `json:"images"`
		Videos        []string `json:"videos"`
		Stock         int      `json:"stock"`
	}

	HealthCheckResponse struct {
		ServiceName string `json:"service_name"`
		Version     string `json:"version"`
		Hostname    string `json:"hostname"`
		Timelife    string `json:"time_life"`
	}
)
