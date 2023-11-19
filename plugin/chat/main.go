package chat

import (
	`fmt`
	`strconv`
	`time`

	`github.com/FloatTech/ReiBot-Plugin/utils/toolchain`
	ctrl `github.com/FloatTech/zbpctrl`
	rei `github.com/fumiama/ReiBot`
)

var (
	engine = rei.Register("chat", &ctrl.Options[*rei.Ctx]{
		DisableOnDefault: false,
		Help:             "chat",
	})
)

func init() {
	engine.OnMessage(rei.OnlyToMe).SetBlock(false).Handle(func(ctx *rei.Ctx) {
		nickname := "Lucy" // hardcoded is a good choice ( I will fix it later.(
		if ctx.Message.Text != "" {
			fmt.Print(ctx.Message.Text)
			return
		}
		time.Sleep(time.Second * 1)
		toolchain.FastSendRandMuiltText(ctx, "这里是"+nickname+"(っ●ω●)っ", nickname+"不在呢~", "哼！"+nickname+"不想理你~")
	})
	engine.OnMessageCommand("callname").SetBlock(true).Handle(func(ctx *rei.Ctx) {
		_, results := toolchain.SplitCommandTo(ctx.Message.Text, 2)
		if len(results) <= 1 {
			return
		}
		texts := results[1]
		if texts == "" {
			return
		}
		if toolchain.StringInArray(texts, []string{"Lucy", "笨蛋", "老公", "猪", "夹子", "主人"}) {
			ctx.SendPlainMessage(true, "这些名字可不好哦(敲)")
			return
		}
		getID, _ := toolchain.GetChatUserInfoID(ctx)
		userID := strconv.FormatInt(getID, 10)
		err := toolchain.StoreUserNickname(userID, texts)
		if err != nil {
			ctx.SendPlainMessage(true, "发生了一些不可预料的问题 请稍后再试, ERR: ", err)
			return
		}
		ctx.SendPlainMessage(true, "好哦~ ", texts, " ちゃん~~~")
	})
}
