// Package emojimix 合成emoji
package emojimix

import (
	"fmt"
	"net/http"

	rei "github.com/fumiama/ReiBot"
	tgba "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/sirupsen/logrus"

	"github.com/FloatTech/ReiBot-Plugin/utils/ctxext"
)

const bed = "https://www.gstatic.com/android/keyboard/emojikitchen/%d/u%x/u%x_u%x.png"

func init() {
	rei.Register("emojimix", &ctrl.Options[*rei.Ctx]{
		DisableOnDefault: false,
		Help: "合成emoji\n" +
			"- [emoji][emoji]",
	}).OnMessage(match).SetBlock(true).Limit(ctxext.LimitByUser).
		Handle(func(ctx *rei.Ctx) {
			r := ctx.State["emojimix"].([]rune)
			logrus.Debugln("[emojimix] match:", r)
			r1, r2 := r[0], r[1]
			u1 := fmt.Sprintf(bed, emojis[r1], r1, r1, r2)
			u2 := fmt.Sprintf(bed, emojis[r2], r2, r2, r1)
			resp1, err := http.Head(u1)
			if err == nil {
				resp1.Body.Close()
				if resp1.StatusCode == http.StatusOK {
					_, _ = ctx.Caller.Send(tgba.NewPhoto(ctx.Message.Chat.ID, tgba.FileURL(u1)))
					return
				}
			}
			resp2, err := http.Head(u2)
			if err == nil {
				resp2.Body.Close()
				if resp2.StatusCode == http.StatusOK {
					_, _ = ctx.Caller.Send(tgba.NewPhoto(ctx.Message.Chat.ID, tgba.FileURL(u2)))
					return
				}
			}
		})
}

func match(ctx *rei.Ctx) bool {
	r := []rune(ctx.Message.Text)
	if len(r) == 2 {
		if _, ok := emojis[r[0]]; !ok {
			return false
		}
		if _, ok := emojis[r[1]]; !ok {
			return false
		}
		ctx.State["emojimix"] = r
		return true
	}
	return false
}
