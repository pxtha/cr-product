package model

import "github.com/google/uuid"

type (
	MessageReceive struct {
		VendorID uuid.UUID `json:"vendor_id"`
		Shop     string    `json:"shop_name"`
		CateID   uuid.UUID `json:"cate_id"`
		Link     string    `json:"link"`
	}

	RawProduct struct {
		EcProductID string    `json:"ec_product_id"`
		Title       string    `json:"title"`
		Description string    `json:"description"`
		CateID      string    `json:"category_id"`
		VendorID    string    `json:"vendor_id"`
		MadeIn      string    `json:"made_in"`
		Shop        string    `json:"shop_name"`
		Detail      string    `json:"detail"`
		Variant     []Variant `json:"variant"`
	}

	Variant struct {
		SKU           string   `json:"sku"`
		Link          string   `json:"link"`
		Price         float64  `json:"price"`
		DiscountPrice float64  `json:"discount_price"`
		Name          string   `json:"name"`
		Color         string   `json:"color"`
		Size          string   `json:"size"`
		Images        []string `json:"images"`
		Videos        []string `json:"videos"`
		IsAvailable   bool     `json:"is_available"`
		Stock         int      `json:"stock"`
	}

	HealthCheckResponse struct {
		ServiceName string `json:"service_name"`
		Version     string `json:"version"`
		HostName    string `json:"host_name"`
		TimeLife    string `json:"time_life"`
	}
)
