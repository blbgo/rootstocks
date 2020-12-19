package rootstocks

import (
	"encoding/json"

	"github.com/blbgo/record/root"
)

// Stock represents a stock in the database, provides access to information and data points about
// the stock
type Stock interface {
	Ticker() string
	Details() (*StockDetails, error)
}

type stock struct {
	root.Item
}

// StockDetails is a collection of data about a stock
type StockDetails struct {
	Name          string
	Sector        string
	WeightPercent float32
	SharesInETF   uint64
	CUSIP         string
	ISIN          string
	SEDOL         string
	//Price
	Location string
	Exchange string
	Currency string
	//FX Rate
	MarketCurrency string
}

func (r stock) Ticker() string {
	return string(r.Item.CopyKey(nil))
}

func (r stock) Details() (*StockDetails, error) {
	result := &StockDetails{}
	err := json.Unmarshal(r.Item.Value(), result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
