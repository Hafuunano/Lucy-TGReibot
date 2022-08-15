// Package hyaku 百人一首
package hyaku

import (
	"encoding/csv"
	"fmt"
	"math/rand"
	"os"
	"reflect"
	"strconv"
	"unsafe"

	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/binary"
	"github.com/FloatTech/zbputils/file"
	rei "github.com/fumiama/ReiBot"
	tgba "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/FloatTech/ReiBot-Plugin/utils/ctxext"
)

const bed = "https://gitcode.net/u011570312/OguraHyakuninIsshu/-/raw/master/"

//nolint: asciicheck
type line struct {
	番号, 歌人, 上の句, 下の句, 上の句ひらがな, 下の句ひらがな string
}

func (l *line) String() string {
	b := binary.NewWriterF(func(w *binary.Writer) {
		r := reflect.ValueOf(l).Elem().Type()
		for i := 0; i < r.NumField(); i++ {
			switch i {
			case 0:
				w.WriteString("●")
			case 1:
				w.WriteString("◉")
			case 2, 3:
				w.WriteString("○")
			case 4, 5:
				w.WriteString("◎")
			}
			w.WriteString(r.Field(i).Name)
			w.WriteString(": ")
			w.WriteString((*[6]string)(unsafe.Pointer(l))[i])
			w.WriteString("\n")
		}
	})
	return binary.BytesToString(b)
}

var lines [100]*line

func init() {
	engine := rei.Register("hyaku", &ctrl.Options[*rei.Ctx]{
		DisableOnDefault: false,
		Help: "百人一首\n" +
			"- 百人一首(随机发一首)\n" +
			"- 百人一首之n",
		PrivateDataFolder: "hyaku",
	})
	csvfile := engine.DataFolder() + "hyaku.csv"
	go func() {
		if file.IsNotExist(csvfile) {
			err := file.DownloadTo(bed+"小倉百人一首.csv", csvfile, true)
			if err != nil {
				_ = os.Remove(csvfile)
				panic(err)
			}
		}
		f, err := os.Open(csvfile)
		if err != nil {
			panic(err)
		}
		records, err := csv.NewReader(f).ReadAll()
		if err != nil {
			panic(err)
		}
		_ = f.Close()
		records = records[1:] // skip title
		if len(records) != 100 {
			panic("invalid csvfile")
		}
		for j, r := range records {
			if len(r) != 6 {
				panic("invalid csvfile")
			}
			i, err := strconv.Atoi(r[0])
			if err != nil {
				panic(err)
			}
			i--
			if j != i {
				panic("invalid csvfile")
			}
			lines[i] = (*line)(*(*unsafe.Pointer)(unsafe.Pointer(&r)))
		}
	}()
	engine.OnMessageFullMatch("百人一首").SetBlock(true).Limit(ctxext.LimitByUser).Handle(func(ctx *rei.Ctx) {
		i := rand.Intn(100)
		mg := tgba.NewMediaGroup(ctx.Message.Chat.ID, []any{
			tgba.InputMediaPhoto{
				BaseInputMedia: tgba.BaseInputMedia{
					Type:  "photo",
					Media: tgba.FileURL(fmt.Sprintf(bed+"img/%03d.jpg", i+1)),
				},
			},
			tgba.InputMediaPhoto{
				BaseInputMedia: tgba.BaseInputMedia{
					Type:    "photo",
					Media:   tgba.FileURL(fmt.Sprintf(bed+"img/%03d.png", i+1)),
					Caption: lines[i].String(),
				},
			},
		})
		_, _ = ctx.Caller.SendMediaGroup(mg)
	})
	engine.OnMessageRegex(`^百人一首之\s?(\d+)$`).SetBlock(true).Limit(ctxext.LimitByUser).Handle(func(ctx *rei.Ctx) {
		i, err := strconv.Atoi(ctx.State["regex_matched"].([]string)[1])
		if err != nil {
			_, _ = ctx.Caller.Send(tgba.NewMessage(ctx.Message.Chat.ID, "ERROR:"+err.Error()))
			return
		}
		if i > 100 || i < 1 {
			_, _ = ctx.Caller.Send(tgba.NewMessage(ctx.Message.Chat.ID, "ERROR:超出范围"))
			return
		}
		mg := tgba.NewMediaGroup(ctx.Message.Chat.ID, []any{
			tgba.InputMediaPhoto{
				BaseInputMedia: tgba.BaseInputMedia{
					Type:  "photo",
					Media: tgba.FileURL(fmt.Sprintf(bed+"img/%03d.jpg", i)),
				},
			},
			tgba.InputMediaPhoto{
				BaseInputMedia: tgba.BaseInputMedia{
					Type:    "photo",
					Media:   tgba.FileURL(fmt.Sprintf(bed+"img/%03d.png", i)),
					Caption: lines[i-1].String(),
				},
			},
		})
		_, _ = ctx.Caller.SendMediaGroup(mg)
	})
}
