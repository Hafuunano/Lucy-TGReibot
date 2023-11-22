// Package toolchain can make it more easily to code.
package toolchain

import (
	"encoding/json"
	"fmt"
	"hash/crc64"
	"image"
	"io"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unsafe"

	"github.com/FloatTech/ReiBot-Plugin/utils/CoreFactory"
	"github.com/FloatTech/ReiBot-Plugin/utils/userlist"
	"github.com/FloatTech/floatbox/binary"
	"github.com/FloatTech/floatbox/file"
	rei "github.com/fumiama/ReiBot"
	tgba "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/wdvxdr1123/ZeroBot/extension/rate"
)

var OnHoldSaver = rate.NewManager[int64](time.Hour*24, 1)  // only update once.
var OnGroupSaver = rate.NewManager[int64](time.Hour*24, 1) // only update once.

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
		return nil
	}
	datas, err := http.Get(avatar)
	// avatar
	avatarByteUni, _, _ := image.Decode(datas.Body)
	return avatarByteUni
}

func GetReferTargetAvatar(ctx *rei.Ctx, uid int64) string {
	getGroupChatConfig := tgba.ChatInfoConfig{ChatConfig: tgba.ChatConfig{ChatID: uid}}
	chatGroupInfo, err := ctx.Caller.GetChat(getGroupChatConfig)
	if err != nil {
		return ""
	}
	if chatGroupInfo.Photo == nil {
		return ""
	}
	chatPhoto := chatGroupInfo.Photo.SmallFileID
	avatar, err := ctx.Caller.GetFileDirectURL(chatPhoto)
	return avatar
}

// GetChatUserInfoID GetID and UserName, support Channel | Userself and Annoy Group
func GetChatUserInfoID(ctx *rei.Ctx) (id int64, name string) {
	getUserName := ctx.Event.Value.(*tgba.Message).From.FirstName
	switch {
	case getUserName == "Group" || getUserName == "Channel":
		getGroupChatConfig := tgba.ChatInfoConfig{ChatConfig: tgba.ChatConfig{ChatID: ctx.Message.Chat.ID}}
		chatGroupInfo, err := ctx.Caller.GetChat(getGroupChatConfig)
		if err != nil {
			return
		}
		return chatGroupInfo.ID, chatGroupInfo.Title
	default:
		return ctx.Event.Value.(*tgba.Message).From.ID, getUserName + " " + ctx.Event.Value.(*tgba.Message).From.LastName
	}

}

// GetThisGroupID Get Group ID
func GetThisGroupID(ctx *rei.Ctx) (id int64) {
	if !ctx.Message.Chat.IsGroup() && !ctx.Message.Chat.IsSuperGroup() {
		return 0
	}
	return ctx.Message.Chat.ID
}

// GetNickNameFromUsername Use Sniper, not api.
func GetNickNameFromUsername(username string) (name string) {
	// https://github.com/XiaoMengXinX/Telegram_QuoteReply_Bot-Go/blob/master/api/bot.go
	if strings.HasPrefix(username, "@") {
		username = strings.Replace(username, "@", "", 1)
	}
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

func GetNickNameFromUserid(ctx *rei.Ctx, userid int64) string {
	data := CoreFactory.GetUserSampleUserinfobyid(userid)
	if data == nil {
		return ""
	}
	return GetNickNameFromUsername(data.UserName)
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

// RequestImageTo Request Image and return PhotoSize To handle.
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

// FastSendRandMuiltText Send Muilt Text to help/
func FastSendRandMuiltText(ctx *rei.Ctx, raw ...string) error {
	_, err := ctx.SendPlainMessage(true, raw[rand.Intn(len(raw))])
	return err
}

// FastSendRandMuiltPic Send Multi picture to help/
func FastSendRandMuiltPic(ctx *rei.Ctx, raw ...string) error {
	_, err := ctx.SendPhoto(tgba.FilePath(raw[rand.Intn(len(raw))]), true, "")
	return err
}

// StringInArray 检查列表是否有关键词 https://github.com/Kyomotoi/go-ATRI
func StringInArray(aim string, list []string) bool {
	for _, i := range list {
		if i == aim {
			return true
		}
	}
	return false
}

// StoreUserNickname store names in jsons
func StoreUserNickname(userID string, nickname string) error {
	var userNicknameData map[string]interface{}
	filePath := file.BOTPATH + "/data/rbp/users.json"
	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			_ = os.WriteFile(filePath, []byte("{}"), 0777)
		} else {
			return err
		}
	}
	_ = json.Unmarshal(data, &userNicknameData)
	userNicknameData[userID] = nickname // setdata.
	newData, err := json.Marshal(userNicknameData)
	if err != nil {
		return err
	}
	_ = os.WriteFile(filePath, newData, 0777)
	return nil
}

// LoadUserNickname Load UserNames, it will work on simai plugin
func LoadUserNickname(userID string) string {
	var d map[string]string
	// read main files
	filePath := file.BOTPATH + "/data/rbp/users.json"
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "你"
	}
	err = json.Unmarshal(data, &d)
	if err != nil {
		return "你"
	}
	result := d[userID]
	if result == "" {
		result = "你"
	}
	return result
}

func GetBotIsAdminInThisGroup(ctx *rei.Ctx) bool {
	getSelfMember, err := ctx.Caller.GetChatMember(
		tgba.GetChatMemberConfig{
			ChatConfigWithUser: tgba.ChatConfigWithUser{
				ChatID: ctx.Message.Chat.ID,
				UserID: ctx.Caller.Self.ID,
			},
		},
	)
	if err != nil {
		return false
	}
	return getSelfMember.IsCreator() || getSelfMember.IsAdministrator()
}

func GetTheTargetIsNormalUser(ctx *rei.Ctx) bool {
	// stop channel to take part in this.
	getUserChannelStatus := ctx.Event.Value.(*tgba.Message).From.FirstName
	if getUserChannelStatus == "Group" || getUserChannelStatus == "Channel" || ctx.Message.From.ID == 777000 { // unknownUser.
		return false
	}
	return true
}

func IsTargetSettedUserName(ctx *rei.Ctx) bool {
	return ctx.Message.From.UserName != ""
}

// FastSaveUserStatus I hope this will not ruin my machine. (
func FastSaveUserStatus(ctx *rei.Ctx) bool {
	// only save normal user

	if !OnHoldSaver.Load(ctx.Message.From.ID).Acquire() || !GetTheTargetIsNormalUser(ctx) || !IsTargetSettedUserName(ctx) {
		// save database onload time.
		return false
	}
	CoreFactory.StoreUserDataBase(ctx.Message.From.ID, ctx.Message.From.UserName)
	return true
}

func FastSaveUserGroupList(ctx *rei.Ctx) {
	if !OnGroupSaver.Load(ctx.Message.From.ID).Acquire() || !GetTheTargetIsNormalUser(ctx) || GetThisGroupID(ctx) == 0 {
		return
	}
	userlist.SaveUserOnList(strconv.FormatInt(ctx.Message.From.ID, 10), strconv.FormatInt(ctx.Message.Chat.ID, 10), ctx.Message.From.UserName)

}

// CheckIfthisUserInThisGroup Check the user if in this group.
func CheckIfthisUserInThisGroup(userID int64, ctx *rei.Ctx) bool {
	group := GetThisGroupID(ctx)
	if group == 0 {
		// not a group.
		return false
	}
	getResult, err := ctx.Caller.GetChatMember(
		tgba.GetChatMemberConfig{
			ChatConfigWithUser: tgba.ChatConfigWithUser{
				ChatID: group,
				UserID: userID,
			},
		},
	)
	if err != nil {
		return false
	}
	if getResult.User != nil {
		return true
	}
	return false
}

// ListEntitiesMention List Entities and return a simple list with user.
func ListEntitiesMention(ctx *rei.Ctx) []string {
	var tempList []string
	msg := ctx.Message.Text
	for _, entity := range ctx.Message.Entities {
		if entity.Type == "mention" {
			mentionText := msg[entity.Offset : entity.Offset+entity.Length]
			tempList = append(tempList, mentionText)
		}
	}
	return tempList
}

// GetUserIDFromUserName with @, only works when the data saved.
func GetUserIDFromUserName(ctx *rei.Ctx, userName string) int64 {
	getUserData := CoreFactory.GetUserSampleUserinfo(strings.Replace(userName, "@", "", 1))
	if getUserData.UserID == 0 {
		return 0
	}
	// check the user is in group?
	if !CheckIfthisUserInThisGroup(getUserData.UserID, ctx) {
		return 0
	}
	return getUserData.UserID
}

// ExtractNumbers Extract Numbers by using regexp.
func ExtractNumbers(text string) int64 {
	re := regexp.MustCompile(`\d+`)
	numbers := re.FindAllString(text, 1)
	num, _ := strconv.ParseInt(numbers[0], 10, 64)
	return num
}

func GetUserNickNameByIDInGroup(ctx *rei.Ctx, id int64) string {
	if !CheckIfthisUserInThisGroup(id, ctx) {
		return ""
	}
	chatPrefer, err := ctx.Caller.GetChatMember(tgba.GetChatMemberConfig{ChatConfigWithUser: tgba.ChatConfigWithUser{ChatID: ctx.Message.Chat.ID, UserID: id}})
	if err != nil {
		panic(err)
	}
	return chatPrefer.User.FirstName + chatPrefer.User.LastName

}
