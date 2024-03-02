// Package CoreFactory add simple event listener to database, record user id and some other info to judge user is correct.
package CoreFactory

import (
	"strconv"
	"sync"
	"time"

	"github.com/FloatTech/floatbox/file"
	sql "github.com/FloatTech/sqlite"
)

type Data struct {
	UserID   int64  `db:"userid"`
	UserName string `db:"username"` // refer: username means @xx, can be direct to user by t.me/@username || channel / group user use save.
}

var (
	coreSaver  = &sql.Sqlite{}
	coreLocker = sync.Mutex{}
)

// handle when receiving user data, save it to database KV:{username : UserID}

func init() {
	// when using Lucy's functions (need to sure the trigger id.)
	filePath := file.BOTPATH + "/data/rbp/lib.db"
	coreSaver.DBPath = filePath
	err := coreSaver.Open(time.Hour * 24)
	if err != nil {
		return
	}
	_ = initDatabase()

}

// GetUserSampleUserinfo User Info.
func GetUserSampleUserinfo(username string) *Data {
	coreLocker.Lock()
	defer coreLocker.Unlock()
	var ResultData Data
	coreSaver.Find("userinfo", &ResultData, "Where username is '"+username+"'")
	if &ResultData == nil {
		return nil
	}
	return &ResultData
}

// GetUserSampleUserinfobyid User Info.
func GetUserSampleUserinfobyid(userid int64) *Data {
	coreLocker.Lock()
	defer coreLocker.Unlock()
	var ResultData Data
	coreSaver.Find("userinfo", &ResultData, "Where userid is "+strconv.FormatInt(userid, 10))
	if &ResultData == nil {
		return nil
	}
	return &ResultData
}

func StoreUserDataBase(userid int64, userName string) error {
	coreLocker.Lock()
	defer coreLocker.Unlock()
	return coreSaver.Insert("userinfo", &Data{UserID: userid, UserName: userName})
}

func initDatabase() error {
	coreLocker.Lock()
	defer coreLocker.Unlock()
	return coreSaver.Create("userinfo", &Data{})
}
