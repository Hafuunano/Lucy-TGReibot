package novelai

import (
	"hash/crc64"
	"sync"
	"time"

	"github.com/FloatTech/floatbox/binary"
	sql "github.com/FloatTech/sqlite"
)

type imgstorage struct {
	ID   int64  `db:"id"`
	Seed int32  `db:"seed"`
	Tags string `db:"tags"`
}

type keystorage struct {
	Sender int64  `db:"sender"`
	OnlyMe bool   `db:"onlyme"`
	Key    string `db:"key"`
}

var (
	ims *sql.Sqlite
	mu  sync.RWMutex
	iso = crc64.MakeTable(crc64.ISO)
)

func newims(dbpath string) *sql.Sqlite {
	ims = &sql.Sqlite{}
	ims.DBPath = dbpath
	err := ims.Open(time.Hour)
	if err != nil {
		panic(err)
	}
	err = ims.Create("s", &imgstorage{})
	if err != nil {
		panic(err)
	}
	err = ims.Create("k", &keystorage{})
	if err != nil {
		panic(err)
	}
	return ims
}

func idof(fn string) int64 {
	return int64(crc64.Checksum(binary.StringToBytes(fn), iso))
}
