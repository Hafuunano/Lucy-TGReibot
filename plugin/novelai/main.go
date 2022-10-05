// Package novelai NovelAI作画
package novelai

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	rei "github.com/fumiama/ReiBot"
	tgba "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	nvai "github.com/FloatTech/AnimeAPI/novelai"
	"github.com/FloatTech/floatbox/binary"
	"github.com/FloatTech/floatbox/file"
	ctrl "github.com/FloatTech/zbpctrl"

	"github.com/FloatTech/ReiBot-Plugin/utils/ctxext"
)

var nv *nvai.NovalAI

func init() {
	en := rei.Register("novelai", &ctrl.Options[*rei.Ctx]{
		DisableOnDefault: false,
		Help: "novelai\n" +
			"- novelai作图 tag1 tag2...\n" +
			"- 设置 novelai key [key]",
		PrivateDataFolder: "novelai",
	}).ApplySingle(ctxext.DefaultSingle)
	keyfile := en.DataFolder() + "key.txt"
	if file.IsExist(keyfile) {
		key, err := os.ReadFile(keyfile)
		if err != nil {
			panic(err)
		}
		nv = nvai.NewNovalAI(binary.BytesToString(key), nvai.NewDefaultPayload())
		err = nv.Login()
		if err != nil {
			panic(err)
		}
	}
	ims.DBPath = en.DataFolder() + "images.db"
	err := ims.Open(time.Hour)
	if err != nil {
		panic(err)
	}
	err = ims.Create("s", &imgstorage{})
	if err != nil {
		panic(err)
	}
	en.OnMessagePrefix("novelai作图").Limit(ctxext.LimitByGroup).SetBlock(true).
		Handle(func(ctx *rei.Ctx) {
			if nv == nil {
				_, _ = ctx.SendPlainMessage(false, "请私聊发送 设置 novelai key [key] 以启用 novelai 作图 (方括号不需要输入)")
				return
			}
			seed, tags, img, err := nv.Draw(strings.TrimSpace(ctx.State["args"].(string)))
			if err != nil {
				_, _ = ctx.SendPlainMessage(false, "ERROR: ", err)
				return
			}
			seedtext := strconv.Itoa(seed)
			fn := tags + " " + seedtext
			err = os.WriteFile(en.DataFolder()+fn+".png", img, 0755)
			if err != nil {
				_, _ = ctx.SendPlainMessage(false, "ERROR: ", err)
				return
			}
			id := idof(fn)
			mu.Lock()
			err = ims.Insert("s", &imgstorage{
				ID:   id,
				Seed: int32(seed),
				Tags: tags,
			})
			mu.Unlock()
			pho := &tgba.PhotoConfig{
				BaseFile: tgba.BaseFile{
					BaseChat: tgba.BaseChat{
						ChatID:           ctx.Message.Chat.ID,
						ReplyToMessageID: ctx.Message.MessageID,
					},
					File: tgba.FileBytes{Bytes: img},
				},
				Caption: "seed: " + seedtext + "\ntags: " + tags,
				CaptionEntities: []tgba.MessageEntity{
					{Type: "bold", Offset: 0, Length: 5},
					{Type: "bold", Offset: 5 + 1 + len(seedtext) + 1, Length: 5},
				},
			}
			if err == nil {
				pho.ReplyMarkup = tgba.NewInlineKeyboardMarkup(
					tgba.NewInlineKeyboardRow(
						tgba.NewInlineKeyboardButtonData(
							"发送原图",
							"nvaiorg"+fmt.Sprintf("%016x", uint64(id)),
						),
					),
				)
			}
			_, _ = ctx.Caller.Send(pho)
		})
	en.OnCallbackQueryRegex(`^nvaiorg([0-9a-f]{16})$`).SetBlock(true).
		Handle(func(ctx *rei.Ctx) {
			fn := ctx.State["regex_matched"].([]string)[1]
			s := imgstorage{}
			id, err := strconv.ParseUint(fn, 16, 64)
			if err != nil {
				_, _ = ctx.Caller.Send(tgba.NewCallbackWithAlert(ctx.Value.(*tgba.CallbackQuery).ID, "ERROR: "+err.Error()))
				return
			}
			ids := strconv.FormatInt(int64(id), 10)
			mu.RLock()
			err = ims.Find("s", &s, "WHERE id="+ids)
			mu.RUnlock()
			if err != nil {
				_, _ = ctx.Caller.Send(tgba.NewCallbackWithAlert(ctx.Value.(*tgba.CallbackQuery).ID, "ERROR: 找不到文件"))
				return
			}
			fn = s.Tags + " " + strconv.Itoa(int(s.Seed))
			f := tgba.NewDocument(ctx.Message.Chat.ID, tgba.FilePath(en.DataFolder()+fn+".png"))
			f.ReplyToMessageID = ctx.Message.MessageID
			_, err = ctx.Caller.Send(&f)
			if err != nil {
				_, _ = ctx.Caller.Send(tgba.NewCallbackWithAlert(ctx.Value.(*tgba.CallbackQuery).ID, "ERROR: "+err.Error()))
				return
			}
			if len(ctx.Message.ReplyMarkup.InlineKeyboard) > 0 {
				_, _ = ctx.Caller.Send(&tgba.EditMessageReplyMarkupConfig{
					BaseEdit: tgba.BaseEdit{
						ChatID:    ctx.Message.Chat.ID,
						MessageID: ctx.Message.MessageID,
					},
				})
			}
			_, _ = ctx.Caller.Send(tgba.NewCallbackWithAlert(ctx.Value.(*tgba.CallbackQuery).ID, "已发送"))
			mu.Lock()
			_ = ims.Del("s", "WHERE id="+ids)
			mu.Unlock()
		})
	en.OnMessageRegex(`^设置\s?novelai\s?key\s?([0-9A-Za-z_]{64})$`, rei.SuperUserPermission, rei.OnlyPrivate).SetBlock(true).
		Handle(func(ctx *rei.Ctx) {
			key := ctx.State["regex_matched"].([]string)[1]
			err := os.WriteFile(keyfile, binary.StringToBytes(key), 0644)
			if err != nil {
				_, _ = ctx.SendPlainMessage(false, "ERROR: ", err)
				return
			}
			nnv := nvai.NewNovalAI(key, nvai.NewDefaultPayload())
			err = nnv.Login()
			if err != nil {
				_, _ = ctx.SendPlainMessage(false, "ERROR: ", err)
				return
			}
			nv = nnv
			_, _ = ctx.SendPlainMessage(false, "成功!")
		})
}
