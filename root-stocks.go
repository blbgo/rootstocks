package rootstocks

import (
	"encoding/json"
	"errors"

	"github.com/blbgo/record/root"
)

// RootStocks provides access to a database of stocks
type RootStocks interface {
	CreateStock(ticker string, details *StockDetails) (Stock, error)
	ReadStock(ticker string) (Stock, error)
	RangeStocks(startTicker string, reverse bool, cb func(stock Stock) bool) error
	RangeStockTickers(startTicker string, reverse bool, cb func(ticker string) bool) error
}

// ErrBarIndexWrongLength bar record index is not 9 bytes long
var ErrBarIndexWrongLength = errors.New("Bar record index is not 9 bytes long")

// ErrBarValueWrongLength bar record value is wrong length
var ErrBarValueWrongLength = errors.New("Bar record value is wrong length")

// ErrNilArgument bar record value is wrong length
var ErrNilArgument = errors.New("Argument in nil")

// ErrInvalidNotePrefix note key prefix is invalid
var ErrInvalidNotePrefix = errors.New("Note key prefix invalid")

// ErrNoteWithTimeAlreadyExists Note with time already exists
var ErrNoteWithTimeAlreadyExists = errors.New("Note with time already exists")

type rootStocks struct {
	root.Item
}

// New creates a RootStocks implemented by rootStocks
func New(theRoot root.Root) (RootStocks, error) {
	item, err := theRoot.RootItem(
		"github.com/blbgo/rootstocks",
		"github.com/blbgo/rootstocks root item",
	)
	if err != nil {
		return nil, err
	}
	return rootStocks{Item: item}, nil
}

// **************** implement RootStocks

func (r rootStocks) CreateStock(ticker string, details *StockDetails) (Stock, error) {
	value, err := json.Marshal(details)
	if err != nil {
		return nil, err
	}
	item, err := r.Item.CreateChild(tickerToBytes(ticker), value, nil)
	if err != nil {
		return nil, err
	}
	return stock{Item: item}, err
}

func (r rootStocks) ReadStock(ticker string) (Stock, error) {
	item, err := r.Item.ReadChild(tickerToBytes(ticker))
	if err != nil {
		return nil, err
	}
	return stock{Item: item}, err
}

func (r rootStocks) RangeStocks(
	startTicker string,
	reverse bool,
	cb func(stock Stock) bool,
) error {
	return r.Item.RangeChildren(
		[]byte(startTicker),
		0,
		reverse,
		func(item root.Item) bool {
			return cb(stock{Item: item})
		},
	)
}

func (r rootStocks) RangeStockTickers(
	startTicker string,
	reverse bool,
	cb func(ticker string) bool,
) error {
	return r.Item.RangeChildKeys(
		[]byte(startTicker),
		0,
		reverse,
		func(key []byte) bool {
			return cb(tickerFromBytes(key))
		},
	)
}

func tickerToBytes(ticker string) []byte {
	tickerBytes := []byte(ticker)
	if len(tickerBytes) == 1 {
		return append(tickerBytes, 0)
	}
	return tickerBytes
}

func tickerFromBytes(key []byte) string {
	if len(key) == 2 && key[1] == 0 {
		return string(key[:1])
	}
	return string(key)
}
