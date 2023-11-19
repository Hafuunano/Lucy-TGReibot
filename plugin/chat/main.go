package chat

import (
	`fmt`
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

}
