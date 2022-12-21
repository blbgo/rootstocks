package rootstocks

import (
	"encoding/binary"
	"math"
	"time"

	"github.com/blbgo/record/root"
)

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
	Duration  BarDuration
	Timestamp time.Time
	Open      float32
	High      float32
	Low       float32
	Close     float32
	Volume    uint64
	Status    uint32
}

const barDetailsBinaryLength = 4 + 4 + 4 + 4 + 8 + 4

func (r *BarDetails) fromItem(item root.Item) error {
	err := r.fromBytes(item.Value())
	//err := json.Unmarshal(r.Item.Value(), result)
	if err != nil {
		return err
	}
	key := item.CopyKey(nil)
	if len(key) != 9 {
		return ErrBarIndexWrongLength
	}
	r.Duration = BarDuration(key[0])
	epochTime := binary.BigEndian.Uint64(key[1:])
	r.Timestamp = time.Unix(int64(epochTime), 0)
	return nil
}

func (r *BarDetails) toBytes() []byte {
	result := make([]byte, barDetailsBinaryLength)
	binary.BigEndian.PutUint32(result, math.Float32bits(r.Open))
	binary.BigEndian.PutUint32(result[4:], math.Float32bits(r.High))
	binary.BigEndian.PutUint32(result[8:], math.Float32bits(r.Low))
	binary.BigEndian.PutUint32(result[12:], math.Float32bits(r.Close))
	binary.BigEndian.PutUint64(result[16:], r.Volume)
	binary.BigEndian.PutUint32(result[24:], r.Status)
	return result
}

func (r *BarDetails) fromBytes(bytes []byte) error {
	if len(bytes) < barDetailsBinaryLength {
		return ErrBarValueWrongLength
	}
	r.Open = math.Float32frombits(binary.BigEndian.Uint32(bytes))
	r.High = math.Float32frombits(binary.BigEndian.Uint32(bytes[4:]))
	r.Low = math.Float32frombits(binary.BigEndian.Uint32(bytes[8:]))
	r.Close = math.Float32frombits(binary.BigEndian.Uint32(bytes[12:]))
	r.Volume = binary.BigEndian.Uint64(bytes[16:])
	r.Status = binary.BigEndian.Uint32(bytes[24:])
	return nil
}
