package wife

import (
	"strconv"
	"sync"
	"time"

	sql "github.com/FloatTech/sqlite"
)

var (
	db = &database{db: &sql.Sqlite{}}
)

type database struct {
	db *sql.Sqlite
	sync.RWMutex
}

// 群配置
type groupinfo struct {
	GroupID     int64
	SwitchMarry int
	SwitchNTR   int
	LimitTime   float64
	UpdateTime  string
}

// 信息
type certificates struct {
	ManID      int64
	WomanID    int64
	ManName    string
	WomanName  string
	UpdateTime string
}

// 好感度
type favorability struct {
	UserInfo string
	Favor    int
}

// 技能CD记录表
type cooling struct {
	UnixTime int64  // 时间
	GroupID  int64  // 群号
	UserID   int64  // 用户
	ModeID   string // 技能类型
}

func nowday() string {
	return time.Now().Format("2006/01/02")
}

func nowtime() string {
	return time.Now().Format("15:04:05")
}

func (sql *database) checktime(gid int64) error {
	gpinfo, err := sql.watchsetting(gid)
	if err != nil {
		return err
	}
	strgid := strconv.FormatInt(gid, 10)
	sql.Lock()
	defer sql.Unlock()
	if nowday() != gpinfo.UpdateTime {
		_ = sql.db.Drop("group" + strgid)
		gpinfo.UpdateTime = nowday()
		return sql.db.Insert("groupinfo", &gpinfo)
	}
	return nil
}

func (sql *database) watchsetting(gid int64) (info groupinfo, err error) {
	sql.Lock()
	defer sql.Unlock()
	strgid := strconv.FormatInt(gid, 10)
	err = sql.db.Find("groupinfo", &info, "where gid is "+strgid)
	if err == nil {
		return
	}
	info = groupinfo{
		GroupID:     gid,
		SwitchMarry: 1,
		SwitchNTR:   1,
		LimitTime:   12,
	}
	err = sql.db.Insert("groupinfo", &info)
	return
}

func (sql *database) findcertificates(gid, uid int64) (info certificates, err error) {
	sql.Lock()
	defer sql.Unlock()
	strgid := "group" + strconv.FormatInt(gid, 10)
	ok := sql.db.CanFind(strgid, "where id = 0")
	if !ok {
		err = sql.db.Create(strgid, &certificates{})
		if err != nil {
			return
		}
	}
	struid := strconv.FormatInt(uid, 10)
	err = sql.db.Find(strgid, &info, "where ManID = "+struid)
	if err != nil {
		err = sql.db.Find(strgid, &info, "where WomanID = "+struid)
	}
	return
}

func (sql *database) updatecertificates(gid int64, info *certificates) error {
	sql.Lock()
	defer sql.Unlock()
	strgid := "group" + strconv.FormatInt(gid, 10)
	return sql.db.Insert(strgid, info)
}
