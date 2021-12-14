package model

import (
	"time"

	"github.com/google/uuid"
)

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
		Variant     []Variant `json:"variant"`
	}

	Variant struct {
		SKU           string   `json:"sku"`
		Link          string   `json:"link"`
		Price         string   `json:"price"`
		DiscountPrice string   `json:"discount_price"`
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

	RawJunoJson struct {
		Available            bool        `json:"available"`
		CompareAtPriceMax    int         `json:"compare_at_price_max"`
		CompareAtPriceMin    int         `json:"compare_at_price_min"`
		CompareAtPriceVaries bool        `json:"compare_at_price_varies"`
		CompareAtPrice       int         `json:"compare_at_price"`
		Content              interface{} `json:"content"`
		Description          string      `json:"description"`
		FeaturedImage        string      `json:"featured_image"`
		Handle               string      `json:"handle"`
		ID                   int         `json:"id"`
		Images               []string    `json:"images"`
		Options              []struct {
			Name      string   `json:"name"`
			Position  int      `json:"position"`
			ProductID int      `json:"product_id"`
			Values    []string `json:"values"`
		} `json:"options"`
		Price           int         `json:"price"`
		PriceMax        int         `json:"price_max"`
		PriceMin        int         `json:"price_min"`
		PriceVaries     bool        `json:"price_varies"`
		Tags            []string    `json:"tags"`
		TemplateSuffix  interface{} `json:"template_suffix"`
		Title           string      `json:"title"`
		Type            string      `json:"type"`
		URL             string      `json:"url"`
		Pagetitle       string      `json:"pagetitle"`
		Metadescription string      `json:"metadescription"`
		Variants        []struct {
			ID                   int         `json:"id"`
			Barcode              string      `json:"barcode"`
			Available            bool        `json:"available"`
			Price                int         `json:"price"`
			Sku                  string      `json:"sku"`
			Option1              string      `json:"option1"`
			Option2              string      `json:"option2"`
			Option3              string      `json:"option3"`
			Options              []string    `json:"options"`
			InventoryQuantity    int         `json:"inventory_quantity"`
			OldInventoryQuantity int         `json:"old_inventory_quantity"`
			Title                string      `json:"title"`
			Weight               int         `json:"weight"`
			CompareAtPrice       int         `json:"compare_at_price"`
			InventoryManagement  string      `json:"inventory_management"`
			InventoryPolicy      string      `json:"inventory_policy"`
			Selected             bool        `json:"selected"`
			URL                  interface{} `json:"url"`
			FeaturedImage        struct {
				ID         int    `json:"id"`
				CreatedAt  string `json:"created_at"`
				Position   int    `json:"position"`
				ProductID  int    `json:"product_id"`
				UpdatedAt  string `json:"updated_at"`
				Src        string `json:"src"`
				VariantIds []int  `json:"variant_ids"`
			} `json:"featured_image"`
		} `json:"variants"`
		Vendor            string    `json:"vendor"`
		PublishedAt       time.Time `json:"published_at"`
		CreatedAt         time.Time `json:"created_at"`
		NotAllowPromotion bool      `json:"not_allow_promotion"`
	}
)
