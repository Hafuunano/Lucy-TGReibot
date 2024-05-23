package chun

import (
	"encoding/json"
	"github.com/FloatTech/gg"
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/MoYoez/Lucy_reibot/plugin/mai"
	"github.com/MoYoez/Lucy_reibot/utils/toolchain"
	rei "github.com/fumiama/ReiBot"
	tgba "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
)

var engine = rei.Register("chun", &ctrl.Options[*rei.Ctx]{
	DisableOnDefault:  false,
	Help:              "chun for Lucy",
	PrivateDataFolder: "chun",
})

func init() {
	engine.OnMessageCommand("chun").SetBlock(true).Handle(func(ctx *rei.Ctx) {
		getSplitLength, GetSplitInfo := toolchain.SplitCommandTo(ctx.Message.Text, 3)
		if getSplitLength >= 2 {
			switch {
			case GetSplitInfo[1] == "raw":
				// baserender
				ChunRender(ctx, true)
			case GetSplitInfo[1] == "bind":
				if getSplitLength < 3 {
					ctx.SendPlainMessage(true, "w? 绑定错了啊 是 /chun bind <username> ")
					return
				}
				getUsername := GetSplitInfo[2]
				mai.BindUserToMaimai(ctx, getUsername)
				return
			}

		} else {
			ChunRender(ctx, false)
		}
		
	})
}

func ChunRender(ctx *rei.Ctx, israw bool) {
	// check the user using.
	getUserID, _ := toolchain.GetChatUserInfoID(ctx)
	getUsername := mai.GetUserInfoNameFromDatabase(getUserID)
	if getUsername == "" {
		ctx.SendPlainMessage(true, "你还没有绑定呢！使用/chun bind <UserName> 以绑定")
		return
	}
	getUserData, err := QueryChunDataFromUserName(getUsername)
	if err != nil {
		ctx.SendPlainMessage(true, err)
		return
	}
	var data ChunData
	_ = json.Unmarshal(getUserData, &data)
	renderImg := BaseRender(data, ctx)
	_ = gg.NewContextForImage(renderImg).SaveJPG(engine.DataFolder()+"save/"+strconv.Itoa(int(getUserID))+".png", 80)

	if israw {
		getDocumentType := &tgba.DocumentConfig{
			BaseFile: tgba.BaseFile{BaseChat: tgba.BaseChat{
				ChatConfig: tgba.ChatConfig{ChatID: ctx.Message.Chat.ID},
			},
				File: tgba.FilePath(engine.DataFolder() + "save/" + strconv.Itoa(int(getUserID)) + ".png")},
			Caption:         "",
			CaptionEntities: nil,
		}
		ctx.Send(true, getDocumentType)
	} else {
		ctx.SendPhoto(tgba.FilePath(engine.DataFolder()+"save/"+strconv.Itoa(int(getUserID))+".png"), true, "")
	}
}
