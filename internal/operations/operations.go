package operations

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"coinbase-scalper/internal/jwtgenerator"
	"coinbase-scalper/internal/models"
	"coinbase-scalper/internal/utils"
)


func GetAssetDetails(name string) (price float64, baseName string, quoteName string, ppc24h string){
	keyName := os.Getenv("KEY_NAME")
	keySecret := os.Getenv("KEY_SECRET")
	uri := fmt.Sprintf("%s %s%s%s", "GET", "api.coinbase.com", "/api/v3/brokerage/products/", name)
	url := fmt.Sprintf("%s%s", "https://api.coinbase.com/api/v3/brokerage/products/", name)
	method := "GET"

	jwt, err := jwtgenerator.BuildJWT(uri, keyName, keySecret)
	if err != nil {
		utils.LogAndPrintError("Error generating JWT: %v", err)
	}

	req, err := http.NewRequest(method, url, nil)
	
	if err != nil {
		utils.LogAndPrintError("Error creating request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer " + jwt)
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}

	res, err := client.Do(req)
	if err!= nil {
		utils.LogAndPrintError("Error in the response: %v", err)
	}

	if res.StatusCode != http.StatusOK {
		utils.LogAndPrintError("Unexpected status code: %d", err)
	}

	var product models.Product
	err = json.NewDecoder(res.Body).Decode(&product)
	if err != nil {
		utils.LogAndPrintError("Error parsing response to model: %v", err)
	}

	price, err = strconv.ParseFloat(product.Price, 64)
	if err != nil {
		utils.LogAndPrintError("Error parsing price: %v", err)
	}

	baseName = product.BaseName
	quoteName = product.QuoteName
	ppc24h = product.PricePercentageChange24h
	return
}

// simulate buy and sell orders for now

func SimulateBuy(
	price float64,
	balance *float64,
	btcHolding *float64,
	entryPrice *float64,
	position *string,
	) {
		tradeAmount := *balance * 0.5
		btcAmount := tradeAmount / price
		*balance -= tradeAmount
		*btcHolding += btcAmount
		*entryPrice = price
		*position = "long"
		fmt.Printf("Simulated buy: %.8f BTC at $%.2f\n", btcAmount, price)
}

func SimulateSell(
	price float64,
	balance *float64,
	btcHolding *float64,
	entryPrice *float64,
	position *string,
	) {
		soldAmount := *btcHolding * price
		*balance += soldAmount
		*btcHolding = 0
		*position = ""
		profit := soldAmount - (*entryPrice * *btcHolding)
		fmt.Printf("Simulated sell: %.8f BTC %.2f, Profit: %.2f\n", *btcHolding, price, profit)
}