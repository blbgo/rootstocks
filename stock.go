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
	Details(details *StockDetails) error
	Update(details *StockDetails) error
	Delete() error
	WriteBar(details *BarDetails) error
	RangeBars(
		duration BarDuration,
		start time.Time,
		reverse bool,
		cb func(details *BarDetails) bool,
	) error
	DeleteBars(duration BarDuration) error
	WriteNote(details *NoteDetails) error
	RangeNotes(
		start time.Time,
		reverse bool,
		cb func(details *NoteDetails) bool,
	) error
	DeleteNote(timestamp time.Time) error
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
	Name   string
	Sector string
	//CUSIP    string
	//ISIN     string
	//SEDOL    string
	Location string
	Exchange string
	Currency string
	R3kRank  uint32
	MemberOf []string
}

func (r stock) Ticker() string {
	return tickerFromBytes(r.Item.CopyKey(nil))
}

func (r stock) Details(details *StockDetails) error {
	if details == nil {
		return ErrNilArgument
	}
	err := json.Unmarshal(r.Item.Value(), details)
	if err != nil {
		return err
	}
	return nil
}

func (r stock) Update(details *StockDetails) error {
	value, err := json.Marshal(details)
	if err != nil {
		return err
	}
	return r.Item.UpdateValue(value)
}

func (r stock) Delete() error {
	return r.Item.Delete()
}

func (r stock) WriteBar(details *BarDetails) error {
	value := details.toBytes()
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
	cb func(details *BarDetails) bool,
) error {
	startKey := make([]byte, 9)
	startKey[0] = byte(duration)
	binary.BigEndian.PutUint64(startKey[1:], uint64(start.Unix()))
	barDetails := &BarDetails{}
	var err error
	rangeErr := r.Item.RangeChildren(startKey, 1, reverse, func(item root.Item) bool {
		err = barDetails.fromItem(item)
		if err != nil {
			return false
		}
		return cb(barDetails)
	})
	if rangeErr != nil {
		return rangeErr
	}
	return err
}

func (r stock) DeleteBars(duration BarDuration) error {
	startKey := make([]byte, 1)
	startKey[0] = byte(duration)
	var err error
	rangeErr := r.Item.RangeChildren(startKey, 1, false, func(item root.Item) bool {
		err = item.Delete()
		return err == nil
	})
	if rangeErr != nil {
		return rangeErr
	}
	return err

}

func (r stock) WriteNote(details *NoteDetails) error {
	key := make([]byte, 9)
	key[0] = byte(notePrefix)
	binary.BigEndian.PutUint64(key[1:], uint64(details.Timestamp.Unix()))
	_, err := r.Item.ReadChild(key)
	if err == root.ErrItemNotFound {
		return r.Item.QuickChild(key, []byte(details.Note))
	}
	if err != nil {
		return err
	}
	return ErrNoteWithTimeAlreadyExists
}

func (r stock) RangeNotes(
	start time.Time,
	reverse bool,
	cb func(details *NoteDetails) bool,
) error {
	startKey := make([]byte, 9)
	startKey[0] = byte(notePrefix)
	binary.BigEndian.PutUint64(startKey[1:], uint64(start.Unix()))
	noteDetails := &NoteDetails{}
	var err error
	rangeErr := r.Item.RangeChildren(startKey, 1, reverse, func(item root.Item) bool {
		err = noteDetails.fromItem(item)
		if err != nil {
			return false
		}
		return cb(noteDetails)
	})
	if rangeErr != nil {
		return rangeErr
	}
	return err
}

func (r stock) DeleteNote(timestamp time.Time) error {
	key := make([]byte, 9)
	key[0] = byte(notePrefix)
	binary.BigEndian.PutUint64(key[1:], uint64(timestamp.Unix()))
	noteItem, err := r.Item.ReadChild(key)
	if err != nil {
		return err
	}
	return noteItem.Delete()
}
