package rootstocks

import (
	"encoding/binary"
	"encoding/json"
	"time"

	"github.com/blbgo/record/root"
)

// Stock represents a stock in the database, provides access to information and data points about
// the stock
type Stock interface {
	Ticker() string
	Details() (*StockDetails, error)
	Update(details *StockDetails) error
	WriteBar(details BarDetails) error
	RangeBars(duration BarDuration, start time.Time, reverse bool, cb func(bar Bar) bool) error
}

type stock struct {
	root.Item
}

const (
	childTypeDayBar byte = iota
	childTypeMinuteBar
)

// StockDetails is a collection of data about a stock
type StockDetails struct {
	Name     string
	Sector   string
	CUSIP    string
	ISIN     string
	SEDOL    string
	Location string
	Exchange string
	Currency string
	R3kRank  uint32
	MemberOf []string
}

func (r stock) Ticker() string {
	return tickerFromBytes(r.Item.CopyKey(nil))
}

func (r stock) Details() (*StockDetails, error) {
	result := &StockDetails{}
	err := json.Unmarshal(r.Item.Value(), result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r stock) Update(details *StockDetails) error {
	value, err := json.Marshal(details)
	if err != nil {
		return err
	}
	return r.Item.UpdateValue(value)
}

func (r stock) WriteBar(details BarDetails) error {
	value, err := json.Marshal(details)
	if err != nil {
		return err
	}
	key := make([]byte, 9)
	key[0] = byte(details.Duration)
	binary.BigEndian.PutUint64(key[1:], uint64(details.Timestamp.Unix()))
	item, err := r.Item.ReadChild(key)
	if err == root.ErrItemNotFound {
		return r.Item.QuickChild(key, value)
	}
	if err != nil {
		return err
	}
	return item.UpdateValue(value)
}

func (r stock) RangeBars(
	duration BarDuration,
	start time.Time,
	reverse bool,
	cb func(bar Bar) bool,
) error {
	startKey := make([]byte, 9)
	startKey[0] = byte(duration)
	binary.BigEndian.PutUint64(startKey[1:], uint64(start.Unix()))
	return r.Item.RangeChildren(startKey, 1, reverse, func(item root.Item) bool {
		return cb(bar{Item: item})
	})
}
