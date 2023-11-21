package userlist

import (
	"sync"
	"time"

	"github.com/FloatTech/floatbox/file"
	sql "github.com/FloatTech/sqlite"
)

type UserList struct {
	UserID   string `db:"userid"`
	UserName string `db:"username"`
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
		panic(err)
	}
}

// InitDataGroup Init A group if we cannot find it.
func InitDataGroup(groupID string) error {

	return groupSaver.Create("group_"+groupID, &UserList{})
}

func SaveUserOnList(userid string, groupID string, username string) {

	// check the table is existed? if not, create it.
	InitDataGroup(groupID)
	groupSaver.Insert("group_"+groupID, &UserList{UserID: userid, UserName: username})
}

func RemoveUserOnList(userid string, groupID string) {

	// check the table is existed? if not, create it.
	groupSaver.Del("group_"+groupID, "WHERE userid is "+userid)
	InitDataGroup(groupID)
}

func PickUserOnGroup(gid string) string {

	var PickerAxe UserList
	groupSaver.Pick("group_"+gid, &PickerAxe)
	InitDataGroup(gid)
	return PickerAxe.UserID
}

func GetThisGroupList(gid string) []string {

	getNum, _ := groupSaver.Count("group_" + gid)
	if getNum == 0 {
		return nil
	}
	var list []string
	var onTemploader UserList

	_ = groupSaver.FindFor("group_"+gid, &onTemploader, "WHERE id = 0", func() error {
		list = append(list, onTemploader.UserID)
		return nil
	})
	return list
}
