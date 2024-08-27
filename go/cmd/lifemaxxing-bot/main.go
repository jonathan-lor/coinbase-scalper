package main

import (
	"fmt"

	"github.com/joho/godotenv"

	"coinbase-scalper/internal/operations"
	"coinbase-scalper/internal/utils"
	"sync"
)

const (
	profitTarget   = 0.00025 // .005%
	stopLoss       = 0.00025 // .0025%
	initialBalance = 10000
)

var (
	lastPrice  float64
	position   bool
	mu         sync.Mutex
	balance    float64
	btcHolding float64
	entryPrice float64
)

func main() {
	balance = initialBalance

	err := godotenv.Load()
	if err != nil {
		utils.LogAndPrintError("Error loading .env file: %v", err)
	}

	for {
		initialPosition := position
		price, _, _, _ := operations.GetAssetDetails("BTC-USD")

		mu.Lock() // not actually used rn because no concurrency but putting here for later

		if !position && (lastPrice > 0 && price < lastPrice*(1-stopLoss)) {
			operations.SimulateBuy(price, &balance, &btcHolding, &entryPrice, &position)
		} else if position && (price >= entryPrice*(1+profitTarget) || price <= entryPrice*(1-stopLoss)) {
			operations.SimulateSell(price, &balance, &btcHolding, &entryPrice, &position)
		}

		lastPrice = price
		mu.Unlock() // see above mu.Lock()

		if initialPosition != position {
			fmt.Printf("Balance: $%.2f\n", balance)
			fmt.Printf("BTC Holding: %.8f\n", btcHolding)
			fmt.Println()
		}
	}
}
