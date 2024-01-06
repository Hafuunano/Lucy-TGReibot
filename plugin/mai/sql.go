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

type UserSwitcherService struct {
	TGId   int64 `db:"tgid"`
	IsUsed bool  `db:"isused"` // true == lxns service \ false == Diving Fish.
}

type UserIDToMaimaiFriendCode struct {
	TelegramId int64 `db:"telegramid"`
	MaimaiID   int64 `db:"friendid"`
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

func FormatUserSwitcher(tgid int64, isSwitcher bool) *UserSwitcherService {
	return &UserSwitcherService{TGId: tgid, IsUsed: isSwitcher}
}

func FormatMaimaiFriendCode(friendCode int64, tgid int64) *UserIDToMaimaiFriendCode {
	return &UserIDToMaimaiFriendCode{TelegramId: tgid, MaimaiID: friendCode}
}

func InitDataBase() error {
	maiLocker.Lock()
	defer maiLocker.Unlock()
	maiDatabase.Create("userinfo", &DataHostSQL{})
	maiDatabase.Create("useridinfo", &UserIDToQQ{})
	maiDatabase.Create("userswitchinfo", &UserSwitcherService{})
	maiDatabase.Create("usermaifriendinfo", &UserIDToMaimaiFriendCode{})
	return nil
}

func GetUserMaiFriendID(userid int64) UserIDToMaimaiFriendCode {
	maiLocker.Lock()
	defer maiLocker.Unlock()
	var infosql UserIDToMaimaiFriendCode
	userIDStr := strconv.FormatInt(userid, 10)
	_ = maiDatabase.Find("usermaifriendinfo", &infosql, "where telegramid is "+userIDStr)
	return infosql
}

func GetUserSwitcherInfoFromDatabase(userid int64) bool {
	maiLocker.Lock()
	defer maiLocker.Unlock()
	var info UserSwitcherService
	userIDStr := strconv.FormatInt(userid, 10)
	err := maiDatabase.Find("userswitchinfo", &info, "where tgid is "+userIDStr)
	if err != nil {
		return false
	}
	return info.IsUsed
}

func (info *UserSwitcherService) ChangeUserSwitchInfoFromDataBase() error {
	maiLocker.Lock()
	defer maiLocker.Unlock()
	return maiDatabase.Insert("userswitchinfo", info)
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

func (info *UserIDToMaimaiFriendCode) BindUserFriendCode() error {
	maiLocker.Lock()
	defer maiLocker.Unlock()
	return maiDatabase.Insert("usermaifriendinfo", info)
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
