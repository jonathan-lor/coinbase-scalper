package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"coinbase-scalper/internal/jwtgenerator"
	"coinbase-scalper/internal/operations"
	"coinbase-scalper/internal/utils"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		utils.LogAndPrintError("Error loading .env file: %v", err)
	}

	keyName := os.Getenv("KEY_NAME")
	keySecret := os.Getenv("KEY_SECRET")

	// getAccounts(keyName, keySecret)
	//getShib(keyName, keySecret)
	getPerms(keyName, keySecret)
	getAccounts(keyName, keySecret)
	fmt.Println(operations.GetAssetDetails("BTC-USD"))
	price, _, _, _ := operations.GetAssetDetails("BTC-USD")
	err = placeBuyOrder(price, .000001)
	if err != nil {
		utils.LogAndPrintError("Error placing buy order: %v", err)
	}

}

const (
    baseURL = "https://api.coinbase.com"
    orderEndpoint = "/api/v3/brokerage/orders"
)

type OrderRequest struct {
    ClientOrderID string `json:"client_order_id"`
    ProductID     string `json:"product_id"`
    Side          string `json:"side"`
    OrderConfiguration OrderConfiguration `json:"order_configuration"`
}

type OrderConfiguration struct {
    MarketMarketIoc MarketMarketIoc `json:"market_market_ioc"`
}

type MarketMarketIoc struct {
    QuoteSize string `json:"quote_size"`
    BaseSize  string `json:"base_size"`
}

func placeBuyOrder(price, amount float64) error {
    return placeOrder("BUY", price, amount)
}

func placeSellOrder(price, amount float64) error {
    return placeOrder("SELL", price, amount)
}

func placeOrder(side string, price, amount float64) error {
    clientOrderID := fmt.Sprintf("order_%d", time.Now().UnixNano())
    
    var orderConfig OrderConfiguration
    if side == "BUY" {
        orderConfig = OrderConfiguration{
            MarketMarketIoc: MarketMarketIoc{
                QuoteSize: strconv.FormatFloat(amount, 'f', 2, 64),
            },
        }
    } else {
        orderConfig = OrderConfiguration{
            MarketMarketIoc: MarketMarketIoc{
                BaseSize: strconv.FormatFloat(amount/price, 'f', 8, 64),
            },
        }
    }

    orderRequest := OrderRequest{
        ClientOrderID: clientOrderID,
        ProductID:     "BTC-USD",
        Side:          side,
        OrderConfiguration: orderConfig,
    }

    jsonData, err := json.Marshal(orderRequest)
    if err != nil {
        return fmt.Errorf("error marshalling order request: %v", err)
    }

    uri := fmt.Sprintf("POST %s%s", baseURL, orderEndpoint)
    jwt, err := jwtgenerator.BuildJWT(uri, os.Getenv("KEY_NAME"), os.Getenv("KEY_SECRET"))
    if err != nil {
        return fmt.Errorf("error generating JWT: %v", err)
    }

    req, err := http.NewRequest("POST", baseURL+orderEndpoint, bytes.NewBuffer(jsonData))
    if err != nil {
        return fmt.Errorf("error creating request: %v", err)
    }

    req.Header.Set("Authorization", "Bearer "+jwt)
    req.Header.Set("Content-Type", "application/json")

    // Print request details for debugging
    fmt.Printf("Request URL: %s\n", req.URL)
    fmt.Printf("Request Headers: %v\n", req.Header)
    fmt.Printf("Request Body: %s\n", string(jsonData))

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return fmt.Errorf("error sending request: %v", err)
    }
    defer resp.Body.Close()

    // Read and print response body
    body, _ := io.ReadAll(resp.Body)
    fmt.Printf("Response Status: %s\n", resp.Status)
    fmt.Printf("Response Body: %s\n", string(body))

    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
    }
	return nil
}

func getShib(keyName, keySecret string) {
	uri := fmt.Sprintf("%s %s%s", "GET", "api.coinbase.com", "/api/v3/brokerage/products/SHIB-USD")
	url := "https://api.coinbase.com/api/v3/brokerage/products/SHIB-USD"
	method := "GET"

	jwt, err := jwtgenerator.BuildJWT(uri, keyName, keySecret)
	if err != nil {
		utils.LogAndPrintError("Error generating JWT: %v", err)
	}

	client := &http.Client {
	}
	req, err := http.NewRequest(method, url, nil)
  
	if err != nil {
	  fmt.Println(err)
	  return
	}
	req.Header.Set("Authorization", "Bearer " + jwt)
	req.Header.Add("Content-Type", "application/json")
  
	res, err := client.Do(req)
	if err != nil {
	  fmt.Println(err)
	  return
	}
	defer res.Body.Close()
  
	body, err := io.ReadAll(res.Body)
	if err != nil {
	  fmt.Println(err)
	  return
	}
	fmt.Println(string(body))
}

func getAccounts(keyName, keySecret string) {
	uri := fmt.Sprintf("%s %s%s", "GET", "api.coinbase.com", "/api/v3/brokerage/accounts")
	url := "https://api.coinbase.com/api/v3/brokerage/accounts"
  	method := "GET"

	jwt, err := jwtgenerator.BuildJWT(uri, keyName, keySecret)
	if err != nil {
		utils.LogAndPrintError("Error generating JWT: %v", err)
	}

	client := &http.Client {
	}
	req, err := http.NewRequest(method, url, nil)
  
	if err != nil {
	  fmt.Println(err)
	  return
	}
	req.Header.Set("Authorization", "Bearer " + jwt)
	req.Header.Add("Content-Type", "application/json")
  
	res, err := client.Do(req)
	if err != nil {
	  fmt.Println(err)
	  return
	}
	defer res.Body.Close()
  
	body, err := io.ReadAll(res.Body)
	if err != nil {
	  fmt.Println(err)
	  return
	}
	fmt.Println(string(body))
}

func getPerms(keyName, keySecret string) {
	url := "https://api.coinbase.com/api/v3/brokerage/key_permissions"
	uri := fmt.Sprintf("%s %s%s", "GET", "api.coinbase.com", "/api/v3/brokerage/key_permissions/")
	method := "GET"
  
	client := &http.Client {
	}
	req, err := http.NewRequest(method, url, nil)
  
	if err != nil {
	  fmt.Println(err)
	  return
	}

	jwt, err := jwtgenerator.BuildJWT(uri, keyName, keySecret)
	if err != nil {
		utils.LogAndPrintError("Error generating JWT: %v", err)
	}
	req.Header.Set("Authorization", "Bearer " + jwt)
	req.Header.Add("Content-Type", "application/json")
  
	res, err := client.Do(req)
	if err != nil {
	  fmt.Println(err)
	  return
	}
	defer res.Body.Close()
  
	body, err := io.ReadAll(res.Body)
	if err != nil {
	  fmt.Println(err)
	  return
	}
	fmt.Println(string(body))
}