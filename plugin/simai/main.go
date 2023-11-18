// Package simai refactory From Lucy For Onebot. (origin github.com/FloatTech/Zerobot-Plugin)
package simai

import (
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/FloatTech/ReiBot-Plugin/utils/toolchain"
	"github.com/FloatTech/ReiBot-Plugin/utils/transform"
	ctrl "github.com/FloatTech/zbpctrl"
	rei "github.com/fumiama/ReiBot"
	"github.com/wdvxdr1123/ZeroBot/extension/rate"
	"gopkg.in/yaml.v3"
)

// SimPackData simai Data
type SimPackData struct {
	Proud  map[string][]string `yaml:"傲娇"`
	Kawaii map[string][]string `yaml:"可爱"`
}

var limit = rate.NewManager[int64](time.Minute*3, 28) // 回复限制

func init() {
	engine := rei.Register("simai", &ctrl.Options[*rei.Ctx]{
		DisableOnDefault:  false,
		PrivateDataFolder: "simai",
		Help:              "simai - Use simia pre-render dict to make it more clever",
	})
	// onload simai dict
	dictLoaderLocation := transform.ReturnLucyMainDataIndex("simai") + "simai.yml"
	dictLoader, err := os.ReadFile(dictLoaderLocation)
	if err != nil {
		panic(err)
	}
	var data SimPackData
	_ = yaml.Unmarshal(dictLoader, &data)
	engine.OnMessage(rei.OnlyToMe).SetBlock(false).Handle(func(ctx *rei.Ctx) {
		msg := ctx.Message.Text
		var getChartReply []string
		if GetTiredToken(ctx) < 4 {
			getChartReply = data.Proud[msg]
			// if no data
			if getChartReply == nil {
				getChartReply = data.Kawaii[msg]
				if getChartReply == nil {
					// no reply
					return
				}
			}
		} else {
			getChartReply = data.Kawaii[msg]
			// if no data
			if getChartReply == nil {
				getChartReply = data.Proud[msg]
				if getChartReply == nil {
					// no reply
					return
				}
			}
		}
		if GetTiredToken(ctx) < 4 {
			ctx.SendPlainMessage(true, "咱不想说话 好累awww")
			return
		} else {
			GetCostTiredToken(ctx)
		}
		getReply := getChartReply[rand.Intn(len(getChartReply))]
		getLucyName := []string{"Lucy", "Lucy酱"}[rand.Intn(2)]
		getReply = strings.ReplaceAll(getReply, "{segment}", " ")
		getReply = strings.ReplaceAll(getReply, "{name}", "你")
		getReply = strings.ReplaceAll(getReply, "{me}", getLucyName)
		ctx.SendPlainMessage(true, getReply)
	})
}

func GetTiredToken(ctx *rei.Ctx) float64 {
	getID, _ := toolchain.GetChatUserInfoID(ctx)
	return limit.Load(getID).Tokens()
}

func GetCostTiredToken(ctx *rei.Ctx) bool {
	getID, _ := toolchain.GetChatUserInfoID(ctx)
	return limit.Load(getID).AcquireN(3)
}
