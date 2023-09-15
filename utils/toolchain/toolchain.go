// Package toolchain can make it more easily to code.
package toolchain

import (
	"fmt"
	"image"
	"io"
	"net/http"
	"regexp"

	rei "github.com/fumiama/ReiBot"
	tgba "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// GetTargetAvatar GetUserTarget ID
func GetTargetAvatar(ctx *rei.Ctx) image.Image {
	getUserName := ctx.Event.Value.(*tgba.Message).From.FirstName
	switch {
	case getUserName == "Group" || getUserName == "Channel":
		getGroupChatConfig := tgba.ChatInfoConfig{ChatConfig: tgba.ChatConfig{ChatID: ctx.Message.Chat.ID}}
		chatGroupInfo, err := ctx.Caller.GetChat(getGroupChatConfig)
		if err != nil {
			panic(err)
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
	default:
		userProfilePhotosConfig := tgba.UserProfilePhotosConfig{
			UserID: ctx.Event.Value.(*tgba.Message).From.ID,
		}
		userProfilePhotos, err := ctx.Caller.GetUserProfilePhotos(userProfilePhotosConfig)
		if err != nil {
			return nil
		}
		getLengthImage := len(userProfilePhotos.Photos)
		// WHY TELEGRAM CAN SET NO TO PUBLIC ADMISSION ON AVATAR????
		if getLengthImage != 0 {
			// offset draw
			photo := userProfilePhotos.Photos[0][0]
			avatar, err := ctx.Caller.GetFileDirectURL(photo.FileID)
			if err != nil {
				panic(err)
			}
			datas, err := http.Get(avatar)
			// avatar
			avatarByteUni, _, _ := image.Decode(datas.Body)
			return avatarByteUni
		}
	}
	return nil
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
