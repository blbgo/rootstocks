package rootstocks

import (
	"encoding/binary"
	"encoding/json"
	"time"

	"github.com/blbgo/record/root"
)

// Bar represents price information of a stock over a period of time
type Bar interface {
	Details() (*BarDetails, error)
}

type bar struct {
	root.Item
}

// BarDuration indicates the duration of a bar
type BarDuration byte

const (
	// DayBar represents a bar covering a whole day
	DayBar BarDuration = BarDuration(childTypeDayBar)
	// MinuteBar represents a bar covering one minuet
	MinuteBar BarDuration = BarDuration(childTypeMinuteBar)
)

// BarDetails is a collection of data about a bar
type BarDetails struct {
	Duration        BarDuration `json:"-"`
	Timestamp       time.Time   `json:"-"`
	Open            float64
	High            float64
	Low             float64
	Close           float64
	UpTicks         uint64
	UpVolume        uint64
	DownTicks       uint64
	DownVolume      uint64
	UnchangedTicks  uint64
	UnchangedVolume uint64
	TotalTicks      uint64
	TotalVolume     uint64
	Status          uint32
	// not needed for stock
	//OpenInterest         uint64
}

func (r bar) Details() (*BarDetails, error) {
	result := &BarDetails{}
	err := json.Unmarshal(r.Item.Value(), result)
	if err != nil {
		return nil, err
	}
	key := r.Item.CopyKey(nil)
	if len(key) != 9 {
		return nil, ErrBarIndexWrongLength
	}
	result.Duration = BarDuration(key[0])
	epochTime := binary.BigEndian.Uint64(key[1:])
	result.Timestamp = time.Unix(int64(epochTime), 0)
	return result, nil
}
