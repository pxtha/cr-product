package model

import "github.com/google/uuid"

type (
	MessageReceive struct {
		Vendor_ID uuid.UUID `json:"vendor_id"`
		Shop      string    `json:"shop_name"`
		Cate_ID   uuid.UUID `json:"cate_id"`
		Link      string    `json:"link"`
	}

	RawProduct struct {
<<<<<<< HEAD
		SKU           string `json:"sku"`
		Description   string `json:"description"`
		SlugName      string `json:"slug_name"`
		Link          string `json:"link"`
		EC_Product_ID string `json:"ec_product_id"`
		CateId        string `json:"category_id"`
		Title         string `json:"title"`
		VendorId      string `json:"vendor_id"`
		MadeIn        string `json:"made_in"`
		Detail        string `json:"detail"`
		Variant       []Variant
=======
		EcProductId string `json:"ec_product_id"`
		Title       string `json:"title"`
		Description string `json:"description"`
		CateId      string `json:"category_id"`
		VendorId    string `json:"vendor_id"`
		MadeIn      string `json:"made_in"`
		Detail      string `json:"detail"`
		Variant     []Variant
>>>>>>> develop
	}

	Variant struct {
		SKU           string   `json:"sku"`
		Link          string   `json:"link"`
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
		HostName    string `json:"host_name"`
		TimeLife    string `json:"time_life"`
	}
)
