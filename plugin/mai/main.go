package mai

import (
	"bytes"
	"encoding/json"
	"image"
	"os"
	"strconv"

	"github.com/FloatTech/ReiBot-Plugin/utils/toolchain"
	"github.com/FloatTech/floatbox/web"
	"github.com/FloatTech/gg"
	ctrl "github.com/FloatTech/zbpctrl"
	rei "github.com/fumiama/ReiBot"
	tgba "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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
				userID := GetUserMaiUserid(getSplitStringList[2])
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
				getCode := FastUnlockerMai15mins(getMaiID.Userid)
				if getCode == 200 {
					ctx.SendPlainMessage(true, "发信成功，如果未生效请重新尝试")
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
	if getBool == false {
		getEventText = "Lxns查分"
	} else {
		getEventText = "Diving Fish查分"
	}
	ctx.SendPlainMessage(true, "已经修改为"+getEventText)
}
