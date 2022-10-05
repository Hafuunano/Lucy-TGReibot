// Package lolicon 基于 https://api.lolicon.app 随机图片
package lolicon

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	rei "github.com/fumiama/ReiBot"
	tgba "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	base14 "github.com/fumiama/go-base16384"

	"github.com/FloatTech/floatbox/binary"
	"github.com/FloatTech/floatbox/math"
	"github.com/FloatTech/floatbox/web"
	ctrl "github.com/FloatTech/zbpctrl"
)

const (
	api      = "https://api.lolicon.app/setu/v2?r18=2&proxy=i.pixiv.cat"
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
	en := rei.Register("lolicon", &ctrl.Options[*rei.Ctx]{
		DisableOnDefault: false,
		Help: "lolicon\n" +
			"- 来份萝莉",
	}).ApplySingle(rei.NewSingle(
		rei.WithKeyFn(func(ctx *rei.Ctx) int64 {
			return ctx.Message.Chat.ID
		}),
		rei.WithPostFn[int64](func(ctx *rei.Ctx) {
			_, _ = ctx.SendPlainMessage(false, "有其他萝莉操作正在执行中, 不要着急哦")
		})))
	en.OnMessageFullMatch("来份萝莉").SetBlock(true).
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
					caption.WriteByte('\n')
					for _, t := range r.Data[0].Tags {
						caption.WriteByte(' ')
						caption.WriteString(t)
					}
					uidlink := "https://pixiv.net/u/" + strconv.Itoa(r.Data[0].UID)
					pidlink := "https://pixiv.net/i/" + strconv.Itoa(r.Data[0].Pid)
					title16, err := base14.UTF82UTF16BE(binary.StringToBytes(r.Data[0].Title))
					if err != nil {
						continue
					}
					auth16, err := base14.UTF82UTF16BE(binary.StringToBytes(r.Data[0].Author))
					if err != nil {
						continue
					}
					_, imgcallbackdata, _ := strings.Cut(r.Data[0].Urls.Original, "/img-original/img/")
					queue <- &tgba.PhotoConfig{
						BaseFile: tgba.BaseFile{
							BaseChat: tgba.BaseChat{
								ReplyMarkup: tgba.NewInlineKeyboardMarkup(
									tgba.NewInlineKeyboardRow(
										tgba.NewInlineKeyboardButtonURL(
											"UID "+strconv.Itoa(r.Data[0].UID),
											uidlink,
										),
										tgba.NewInlineKeyboardButtonURL(
											"PID "+strconv.Itoa(r.Data[0].Pid),
											pidlink,
										),
									),
									tgba.NewInlineKeyboardRow(
										tgba.NewInlineKeyboardButtonData(
											"发送原图",
											imgcallbackdata,
										),
									),
								),
							},
							File: tgba.FileURL(r.Data[0].Urls.Original),
						},
						Caption: caption.String(),
						CaptionEntities: []tgba.MessageEntity{
							{
								Type:   "bold",
								Offset: 0,
								Length: len(title16) / 2,
							},
							{
								Type:   "underline",
								Offset: len(title16)/2 + 1,
								Length: len(auth16)/2 + 1,
							},
						},
					}
				}
			}()
			select {
			case <-time.After(time.Minute):
				_, _ = ctx.SendPlainMessage(false, "ERROR: 等待填充，请稍后再试...")
			case img := <-queue:
				_, err := ctx.Send(false, img)
				if err != nil {
					_, _ = ctx.SendPlainMessage(false, "ERROR: ", err)
					return
				}
			}
		})
	en.OnCallbackQueryRegex(`^(\d{4}/\d{2}/\d{2}/\d{2}/\d{2}/\d{2}/\d+_p\d+.\w+){1}$`).SetBlock(true).
		Handle(func(ctx *rei.Ctx) {
			if len(ctx.Message.ReplyMarkup.InlineKeyboard) > 1 {
				ctx.Message.ReplyMarkup.InlineKeyboard = ctx.Message.ReplyMarkup.InlineKeyboard[:1]
				_, _ = ctx.Caller.Send(tgba.NewEditMessageReplyMarkup(ctx.Message.Chat.ID, ctx.Message.MessageID, *ctx.Message.ReplyMarkup))
			}
			f := tgba.NewDocument(ctx.Message.Chat.ID, tgba.FileURL("https://i.pixiv.cat/img-original/img/"+ctx.State["regex_matched"].([]string)[1]))
			f.ReplyToMessageID = ctx.Message.MessageID
			_, err := ctx.Caller.Send(&f)
			if err != nil {
				_, _ = ctx.Caller.Send(tgba.NewCallbackWithAlert(ctx.Value.(*tgba.CallbackQuery).ID, "ERROR: "+err.Error()))
				return
			}
			_, _ = ctx.Caller.Send(tgba.NewCallbackWithAlert(ctx.Value.(*tgba.CallbackQuery).ID, "已发送"))
		})
}
