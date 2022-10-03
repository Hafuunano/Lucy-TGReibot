// Package base64gua base64卦 与 tea 加解密
package base64gua

import (
	"github.com/FloatTech/floatbox/crypto"
	ctrl "github.com/FloatTech/zbpctrl"
	rei "github.com/fumiama/ReiBot"
	"github.com/fumiama/unibase2n"
	tgba "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/wdvxdr1123/ZeroBot/utils/helper"
)

func init() {
	en := rei.Register("base64gua", &ctrl.Options[*rei.Ctx]{
		DisableOnDefault: false,
		Help: "base64gua加解密\n" +
			"- 六十四卦加密xxx\n- 六十四卦解密xxx\n- 六十四卦用yyy加密xxx\n- 六十四卦用yyy解密xxx",
	})
	en.OnMessageRegex(`^六十四卦加密\s*(.+)$`).SetBlock(true).
		Handle(func(ctx *rei.Ctx) {
			str := ctx.State["regex_matched"].([]string)[1]
			es := unibase2n.Base64Gua.EncodeString(str)
			if es != "" {
				_, _ = ctx.Caller.Send(tgba.NewMessage(ctx.Message.Chat.ID, es))
			} else {
				_, _ = ctx.Caller.Send(tgba.NewMessage(ctx.Message.Chat.ID, "加密失败!"))
			}
		})
	en.OnMessageRegex(`^六十四卦解密\s*([䷀-䷿]+[☰☱]?)$`).SetBlock(true).
		Handle(func(ctx *rei.Ctx) {
			str := ctx.State["regex_matched"].([]string)[1]
			es := unibase2n.Base64Gua.DecodeString(str)
			if es != "" {
				_, _ = ctx.Caller.Send(tgba.NewMessage(ctx.Message.Chat.ID, es))
			} else {
				_, _ = ctx.Caller.Send(tgba.NewMessage(ctx.Message.Chat.ID, "解密失败!"))
			}
		})
	en.OnMessageRegex(`^六十四卦用(.+)加密\s*(.+)$`).SetBlock(true).
		Handle(func(ctx *rei.Ctx) {
			key, str := ctx.State["regex_matched"].([]string)[1], ctx.State["regex_matched"].([]string)[2]
			t := crypto.GetTEA(key)
			es, err := unibase2n.UTF16BE2UTF8(unibase2n.Base64Gua.Encode(t.Encrypt(helper.StringToBytes(str))))
			if err == nil {
				_, _ = ctx.Caller.Send(tgba.NewMessage(ctx.Message.Chat.ID, helper.BytesToString(es)))
			} else {
				_, _ = ctx.Caller.Send(tgba.NewMessage(ctx.Message.Chat.ID, "加密失败!"))
			}
		})
	en.OnMessageRegex(`^六十四卦用(.+)解密\s*([䷀-䷿]+[☰☱]?)$`).SetBlock(true).
		Handle(func(ctx *rei.Ctx) {
			key, str := ctx.State["regex_matched"].([]string)[1], ctx.State["regex_matched"].([]string)[2]
			t := crypto.GetTEA(key)
			es, err := unibase2n.UTF82UTF16BE(helper.StringToBytes(str))
			if err == nil {
				_, _ = ctx.Caller.Send(tgba.NewMessage(ctx.Message.Chat.ID, helper.BytesToString(t.Decrypt(unibase2n.Base64Gua.Decode(es)))))
			} else {
				_, _ = ctx.Caller.Send(tgba.NewMessage(ctx.Message.Chat.ID, "解密失败!"))
			}
		})
}
