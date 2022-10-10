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
	cnapi = "https://genshin.azurewebsites.net/api/speak?format=mp3&text=%s&id=%d"
)

var speakers = map[string]uint{
	"宁宁": 0, "爱瑠": 1, "芳乃": 2, "茉子": 3, "丛雨": 4, "小春": 5, "七海": 6,
	"수아": 0, "미미르": 1, "아린": 2, "연화": 3, "유화": 4, "선배": 5,
	"派蒙": 0, "凯亚": 1, "安柏": 2, "丽莎": 3, "琴": 4, "香菱": 5, "枫原万叶": 6, "迪卢克": 7, "温迪": 8, "可莉": 9, "早柚": 10, "托马": 11, "芭芭拉": 12, "优菈": 13, "云堇": 14, "钟离": 15, "魈": 16, "凝光": 17, "雷电将军": 18, "北斗": 19, "甘雨": 20, "七七": 21, "刻晴": 22, "神里绫华": 23, "戴因斯雷布": 24, "雷泽": 25, "神里绫人": 26, "罗莎莉亚": 27, "阿贝多": 28, "八重神子": 29, "宵宫": 30, "荒泷一斗": 31, "九条裟罗": 32, "夜兰": 33, "珊瑚宫心海": 34, "五郎": 35, "散兵": 36, "女士": 37, "达达利亚": 38, "莫娜": 39, "班尼特": 40, "申鹤": 41, "行秋": 42, "烟绯": 43, "久岐忍": 44, "辛焱": 45, "砂糖": 46, "胡桃": 47, "重云": 48, "菲谢尔": 49, "诺艾尔": 50, "迪奥娜": 51, "鹿野院平藏": 52,
}

func init() {
	en := rei.Register("moegoe", &ctrl.Options[*rei.Ctx]{
		DisableOnDefault: false,
		Help: "moegoe\n" +
			"- 让[宁宁|爱瑠|芳乃|茉子|丛雨|小春|七海]说(日语)\n" +
			"- 让[수아|미미르|아린|연화|유화|선배]说(韩语)\n" +
			"- 让[派蒙|凯亚|安柏|丽莎|琴|香菱|枫原万叶|迪卢克|温迪|可莉|早柚|托马|芭芭拉|优菈|云堇|钟离|魈|凝光|雷电将军|北斗|甘雨|七七|刻晴|神里绫华|雷泽|神里绫人|罗莎莉亚|阿贝多|八重神子|宵宫|荒泷一斗|九条裟罗|夜兰|珊瑚宫心海|五郎|达达利亚|莫娜|班尼特|申鹤|行秋|烟绯|久岐忍|辛焱|砂糖|胡桃|重云|菲谢尔|诺艾尔|迪奥娜|鹿野院平藏]说(中文)",
	}).ApplySingle(ctxext.DefaultSingle)
	en.OnMessageRegex("^让(宁宁|爱瑠|芳乃|茉子|丛雨|小春|七海)说([A-Za-z\\s\\d\u3005\u3040-\u30ff\u4e00-\u9fff\uff11-\uff19\uff21-\uff3a\uff41-\uff5a\uff66-\uff9d\\pP]+)$").Limit(ctxext.LimitByGroup).SetBlock(true).
		Handle(func(ctx *rei.Ctx) {
			speaker := ctx.State["regex_matched"].([]string)[1]
			text := ctx.State["regex_matched"].([]string)[2]
			id := speakers[speaker]
			_, err := ctx.SendAudio(
				tgba.FileURL(fmt.Sprintf(jpapi, url.QueryEscape(text), id)),
				false, speaker+": "+text,
				tgba.MessageEntity{Type: "bold", Length: len([]rune(speaker))},
			)
			if err != nil {
				_, _ = ctx.SendPlainMessage(false, "ERROR: ", err)
				return
			}
		})
	en.OnMessageRegex("^让(수아|미미르|아린|연화|유화|선배)说([A-Za-z\\s\\d\u3131-\u3163\uac00-\ud7ff\\pP]+)$").Limit(ctxext.LimitByGroup).SetBlock(true).
		Handle(func(ctx *rei.Ctx) {
			speaker := ctx.State["regex_matched"].([]string)[1]
			text := ctx.State["regex_matched"].([]string)[2]
			id := speakers[speaker]
			_, err := ctx.SendAudio(
				tgba.FileURL(fmt.Sprintf(krapi, url.QueryEscape(text), id)),
				false, speaker+": "+text,
				tgba.MessageEntity{Type: "bold", Length: len([]rune(speaker))},
			)
			if err != nil {
				_, _ = ctx.SendPlainMessage(false, "ERROR: ", err)
				return
			}
		})
	en.OnMessageRegex("^让(派蒙|凯亚|安柏|丽莎|琴|香菱|枫原万叶|迪卢克|温迪|可莉|早柚|托马|芭芭拉|优菈|云堇|钟离|魈|凝光|雷电将军|北斗|甘雨|七七|刻晴|神里绫华|雷泽|神里绫人|罗莎莉亚|阿贝多|八重神子|宵宫|荒泷一斗|九条裟罗|夜兰|珊瑚宫心海|五郎|达达利亚|莫娜|班尼特|申鹤|行秋|烟绯|久岐忍|辛焱|砂糖|胡桃|重云|菲谢尔|诺艾尔|迪奥娜|鹿野院平藏)说([\\s\u4e00-\u9fa5\\pP]+)$").Limit(ctxext.LimitByGroup).SetBlock(true).
		Handle(func(ctx *rei.Ctx) {
			speaker := ctx.State["regex_matched"].([]string)[1]
			text := ctx.State["regex_matched"].([]string)[2]
			id := speakers[speaker]
			_, err := ctx.SendAudio(
				tgba.FileURL(fmt.Sprintf(cnapi, url.QueryEscape(text), id)),
				false, speaker+": "+text,
				tgba.MessageEntity{Type: "bold", Length: len([]rune(speaker))},
			)
			if err != nil {
				_, _ = ctx.SendPlainMessage(false, "ERROR: ", err)
				return
			}
		})
}
