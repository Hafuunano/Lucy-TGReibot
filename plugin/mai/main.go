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
	engine.OnMessageRegex(`^[! /]mai\sbind\s(.*)$`).SetBlock(true).Handle(func(ctx *rei.Ctx) {
		matched := ctx.State["regex_matched"].([]string)[1]
		getUserID, _ := toolchain.GetChatUserInfoID(ctx)
		FormatUserDataBase(getUserID, GetUserPlateInfoFromDatabase(getUserID), GetUserDefaultBackgroundDataFromDatabase(getUserID), matched).BindUserDataBase()
		ctx.SendPlainMessage(true, "绑定成功~！")
	})
	engine.OnMessageRegex(`^[! /]mai\splate\s(.*)$`).SetBlock(true).Handle(func(ctx *rei.Ctx) {
		getPlateInfo := ctx.State["regex_matched"].([]string)[1]
		getUserID, _ := toolchain.GetChatUserInfoID(ctx)
		FormatUserDataBase(getUserID, getPlateInfo, GetUserDefaultBackgroundDataFromDatabase(getUserID), GetUserInfoNameFromDatabase(getUserID)).BindUserDataBase()
		ctx.SendPlainMessage(true, "好哦~ 是个好名称w")
	})
	engine.OnMessageRegex(`^[! /]mai\supload`, rei.MustProvidePhoto("请提供一张图片，图片大小比例适应为6:1 (1260x210) ,如果图片不适应将会自动剪辑到合适大小", "没有拿到图片呢 awa")).SetBlock(true).Handle(func(ctx *rei.Ctx) {
		getUserID, _ := toolchain.GetChatUserInfoID(ctx)
		ps := ctx.State["photos"].([]tgba.PhotoSize)
		pic := ps[len(ps)-1]
		picu, err := ctx.Caller.GetFileDirectURL(pic.FileID)
		imageData, err := web.GetData(picu)
		if err != nil {
			return
		}
		getRaw, _, err := image.Decode(bytes.NewReader(imageData))
		if err != nil {
			panic(err)
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
	})
	engine.OnMessageRegex(`^[! /]mai\sremove`).SetBlock(true).Handle(func(ctx *rei.Ctx) {
		getUserID, _ := toolchain.GetChatUserInfoID(ctx)
		_ = os.Remove(userPlate + strconv.Itoa(int(getUserID)) + ".png")
		ctx.SendPlainMessage(true, "已经移除了~ ")
	})
	engine.OnMessageRegex(`^[! ！/]mai\sdefault\splate\s(.*)$`).SetBlock(true).Handle(func(ctx *rei.Ctx) {
		getUserID, _ := toolchain.GetChatUserInfoID(ctx)
		getDefaultInfo := ctx.State["regex_matched"].([]string)[1]
		_, err := GetDefaultPlate(getDefaultInfo)
		if err != nil {
			ctx.SendPlainMessage(true, "设定的预设不正确")
			return
		}
		FormatUserDataBase(getUserID, GetUserPlateInfoFromDatabase(getUserID), getDefaultInfo, GetUserInfoNameFromDatabase(getUserID)).BindUserDataBase()
		ctx.SendPlainMessage(true, "已经设定好了哦w~ ")
	})
	engine.OnMessageCommand("mai").SetBlock(true).Handle(func(ctx *rei.Ctx) {
		// query data from sql
		getUserID, _ := toolchain.GetChatUserInfoID(ctx)
		getUsername := GetUserInfoNameFromDatabase(getUserID)
		if getUsername == "" {
			ctx.SendPlainMessage(true, "你还没有绑定呢！")
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
	})

}
