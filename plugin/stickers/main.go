package stickers

import (
	"math/rand"
	"time"

	ctrl "github.com/FloatTech/zbpctrl"
	rei "github.com/fumiama/ReiBot"
	tgba "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/wdvxdr1123/ZeroBot/extension/rate"
)

var (
	limitedManager = rate.NewManager[int64](time.Minute*10, 8)
	engine         = rei.Register("stickers", &ctrl.Options[*rei.Ctx]{
		DisableOnDefault: true,
		Help:             "stickers",
	})
)

func init() {
	engine.OnMessage().SetBlock(false).Handle(func(ctx *rei.Ctx) {
		if rand.Intn(10) > 1 {
			return
		}
		if !limitedManager.Load(ctx.Message.Chat.ID).Acquire() {
			return
		}
		if ctx.Message.Sticker != nil {
			getStickerPack, err := ctx.Caller.GetStickerSet(tgba.GetStickerSetConfig{Name: ctx.Message.Sticker.SetName})
			if err != nil {
				return
			}
			ctx.Caller.Request(tgba.StickerConfig{BaseFile: tgba.BaseFile{BaseChat: tgba.BaseChat{ChatConfig: tgba.ChatConfig{ChatID: ctx.Message.Chat.ID}}, File: tgba.FileID(getStickerPack.Stickers[rand.Intn(len(getStickerPack.Stickers))].FileID)}})
		}

	})
}
