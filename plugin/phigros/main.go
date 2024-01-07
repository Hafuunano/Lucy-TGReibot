package phigros

import (
	"encoding/json"
	"time"
	"unicode/utf8"

	"image/color"

	"strconv"
	"sync"

	"github.com/FloatTech/ReiBot-Plugin/utils/toolchain"
	"github.com/FloatTech/floatbox/web"
	"github.com/FloatTech/gg"
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/disintegration/imaging"
	rei "github.com/fumiama/ReiBot"
	tgba "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	engine = rei.Register("phigros", &ctrl.Options[*rei.Ctx]{
		DisableOnDefault:  false,
		Help:              "phigros - | generate phigros b19 | roll ",
		PrivateDataFolder: "phigros",
	})
	router = "http://221.131.165.69:51602/"
)

func init() {
	engine.OnMessageCommand("pgr").SetBlock(true).Handle(func(ctx *rei.Ctx) {
		getMsg := ctx.Message.Text
		getSplitLength, getSplitStringList := toolchain.SplitCommandTo(getMsg, 3)
		if getSplitLength >= 2 {
			switch {
			case getSplitStringList[1] == "bind":
				if getSplitLength < 3 {
					ctx.SendPlainMessage(true, "参数提供不足")
					return
				}
				PhiBind(ctx, getSplitStringList[2])
			case getSplitStringList[1] == "roll":
				if getSplitLength < 3 {
					ctx.SendPlainMessage(true, "参数提供不足")
					return
				}
				RollToRenderPhigros(ctx, getSplitStringList[2])
			default:
				ctx.SendPlainMessage(true, "未知的指令或者指令出现错误~")
			}
		} else {
			RenderPhi(ctx)
		}
	})

}

func RenderPhi(ctx *rei.Ctx) {
	getUserId, _ := toolchain.GetChatUserInfoID(ctx)
	data := GetUserInfoFromDatabase(getUserId)
	getDataSession := data.PhiSession
	if getDataSession == "" {
		_, _ = ctx.SendPlainMessage(true, "请前往 https://phi.lemonkoi.one 获取绑定码进行绑定 ")
		return
	}
	userData := GetUserInfoFromDatabase(getUserId)
	ctx.SendPlainMessage(true, "好哦~正在帮你请求，请稍等一下啦w~大约需要1-2分钟")
	var dataWaiter sync.WaitGroup
	var AvatarWaiter sync.WaitGroup
	var getAvatarFormat *gg.Context
	var phidata []byte
	var setGlobalStat = true
	AvatarWaiter.Add(1)
	dataWaiter.Add(1)
	go func() {
		defer dataWaiter.Done()
		getFullLink := router + "/api/phi/bests?session=" + userData.PhiSession + "&overflow=2"
		phidata, _ = web.GetData(getFullLink)
		if phidata == nil {
			ctx.SendPlainMessage(true, "目前 Unoffical Phigros API 暂时无法工作 请过一段时候尝试")
			setGlobalStat = false
			return
		}
		err := json.Unmarshal(phidata, &phigrosB19)
		if !phigrosB19.Status || err != nil {
			ctx.SendPlainMessage(true, "w? 貌似出现了一些问题x")
			return
		}
	}()
	go func() {
		defer AvatarWaiter.Done()
		avatarByteUni := toolchain.GetTargetAvatar(ctx)
		if avatarByteUni == nil {
			return
		}
		// avatar
		showUserAvatar := imaging.Resize(avatarByteUni, 250, 250, imaging.Lanczos)
		getAvatarFormat = gg.NewContext(250, 250)
		getAvatarFormat.DrawRoundedRectangle(0, 0, 248, 248, 20)
		getAvatarFormat.Clip()
		getAvatarFormat.DrawImage(showUserAvatar, 0, 0)
		getAvatarFormat.Fill()
	}()
	getRawBackground, _ := gg.LoadImage(backgroundRender)
	getMainBgRender := gg.NewContextForImage(imaging.Resize(getRawBackground, 2750, 5500, imaging.Lanczos))
	_ = getMainBgRender.LoadFontFace(font, 30)
	// header background
	drawTriAngle(getMainBgRender, a, 0, 166, 1324, 410)
	getMainBgRender.SetRGBA255(0, 0, 0, 160)
	getMainBgRender.Fill()
	drawTriAngle(getMainBgRender, a, 1318, 192, 1600, 350)
	getMainBgRender.SetRGBA255(0, 0, 0, 160)
	getMainBgRender.Fill()
	drawTriAngle(getMainBgRender, a, 1320, 164, 6, 414)
	getMainBgRender.SetColor(color.White)
	getMainBgRender.Fill()
	// header background end.
	// load icon with other userinfo.
	getMainBgRender.SetColor(color.White)
	logo, _ := gg.LoadPNG(icon)
	getImageLogo := imaging.Resize(logo, 290, 290, imaging.Lanczos)
	getMainBgRender.DrawImage(getImageLogo, 50, 216)
	fontface, _ := gg.LoadFontFace(font, 90)
	getMainBgRender.SetFontFace(fontface)
	getMainBgRender.DrawString("Phigros", 422, 336)
	getMainBgRender.DrawString("RankingScore查询", 422, 462)
	// draw userinfo path
	renderHeaderText, _ := gg.LoadFontFace(font, 54)
	getMainBgRender.SetFontFace(renderHeaderText)
	dataWaiter.Wait()
	if !setGlobalStat {
		return
	}
	getMainBgRender.DrawString("Player: "+phigrosB19.Content.PlayerID, 1490, 300)
	getMainBgRender.DrawString("RankingScore: "+strconv.FormatFloat(phigrosB19.Content.RankingScore, 'f', 3, 64), 1490, 380)
	getMainBgRender.DrawString("ChanllengeMode: ", 1490, 460) // +56
	getColor, getLink := GetUserChallengeMode(phigrosB19.Content.ChallengeModeRank)
	if getColor != "" {
		getColorLink := ChanllengeMode + getColor + ".png"
		getColorImage, _ := gg.LoadImage(getColorLink)
		getMainBgRender.DrawImage(imaging.Resize(getColorImage, 238, 130, imaging.Lanczos), 1912, 390)
	}
	renderHeaderTextNumber, _ := gg.LoadFontFace(font, 65)
	getMainBgRender.SetFontFace(renderHeaderTextNumber)
	// white glow render
	getMainBgRender.SetRGB(1, 1, 1)
	getMainBgRender.DrawStringAnchored(getLink, 2021, 430, 0.4, 0.4)
	// avatar
	AvatarWaiter.Wait()
	if getAvatarFormat != nil {
		getMainBgRender.DrawImage(getAvatarFormat.Image(), 2321, 230)
		getMainBgRender.Fill()
	}
	// render
	CardRender(getMainBgRender, phidata)
	// draw bottom
	_ = getMainBgRender.LoadFontFace(font, 40)
	getMainBgRender.SetColor(color.White)
	getMainBgRender.Fill()
	getMainBgRender.DrawString("Generated By Lucy (HiMoYoBOT) | Designed By Eastown | Data From Phigros Library Project", 10, 5480)
	_ = getMainBgRender.SavePNG(engine.DataFolder() + "save/" + strconv.Itoa(int(getUserId)) + ".png")
	ctx.SendPhoto(tgba.FilePath(engine.DataFolder()+"save/"+strconv.Itoa(int(getUserId))+".png"), true, "")
}

func RollToRenderPhigros(ctx *rei.Ctx, num string) {
	var wg sync.WaitGroup
	var avatarWaitGroup sync.WaitGroup
	var dataWaiter sync.WaitGroup
	var getMainBgRender *gg.Context
	var getAvatarFormat *gg.Context
	var setGlobalStat = true
	var phidata []byte
	wg.Add(1)
	avatarWaitGroup.Add(1)
	dataWaiter.Add(1)
	// get Session From Database.
	userid, _ := toolchain.GetChatUserInfoID(ctx)
	data := GetUserInfoFromDatabase(userid)
	getDataSession := data.PhiSession
	if getDataSession == "" {
		ctx.SendPlainMessage(true, "由于Session特殊性，请前往 https://phi.lemonkoi.one 获取绑定码进行绑定")
		return
	}
	userData := GetUserInfoFromDatabase(userid)
	getRoll := num
	getRollInt, err := strconv.ParseInt(getRoll, 10, 64)
	if err != nil {
		ctx.SendPlainMessage(true, "请求roll不合法")
		return
	}
	if getRollInt > 40 {
		ctx.SendPlainMessage(true, "限制查询数为小于40")
		return
	}
	getOverFlowNumber := getRollInt - 19
	if getOverFlowNumber <= 0 {
		getOverFlowNumber = 0
	}
	ctx.SendPlainMessage(true, "好哦~正在帮你请求，请稍等一下啦w~大约需要1-2分钟")
	// data handling.
	go func() {
		defer dataWaiter.Done()
		getFullLink := router + "/api/phi/bests?session=" + userData.PhiSession + "&overflow=" + strconv.Itoa(int(getOverFlowNumber))
		phidata, _ = web.GetData(getFullLink)
		if phidata == nil {
			ctx.SendPlainMessage(true, "目前 Unoffical Phigros Library 暂时无法工作 请过一段时候尝试")
			setGlobalStat = false
			return
		}
		err = json.Unmarshal(phidata, &phigrosB19)
		if err != nil {
			ctx.SendPlainMessage(true, "发生解析错误: \n", err)
			setGlobalStat = false
			return
		}
		if !phigrosB19.Status {
			ctx.SendPlainMessage(true, "w? 貌似出现了一些问题x\n", phigrosB19.Message)
			setGlobalStat = false
			return
		}
	}()
	go func() {
		defer wg.Done()
		getRawBackground, _ := gg.LoadImage(backgroundRender)
		getMainBgRender = gg.NewContextForImage(imaging.Resize(getRawBackground, 2750, int(5250+getOverFlowNumber*200), imaging.Lanczos))
	}()
	go func() {
		defer avatarWaitGroup.Done()
		avatarByteUni := toolchain.GetTargetAvatar(ctx)
		if avatarByteUni == nil {
			return
		}
		// avatar
		showUserAvatar := imaging.Resize(avatarByteUni, 250, 250, imaging.Lanczos)
		getAvatarFormat = gg.NewContext(250, 250)
		getAvatarFormat.DrawRoundedRectangle(0, 0, 248, 248, 20)
		getAvatarFormat.Clip()
		getAvatarFormat.DrawImage(showUserAvatar, 0, 0)
		getAvatarFormat.Fill()
	}()
	wg.Wait()
	_ = getMainBgRender.LoadFontFace(font, 30)
	// header background
	drawTriAngle(getMainBgRender, a, 0, 166, 1324, 410)
	getMainBgRender.SetRGBA255(0, 0, 0, 160)
	getMainBgRender.Fill()
	drawTriAngle(getMainBgRender, a, 1318, 192, 1600, 350)
	getMainBgRender.SetRGBA255(0, 0, 0, 160)
	getMainBgRender.Fill()
	drawTriAngle(getMainBgRender, a, 1320, 164, 6, 414)
	getMainBgRender.SetColor(color.White)
	getMainBgRender.Fill()
	// header background end.
	// load icon with other userinfo.
	getMainBgRender.SetColor(color.White)
	logo, _ := gg.LoadPNG(icon)
	getImageLogo := imaging.Resize(logo, 290, 290, imaging.Lanczos)
	getMainBgRender.DrawImage(getImageLogo, 50, 216)
	fontface, _ := gg.LoadFontFace(font, 90)
	getMainBgRender.SetFontFace(fontface)
	getMainBgRender.DrawString("Phigros", 422, 336)
	getMainBgRender.DrawString("RankingScore查询", 422, 462)
	dataWaiter.Wait()
	if !setGlobalStat {
		return
	}
	// draw userinfo path
	renderHeaderText, _ := gg.LoadFontFace(font, 54)
	getMainBgRender.SetFontFace(renderHeaderText)
	// wait data until fine.
	getMainBgRender.DrawString("Player: "+phigrosB19.Content.PlayerID, 1490, 300)
	getMainBgRender.DrawString("RankingScore: "+strconv.FormatFloat(phigrosB19.Content.RankingScore, 'f', 3, 64), 1490, 380)
	getMainBgRender.DrawString("ChanllengeMode: ", 1490, 460) // +56
	getColor, getLink := GetUserChallengeMode(phigrosB19.Content.ChallengeModeRank)
	if getColor != "" {
		getColorLink := ChanllengeMode + getColor + ".png"
		getColorImage, _ := gg.LoadImage(getColorLink)
		getMainBgRender.DrawImage(imaging.Resize(getColorImage, 238, 130, imaging.Lanczos), 1912, 390)
	}
	renderHeaderTextNumber, _ := gg.LoadFontFace(font, 65)
	getMainBgRender.SetFontFace(renderHeaderTextNumber)
	// white glow render
	getMainBgRender.SetRGB(1, 1, 1)
	getMainBgRender.DrawStringAnchored(getLink, 2021, 430, 0.4, 0.4)
	avatarWaitGroup.Wait()
	getMainBgRender.DrawImage(getAvatarFormat.Image(), 2321, 230)
	getMainBgRender.Fill()
	// render
	CardRender(getMainBgRender, phidata)
	// draw bottom
	_ = getMainBgRender.LoadFontFace(font, 40)
	getMainBgRender.SetColor(color.White)
	getMainBgRender.Fill()
	getMainBgRender.DrawString("Generated By Lucy (HiMoYoBOT) | Designed By Eastown | Data From Phigros Library Project", 10, float64(5110+getOverFlowNumber*200+100))
	_ = getMainBgRender.SavePNG(engine.DataFolder() + "save/" + "roll" + strconv.Itoa(int(userid)) + ".png")
	ctx.SendPhoto(tgba.FilePath(engine.DataFolder()+"save/"+"roll"+strconv.Itoa(int(userid))+".png"), true, "")

}

func PhiBind(ctx *rei.Ctx, bindAcc string) {
	hash := bindAcc
	getUserId, _ := toolchain.GetChatUserInfoID(ctx)
	userInfo := GetUserInfoTimeFromDatabase(getUserId)
	if userInfo+(12*60*60) > time.Now().Unix() {
		ctx.SendPlainMessage(true, "12小时内仅允许绑定一次哦")
		return
	}
	indexReply := DecHashToRaw(hash)
	// get session.
	if indexReply == "" {
		ctx.SendPlainMessage(true, "请前往 https://phi.lemonkoi.one 获取绑定码进行绑定")
		return
	}
	getQQID, getSessionID := RawJsonParse(indexReply)
	if getQQID != getUserId {
		ctx.SendPlainMessage(true, "请求Hash中Telegram ID不一致，请使用自己的号重新申请")
		return
	}
	if utf8.RuneCountInString(getSessionID) != 25 {
		ctx.SendPlainMessage(true, "Session 传入数值出现错误，请重新绑定")
		return
	}
	_ = FormatUserDataBase(getQQID, getSessionID, time.Now().Unix()).BindUserDataBase()
	ctx.SendPlainMessage(true, "绑定成功～")
}
