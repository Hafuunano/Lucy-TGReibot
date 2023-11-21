// Package wife From "github.com/MoYoez/Lucy_zerobot"
package wife

import (
	"time"

	coins "github.com/FloatTech/ReiBot-Plugin/utils/coins"
	ctrl "github.com/FloatTech/zbpctrl"
	rei "github.com/fumiama/ReiBot"
	tgba "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/wdvxdr1123/ZeroBot/extension/rate"
)

var (
	MessageTickerLimiter = rate.NewManager[int64](time.Minute*1, 2)
	engine               = rei.Register("wife", &ctrl.Options[*rei.Ctx]{
		DisableOnDefault:  false,
		Help:              "Hi NekoPachi!",
		PrivateDataFolder: "wife",
	}).ApplySingle(ReverseSingle)
	ReverseSingle = rei.NewSingle(
		rei.WithKeyFn(func(ctx *rei.Ctx) int64 {
			switch msg := ctx.Value.(type) {
			case *tgba.Message:
				return msg.From.ID
			case *tgba.CallbackQuery:
				return msg.From.ID
			}
			return 0
		}), rei.WithPostFn[int64](func(ctx *rei.Ctx) {
			if !MessageTickerLimiter.Load(ctx.Message.Chat.ID).Acquire() {
				return
			}
			_, _ = ctx.SendPlainMessage(false, "正在操作哦～")
		}),
	)
)

/*
StatusID:

Type 1: Normal Mode, nothing happened.

Type 2: Cannot be the Target, Target became initiative, so reverse.

(However the target and the initiative should be in their position, DO NOT CHANGE. )

Type 3: Something is wrong, you are Target == initiative Person. (Drop The Person Before.)

Type 4: Removed.
(When User get others person. || IF REMARRIED, CHANGE IT TO TYPE1.) || (Be check more Time to reduce to err.)

Type 5: NTR Mode
(Tips: NTR means changed their pairkey & TargetID || UserID, need to do some changes. ) ||
(Attempt to do once more every person.)

Type 6: No wife Mod?
Fake - Invisible person here.
(Lucy Hides this and shows it in the next Time if a person uses NTR,
shows nothing, and Lucy will make it for joke. LMAO)

Type 7: NTRED BY SOMEONE.
*/

func init() {
	sdb := coins.Initialize("./data/score/score.db")
	dict := make(map[string][]string) // this dict is used to reply
	// dict path.
	dict["block"] = []string{"嗯哼？貌似没有找到哦w", "再试试哦w，或许有帮助w", "运气不太好哦，想一下办法呢x"}
	dict["success"] = []string{"Lucky For You~", "恭喜哦ww~ ", "这边来恭喜一下哦w～", "貌似很成功的一次尝试呢w~"}
	dict["failed"] = []string{"今天的运气有一点背哦~这一次没有成功呢x", "_(:з」∠)_下次还有机会 抱抱w", "没关系哦，虽然失败了但还有机会呢x"}
	dict["ntr"] = []string{"嗯哼～这位还是成功了呢x", "aaa 好怪 不过还是让你通过了 ^^ "}
	dict["lost_failed"] = []string{"为什么要分呢? 让咱捏捏w", "太坏了啦！不许！"}
	dict["lost_success"] = []string{"好呢w 就这样呢(", "已经成功了哦w"}
	dict["hide_mode"] = []string{"哼哼～ 哼唧", "喵喵喵？！"}

	engine.OnMessageCommand("marry").SetBlock(true).Handle(func(ctx *rei.Ctx) {

	})
}
