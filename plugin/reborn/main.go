package reborn

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/FloatTech/ReiBot-Plugin/utils/toolchain"
	"github.com/FloatTech/ReiBot-Plugin/utils/transform"
	ctrl "github.com/FloatTech/zbpctrl"
	rei "github.com/fumiama/ReiBot"
	wr "github.com/mroth/weightedrand"
	"github.com/sirupsen/logrus"
	"github.com/wdvxdr1123/ZeroBot/extension/rate"
)

var (
	areac     *wr.Chooser
	gender, _ = wr.NewChooser(
		wr.Choice{Item: "男孩子", Weight: 23707},
		wr.Choice{Item: "女孩子", Weight: 49292},
		wr.Choice{Item: "雌雄同体", Weight: 1001},
		wr.Choice{Item: "猫猫!", Weight: 15000},
		wr.Choice{Item: "狗狗!", Weight: 5000},
		wr.Choice{Item: "龙猫~", Weight: 6000},
	)
	rebornTimerManager = rate.NewManager[int64](time.Minute*2, 5)
)

type ratego []struct {
	Name   string  `json:"name"`
	Weight float64 `json:"weight"`
}

func init() {
	engine := rei.Register("reborn", &ctrl.Options[*rei.Ctx]{
		DisableOnDefault: false,
		Help:             "reborn",
	})
	go func() {
		datapath := transform.ReturnLucyMainDataIndex("funwork")
		jsonfile := datapath + "ratego.json"
		area := make(ratego, 226)
		err := load(&area, jsonfile)
		if err != nil {
			panic(err)
		}
		choices := make([]wr.Choice, len(area))
		for i, a := range area {
			choices[i].Item = a.Name
			choices[i].Weight = uint(a.Weight * 1e9)
		}
		areac, err = wr.NewChooser(choices...)
		if err != nil {
			panic(err)
		}
		logrus.Printf("[Reborn]读取%d个国家/地区", len(area))
	}()
	engine.OnMessageCommand("reborn").SetBlock(true).Handle(func(ctx *rei.Ctx) {
		if !rebornTimerManager.Load(toolchain.GetThisGroupID(ctx)).Acquire() {
			ctx.SendPlainMessage(true, "太快了哦，麻烦慢一点~")
			return
		}
		if rand.Int31() > 1<<27 {
			ctx.SendPlainMessage(true, fmt.Sprintf("投胎成功！\n您出生在 %s, 是 %s。", randcoun(), randgen()))
		} else {
			ctx.SendPlainMessage(true, "投胎失败！\n您没能活到出生，希望下次运气好一点呢~！")
		}

	})
}

// load 加载rate数据
func load(area *ratego, jsonfile string) error {
	data, err := os.ReadFile(jsonfile)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, area)
}

func randcoun() string {
	return areac.Pick().(string)
}

func randgen() string {
	return gender.Pick().(string)
}
