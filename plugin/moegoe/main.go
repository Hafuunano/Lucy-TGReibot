// Package moegoe 日韩 VITS 模型拟声
package moegoe

import (
	"fmt"
	"net/url"

	rei "github.com/fumiama/ReiBot"
	tgba "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	ctrl "github.com/FloatTech/zbpctrl"

	"github.com/FloatTech/ReiBot-Plugin/utils/ctxext"
)

const (
	jpapi = "https://moegoe.azurewebsites.net/api/speak?format=mp3&text=%s&id=%d"
	krapi = "https://moegoe.azurewebsites.net/api/speakkr?format=mp3&text=%s&id=%d"
)

var speakers = map[string]uint{
	"宁宁": 0, "爱瑠": 1, "芳乃": 2, "茉子": 3, "丛雨": 4, "小春": 5, "七海": 6,
	"수아": 0, "미미르": 1, "아린": 2, "연화": 3, "유화": 4, "선배": 5,
}

func init() {
	en := rei.Register("moegoe", &ctrl.Options[*rei.Ctx]{
		DisableOnDefault: false,
		Help: "moegoe\n" +
			"- 让[宁宁|爱瑠|芳乃|茉子|丛雨|小春|七海]说(日语)\n" +
			"- 让[수아|미미르|아린|연화|유화|선배]说(韩语)",
	}).ApplySingle(ctxext.DefaultSingle)
	en.OnMessageRegex("^让(宁宁|爱瑠|芳乃|茉子|丛雨|小春|七海)说([A-Za-z\\s\\d\u3005\u3040-\u30ff\u4e00-\u9fff\uff11-\uff19\uff21-\uff3a\uff41-\uff5a\uff66-\uff9d.。,，、:：;；!！?？]+)$").Limit(ctxext.LimitByGroup).SetBlock(true).
		Handle(func(ctx *rei.Ctx) {
			speaker := ctx.State["regex_matched"].([]string)[1]
			text := ctx.State["regex_matched"].([]string)[2]
			id := speakers[speaker]
			a := tgba.NewAudio(ctx.Message.Chat.ID, tgba.FileURL(fmt.Sprintf(jpapi, url.QueryEscape(text), id)))
			a.Caption = speaker + ": " + text
			a.CaptionEntities = []tgba.MessageEntity{{Type: "bold", Length: len([]rune(speaker))}}
			_, err := ctx.Caller.Send(&a)
			if err != nil {
				_, _ = ctx.Caller.Send(tgba.NewMessage(ctx.Message.Chat.ID, "ERROR:"+err.Error()))
				return
			}
		})
	en.OnMessageRegex("^让(수아|미미르|아린|연화|유화|선배)说([A-Za-z\\s\\d\u3131-\u3163\uac00-\ud7ff.。,，、:：;；!！?？]+)$").Limit(ctxext.LimitByGroup).SetBlock(true).
		Handle(func(ctx *rei.Ctx) {
			speaker := ctx.State["regex_matched"].([]string)[1]
			text := ctx.State["regex_matched"].([]string)[2]
			id := speakers[speaker]
			a := tgba.NewAudio(ctx.Message.Chat.ID, tgba.FileURL(fmt.Sprintf(krapi, url.QueryEscape(text), id)))
			a.Caption = speaker + ": " + text
			a.CaptionEntities = []tgba.MessageEntity{{Type: "bold", Length: len([]rune(speaker))}}
			_, err := ctx.Caller.Send(&a)
			if err != nil {
				_, _ = ctx.Caller.Send(tgba.NewMessage(ctx.Message.Chat.ID, "ERROR:"+err.Error()))
				return
			}
		})
}
