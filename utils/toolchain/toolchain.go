// Package toolchain can make it more faster.
package toolchain

import (
	rei "github.com/fumiama/ReiBot"
	tgba "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"image"
	"net/http"
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
