package mai

import (
	"encoding/json"
	"github.com/FloatTech/ReiBot-Plugin/utils/toolchain"
	"github.com/FloatTech/gg"
	ctrl "github.com/FloatTech/zbpctrl"
	rei "github.com/fumiama/ReiBot"
	tgba "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
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
		FormatUserDataBase(getUserID, "", "", matched).BindUserDataBase()
		ctx.SendPlainMessage(true, "绑定成功~！")
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
