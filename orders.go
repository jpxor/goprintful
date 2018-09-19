package goprintful

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func (c Client) DraftOrder(request OrderRequest) (OrderResponse, error) {
	const url = "https://api.printful.com/orders"

	//track costs
	var costs float32
	for _, item := range request.Items {
		request.RetailCosts.Subtotal += float32(item.Quantity) * item.RetailPrice
		costs += float32(item.Quantity) * item.Price
	}
	//assumes free shipping and no taxes
	request.RetailCosts.Total = request.RetailCosts.Subtotal
	costs += request.ShippingCost

	if costs > request.RetailCosts.Total {
		log.Println("Warning: costs exceed income:", request)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(request.toJSON()))
	req.Header.Set("Authorization", "Basic "+c.APIKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return OrderResponse{}, err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	response := OrderResponse{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Println(err)
		return response, err
	}
	return response, nil
}

func (c Client) Fulfill(orderID string) (OrderResponse, error) {
	const url = "https://api.printful.com/orders/{id}/confirm"
	requestURL := strings.Replace(url, "{id}", orderID, 1)

	req, err := http.NewRequest("POST", requestURL, nil)
	req.Header.Set("Authorization", "Basic "+c.APIKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return OrderResponse{}, err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	response := OrderResponse{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Println(err)
		return response, err
	}
	return response, nil
}
