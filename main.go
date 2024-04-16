package main

import (
	"flag"
	"fmt"

	"github.com/MoYoez/Lucy_reibot/kanban"
	rei "github.com/fumiama/ReiBot"

	_ "github.com/MoYoez/Lucy_reibot/plugin/chat"
	_ "github.com/MoYoez/Lucy_reibot/plugin/fortune"
	_ "github.com/MoYoez/Lucy_reibot/plugin/lolicon"
	_ "github.com/MoYoez/Lucy_reibot/plugin/mai"
	_ "github.com/MoYoez/Lucy_reibot/plugin/phigros"
	_ "github.com/MoYoez/Lucy_reibot/plugin/reborn"
	_ "github.com/MoYoez/Lucy_reibot/plugin/score"
	_ "github.com/MoYoez/Lucy_reibot/plugin/tools"
	_ "github.com/MoYoez/Lucy_reibot/plugin/tracemoe"
	_ "github.com/MoYoez/Lucy_reibot/plugin/what2eat"
	_ "github.com/MoYoez/Lucy_reibot/plugin/wife"

	_ "github.com/MoYoez/Lucy_reibot/plugin/action"
	_ "github.com/MoYoez/Lucy_reibot/plugin/simai"
	_ "github.com/MoYoez/Lucy_reibot/plugin/slash" // slash should be the last
	_ "github.com/MoYoez/Lucy_reibot/plugin/stickers"

	"os"
	"strconv"

	"github.com/joho/godotenv"

	tgba "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
)

func main() {
	_ = godotenv.Load()
	token := flag.String("t", os.Getenv("tgbot"), "telegram api token")
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

	rei.OnMessageCommand("help").SetBlock(true).
		Handle(func(ctx *rei.Ctx) {
			ctx.SendPlainMessage(true, kanban.Banner)
		})
	rei.Run(rei.Bot{
		Token:   *token,
		Botname: "Lucy",
		Buffer:  *buffer,
		UpdateConfig: tgba.UpdateConfig{
			Offset:  *offset,
			Limit:   0,
			Timeout: *timeout,
			//	AllowedUpdates: []string{"message", "edited_message", "message_reaction", "message_reaction_count", "inline_query", "chosen_inline_result", "callback_query", "shipping_query", "pre_checkout_query", "poll", "poll_answer", "my_chat_member", "chat_member", "chat_join_request", "chat_boost", "removed_chat_boost"},
			AllowedUpdates: []string{"message", "edited_message", "inline_query", "chosen_inline_result", "callback_query", "shipping_query", "pre_checkout_query", "poll", "poll_answer", "my_chat_member", "chat_member", "chat_join_request", "chat_boost", "removed_chat_boost"},
		},
		SuperUsers: sus,
		Debug:      *debug,
	})
}
