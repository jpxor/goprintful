package goprintful

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

func (c Client) LiveShipRates(request OrderRequest) []ShippingMethod {
	const url = "https://api.printful.com/shipping/rates"

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(request.toJSON()))
	req.Header.Set("Authorization", "Basic "+c.APIKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		log.Println("failed to get live shipping rates:", resp.Status)
	}

	response := ShipRateResponse{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Println(err)
	}
	return response.Result
}

// BestShipMethod is subjective, but this function returns the
// most expensive shipping method available which comes under the
// budgeted shipping expenses, which are defined per item. It assumes
// more expensive shipping is better/faster
func (c Client) BestShipMethod(order *OrderRequest) ShippingMethod {
	shipMethods := c.LiveShipRates(*order)

	var shippingBudget float32
	for _, item := range order.Items {
		if item.Quantity > 0 {
			shippingBudget += item.ShippingBudget.First
			shippingBudget += float32(item.Quantity-1) * item.ShippingBudget.Remainder
		}
	}
	budget := fmt.Sprintf("%.2f", shippingBudget)

	best := ShippingMethod{ID: "", Rate: "0"}
	min := shipMethods[0]

	for _, method := range shipMethods {
		if method.Rate < min.Rate {
			min = method
		}
		if method.Rate <= budget && method.Rate > best.Rate {
			best = method
		}
	}
	if best.ID == "" {
		log.Printf("Warning: shipping costs more than budgeted: items: %+v, dest: %+v %+v, cost: %+v", order.Items, order.Recipient.CountryCode, order.Recipient.StateCode, min.Rate)
		val, _ := strconv.ParseFloat(min.Rate, 32)
		order.ShippingCost = float32(val)
		return min
	}
	val, _ := strconv.ParseFloat(best.Rate, 32)
	order.ShippingCost = float32(val)
	return best
}
