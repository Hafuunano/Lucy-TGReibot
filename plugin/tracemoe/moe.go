// Package tracemoe 搜番
package tracemoe

import (
	"fmt"
	"strconv"

	rei "github.com/fumiama/ReiBot"
	tgba "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	trmoe "github.com/fumiama/gotracemoe"

	"github.com/FloatTech/floatbox/binary"
	ctrl "github.com/FloatTech/zbpctrl"
)

var (
	moe = trmoe.NewMoe("")
)

func init() { // 插件主体
	engine := rei.Register("tracemoe", &ctrl.Options[*rei.Ctx]{
		DisableOnDefault: false,
		Help:             "tracemoe\n- 搜番 | 搜索番剧[图片]",
	})
	// 以图搜图
	engine.OnMessageCommand("tracemoe", rei.MustProvidePhoto("请发送一张图片", "获取图片失败!")).SetBlock(true).
		Handle(func(ctx *rei.Ctx) {
			// 开始搜索图片
			_, _ = ctx.SendPlainMessage(false, "少女祈祷中...")
			ps := ctx.State["photos"].([]tgba.PhotoSize)
			pic := ps[len(ps)-1]
			picu, err := ctx.Caller.GetFileDirectURL(pic.FileID)
			if err != nil {
				_, _ = ctx.SendPlainMessage(false, "ERROR: ", err)
				return
			}
			if result, err := moe.Search(picu, true, true); err != nil {
				_, _ = ctx.SendPlainMessage(false, "ERROR: ", err)
			} else if len(result.Result) > 0 {
				r := result.Result[0]
				hint := "我有把握是这个！"
				if r.Similarity < 80 {
					hint = "大概是这个？"
				}
				mf := int(r.From / 60)
				mt := int(r.To / 60)
				sf := r.From - float32(mf*60)
				st := r.To - float32(mt*60)
				_, _ = ctx.Caller.Send(&tgba.PhotoConfig{
					BaseFile: tgba.BaseFile{
						BaseChat: tgba.BaseChat{
							ChatID:           ctx.Message.Chat.ID,
							ReplyToMessageID: ctx.Event.Value.(*tgba.Message).MessageID,
						},
						File: tgba.FileURL(r.Image),
					},
					Caption: binary.BytesToString(binary.NewWriterF(func(m *binary.Writer) {
						m.WriteString(hint)
						_ = m.WriteByte('\n')
						m.WriteString("番剧名: ")
						m.WriteString(r.Anilist.Title.Native)
						_ = m.WriteByte('\n')
						m.WriteString("话数: ")
						m.WriteString(strconv.Itoa(r.Episode))
						_ = m.WriteByte('\n')
						m.WriteString(fmt.Sprint("时间：", mf, ":", sf, "-", mt, ":", st))
					})),
				})
			}
		})
}
