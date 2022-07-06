// Package b14coder base16384 与 tea 加解密
package b14coder

import (
	"unsafe"

	ctrl "github.com/FloatTech/zbpctrl"
	rei "github.com/fumiama/ReiBot"
	base14 "github.com/fumiama/go-base16384"
	tea "github.com/fumiama/gofastTEA"
	tgba "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/wdvxdr1123/ZeroBot/utils/helper"
)

func init() {
	en := rei.Register("base16384", &ctrl.Options[*rei.Ctx]{
		DisableOnDefault: false,
		Help: "base16384加解密\n" +
			"- 加密xxx\n- 解密xxx\n- 用yyy加密xxx\n- 用yyy解密xxx",
	})
	en.OnMessageRegex(`^加密\s*(.+)$`).SetBlock(true).
		Handle(func(ctx *rei.Ctx) {
			str := ctx.State["regex_matched"].([]string)[1]
			es := base14.EncodeString(str)
			if es != "" {
				_, _ = ctx.Caller.Send(tgba.NewMessage(ctx.Message.Chat.ID, es))
			} else {
				_, _ = ctx.Caller.Send(tgba.NewMessage(ctx.Message.Chat.ID, "加密失败!"))
			}
		})
	en.OnMessageRegex(`^解密\s*([一-踀]+[㴁-㴆]?)$`).SetBlock(true).
		Handle(func(ctx *rei.Ctx) {
			str := ctx.State["regex_matched"].([]string)[1]
			es := base14.DecodeString(str)
			if es != "" {
				_, err := ctx.Caller.Send(tgba.NewMessage(ctx.Message.Chat.ID, es))
				if err != nil {
					_, _ = ctx.Caller.Send(tgba.NewMessage(ctx.Message.Chat.ID, "ERROR: "+err.Error()))
				}
			} else {
				_, _ = ctx.Caller.Send(tgba.NewMessage(ctx.Message.Chat.ID, "解密失败!"))
			}
		})
	en.OnMessageRegex(`^用(.+)加密\s*(.+)$`).SetBlock(true).
		Handle(func(ctx *rei.Ctx) {
			key, str := ctx.State["regex_matched"].([]string)[1], ctx.State["regex_matched"].([]string)[2]
			t := getea(key)
			es, err := base14.UTF16BE2UTF8(base14.Encode(t.Encrypt(helper.StringToBytes(str))))
			if err == nil {
				_, _ = ctx.Caller.Send(tgba.NewMessage(ctx.Message.Chat.ID, helper.BytesToString(es)))
			} else {
				_, _ = ctx.Caller.Send(tgba.NewMessage(ctx.Message.Chat.ID, "加密失败!"))
			}
		})
	en.OnMessageRegex(`^用(.+)解密\s*([一-踀]+[㴁-㴆]?)$`).SetBlock(true).
		Handle(func(ctx *rei.Ctx) {
			key, str := ctx.State["regex_matched"].([]string)[1], ctx.State["regex_matched"].([]string)[2]
			t := getea(key)
			es, err := base14.UTF82UTF16BE(helper.StringToBytes(str))
			if err == nil {
				_, err := ctx.Caller.Send(tgba.NewMessage(ctx.Message.Chat.ID, helper.BytesToString(t.Decrypt(base14.Decode(es)))))
				if err != nil {
					_, _ = ctx.Caller.Send(tgba.NewMessage(ctx.Message.Chat.ID, "ERROR: "+err.Error()))
				}
			} else {
				_, _ = ctx.Caller.Send(tgba.NewMessage(ctx.Message.Chat.ID, "解密失败!"))
			}
		})
}

func getea(key string) tea.TEA {
	kr := []rune(key)
	if len(kr) > 4 {
		kr = kr[:4]
	} else {
		for len(kr) < 4 {
			kr = append(kr, rune(4-len(kr)))
		}
	}
	return *(*tea.TEA)(*(*unsafe.Pointer)(unsafe.Pointer(&kr)))
}
