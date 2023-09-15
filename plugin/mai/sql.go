package mai

import (
	"strconv"
	"sync"
	"time"

	sql "github.com/FloatTech/sqlite"
)

type DataHostSQL struct {
	TelegramId int64  `db:"telegramid"` // telegramid
	Username   string `db:"username"`   // maimai user id ,query from this one.
	Plate      string `db:"plate"`      // plate
	Background string `db:"bg"`         // bg
}

var (
	maiDatabase = &sql.Sqlite{}
	maiLocker   = sync.Mutex{}
)

func init() {
	maiDatabase.DBPath = engine.DataFolder() + "maisql.db"
	err := maiDatabase.Open(time.Hour * 24)
	if err != nil {
		panic(err)
	}
	_ = InitDataBase()
}

func FormatUserDataBase(tgid int64, plate string, bg string, username string) *DataHostSQL {
	return &DataHostSQL{TelegramId: tgid, Plate: plate, Background: bg, Username: username}
}

func InitDataBase() error {
	maiLocker.Lock()
	defer maiLocker.Unlock()
	return maiDatabase.Create("userinfo", &DataHostSQL{})
}

// GetUserPlateInfoFromDatabase Get plate data
func GetUserPlateInfoFromDatabase(userID int64) string {
	maiLocker.Lock()
	defer maiLocker.Unlock()
	var infosql DataHostSQL
	userIDStr := strconv.FormatInt(userID, 10)
	_ = maiDatabase.Find("userinfo", &infosql, "where telegramid is "+userIDStr)
	return infosql.Plate
}

// GetUserInfoNameFromDatabase GetUserName
func GetUserInfoNameFromDatabase(userID int64) string {
	maiLocker.Lock()
	defer maiLocker.Unlock()
	var infosql DataHostSQL
	userIDStr := strconv.FormatInt(userID, 10)
	_ = maiDatabase.Find("userinfo", &infosql, "where telegramid is "+userIDStr)
	if infosql.Username == "" {
		return ""
	}
	return infosql.Username
}

// GetUserDefaultBackgroundDataFromDatabase Get Default Background.
func GetUserDefaultBackgroundDataFromDatabase(userID int64) string {
	maiLocker.Lock()
	defer maiLocker.Unlock()
	var infosql DataHostSQL
	userIDStr := strconv.FormatInt(userID, 10)
	_ = maiDatabase.Find("userinfo", &infosql, "where telegramid is "+userIDStr)
	return infosql.Background
}

// BindUserDataBase Bind Database only for DataHost Inline code.
func (info *DataHostSQL) BindUserDataBase() error {
	maiLocker.Lock()
	defer maiLocker.Unlock()
	return maiDatabase.Insert("userinfo", info)
}
