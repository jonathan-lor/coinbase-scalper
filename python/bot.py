from coinbase.rest import RESTClient
from json import dumps
import os
import time
from decimal import Decimal
import sys
import uuid
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
        # self.mock_balance = Decimal("10000")
        self.balance = self.get_current_usd_balance()
        self.entry_price = Decimal('0')
        self.last_price = Decimal('0')
        self.btc_holding = self.get_current_btc_wallet_balance_dec()
        self.position = "long" if self.btc_holding > Decimal("0.000001") else ""

        # print(self.mock_balance)
        print(self.balance)
        print(self.btc_holding)
        print(self.entry_price)
        print(self.last_price)
        print(self.position)

    def get_current_price(self):
        product = client.get_product('BTC-USD')
        return Decimal(product['price'])
    
    def get_current_usd_balance(self):
        return Decimal(client.get_account(os.getenv('USD_ACCOUNT_UUID'))['account']['available_balance']['value'])
    
    def get_current_btc_wallet_balance_dec(self):
        btcHolding = Decimal(client.get_account(os.getenv('BTC_WALLET_ACCOUNT_UUID'))['account']['available_balance']['value'])
        btcHoldingFormatted = f"{btcHolding:.8f}"
        return Decimal(btcHoldingFormatted)
    
    def get_current_btc_wallet_balance_str(self):
        btcHolding = Decimal(client.get_account(os.getenv('BTC_WALLET_ACCOUNT_UUID'))['account']['available_balance']['value'])
        btcHoldingFormatted = f"{btcHolding:.8f}"
        return btcHoldingFormatted
    
    def generate_client_order_id(self):
        timestamp = int(time.time() * 1000) # get current time in ms
        unique_id = str(uuid.uuid4().hex[:6])
        return f"order_{timestamp}_{unique_id}"
    
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

    # actual buy sell functions
    def buy(self, current_price):
        trade_amount = self.balance * Decimal('0.1') # money to spend in usd
        # btc_amount = trade_amount / current_price
        formatted_trade_amount = "{:.2f}".format(trade_amount)
        #print(formatted_trade_amount)
        #sys.exit(1)

        order = print(client.market_order_buy(
            client_order_id = self.generate_client_order_id(),
            product_id="BTC-USD",
            quote_size = formatted_trade_amount
        ))

        # order_id = order["order_id"]

        # fills = client.get_fills(order_id=order_id)
        # print(dumps(fills, indent=2))
        
        self.btc_holding = self.get_current_btc_wallet_balance_dec()
        self.position = "long"

        print("BOUGHT SOME BITCOIN CUHHHHHHHHHH")


    def sell(self, current_price):
        try:
            # Get the current BTC balance
            btc_account = client.get_account(os.getenv('BTC_WALLET_ACCOUNT_UUID'))
            btc_balance = btc_account['account']['available_balance']['value']

            btc_balance_formatted = "{:.4f}".format(Decimal(btc_balance))
        
            print(f"Attempting to sell all available BTC: {btc_balance}")
        
            order = client.market_order_sell(
                client_order_id=self.generate_client_order_id(),
                product_id="BTC-USD",
                base_size=btc_balance_formatted  # This will sell all available BTC
            )
        
            print("Sell order response:", order)
        
            if order.get('success', False):
                self.btc_holding = self.get_current_btc_wallet_balance_dec()
                self.position = ""
                print(f"SOLD ALL BITCOIN: {btc_balance} BTC")
            else:
                print("Failed to sell Bitcoin. Error:", order.get('error_response', {}).get('message', 'Unknown error'))
    
        except Exception as e:
            print(f"Error while placing sell order: {e}")

    
    def run(self):
        while True:
            if(self.get_current_usd_balance() < Decimal('150')): # account balance stop loss
                print("WE HIT USD BALANCE STOP LOSS!!!!")
                sys.exit(1)
            try:
                current_price = self.get_current_price() # get current BTC price
                print(f"Last BTC price: ${self.last_price:.2f}")
                print(f"Current BTC price: ${current_price:.2f}")
                print(f"Position: {self.position}")
                print(f"BTC Holding: {self.btc_holding:.8f}")
                # print(f"Current (Mock) Balance: ${self.mock_balance:.2f}")
                print(f"Current Balance: ${self.balance:.2f}\n")

                if self.position == "":
                    if self.last_price > 0 and (current_price < self.last_price * (1 - STOP_LOSS)):
                        # stuff here
                        #self.simulate_buy(current_price)
                        self.buy(current_price)
                elif self.position == "long":
                    if (current_price >= self.entry_price * (1 + PROFIT_TARGET)) or (current_price <= self.entry_price * (1 - STOP_LOSS)):
                        #self.simulate_sell(current_price)
                        self.sell(current_price)
                
                self.last_price = current_price;
            except Exception as e:
                print(f"An error occurred: {e}")
                sys.exit(1)



if __name__ == "__main__":

    bot = Bot()
    bot.run()
    # print(bot.get_current_btc_wallet_balance_str())
    # bot.sell(1)


    # WORKING BUY
    # order = print(client.market_order_buy(
    #     client_order_id=generate_client_order_id(),
    #     product_id="BTC-USD",
    #     quote_size="10"
    # ))

    # order_id = order["order_id"]

    # fills = client.get_fills(order_id=order_id)
    # print(dumps(fills, indent=2))

    # WORKING SELL 
    # order = print(client.market_order_sell(
    #     client_order_id=generate_client_order_id(),
    #     product_id="BTC-USD",
    #     base_size=btcHoldingFormatted
    # ))

    # order_id = order["order_id"]

    # fills = client.get_fills(order_id=order_id)
    # print(dumps(fills, indent=2))
