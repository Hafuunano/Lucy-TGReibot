// Package what2eat Waht 2 Eat Package for group.
package what2eat

import (
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/MoYoez/Lucy_reibot/utils/toolchain"
	"github.com/MoYoez/Lucy_reibot/utils/userlist"
	rei "github.com/fumiama/ReiBot"
	"strconv"
)

var engine = rei.Register("what2eat", &ctrl.Options[*rei.Ctx]{
	DisableOnDefault:  false,
	Help:              "今天吃什么群友",
	PrivateDataFolder: "what2eat",
})

func init() {
	engine.OnMessageFullMatch("今天吃什么群友").SetBlock(true).Handle(func(ctx *rei.Ctx) {
		getInt := GetUserListAndChooseOne(ctx)
		if getInt != 0 && toolchain.GetUserNickNameByIDInGroup(ctx, getInt) != "" {
			ctx.SendPlainMessage(true, "决定了, 今天吃"+toolchain.GetUserNickNameByIDInGroup(ctx, getInt))
		} else {
			ctx.SendPlainMessage(true, "Lucy 正在确认群里的人数, 过段时间试试吧w")
		}
	})
}

// GetUserListAndChooseOne choose people.
func GetUserListAndChooseOne(ctx *rei.Ctx) int64 {
	toint64, _ := strconv.ParseInt(userlist.PickUserOnGroup(strconv.FormatInt(ctx.Message.Chat.ID, 10)), 10, 64)
	if !toolchain.CheckIfthisUserInThisGroup(toint64, ctx) {
		userlist.RemoveUserOnList(strconv.FormatInt(toint64, 10), strconv.FormatInt(ctx.Message.Chat.ID, 10))
		TrackerCallFuncGetUserListAndChooseOne(ctx)
	}
	return toint64
}

func TrackerCallFuncGetUserListAndChooseOne(ctx *rei.Ctx) int64 {
	var toint64 int64
	for i := 0; i < 3; i++ {
		toint64, _ = strconv.ParseInt(userlist.PickUserOnGroup(strconv.FormatInt(ctx.Message.Chat.ID, 10)), 10, 64)
		if toolchain.CheckIfthisUserInThisGroup(toint64, ctx) {
			break
		}
	}
	return toint64
}
