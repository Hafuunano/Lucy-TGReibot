package userlist

import (
	"strconv"
	"sync"
	"time"

	"github.com/FloatTech/floatbox/file"
	sql "github.com/FloatTech/sqlite"
)

type UserList struct {
	UserID int64 `db:"userid"`
}

var (
	groupSaver  = &sql.Sqlite{}
	groupLocker = sync.Mutex{}
)

// handle when receiving user data, save it to database KV:{username : UserID}

func init() {
	// when using Lucy's functions (need to sure the trigger id.)
	filePath := file.BOTPATH + "/data/rbp/group_lib.db"
	groupSaver.DBPath = filePath
	err := groupSaver.Open(time.Hour * 24)
	if err != nil {
		return
	}
}

// InitDataGroup Init A group if we cannot find it.
func InitDataGroup(groupID int64) {
	groupLocker.Lock()
	defer groupLocker.Unlock()
	groupSaver.Create(strconv.FormatInt(groupID, 10), &UserList{})
}

func SaveUserOnList(userid int64, groupID int64) {
	groupLocker.Lock()
	defer groupLocker.Unlock()
	// check the table is existed? if not, create it.
	err := groupSaver.Insert(strconv.FormatInt(groupID, 10), &UserList{UserID: userid})
	if err != nil {
		InitDataGroup(groupID)
		groupSaver.Insert(strconv.FormatInt(groupID, 10), &UserList{UserID: userid})
	}
}

func RemoveUserOnList(userid int64, groupID int64) {
	groupLocker.Lock()
	defer groupLocker.Unlock()
	// check the table is existed? if not, create it.
	err := groupSaver.Del(strconv.FormatInt(groupID, 10), "WHERE userid is "+strconv.FormatInt(userid, 10))
	if err != nil {
		InitDataGroup(groupID)
	}
}

func PickUserOnGroup(gid int64) int64 {
	groupLocker.Lock()
	defer groupLocker.Unlock()
	var PickerAxe UserList
	err := groupSaver.Pick(strconv.FormatInt(gid, 10), &PickerAxe)
	if err != nil {
		InitDataGroup(gid)
		return 0
	}
	return PickerAxe.UserID
}

func GetThisGroupList(gid int64) []int64 {
	groupLocker.Lock()
	defer groupLocker.Unlock()
	getNum, _ := groupSaver.Count(strconv.FormatInt(gid, 10))
	if getNum == 0 {
		return nil
	}
	var list []int64
	var onTemploader UserList
	_ = groupSaver.FindFor(strconv.FormatInt(gid, 10), &onTemploader, "WHERE id = 0", func() error {
		list = append(list, onTemploader.UserID)
		return nil
	})
	return list
}
