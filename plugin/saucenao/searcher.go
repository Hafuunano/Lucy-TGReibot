// Package saucenao P站ID/saucenao/ascii2d搜图
package saucenao

import (
	"fmt"
	"net/http"
	"os"
	"reflect"
	"strconv"

	rei "github.com/fumiama/ReiBot"
	tgba "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/jozsefsallai/gophersauce"
	"github.com/sirupsen/logrus"

	"github.com/FloatTech/AnimeAPI/pixiv"
	"github.com/FloatTech/floatbox/binary"
	"github.com/FloatTech/floatbox/file"
	ctrl "github.com/FloatTech/zbpctrl"
)

var (
	saucenaocli *gophersauce.Client
)

func init() { // 插件主体
	engine := rei.Register("saucenao", &ctrl.Options[*rei.Ctx]{
		DisableOnDefault: false,
		Help: "搜图\n" +
			"- 以图搜图 | 搜索图片 | 以图识图[图片]\n" +
			"- 搜图[P站图片ID]",
		PrivateDataFolder: "saucenao",
	})
	apikeyfile := engine.DataFolder() + "apikey.txt"
	if file.IsExist(apikeyfile) {
		key, err := os.ReadFile(apikeyfile)
		if err != nil {
			panic(err)
		}
		saucenaocli, err = gophersauce.NewClient(&gophersauce.Settings{
			MaxResults: 1,
			APIKey:     binary.BytesToString(key),
		})
		if err != nil {
			panic(err)
		}
	}
	// 根据 PID 搜图
	engine.OnMessageRegex(`^搜图(\d+)$`).SetBlock(true).
		Handle(func(ctx *rei.Ctx) {
			id, _ := strconv.ParseInt(ctx.State["regex_matched"].([]string)[1], 10, 64)
			_, _ = ctx.SendPlainMessage(false, "少女祈祷中......")
			// 获取P站插图信息
			illust, err := pixiv.Works(id)
			if err != nil {
				_, _ = ctx.SendPlainMessage(false, "ERROR: ", err)
				return
			}
			if illust.Pid > 0 {
				name := strconv.FormatInt(illust.Pid, 10)
				txt := fmt.Sprint(
					"标题: ", illust.Title, "\n",
					"插画ID: ", illust.Pid, "\n",
					"画师: ", illust.UserName, "\n",
					"画师ID: ", illust.UserId,
				)
				var imgs []any
				for i := range illust.ImageUrls {
					f := file.BOTPATH + "/" + illust.Path(i)
					n := name + "_p" + strconv.Itoa(i)
					if file.IsNotExist(f) {
						logrus.Debugln("[sausenao]开始下载", n)
						logrus.Debugln("[sausenao]urls:", illust.ImageUrls)
						err := illust.DownloadToCache(i)
						if err != nil {
							logrus.Debugln("[sausenao]下载第", i, "张err:", err)
							continue
						}
					}
					d := tgba.NewInputMediaDocument(tgba.FilePath(f))
					if i == len(illust.ImageUrls)-1 {
						d.Caption = txt
					}
					imgs = append(imgs, d)
				}

				if len(imgs) > 0 {
					_, _ = ctx.Caller.SendMediaGroup(tgba.NewMediaGroup(ctx.Message.Chat.ID, imgs))
				} else {
					// 全部图片下载失败，仅发送文字结果
					_, _ = ctx.SendPlainMessage(true, txt)
				}
			} else {
				_, _ = ctx.SendPlainMessage(false, "图片不存在!")
			}
		})
	// 以图搜图
	engine.OnMessageFullMatchGroup([]string{"以图搜图", "搜索图片", "以图识图"}, rei.MustProvidePhoto("请发送一张图片", "获取图片失败!")).SetBlock(true).
		Handle(func(ctx *rei.Ctx) {
			if saucenaocli == nil {
				_, _ = ctx.SendPlainMessage(false, "请私聊发送 设置 saucenao api key [apikey] 以启用 saucenao 搜图 (方括号不需要输入), key 请前往 https://saucenao.com/user.php?page=search-api 获取")
				return
			}
			// 开始搜索图片
			_, _ = ctx.SendPlainMessage(false, "少女祈祷中...")
			ps := ctx.State["photos"].([]tgba.PhotoSize)
			p := ps[len(ps)-1]
			pic, err := ctx.Caller.GetFileDirectURL(p.FileID)
			if err != nil {
				_, _ = ctx.SendPlainMessage(false, "ERROR: ", err)
				return
			}
			resp, err := saucenaocli.FromURL(pic)
			if err == nil && resp.Count() > 0 {
				result := resp.First()
				s, err := strconv.ParseFloat(result.Header.Similarity, 64)
				if err == nil {
					rr := reflect.ValueOf(&result.Data).Elem()
					b := binary.NewWriterF(func(w *binary.Writer) {
						r := rr.Type()
						for i := 0; i < r.NumField(); i++ {
							if !rr.Field(i).IsZero() {
								w.WriteString("\n")
								w.WriteString(r.Field(i).Name)
								w.WriteString(": ")
								w.WriteString(fmt.Sprint(rr.Field(i).Interface()))
							}
						}
					})
					resp, err := http.Head(result.Header.Thumbnail)
					msg := make([]any, 0, 4)
					if s > 80.0 {
						msg = append(msg, "我有把握是这个!")
					} else {
						msg = append(msg, "也许是这个?")
					}
					var file tgba.RequestFileData
					if err == nil {
						_ = resp.Body.Close()
						if resp.StatusCode == http.StatusOK {
							file = tgba.FileURL(result.Header.Thumbnail)
						} else {
							file = tgba.FileURL(pic)
						}
					} else {
						file = tgba.FileURL(pic)
					}
					msg = append(msg, "\n图源: ", result.Header.IndexName, binary.BytesToString(b))
					_, _ = ctx.SendPhoto(file, false, fmt.Sprint(msg...))
				}
			}
		})
	engine.OnMessageRegex(`^设置\s?saucenao\s?api\s?key\s?([0-9a-f]{40})$`, rei.SuperUserPermission, rei.OnlyPrivate).SetBlock(true).
		Handle(func(ctx *rei.Ctx) {
			var err error
			saucenaocli, err = gophersauce.NewClient(&gophersauce.Settings{
				MaxResults: 1,
				APIKey:     ctx.State["regex_matched"].([]string)[1],
			})
			if err != nil {
				_, _ = ctx.SendPlainMessage(false, "ERROR: ", err)
				return
			}
			err = os.WriteFile(apikeyfile, binary.StringToBytes(saucenaocli.APIKey), 0644)
			if err != nil {
				_, _ = ctx.SendPlainMessage(false, "ERROR: ", err)
				return
			}
			_, _ = ctx.SendPlainMessage(false, "成功!")
		})
}
