// Package toolchain can make it more easily to code.
package toolchain

import (
	"fmt"
	"hash/crc64"
	"image"
	"io"
	"math/rand"
	"net/http"
	"regexp"
	"strings"
	"time"
	"unsafe"

	"github.com/FloatTech/floatbox/binary"
	rei "github.com/fumiama/ReiBot"
	tgba "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// FutureEvent 是 ZeroBot 交互式的核心，用于异步获取指定事件
type FutureEvent struct {
	Type     string
	Priority int
	Block    bool
}

// NewFutureEvent 创建一个FutureEvent, 并返回其指针
func NewFutureEvent(Type string, Priority int, Block bool) *FutureEvent {
	return &FutureEvent{
		Type:     Type,
		Priority: Priority,
		Block:    Block,
	}
}

// GetTargetAvatar GetUserTarget ID
func GetTargetAvatar(ctx *rei.Ctx) image.Image {
	getUserName := ctx.Event.Value.(*tgba.Message).From.FirstName
	userID := ctx.Event.Value.(*tgba.Message).From.ID
	if getUserName == "Group" || getUserName == "Channel" {
		userID = ctx.Message.Chat.ID
	}
	getGroupChatConfig := tgba.ChatInfoConfig{ChatConfig: tgba.ChatConfig{ChatID: userID}}
	chatGroupInfo, err := ctx.Caller.GetChat(getGroupChatConfig)
	if err != nil {
		return nil
	}
	if chatGroupInfo.Photo == nil {
		return nil
	}
	chatPhoto := chatGroupInfo.Photo.SmallFileID
	avatar, err := ctx.Caller.GetFileDirectURL(chatPhoto)
	if err != nil {
		panic(err)
	}
	datas, err := http.Get(avatar)
	// avatar
	avatarByteUni, _, _ := image.Decode(datas.Body)
	return avatarByteUni
}

// GetChatUserInfoID GetID and UserName, support Channel | Userself and Annoy Group
func GetChatUserInfoID(ctx *rei.Ctx) (id int64, name string) {
	getUserName := ctx.Event.Value.(*tgba.Message).From.FirstName
	switch {
	case getUserName == "Group" || getUserName == "Channel":
		getGroupChatConfig := tgba.ChatInfoConfig{ChatConfig: tgba.ChatConfig{ChatID: ctx.Message.Chat.ID}}
		chatGroupInfo, err := ctx.Caller.GetChat(getGroupChatConfig)
		if err != nil {
			panic(err)
		}
		return chatGroupInfo.ID, chatGroupInfo.Title
	default:
		return ctx.Event.Value.(*tgba.Message).From.ID, getUserName + " " + ctx.Event.Value.(*tgba.Message).From.LastName
	}
	return 0, ""
}

// GetNickNameFromUsername Use Sniper, not api.
func GetNickNameFromUsername(username string) (name string) {
	// https://github.com/XiaoMengXinX/Telegram_QuoteReply_Bot-Go/blob/master/api/bot.go
	client := &http.Client{}
	req, _ := http.NewRequest("GET", fmt.Sprintf("https://t.me/%s", username), nil)
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	reName := regexp.MustCompile(`<meta property="og:title" content="([^"]*)"`)
	match := reName.FindStringSubmatch(string(body))
	if len(match) > 1 {
		name = match[1]
	}
	pageTitle := ""
	reTitle1 := regexp.MustCompile(`<title>`)
	reTitle2 := regexp.MustCompile(`</title>`)
	start := reTitle1.FindStringIndex(string(body))
	end := reTitle2.FindStringIndex(string(body))
	if start != nil && end != nil {
		pageTitle = string(body)[start[1]:end[0]]
	}

	if pageTitle == name { // 用户不存在
		name = ""
	}
	return
}

// RandSenderPerDayN 每个用户每天随机数
func RandSenderPerDayN(uid int64, n int) int {
	sum := crc64.New(crc64.MakeTable(crc64.ISO))
	sum.Write(binary.StringToBytes(time.Now().Format("20060102")))
	sum.Write((*[8]byte)(unsafe.Pointer(&uid))[:])
	r := rand.New(rand.NewSource(int64(sum.Sum64())))
	return r.Intn(n)
}

// SplitCommandTo Split Command and Adjust To.
func SplitCommandTo(raw string, setCommandStopper int) (splitCommandLen int, splitInfo []string) {
	rawSplit := strings.SplitN(raw, " ", setCommandStopper)
	return len(rawSplit), rawSplit
}

func RequestImageTo(ctx *rei.Ctx, footpoint string) []tgba.PhotoSize {
	msg, ok := ctx.Value.(*tgba.Message)
	if ok && len(msg.Photo) > 0 {
		ctx.State["photos"] = msg.Photo
		return ctx.State["photos"].([]tgba.PhotoSize)
	} else {
		ctx.SendPlainMessage(true, footpoint)
		return nil
	}

}
