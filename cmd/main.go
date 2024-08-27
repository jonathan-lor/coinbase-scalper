package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"

	"coinbase-scalper/internal/jwtgenerator"
	"coinbase-scalper/internal/models"
	"coinbase-scalper/internal/operations"
	"coinbase-scalper/internal/utils"
	"strconv"
	"sync"
)

const (
	profitTarget = 0.0005 // .005%
	stopLoss = 0.00025 // .0025%
	initialBalance = 10000
)

var (
	lastPrice float64
	position string
	mu sync.Mutex // here for later
	balance float64
	btcHolding float64
	entryPrice float64
)


func main() {
	balance = initialBalance;

	err := godotenv.Load()
	if err != nil {
		fmt.Printf("Error loading .env file: %v", err)
		log.Fatalf("Error loading .env file: %v", err)
	}

	keyName := os.Getenv("KEY_NAME")
	keySecret := os.Getenv("KEY_SECRET")
	uri := fmt.Sprintf("%s %s%s", "GET", "api.coinbase.com", "/api/v3/brokerage/products/BTC-USD")

	url := "https://api.coinbase.com/api/v3/brokerage/products/BTC-USD"
	method := "GET"

	for {
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
		// defer res.Body.Close()
	
		if res.StatusCode != http.StatusOK {
			utils.LogAndPrintError("Unexpected status code: %d", err)
		}

		var product models.Product
		err = json.NewDecoder(res.Body).Decode(&product)
		if err != nil {
			utils.LogAndPrintError("Error parsing response to model: %v", err)
		}

		price, err := strconv.ParseFloat(product.Price, 64)
		if err != nil {
			utils.LogAndPrintError("Error parsing price: %v", err)
		}

		mu.Lock() // not actually used rn because no concurrency but putting here for later

		if position == "" {
			if lastPrice > 0 && price < lastPrice * (1 - stopLoss) {
				operations.SimulateBuy(price, &balance, &btcHolding, &entryPrice, &position)
			}
		} else if position == "long" {
			if price >= entryPrice*(1 + profitTarget) || price <= entryPrice * (1 - stopLoss) {
				operations.SimulateSell(price, &balance, &btcHolding, &entryPrice, &position)
			}
		}
		lastPrice = price
		mu.Unlock() // see above mu.Lock()

		
		// fmt.Printf("Product ID: %s\n", product.ProductID)
        fmt.Printf("Current Time: %s\n", time.Now().Format("2006-01-02 15:04:05"))
        fmt.Printf("%s Price: %s %s\n", product.BaseName, product.Price, product.QuoteName)
        fmt.Printf("24h Price Change: %s%%\n", product.PricePercentageChange24h)
        fmt.Printf("Position: %s\n", position)
        fmt.Printf("Balance: $%.2f\n", balance)
        fmt.Printf("BTC Holding: %.8f\n", btcHolding)
        fmt.Println()
		// time.Sleep(1 * time.Second)
	}
}
