// Package heisi 黑丝
package heisi

import (
	"math/rand"
	"strconv"
	"unsafe"

	rei "github.com/fumiama/ReiBot"
	tgba "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	fbctxext "github.com/FloatTech/floatbox/ctxext"
	ctrl "github.com/FloatTech/zbpctrl"

	"github.com/FloatTech/ReiBot-Plugin/utils/ctxext"
)

var (
	heisiPic []item
	baisiPic []item
	jkPic    []item
	jurPic   []item
	zukPic   []item
	mcnPic   []item
	fileList = [...]string{"heisi.bin", "baisi.bin", "jk.bin", "jur.bin", "zuk.bin", "mcn.bin"}
)

func init() { // 插件主体
	engine := rei.Register("heisi", &ctrl.Options[*rei.Ctx]{
		DisableOnDefault: false,
		Help: "黑丝\n" +
			"- 来点黑丝\n- 来点白丝\n- 来点jk\n- 来点巨乳\n- 来点足控\n- 来点网红",
		PublicDataFolder: "Heisi",
	}).ApplySingle(ctxext.DefaultSingle)

	getbins := fbctxext.DoOnceOnSuccess(filldata(engine.GetLazyData))
	engine.OnMessageFullMatchGroup([]string{"来点黑丝", "来点白丝", "来点jk", "来点巨乳", "来点足控", "来点网红"}, getbins).Limit(ctxext.LimitByGroup).SetBlock(true).
		Handle(func(ctx *rei.Ctx) {
			matched := ctx.State["matched"].(string)
			var pic item
			i := 0
			switch matched {
			case "来点黑丝":
				i = rand.Intn(len(heisiPic))
				pic = heisiPic[i]
			case "来点白丝":
				i = rand.Intn(len(baisiPic))
				pic = baisiPic[i]
			case "来点jk":
				i = rand.Intn(len(jkPic))
				pic = jkPic[i]
			case "来点巨乳":
				i = rand.Intn(len(jurPic))
				pic = jurPic[i]
			case "来点足控":
				i = rand.Intn(len(zukPic))
				pic = zukPic[i]
			case "来点网红":
				i = rand.Intn(len(mcnPic))
				pic = mcnPic[i]
			}
			_, _ = ctx.Caller.Send(&tgba.PhotoConfig{
				BaseFile: tgba.BaseFile{
					BaseChat: tgba.BaseChat{
						ChatID: ctx.Message.Chat.ID,
						ReplyMarkup: tgba.NewInlineKeyboardMarkup(
							tgba.NewInlineKeyboardRow(
								tgba.NewInlineKeyboardButtonData(
									"发送原图",
									matched+strconv.Itoa(i),
								),
							),
						),
					},
					File: tgba.FileURL(pic.String()),
				},
			})
		})

	engine.OnCallbackQueryRegex(`^来点(黑丝|白丝|jk|巨乳|足控|网红)(\d+)$`, ctxext.MustMessageNotNil, getbins).Limit(ctxext.LimitByGroup).SetBlock(true).
		Handle(func(ctx *rei.Ctx) {
			var pic item
			i, err := strconv.Atoi(ctx.State["regex_matched"].([]string)[2])
			if err != nil {
				_, _ = ctx.Caller.Send(tgba.NewCallbackWithAlert(ctx.Value.(*tgba.CallbackQuery).ID, "ERROR: "+err.Error()))
				return
			}
			switch ctx.State["regex_matched"].([]string)[1] {
			case "黑丝":
				pic = heisiPic[i]
			case "白丝":
				pic = baisiPic[i]
			case "jk":
				pic = jkPic[i]
			case "巨乳":
				pic = jurPic[i]
			case "足控":
				pic = zukPic[i]
			case "网红":
				pic = mcnPic[i]
			}
			_, err = ctx.Caller.Send(&tgba.DocumentConfig{
				BaseFile: tgba.BaseFile{
					BaseChat: tgba.BaseChat{
						ChatID:           ctx.Message.Chat.ID,
						ReplyToMessageID: ctx.Message.MessageID,
					},
					File: tgba.FileURL(pic.String()),
				},
			})
			if err != nil {
				_, _ = ctx.Caller.Send(tgba.NewCallbackWithAlert(ctx.Value.(*tgba.CallbackQuery).ID, "ERROR: "+err.Error()))
				return
			}
			if len(ctx.Message.ReplyMarkup.InlineKeyboard) > 0 {
				_, _ = ctx.Caller.Send(tgba.EditMessageReplyMarkupConfig{
					BaseEdit: tgba.BaseEdit{
						ChatID:    ctx.Message.Chat.ID,
						MessageID: ctx.Message.MessageID,
					},
				})
			}
			_, _ = ctx.Caller.Send(tgba.NewCallbackWithAlert(ctx.Value.(*tgba.CallbackQuery).ID, "已发送"))
		})
}

func filldata(getlazy func(string, bool) ([]byte, error)) rei.Rule {
	return func(ctx *rei.Ctx) bool {
		for i, filePath := range fileList {
			data, err := getlazy(filePath, true)
			if err != nil {
				_, _ = ctx.SendPlainMessage(false, "ERROR: ", err)
				return false
			}
			if len(data)%10 != 0 {
				_, _ = ctx.SendPlainMessage(false, "ERROR: invalid data "+strconv.Itoa(i))
				return false
			}
			s := (*slice)(unsafe.Pointer(&data))
			s.len /= 10
			s.cap /= 10
			switch i {
			case 0:
				heisiPic = *(*[]item)(unsafe.Pointer(s))
			case 1:
				baisiPic = *(*[]item)(unsafe.Pointer(s))
			case 2:
				jkPic = *(*[]item)(unsafe.Pointer(s))
			case 3:
				jurPic = *(*[]item)(unsafe.Pointer(s))
			case 4:
				zukPic = *(*[]item)(unsafe.Pointer(s))
			case 5:
				mcnPic = *(*[]item)(unsafe.Pointer(s))
			}
		}
		return true
	}
}
