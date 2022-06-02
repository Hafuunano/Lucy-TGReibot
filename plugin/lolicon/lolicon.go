// Package lolicon 基于 https://api.lolicon.app 随机图片
package lolicon

import (
	"strings"
	"time"

	rei "github.com/fumiama/ReiBot"
	tgba "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/tidwall/gjson"

	"github.com/FloatTech/zbputils/math"
	"github.com/FloatTech/zbputils/web"
)

const (
	api      = "https://api.lolicon.app/setu/v2"
	capacity = 10
)

var (
	queue = make(chan string, capacity)
)

func init() {
	rei.OnMessageFullMatch("来份萝莉").SetBlock(true).SetPriority(10).
		Handle(func(ctx *rei.Ctx) {
			go func() {
				for i := 0; i < math.Min(cap(queue)-len(queue), 2); i++ {
					data, err := web.GetData(api)
					if err != nil {
						continue
					}
					json := gjson.ParseBytes(data)
					if e := json.Get("error").Str; e != "" {
						continue
					}
					url := json.Get("data.0.urls.original").Str
					url = strings.ReplaceAll(url, "i.pixiv.cat", "i.pixiv.re")
					queue <- url
				}
			}()
			msg := ctx.Value.(*tgba.Message)
			select {
			case <-time.After(time.Minute):
				_, _ = ctx.Caller.Send(tgba.NewMessage(msg.Chat.ID, "ERROR: 等待填充，请稍后再试..."))
			case img := <-queue:
				_, _ = ctx.Caller.Send(tgba.NewPhoto(msg.Chat.ID, tgba.FileURL(img)))
			}
		})
}
