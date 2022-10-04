// Package bilibiliparse b站视频链接解析
package bilibiliparse

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strconv"
	"time"

	rei "github.com/fumiama/ReiBot"
	tgba "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	base14 "github.com/fumiama/go-base16384"

	"github.com/FloatTech/floatbox/binary"
	"github.com/FloatTech/floatbox/web"
	ctrl "github.com/FloatTech/zbpctrl"

	"github.com/FloatTech/ReiBot-Plugin/utils/ctxext"
)

type result struct {
	Data struct {
		Bvid      string `json:"bvid"`
		Aid       int    `json:"aid"`
		Copyright int    `json:"copyright"`
		Pic       string `json:"pic"`
		Title     string `json:"title"`
		Pubdate   int    `json:"pubdate"`
		Ctime     int    `json:"ctime"`
		Rights    struct {
			IsCooperation int `json:"is_cooperation"`
		} `json:"rights"`
		Owner struct {
			Mid  int    `json:"mid"`
			Name string `json:"name"`
		} `json:"owner"`
		Stat struct {
			Aid      int `json:"aid"`
			View     int `json:"view"`
			Danmaku  int `json:"danmaku"`
			Reply    int `json:"reply"`
			Favorite int `json:"favorite"`
			Coin     int `json:"coin"`
			Share    int `json:"share"`
			Like     int `json:"like"`
		} `json:"stat"`
		Staff []struct {
			Title    string `json:"title"`
			Name     string `json:"name"`
			Follower int    `json:"follower"`
		} `json:"staff"`
	} `json:"data"`
}

type owner struct {
	Data struct {
		Card struct {
			Fans int `json:"fans"`
		} `json:"card"`
	} `json:"data"`
}

const (
	videoapi = "https://api.bilibili.com/x/web-interface/view?"
	cardapi  = "http://api.bilibili.com/x/web-interface/card?"
	origin   = "https://www.bilibili.com/video/"
)

var (
	reg   = regexp.MustCompile(`https://www.bilibili.com/video/([0-9a-zA-Z]+)`)
	limit = ctxext.NewLimiterManager(time.Second*10, 1)
)

// 插件主体
func init() {
	en := rei.Register("bilibiliparse", &ctrl.Options[*rei.Ctx]{
		DisableOnDefault: false,
		Help: "b站视频链接解析\n" +
			"- https://www.bilibili.com/video/BV1xx411c7BF | https://www.bilibili.com/video/av1605 | https://b23.tv/I8uzWCA | https://www.bilibili.com/video/bv1xx411c7BF",
	})
	en.OnMessageRegex(`(av[0-9]+|BV[0-9a-zA-Z]{10}){1}`).SetBlock(true).Limit(limit.LimitByGroup).
		Handle(func(ctx *rei.Ctx) {
			id := ctx.State["regex_matched"].([]string)[1]
			photo, err := parse(id)
			if err != nil {
				_, _ = ctx.Caller.Send(tgba.NewMessage(ctx.Message.Chat.ID, "ERROR: "+err.Error()))
				return
			}
			photo.ChatID = ctx.Message.Chat.ID
			_, _ = ctx.Caller.Send(photo)
		})
	en.OnMessageRegex(`https://www.bilibili.com/video/([0-9a-zA-Z]+)`).SetBlock(true).Limit(limit.LimitByGroup).
		Handle(func(ctx *rei.Ctx) {
			id := ctx.State["regex_matched"].([]string)[1]
			photo, err := parse(id)
			if err != nil {
				_, _ = ctx.Caller.Send(tgba.NewMessage(ctx.Message.Chat.ID, "ERROR: "+err.Error()))
				return
			}
			photo.ChatID = ctx.Message.Chat.ID
			_, _ = ctx.Caller.Send(photo)
		})
	en.OnMessageRegex(`(https://b23.tv/[0-9a-zA-Z]+)`).SetBlock(true).Limit(limit.LimitByGroup).
		Handle(func(ctx *rei.Ctx) {
			url := ctx.State["regex_matched"].([]string)[1]
			realurl, err := getrealurl(url)
			if err != nil {
				_, _ = ctx.Caller.Send(tgba.NewMessage(ctx.Message.Chat.ID, "ERROR: "+err.Error()))
				return
			}
			photo, err := parse(cuturl(realurl))
			if err != nil {
				_, _ = ctx.Caller.Send(tgba.NewMessage(ctx.Message.Chat.ID, "ERROR: "+err.Error()))
				return
			}
			photo.ChatID = ctx.Message.Chat.ID
			_, _ = ctx.Caller.Send(photo)
		})
}

// parse 解析视频数据
func parse(id string) (*tgba.PhotoConfig, error) {
	var vid string
	switch id[:2] {
	case "av":
		vid = "aid=" + id[2:]
	case "BV":
		vid = "bvid=" + id
	}
	data, err := web.GetData(videoapi + vid)
	if err != nil {
		return nil, err
	}
	var r result
	err = json.Unmarshal(data, &r)
	if err != nil {
		return nil, err
	}
	title16, err := base14.UTF82UTF16BE(binary.StringToBytes("标题: " + r.Data.Title))
	if err != nil {
		return nil, err
	}
	return &tgba.PhotoConfig{
		BaseFile: tgba.BaseFile{
			File: tgba.FileURL(r.Data.Pic),
		},
		Caption: binary.BytesToString(binary.NewWriterF(func(m *binary.Writer) {
			m.WriteString("标题: ")
			m.WriteString(r.Data.Title)
			_ = m.WriteByte('\n')
			if r.Data.Rights.IsCooperation == 1 {
				for i := 0; i < len(r.Data.Staff); i++ {
					m.WriteString(r.Data.Staff[i].Title)
					m.WriteString(": ")
					m.WriteString(r.Data.Staff[i].Name)
					m.WriteString(", 粉丝: ")
					m.WriteString(row(r.Data.Staff[i].Follower))
					_ = m.WriteByte('\n')
				}
			} else {
				o, err := getcard(r.Data.Owner.Mid)
				if err != nil {
					m.WriteString(err.Error())
				} else {
					m.WriteString("UP主: ")
					m.WriteString(r.Data.Owner.Name)
					m.WriteString(", 粉丝: ")
					m.WriteString(row(o.Data.Card.Fans))
				}
				_ = m.WriteByte('\n')
			}
			m.WriteString("播放: ")
			m.WriteString(row(r.Data.Stat.View))
			m.WriteString(", 弹幕: ")
			m.WriteString(row(r.Data.Stat.Danmaku))
			m.WriteString("\n点赞: ")
			m.WriteString(row(r.Data.Stat.Like))
			m.WriteString(", 投币: ")
			m.WriteString(row(r.Data.Stat.Coin))
			m.WriteString("\n收藏: ")
			m.WriteString(row(r.Data.Stat.Favorite))
			m.WriteString(", 分享: ")
			m.WriteString(row(r.Data.Stat.Share))
			_ = m.WriteByte('\n')
			m.WriteString(origin)
			m.WriteString(id)
		})),
		CaptionEntities: []tgba.MessageEntity{
			{
				Type:   "bold",
				Offset: 0,
				Length: len(title16) / 2,
			},
		},
	}, nil
}

// getrealurl 获取跳转后的链接
func getrealurl(url string) (realurl string, err error) {
	data, err := http.Head(url)
	if err != nil {
		return
	}
	_ = data.Body.Close()
	realurl = data.Request.URL.String()
	return
}

// cuturl 获取aid或者bvid
func cuturl(url string) (id string) {
	if !reg.MatchString(url) {
		return
	}
	return reg.FindStringSubmatch(url)[1]
}

// getcard 获取个人信息
func getcard(mid int) (o owner, err error) {
	data, err := web.GetData(cardapi + "mid=" + strconv.Itoa(mid))
	if err != nil {
		return
	}
	err = json.Unmarshal(data, &o)
	return
}

func row(res int) string {
	if res/10000 != 0 {
		return strconv.FormatFloat(float64(res)/10000, 'f', 2, 64) + "万"
	}
	return strconv.Itoa(res)
}
