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

// UserIDToQQ TIPS: Onebot path, actually it refers to Telegram Userid.

type UserIDToQQ struct {
	QQ     int64  `db:"user_qq"` // qq nums
	Userid string `db:"user_id"` // user_id
}

var (
	maiDatabase = &sql.Sqlite{}
	maiLocker   = sync.Mutex{}
)

func init() {
	maiDatabase.DBPath = engine.DataFolder() + "maisql.db"
	err := maiDatabase.Open(time.Hour * 24)
	if err != nil {
		return
	}
	_ = InitDataBase()
}

func FormatUserDataBase(tgid int64, plate string, bg string, username string) *DataHostSQL {
	return &DataHostSQL{TelegramId: tgid, Plate: plate, Background: bg, Username: username}
}

func FormatUserIDDatabase(qq int64, userid string) *UserIDToQQ {
	return &UserIDToQQ{QQ: qq, Userid: userid}
}

func InitDataBase() error {
	maiLocker.Lock()
	defer maiLocker.Unlock()
	maiDatabase.Create("userinfo", &DataHostSQL{})
	maiDatabase.Create("useridinfo", &UserIDToQQ{})
	return nil
}

// GetUserIDFromDatabase Params: user qq id ==> user maimai id.
func GetUserIDFromDatabase(userID int64) UserIDToQQ {
	maiLocker.Lock()
	defer maiLocker.Unlock()
	var infosql UserIDToQQ
	userIDStr := strconv.FormatInt(userID, 10)
	_ = maiDatabase.Find("useridinfo", &infosql, "where user_qq is "+userIDStr)
	return infosql
}

func (info *UserIDToQQ) BindUserIDDataBase() error {
	maiLocker.Lock()
	defer maiLocker.Unlock()
	return maiDatabase.Insert("useridinfo", info)
}

// maimai origin render base.

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
