from coinbase.rest import RESTClient
from json import dumps
import os
import time
from decimal import Decimal
import sys

from dotenv import load_dotenv

load_dotenv()

keyName = os.getenv("KEY_NAME")
keySecret = os.getenv("KEY_SECRET")

client = RESTClient(api_key=keyName, api_secret=keySecret)

# Constants
PROFIT_TARGET = Decimal('0.00025')  # 0.025%
STOP_LOSS = Decimal('0.00025')      # 0.025%
TRADE_AMOUNT = Decimal('10')      # Amount in USD to trade
INITIAL_BALANCE = 10000

# Coinbase client setup
client = RESTClient(api_key=os.getenv('KEY_NAME'), api_secret=os.getenv('KEY_SECRET'))



class Bot:
    def __init__(self):
        self.mock_balance = Decimal("10000")
        self.position = ""
        self.entry_price = Decimal('0')
        self.last_price = Decimal('0')
        self.btc_holding = Decimal('0')

    def get_current_price(self):
        product = client.get_product('BTC-USD')
        return Decimal(product['price'])
    
    # simulating buy and sell functions to test strategy

    def simulate_buy(self, current_price):
        trade_amount = self.mock_balance * Decimal('0.5')
        btc_amount = trade_amount / current_price
        self.mock_balance -= trade_amount
        self.btc_holding += btc_amount
        self.entry_price = current_price
        self.position = "long"
        print(f"Simulated buy: {btc_amount:.8f} BTC at ${current_price:.2f}")

    def simulate_sell(self, current_price):
        sold_amount = self.btc_holding * current_price
        self.mock_balance += sold_amount
        self.btc_holding = 0;
        self.position = ""
        profit = sold_amount - (self.entry_price * self.btc_holding)
        print(f"Simulated sell: {self.btc_holding:.8f} BTC {current_price:.2f}, Profit: {profit:.2f}")

    
    def run(self):
        while True:
            try:
                current_price = self.get_current_price()
                print(f"Last BTC price: ${self.last_price:.2f}")
                print(f"Current BTC price: ${current_price:.2f}")
                print(f"Position: {self.position}")
                print(f"BTC Holding: {self.btc_holding:.8f}")
                print(f"Current (Mock) Balance: ${self.mock_balance:.2f}\n")

                if self.position == "":
                    if self.last_price > 0 and (current_price < self.last_price * (1 - STOP_LOSS)):
                        # stuff here
                        self.simulate_buy(current_price)
                elif self.position == "long":
                    if (current_price >= self.entry_price * (1 + PROFIT_TARGET)) or (current_price <= self.entry_price * (1 - STOP_LOSS)):
                        self.simulate_sell(current_price)
                
                self.last_price = current_price;
            except Exception as e:
                print(f"An error occurred: {e}")
                sys.Exit(1)



if __name__ == "__main__":
    bot = Bot()
    bot.run()
    # print(bot.get_current_price())
