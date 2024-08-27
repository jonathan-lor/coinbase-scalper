package operations

import (
	"fmt"
)

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