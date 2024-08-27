from coinbase.rest import RESTClient
from json import dumps
import os

from dotenv import load_dotenv

load_dotenv()

keyName = os.getenv("KEY_NAME")
keySecret = os.getenv("KEY_SECRET")

client = RESTClient(api_key=keyName, api_secret=keySecret)

accounts = client.get_accounts()
print(dumps(accounts, indent=2))

order = client.market_order_buy(
    client_order_id="00000001",
    product_id="BTC-USD",
        quote_size="10"
    )

order_id = order["order_id"]

fills = client.get_fills(order_id=order_id)
print(dumps(fills, indent=2))