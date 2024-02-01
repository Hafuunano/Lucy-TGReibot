package mai

import (
	"bytes"
	"encoding/json"
	"image"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/FloatTech/ReiBot-Plugin/utils/toolchain"
	"github.com/FloatTech/floatbox/web"
	"github.com/FloatTech/gg"
	ctrl "github.com/FloatTech/zbpctrl"
	rei "github.com/fumiama/ReiBot"
	tgba "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/tidwall/gjson"
	"github.com/wdvxdr1123/ZeroBot/utils/helper"
)

var engine = rei.Register("mai", &ctrl.Options[*rei.Ctx]{
	DisableOnDefault:  false,
	Help:              "maimai - bind Username / maimai b50 render",
	PrivateDataFolder: "mai",
})

func init() {
	engine.OnMessageCommand("mai").SetBlock(true).Handle(func(ctx *rei.Ctx) {
		getMsg := ctx.Message.Text
		getSplitLength, getSplitStringList := toolchain.SplitCommandTo(getMsg, 3)
		if getSplitLength >= 2 {
			switch {
			case getSplitStringList[1] == "bind":
				if getSplitLength < 3 {
					ctx.SendPlainMessage(true, "参数提供不足")
					return
				}
				BindUserToMaimai(ctx, getSplitStringList[2])
			case getSplitStringList[1] == "lxbind":
				if getSplitLength < 3 {
					ctx.SendPlainMessage(true, "参数提供不足")
					return
				}
				Toint64, err := strconv.ParseInt(getSplitStringList[2], 10, 64)
				if err != nil {
					ctx.SendPlainMessage(true, "参数的FriendCode为非法")
					return
				}
				BindFriendCode(ctx, Toint64)
			case getSplitStringList[1] == "userbind":
				if getSplitLength < 3 {
					ctx.SendPlainMessage(true, "参数提供不足, /mai userbind <maiTempID> ")
					return
				}
				getID, _ := toolchain.GetChatUserInfoID(ctx)
				userID := GetWahlapUserID(getSplitStringList[2])
				if userID == -1 {
					ctx.SendPlainMessage(true, "ID 无效或者是过期 ，请使用新的ID或者再次尝试")
					return
				}
				ctx.SendPlainMessage(true, "绑定成功~")
				FormatUserIDDatabase(getID, strconv.FormatInt(userID, 10)).BindUserIDDataBase()
			case getSplitStringList[1] == "unlock":
				getID, _ := toolchain.GetChatUserInfoID(ctx)
				getMaiID := GetUserIDFromDatabase(getID)
				if getMaiID.Userid == "" {
					ctx.SendPlainMessage(true, "没有绑定~ 绑定方式: /mai userbind <maiTempID>")
					return
				}
				getCodeRaw, err := strconv.ParseInt(getMaiID.Userid, 10, 64)
				if err != nil {
					panic(err)
				}
				getCodeStat := Logout(getCodeRaw)
				getCode := gjson.Get(getCodeStat, "returnCode").Int()
				if getCode == 1 {
					ctx.SendPlainMessage(true, "发信成功，服务器返回正常 , 如果未生效请重新尝试")
				} else {
					ctx.SendPlainMessage(true, "发信失败，如果未生效请重新尝试")
				}
			case getSplitStringList[1] == "plate":
				if getSplitLength == 2 {
					SetUserPlateToLocal(ctx, "")
					return
				}
				SetUserPlateToLocal(ctx, getSplitStringList[2])
			case getSplitStringList[1] == "upload":
				// uploadImage
				images := toolchain.RequestImageTo(ctx, "请发送指令同时提供一张图片，图片大小比例适应为6:1 (1260x210) ,如果图片不适应将会自动剪辑到合适大小")
				if images == nil {
					return
				}
				HandlerUserSetsCustomImage(ctx, images)
			case getSplitStringList[1] == "remove":
				RemoveUserLocalCustomImage(ctx)
			case getSplitStringList[1] == "defplate":
				if getSplitLength < 3 {
					SetUserDefaultPlateToDatabase(ctx, "")
					return
				}
				SetUserDefaultPlateToDatabase(ctx, getSplitStringList[2])
			case getSplitStringList[1] == "switch":
				MaimaiSwitcherService(ctx)
			case getSplitStringList[1] == "region":
				getID, _ := toolchain.GetChatUserInfoID(ctx)
				getMaiID := GetUserIDFromDatabase(getID)
				if getMaiID.Userid == "" {
					ctx.SendPlainMessage(true, "没有绑定UserID~ 绑定方式: /mai userbind <maiTempID>")
					return
				}
				getIntID, _ := strconv.ParseInt(getMaiID.Userid, 10, 64)
				getReplyMsg := GetUserRegion(getIntID)
				if strings.Contains(getReplyMsg, "{") == false {
					ctx.SendPlainMessage(true, "返回了错误.png, ERROR:"+getReplyMsg)
					return
				}
				var MixedMagic GetUserRegionStruct
				json.Unmarshal(helper.StringToBytes(getReplyMsg), &MixedMagic)
				var returnText string
				for _, onlistLoader := range MixedMagic.UserRegionList {
					returnText = returnText + MixedRegionWriter(onlistLoader.RegionId-1, onlistLoader.PlayCount, onlistLoader.Created) + "\n\n"
				}
				if returnText == "" {
					ctx.SendPlainMessage(true, "目前 Lucy 没有查到您的游玩记录哦~")
					return
				}
				ctx.SendPlainMessage(true, "目前查询到您的游玩记录如下: \n\n"+returnText)
			case getSplitStringList[1] == "status":
				// getWebStatus
				getWebStatus := ReturnWebStatus()
				getZlibError := ReturnZlibError()
				// 20s one request.
				var getLucyRespHandler int
				if getZlibError.Full.Field3 < 180 {
					getLucyRespHandler = getZlibError.Full.Field3
				} else {
					getLucyRespHandler = getZlibError.Full.Field3 - 180
				}
				getLucyRespHandlerStr := strconv.Itoa(getLucyRespHandler)
				getZlibWord := "Zlib 压缩跳过率: \n" + "10mins (" + ConvertZlib(getZlibError.ZlibError.Field1, getZlibError.Full.Field1) + " Loss)\n" + "30mins (" + ConvertZlib(getZlibError.ZlibError.Field2, getZlibError.Full.Field2) + " Loss)\n" + "60mins (" + ConvertZlib(getZlibError.ZlibError.Field3, getZlibError.Full.Field3) + " Loss)\n"
				getWebStatusCount := "Web Uptime Ping:\n * MaimaiDXCN: " + ConvertFloat(getWebStatus.Details.MaimaiDXCN.Uptime*100) + "%\n * MaimaiDXCN Main Server: " + ConvertFloat(getWebStatus.Details.MaimaiDXCNMain.Uptime*100) + "%\n * MaimaiDXCN Title Server: " + ConvertFloat(float64(getWebStatus.Details.MaimaiDXCNTitle.Uptime*100)) + "%\n * MaimaiDXCN Update Server: " + ConvertFloat(float64(getWebStatus.Details.MaimaiDXCNUpdate.Uptime*100)) + "%\n * MaimaiDXCN NetLogin Server: " + ConvertFloat(getWebStatus.Details.MaimaiDXCNNetLogin.Uptime*100) + "%\n * MaimaiDXCN Net Server: " + ConvertFloat(getWebStatus.Details.MaimaiDXCNDXNet.Uptime*100) + "%\n"
				ctx.SendPlainMessage(true, "* Zlib 压缩跳过率可以很好的反馈当前 MaiNet (Wahlap Service) 当前负载的情况\n* Web Uptime Ping 则可以反馈 MaiNet 在外部原因(DDOS) 下造成的负载详情 ( 100% 即代表服务器为稳定, uptime 越低则代表可用性越差 ) \n* 在 1小时 内，Lucy 共处理了 "+getLucyRespHandlerStr+"次 请求💫，其中详细数据如下:\n\n"+getZlibWord+getWebStatusCount+"\n* Title Server 爆炸 容易造成数据获取失败\n* Zlib 3% Loss 以下则 基本上可以正常游玩\n* 10% Loss 则会有明显断网现象(请准备小黑屋工具)\n* 30% Loss 则无法正常游玩(即使使用小黑屋工具) ")
			case getSplitStringList[1] == "update":
				getID, _ := toolchain.GetChatUserInfoID(ctx)
				getMaiID := GetUserIDFromDatabase(getID)
				if getMaiID.Userid == "" {
					ctx.SendPlainMessage(true, "没有绑定UserID~ 绑定方式: /mai userbind <maiTempID>")
					return
				}
				getTokenId := GetUserToken(strconv.FormatInt(getID, 10))
				if getTokenId == "" {
					ctx.SendPlainMessage(true, "请先 /mai tokenbind <token> 绑定水鱼查分器哦")
					return
				}
				if !CheckTheTicketIsValid(getTokenId) {
					ctx.SendPlainMessage(true, "此 Token 不合法 ，请重新绑定")
					return
				}
				// token is valid, get data.
				getIntID, _ := strconv.ParseInt(getMaiID.Userid, 10, 64)
				// getFullData := GetMusicList(getIntID, 0, 600)
				getFullData := GetMusicList(getIntID, 0, 1000)
				var unmashellData UserMusicListStruct
				json.Unmarshal(helper.StringToBytes(getFullData), &unmashellData)
				getFullDataStruct := convert(unmashellData)
				jsonDumper := getFullDataStruct
				jsonDumperFull, err := json.Marshal(jsonDumper)
				if err != nil {
					panic(err)
				}
				// upload to diving fish api
				req, err := http.NewRequest("POST", "https://www.diving-fish.com/api/maimaidxprober/player/update_records", bytes.NewBuffer(jsonDumperFull))
				if err != nil {
					// Handle error
					panic(err)
				}
				req.Header.Set("Import-Token", getTokenId)
				req.Header.Set("Content-Type", "application/json")

				client := &http.Client{}
				resp, err := client.Do(req)
				if err != nil {
					panic(err)
				}
				//	NewReader, err := io.ReadAll(resp.Body)
				if err != nil {
					panic(err)
				}
				ctx.SendPlainMessage(true, "Update CODE:"+strconv.Itoa(resp.StatusCode))
			case getSplitStringList[1] == "tokenbind":
				if getSplitLength == 2 {
					ctx.SendPlainMessage(true, "缺少参数哦~ qwq")
					return
				}
				getID, _ := toolchain.GetChatUserInfoID(ctx)
				FormatUserToken(strconv.FormatInt(getID, 10), getSplitStringList[2]).BindUserToken()
				ctx.SendPlainMessage(true, "绑定成功~")
			case getSplitStringList[1] == "ticket":
				if getSplitLength == 2 {
					ctx.SendPlainMessage(true, "缺少参数哦~ qwq")
					return
				}
				getID, _ := toolchain.GetChatUserInfoID(ctx)
				getMaiID := GetUserIDFromDatabase(getID)
				if getMaiID.Userid == "" {
					ctx.SendPlainMessage(true, "没有绑定~ 绑定方式: /mai userbind <maiTempID>")
					return
				}
				ticketToFormatNum, err := strconv.ParseInt(getSplitStringList[2], 10, 64)
				if err != nil {
					ctx.SendPlainMessage(true, "传输的数据不合法~")
					return
				}
				getMaiIDInt64, err := strconv.ParseInt(getMaiID.Userid, 10, 64)
				getCode := TicketGain(getMaiIDInt64, int(ticketToFormatNum))
				switch {
				case getCode == 500:
					ctx.SendPlainMessage(true, "TicketID 为非限定内，可使用 2 | 3 | 5 | 20010 | 20020 ")
					return
				case getCode == 102:
					ctx.SendPlainMessage(true, "请在 华立公众号 生成一次二维码 后使用")
					return
				case getCode == 105:
					ctx.SendPlainMessage(true, "已经有了未使用的Ticket了x")
					return
				case getCode == 200:
					ctx.SendPlainMessage(true, "使用成功~将在下一次游戏时自动使用")
					return
				}
			default:
				ctx.SendPlainMessage(true, "未知的指令或者指令出现错误~")
			}
		} else {
			MaimaiRenderBase(ctx)
		}
	})
}

// BindFriendCode Bind FriendCode To Users
func BindFriendCode(ctx *rei.Ctx, bindCode int64) {
	getUserID, _ := toolchain.GetChatUserInfoID(ctx)
	FormatMaimaiFriendCode(bindCode, getUserID).BindUserFriendCode()
	ctx.SendPlainMessage(true, "绑定成功~！")
}

// BindUserToMaimai Bind UserMaiMaiID
func BindUserToMaimai(ctx *rei.Ctx, bindName string) {
	getUserID, _ := toolchain.GetChatUserInfoID(ctx)
	FormatUserDataBase(getUserID, GetUserPlateInfoFromDatabase(getUserID), GetUserDefaultBackgroundDataFromDatabase(getUserID), bindName).BindUserDataBase()
	ctx.SendPlainMessage(true, "绑定成功~！")
}

// SetUserPlateToLocal Set Default Plate to Local
func SetUserPlateToLocal(ctx *rei.Ctx, plateID string) {
	getUserID, _ := toolchain.GetChatUserInfoID(ctx)
	FormatUserDataBase(getUserID, plateID, GetUserDefaultBackgroundDataFromDatabase(getUserID), GetUserInfoNameFromDatabase(getUserID)).BindUserDataBase()
	ctx.SendPlainMessage(true, "好哦~ 是个好名称w")
}

// HandlerUserSetsCustomImage  Handle User Custom Image and Send To Local
func HandlerUserSetsCustomImage(ctx *rei.Ctx, ps []tgba.PhotoSize) {
	getUserID, _ := toolchain.GetChatUserInfoID(ctx)
	pic := ps[len(ps)-1]
	picu, err := ctx.Caller.GetFileDirectURL(pic.FileID)
	imageData, err := web.GetData(picu)
	if err != nil {
		return
	}
	getRaw, _, err := image.Decode(bytes.NewReader(imageData))
	if err != nil {
		return
	}
	// pic Handler
	getRenderPlatePicRaw := gg.NewContext(1260, 210)
	getRenderPlatePicRaw.DrawRoundedRectangle(0, 0, 1260, 210, 10)
	getRenderPlatePicRaw.Clip()
	getHeight := getRaw.Bounds().Dy()
	getLength := getRaw.Bounds().Dx()
	var getHeightHandler, getLengthHandler int
	switch {
	case getHeight < 210 && getLength < 1260:
		getRaw = Resize(getRaw, 1260, 210)
		getHeightHandler = 0
		getLengthHandler = 0
	case getHeight < 210:
		getRaw = Resize(getRaw, getLength, 210)
		getHeightHandler = 0
		getLengthHandler = (getRaw.Bounds().Dx() - 1260) / 3 * -1
	case getLength < 1260:
		getRaw = Resize(getRaw, 1260, getHeight)
		getHeightHandler = (getRaw.Bounds().Dy() - 210) / 3 * -1
		getLengthHandler = 0
	default:
		getLengthHandler = (getRaw.Bounds().Dx() - 1260) / 3 * -1
		getHeightHandler = (getRaw.Bounds().Dy() - 210) / 3 * -1
	}
	getRenderPlatePicRaw.DrawImage(getRaw, getLengthHandler, getHeightHandler)
	getRenderPlatePicRaw.Fill()
	// save.
	_ = getRenderPlatePicRaw.SavePNG(userPlate + strconv.Itoa(int(getUserID)) + ".png")
	ctx.SendPlainMessage(true, "已经存入了哦w~")
}

// RemoveUserLocalCustomImage Remove User Local Image.
func RemoveUserLocalCustomImage(ctx *rei.Ctx) {
	getUserID, _ := toolchain.GetChatUserInfoID(ctx)
	_ = os.Remove(userPlate + strconv.Itoa(int(getUserID)) + ".png")
	ctx.SendPlainMessage(true, "已经移除了~ ")
}

// SetUserDefaultPlateToDatabase Set Default plateID To Database.
func SetUserDefaultPlateToDatabase(ctx *rei.Ctx, plateName string) {
	getUserID, _ := toolchain.GetChatUserInfoID(ctx)
	if plateName == "" {
		FormatUserDataBase(getUserID, GetUserPlateInfoFromDatabase(getUserID), "", GetUserInfoNameFromDatabase(getUserID)).BindUserDataBase()
		ctx.SendPlainMessage(true, "已经删除了预设~")
		return
	}
	getDefaultInfo := plateName
	_, err := GetDefaultPlate(getDefaultInfo)
	if err != nil {
		ctx.SendPlainMessage(true, "设定的预设不正确")
		return
	}
	FormatUserDataBase(getUserID, GetUserPlateInfoFromDatabase(getUserID), getDefaultInfo, GetUserInfoNameFromDatabase(getUserID)).BindUserDataBase()
	ctx.SendPlainMessage(true, "已经设定好了哦w~ ")
}

// MaimaiRenderBase Render Base Maimai B50.
func MaimaiRenderBase(ctx *rei.Ctx) {
	// check the user using.
	getUserID, _ := toolchain.GetChatUserInfoID(ctx)
	if GetUserSwitcherInfoFromDatabase(getUserID) == true {
		// use lxns checker service.
		// check bind first, get user friend id.
		getFriendID := GetUserMaiFriendID(getUserID)
		if getFriendID.MaimaiID == 0 {
			ctx.SendPlainMessage(true, "你还没有绑定呢！使用/mai lxbind <friendcode> 以绑定")
			return
		}
		getUserData := RequestBasicDataFromLxns(getFriendID.MaimaiID)
		if getUserData.Code != 200 {
			ctx.SendPlainMessage(true, "aw 出现了一点小错误~：\n - 请检查你是否有上传过数据\n - 请检查你的设置是否允许了第三方查看")
			return
		}
		getGameUserData := RequestB50DataByFriendCode(getUserData.Data.FriendCode)
		if getGameUserData.Code != 200 {
			ctx.SendPlainMessage(true, "aw 出现了一点小错误~：\n - 请检查你是否有上传过数据\n - 请检查你的设置是否允许了第三方查看")
			return
		}
		getImager, _ := ReFullPageRender(getGameUserData, getUserData, ctx)
		_ = gg.NewContextForImage(getImager).SavePNG(engine.DataFolder() + "save/" + "LXNS_" + strconv.Itoa(int(getUserID)) + ".png")
		ctx.SendPhoto(tgba.FilePath(engine.DataFolder()+"save/"+"LXNS_"+strconv.Itoa(int(getUserID))+".png"), true, "")
	} else {
		// diving fish checker:
		getUsername := GetUserInfoNameFromDatabase(getUserID)
		if getUsername == "" {
			ctx.SendPlainMessage(true, "你还没有绑定呢！使用/mai bind <UserName> 以绑定")
			return
		}
		getUserData, err := QueryMaiBotDataFromUserName(getUsername)
		if err != nil {
			ctx.SendPlainMessage(true, err)
			return
		}
		var data player
		_ = json.Unmarshal(getUserData, &data)
		renderImg := FullPageRender(data, ctx)
		_ = gg.NewContextForImage(renderImg).SavePNG(engine.DataFolder() + "save/" + strconv.Itoa(int(getUserID)) + ".png")
		ctx.SendPhoto(tgba.FilePath(engine.DataFolder()+"save/"+strconv.Itoa(int(getUserID))+".png"), true, "")
	}
}

// MaimaiSwitcherService True == Lxns Service || False == Diving Fish Service.
func MaimaiSwitcherService(ctx *rei.Ctx) {
	getUserID, _ := toolchain.GetChatUserInfoID(ctx)
	getBool := GetUserSwitcherInfoFromDatabase(getUserID)
	err := FormatUserSwitcher(getUserID, !getBool).ChangeUserSwitchInfoFromDataBase()
	if err != nil {
		panic(err)
	}
	var getEventText string
	// due to it changed, so reverse.
	// due to it changed, so reverse.
	if getBool == false {
		getEventText = "Lxns查分"
	} else {
		getEventText = "Diving Fish查分"
	}
	ctx.SendPlainMessage(true, "已经修改为"+getEventText)
}

func CheckTheTicketIsValid(ticket string) bool {
	getData, err := web.GetData("https://www.diving-fish.com/api/maimaidxprober/token_available?token=" + ticket)
	if err != nil {
		panic(err)
	}
	result := gjson.Get(helper.BytesToString(getData), "message").String()
	if result == "ok" {
		return true
	}
	return false
}

func convert(listStruct UserMusicListStruct) []InnerStructChanger {
	getRequest, err := os.ReadFile(engine.DataFolder() + "music_data")
	if err != nil {
		panic(err)
	}
	var divingfishMusicData []DivingFishMusicDataStruct
	err = json.Unmarshal(getRequest, &divingfishMusicData)
	if err != nil {
		panic(err)
	}
	mdMap := make(map[string]DivingFishMusicDataStruct)
	for _, m := range divingfishMusicData {
		mdMap[m.Id] = m
	}
	var dest []InnerStructChanger
	for _, musicList := range listStruct.UserMusicList {
		for _, musicDetailedList := range musicList.UserMusicDetailList {
			/*
				for _, userMusicDetail := range musicList.UserMusicDetailList {
					if _, exists := mdMap[strconv.Itoa(musicDetailedList.MusicId)]; !exists {
						continue
					}
				}
			*/
			level := musicDetailedList.Level
			achievement := math.Min(1010000, float64(musicDetailedList.Achievement))
			fc := []string{"", "fc", "fcp", "ap", "app"}[musicDetailedList.ComboStatus]
			fs := []string{"", "fs", "fsp", "fsd", "fsdp"}[musicDetailedList.SyncStatus]
			dxScore := musicDetailedList.DeluxscoreMax
			dest = append(dest, InnerStructChanger{
				Title:        mdMap[strconv.Itoa(musicDetailedList.MusicId)].Title,
				Type:         mdMap[strconv.Itoa(musicDetailedList.MusicId)].Type,
				LevelIndex:   level,
				Achievements: (achievement) / 10000,
				Fc:           fc,
				Fs:           fs,
				DxScore:      dxScore,
			})
		}
	}
	return dest
}
