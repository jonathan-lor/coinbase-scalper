package main

import (
	"fmt"
	"time"

	"github.com/joho/godotenv"

	"coinbase-scalper/internal/operations"
	"coinbase-scalper/internal/utils"
	"sync"
)

const (
	profitTarget = 0.00025 // .005%
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
		utils.LogAndPrintError("Error loading .env file: %v", err)
	}

	for {
		price, baseName, quoteName, ppc24h := operations.GetAssetDetails("BTC-USD")

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
        fmt.Printf("%s Price: %.2f %s\n", baseName, price, quoteName)
        fmt.Printf("24h Price Change: %s%%\n", ppc24h)
        fmt.Printf("Position: %s\n", position)
        fmt.Printf("Balance: $%.2f\n", balance)
        fmt.Printf("BTC Holding: %.8f\n", btcHolding)
        fmt.Println()
		// time.Sleep(1 * time.Second)
	}
}
