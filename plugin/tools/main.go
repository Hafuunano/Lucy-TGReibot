package tools

import (
	"strconv"
	"strings"

	ctrl "github.com/FloatTech/zbpctrl"
	rei "github.com/fumiama/ReiBot"
	tgba "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var engine = rei.Register("tools", &ctrl.Options[*rei.Ctx]{
	DisableOnDefault: false,
	Help:             "tools for Lucy",
})

func init() {
	engine.OnMessageCommand("leave", rei.SuperUserPermission).SetBlock(true).
		Handle(func(ctx *rei.Ctx) {
			arg := strings.TrimSpace(ctx.State["args"].(string))
			var gid int64
			var err error
			if arg != "" {
				gid, err = strconv.ParseInt(arg, 10, 64)
				if err != nil {
					_, _ = ctx.SendPlainMessage(false, "ERROR: ", err)
					return
				}
			} else {
				gid = ctx.Message.Chat.ID
			}
			_, _ = ctx.Caller.Send(&tgba.LeaveChatConfig{ChatID: gid})
		})
}
