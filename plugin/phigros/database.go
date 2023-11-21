package phigros

import (
	"strconv"
	"sync"
	"time"

	sql "github.com/FloatTech/sqlite"
)

type PhigrosSQL struct {
	Id         int64  `db:"user_tgid"` // tgid
	PhiSession string `db:"session"`   // pgr session
	Time       int64  `db:"time"`      // time.
}

var (
	pgrDatabase = &sql.Sqlite{}
	pgrLocker   = sync.Mutex{}
)

func init() {
	pgrDatabase.DBPath = engine.DataFolder() + "pgrsql.db"
	err := pgrDatabase.Open(time.Hour * 24)
	if err != nil {
		return
	}
	_ = InitDataBase()
}

func FormatUserDataBase(tgid int64, session string, Time int64) *PhigrosSQL {
	return &PhigrosSQL{Id: tgid, PhiSession: session, Time: Time}
}

func InitDataBase() error {
	pgrLocker.Lock()
	defer pgrLocker.Unlock()
	return pgrDatabase.Create("userinfo", &PhigrosSQL{})
}

func GetUserInfoFromDatabase(userID int64) *PhigrosSQL {
	pgrLocker.Lock()
	defer pgrLocker.Unlock()
	var infosql PhigrosSQL
	userIDStr := strconv.FormatInt(userID, 10)
	_ = pgrDatabase.Find("userinfo", &infosql, "where user_tgid is "+userIDStr)
	return &infosql
}

func GetUserInfoTimeFromDatabase(userID int64) int64 {
	pgrLocker.Lock()
	defer pgrLocker.Unlock()
	var infosql PhigrosSQL
	userIDStr := strconv.FormatInt(userID, 10)
	_ = pgrDatabase.Find("userinfo", &infosql, "where user_tgid is "+userIDStr)
	return infosql.Time
}

func (info *PhigrosSQL) BindUserDataBase() error {
	pgrLocker.Lock()
	defer pgrLocker.Unlock()
	return pgrDatabase.Insert("userinfo", info)
}
