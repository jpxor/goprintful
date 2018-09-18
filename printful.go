package goprintful

import (
	"encoding/json"
	"log"
	"os"
)

type Client struct {
	APIKey []byte
}

func NewClient() *Client {
	key := os.Getenv("PRINTFUL_APIKEY")
	return &Client{[]byte(key)}
}

func (r OrderRequest) toJSON() []byte {
	result, err := json.Marshal(&r)
	if err != nil {
		log.Println(err)
	}
	return []byte(result)
}

// API requests

type OrderRequest struct {
	ExternalID   string    `json:external_id`
	Shipping     string    `json:shipping`
	ShippingCost float32   `json:shipping_cost`
	Recipient    Recipient `json:recipient`
	Items        []Item    `json:items`
	RetailCosts  Costs     `json:retail_costs`
	Gift         GiftData  `json:gift`
}

// API Responses

type OrderResponse struct {
	Status        int    `json:status`
	StatusMessage string `json:status_message`
	Costs         Costs  `json:costs`
	RetailCosts   Costs  `json:retail_costs`
}

type ShipRateResponse struct {
	Code   int              `json:code`
	Result []ShippingMethod `json:result`
}

// Data structures

type Recipient struct {
	Name        string `json:name`
	Address1    string `json:address1`
	Address2    string `json:address2`
	City        string `json:city`
	StateCode   string `json:state_code`
	CountryCode string `json:country_code`
	Zip         string `json:zip`
}

type Item struct {
	VariantID      int            `json:variant_id`
	Quantity       int            `json:quantity`
	Price          float32        `json:price`
	RetailPrice    float32        `json:retail_price`
	Name           string         `json:name`
	Files          []PrintFile    `json:files`
	ShippingBudget ShippingBudget `json:shipping_budget`
}

type Costs struct {
	Currency     string  `json:currency`
	Subtotal     float32 `json:subtotal`
	Discount     float32 `json:discount`
	Shipping     float32 `json:shipping`
	Digitization float32 `json:digitization`
	Tax          float32 `json:tax`
	Vat          float32 `json:vat`
	Total        float32 `json:total`
}

type PrintFile struct {
	ID   int    `json:id`
	Type string `json:type`
	URL  string `json:url`
}

type GiftData struct {
	Subject string `json:subject`
	Message string `json:message`
}

type ShippingMethod struct {
	ID       string `json:id`
	Name     string `json:name`
	Rate     string `json:rate`
	Currency string `json:currency`
}

type ShippingBudget struct {
	First     float32 `json:first`
	Remainder float32 `json:remainder`
}
