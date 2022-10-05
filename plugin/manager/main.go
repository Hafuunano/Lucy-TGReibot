// Package manager bot管理相关
package manager

import (
	"strconv"
	"strings"

	rei "github.com/fumiama/ReiBot"
	tgba "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	ctrl "github.com/FloatTech/zbpctrl"
)

func init() {
	en := rei.Register("manager", &ctrl.Options[*rei.Ctx]{
		DisableOnDefault: false,
		Help: "bot管理相关\n" +
			"- [emoji][emoji]",
	})
	en.OnMessageCommand("离开", rei.OnlyToMe, rei.SuperUserPermission).SetBlock(true).
		Handle(func(ctx *rei.Ctx) {
			arg := strings.TrimSpace(ctx.State["args"].(string))
			var gid int64
			var err error
			if arg != "" {
				gid, err = strconv.ParseInt(arg, 10, 64)
				if err != nil {
					_, _ = ctx.Caller.Send(tgba.NewMessage(ctx.Message.Chat.ID, "ERROR: "+err.Error()))
					return
				}
			} else {
				gid = ctx.Message.Chat.ID
			}
			_, _ = ctx.Caller.Send(&tgba.LeaveChatConfig{ChatID: gid})
		})
}
