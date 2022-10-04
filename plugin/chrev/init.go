// Package chrev 英文字符反转
package chrev

import (
	"strings"

	rei "github.com/fumiama/ReiBot"
	tgba "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	ctrl "github.com/FloatTech/zbpctrl"
)

func init() {
	// 初始化engine
	engine := rei.Register("chrev", &ctrl.Options[*rei.Ctx]{
		DisableOnDefault: false,
		Help:             "字符翻转\n- 翻转 I love you",
	})
	// 处理字符翻转指令
	engine.OnMessageRegex(`^翻转\s*([A-Za-z\s]*)$`).SetBlock(true).
		Handle(func(ctx *rei.Ctx) {
			// 获取需要翻转的字符串
			str := ctx.State["regex_matched"].([]string)[1]
			// 将字符顺序翻转
			tmp := strings.Builder{}
			for i := len(str) - 1; i >= 0; i-- {
				tmp.WriteRune(charMap[str[i]])
			}
			// 发送翻转后的字符串
			_, _ = ctx.Caller.Send(tgba.NewMessage(ctx.Message.Chat.ID, tmp.String()))
		})
}
