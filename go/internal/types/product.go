package models

// Product represents the structure of the JSON response from API url: https://api.coinbase.com/api/v3/brokerage/products/BTC-USD
type Product struct {
    ProductID                string  `json:"product_id"`
    Price                    string  `json:"price"`
    PricePercentageChange24h string  `json:"price_percentage_change_24h"`
    Volume24h                string  `json:"volume_24h"`
    VolumePercentageChange24h string  `json:"volume_percentage_change_24h"`
    BaseIncrement            string  `json:"base_increment"`
    QuoteIncrement           string  `json:"quote_increment"`
    QuoteMinSize             string  `json:"quote_min_size"`
    QuoteMaxSize             string  `json:"quote_max_size"`
    BaseMinSize              string  `json:"base_min_size"`
    BaseMaxSize              string  `json:"base_max_size"`
    BaseName                 string  `json:"base_name"`
    QuoteName                string  `json:"quote_name"`
    Watched                  bool    `json:"watched"`
    IsDisabled               bool    `json:"is_disabled"`
    New                      bool    `json:"new"`
    Status                   string  `json:"status"`
    CancelOnly               bool    `json:"cancel_only"`
    LimitOnly                bool    `json:"limit_only"`
    PostOnly                 bool    `json:"post_only"`
    TradingDisabled          bool    `json:"trading_disabled"`
    AuctionMode              bool    `json:"auction_mode"`
    ProductType              string  `json:"product_type"`
    QuoteCurrencyID          string  `json:"quote_currency_id"`
    BaseCurrencyID           string  `json:"base_currency_id"`
    FcmTradingSessionDetails interface{} `json:"fcm_trading_session_details"`
    MidMarketPrice           string  `json:"mid_market_price"`
    Alias                    string  `json:"alias"`
    AliasTo                  []string `json:"alias_to"`
    BaseDisplaySymbol        string  `json:"base_display_symbol"`
    QuoteDisplaySymbol       string  `json:"quote_display_symbol"`
    ViewOnly                 bool    `json:"view_only"`
    PriceIncrement           string  `json:"price_increment"`
    DisplayName              string  `json:"display_name"`
    ProductVenue             string  `json:"product_venue"`
    ApproximateQuote24hVolume string  `json:"approximate_quote_24h_volume"`
}