package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	_ "github.com/FloatTech/ReiBot-Plugin/plugin/bilibili_parse"
	_ "github.com/FloatTech/ReiBot-Plugin/plugin/chrev"
	_ "github.com/FloatTech/ReiBot-Plugin/plugin/fortune"
	_ "github.com/FloatTech/ReiBot-Plugin/plugin/lolicon"
	_ "github.com/FloatTech/ReiBot-Plugin/plugin/tracemoe"

	// -----------------------以下为内置依赖，勿动------------------------ //
	rei "github.com/fumiama/ReiBot"
	tgba "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"

	"github.com/FloatTech/ReiBot-Plugin/kanban"
	// -----------------------以上为内置依赖，勿动------------------------ //
)

func main() {
	rand.Seed(time.Now().UnixNano()) // 全局 seed，其他插件无需再 seed

	token := flag.String("t", "6698655694:AAHyJj6bgketzAOu8GxDTCFOOgOXpNz_lIY", "telegram api token")
	buffer := flag.Int("b", 256, "message sequence length")
	debug := flag.Bool("d", false, "enable debug-level log output")
	offset := flag.Int("o", 0, "the last Update ID to include")
	timeout := flag.Int("T", 60, "timeout")
	help := flag.Bool("h", false, "print this help")
	flag.Parse()
	if *help {
		fmt.Println("Usage:")
		flag.PrintDefaults()
		os.Exit(0)
	}

	if *debug {
		logrus.SetLevel(logrus.DebugLevel)
	}

	sus := make([]int64, 0, 16)
	for _, s := range flag.Args() {
		i, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			continue
		}
		sus = append(sus, i)
	}

	rei.OnMessageCommandGroup([]string{"help"}, rei.OnlyToMe).SetBlock(true).
		Handle(func(ctx *rei.Ctx) {
			_, _ = ctx.SendPlainMessage(false, kanban.Banner)
		})
	rei.Run(rei.Bot{
		Token:  *token,
		Buffer: *buffer,
		UpdateConfig: tgba.UpdateConfig{
			Offset:  *offset,
			Limit:   0,
			Timeout: *timeout,
		},
		SuperUsers: sus,
		Debug:      *debug,
	})
}
