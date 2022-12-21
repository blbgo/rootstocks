package rootstocks

import (
	"encoding/binary"
	"time"

	"github.com/blbgo/record/root"
)

const notePrefix = 32

// NoteDetails is a note about a stock
type NoteDetails struct {
	Timestamp time.Time
	Note      string
}

func (r *NoteDetails) fromItem(item root.Item) error {
	key := item.CopyKey(nil)
	if len(key) != 9 {
		return ErrBarIndexWrongLength
	}
	if key[0] != notePrefix {
		return ErrInvalidNotePrefix
	}
	epochTime := binary.BigEndian.Uint64(key[1:])
	r.Timestamp = time.Unix(int64(epochTime), 0)
	r.Note = string(item.Value())
	return nil
}
