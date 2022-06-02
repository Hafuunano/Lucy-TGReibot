package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	_ "github.com/FloatTech/ReiBot-Plugin/plugin/lolicon"

	// -----------------------以下为内置依赖，勿动------------------------ //
	rei "github.com/fumiama/ReiBot"
	tgba "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/FloatTech/ReiBot-Plugin/kanban"
	// -----------------------以上为内置依赖，勿动------------------------ //
)

func main() {
	rand.Seed(time.Now().UnixNano()) // 全局 seed，其他插件无需再 seed

	token := flag.String("t", "", "telegram api token")
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

	sus := make([]int64, 0, 16)
	for _, s := range flag.Args() {
		i, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			continue
		}
		sus = append(sus, i)
	}

	rei.OnMessageFullMatchGroup([]string{"help", "帮助", "menu", "菜单"}).SetBlock(true).
		Handle(func(ctx *rei.Ctx) {
			msg := ctx.Value.(*tgba.Message)
			_, _ = ctx.Caller.Send(tgba.NewMessage(msg.Chat.ID, kanban.Banner))
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
