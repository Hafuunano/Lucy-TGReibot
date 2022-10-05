package novelai

import (
	"hash/crc64"
	"sync"

	"github.com/FloatTech/floatbox/binary"
	sql "github.com/FloatTech/sqlite"
)

type imgstorage struct {
	ID   int64  `db:"id"`
	Seed int32  `db:"seed"`
	Tags string `db:"tags"`
}

var (
	ims = &sql.Sqlite{}
	mu  sync.RWMutex
	iso = crc64.MakeTable(crc64.ISO)
)

func idof(fn string) int64 {
	return int64(crc64.Checksum(binary.StringToBytes(fn), iso))
}
