package action

import (
	"math/rand"
	"time"

	"github.com/FloatTech/ReiBot-Plugin/utils/ctxext"
	"github.com/FloatTech/ReiBot-Plugin/utils/toolchain"
	ctrl "github.com/FloatTech/zbpctrl"
	rei "github.com/fumiama/ReiBot"
	tgba "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/wdvxdr1123/ZeroBot/extension/rate"
)

var (
	limit   = rate.NewManager[int64](time.Minute*10, 15)
	LucyImg = "/root/Lucy_Project/memes/" // LucyImg for Lucy的meme表情包地址
)

func init() {
	engine := rei.Register("action", &ctrl.Options[*rei.Ctx]{
		DisableOnDefault: false,
		Help:             "Lucy容易被动触发语言\n",
	})
	engine.OnMessageFullMatchGroup([]string{"喵", "喵喵", "喵喵喵"}).SetBlock(true).Handle(func(ctx *rei.Ctx) {
		if !limit.Load(toolchain.GetThisGroupID(ctx)).Acquire() {
			return
		}
		switch rand.Intn(6) {
		case 2, 3:
			ctx.SendPhoto(tgba.FilePath(RandImage("6152277811454.jpg", "meow.jpg", "file_3491851.jpg", "file_3492320.jpg")), true, "")
		case 4, 5:
			ctx.SendPlainMessage(true, []string{"喵喵~", "喵w~"}[rand.Intn(2)])
		}
	})
	engine.OnMessageFullMatch("咕咕").SetBlock(true).Limit(ctxext.LimitByUser).Handle(func(ctx *rei.Ctx) {
		if !limit.Load(toolchain.GetThisGroupID(ctx)).Acquire() {
			return
		}
		ctx.SendPlainMessage(true, []string{"炖了~鸽子都要恰掉w", "咕咕咕", "不许咕咕咕"}[rand.Intn(3)])
	})

}

func RandImage(file ...string) string {
	return LucyImg + file[rand.Intn(len(file))]
}
