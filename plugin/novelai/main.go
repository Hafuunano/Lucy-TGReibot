// Package novelai 日韩 VITS 模型拟声
package novelai

import (
	"os"
	"strconv"
	"strings"

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
			_, _ = ctx.SendPhoto(
				tgba.FileBytes{Bytes: img},
				true, "seed: "+seedtext+"\ntags: "+tags,
				tgba.MessageEntity{Type: "bold", Offset: 0, Length: 5},
				tgba.MessageEntity{Type: "bold", Offset: 5 + 1 + len(seedtext) + 1, Length: 5},
			)
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
