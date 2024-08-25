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
	"coinbase-scalper/internal/utils"
)



func main() {
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
			utils.LogAndPrintError("Error parsing response to model %v", err)
		}
		
		// fmt.Printf("Product ID: %s\n", product.ProductID)
		fmt.Printf("Current Time: %s\n", time.Now().Format("2006-01-02 15:04:05"))
		fmt.Printf("%s Price: %s %ss\n", product.BaseName, product.Price, product.QuoteName)
		fmt.Printf("24h Price Change: %s%%\n", product.PricePercentageChange24h)
		fmt.Println()
		// time.Sleep(1 * time.Second)
	}
}
