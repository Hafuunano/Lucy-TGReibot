// Package lolicon 基于 https://api.lolicon.app 随机图片
package lolicon

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	rei "github.com/fumiama/ReiBot"
	tgba "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/FloatTech/zbputils/math"
	"github.com/FloatTech/zbputils/web"
)

const (
	api      = "https://api.lolicon.app/setu/v2?r18=2"
	capacity = 10
)

type lolire struct {
	Error string `json:"error"`
	Data  []struct {
		Pid        int      `json:"pid"`
		P          int      `json:"p"`
		UID        int      `json:"uid"`
		Title      string   `json:"title"`
		Author     string   `json:"author"`
		R18        bool     `json:"r18"`
		Width      int      `json:"width"`
		Height     int      `json:"height"`
		Tags       []string `json:"tags"`
		Ext        string   `json:"ext"`
		UploadDate int64    `json:"uploadDate"`
		Urls       struct {
			Original string `json:"original"`
		} `json:"urls"`
	} `json:"data"`
}

var (
	queue = make(chan *tgba.PhotoConfig, capacity)
)

func init() {
	rei.OnMessageFullMatch("来份萝莉").SetBlock(true).
		Handle(func(ctx *rei.Ctx) {
			go func() {
				for i := 0; i < math.Min(cap(queue)-len(queue), 2); i++ {
					data, err := web.GetData(api)
					if err != nil {
						continue
					}
					var r lolire
					err = json.Unmarshal(data, &r)
					if err != nil {
						continue
					}
					if r.Error != "" {
						continue
					}
					caption := strings.Builder{}
					caption.WriteString(r.Data[0].Title)
					caption.WriteString(" @")
					caption.WriteString(r.Data[0].Author)
					entities := []tgba.MessageEntity{
						{
							Type:   "text_link",
							Offset: 0,
							Length: len([]rune(r.Data[0].Title)),
							URL:    "https://pixiv.net/i/" + strconv.Itoa(r.Data[0].Pid),
						},
						{
							Type:   "text_link",
							Offset: len([]rune(r.Data[0].Title)) + 1,
							Length: len([]rune(r.Data[0].Author)) + 1,
							URL:    "https://pixiv.net/u/" + strconv.Itoa(r.Data[0].UID),
						},
					}
					offset := entities[1].Offset + entities[1].Length
					for _, t := range r.Data[0].Tags {
						caption.WriteByte(' ')
						caption.WriteString(t)
						l := len([]rune(t))
						entities = append(entities, tgba.MessageEntity{
							Type:   "code",
							Offset: offset + 1,
							Length: l,
						})
						offset += l + 1
					}
					queue <- &tgba.PhotoConfig{
						BaseFile: tgba.BaseFile{
							File: tgba.FileURL(r.Data[0].Urls.Original),
						},
						Caption:         caption.String(),
						CaptionEntities: entities,
					}
				}
			}()
			msg := ctx.Value.(*tgba.Message)
			select {
			case <-time.After(time.Minute):
				_, _ = ctx.Caller.Send(tgba.NewMessage(msg.Chat.ID, "ERROR: 等待填充，请稍后再试..."))
			case img := <-queue:
				img.ChatID = msg.From.ID
				m, err := ctx.Caller.Send(img)
				if err != nil {
					return
				}
				_, _ = ctx.Caller.Send(tgba.NewEditMessageReplyMarkup(msg.From.ID, m.MessageID, tgba.NewInlineKeyboardMarkup(tgba.NewInlineKeyboardRow(
					tgba.NewInlineKeyboardButtonURL(
						"UID "+strings.TrimLeft(img.CaptionEntities[1].URL, "https://pixiv.net/u/"),
						img.CaptionEntities[1].URL,
					),
					tgba.NewInlineKeyboardButtonURL(
						"PID "+strings.TrimLeft(img.CaptionEntities[0].URL, "https://pixiv.net/i/"),
						img.CaptionEntities[0].URL,
					)))))
			}
		})
}
