package kanban

import (
	"fmt"
	"strings"
)

var (
	info = [...]string{
		"* Telegram + ReiBot + Golang",
		"* Version 0.0.4 - 2022-08-16 15:09:12 +0800 CST",
		"* Copyright © 2020 - 2022 FloatTech. All Rights Reserved.",
		"* Project: https://github.com/FloatTech/ReiBot-Plugin",
	}
	// Banner ...
	Banner = strings.Join(info[:], "\n")
)

// PrintBanner ...
func PrintBanner() {
	fmt.Print(
		"\n======================[ReiBot-Plugin]======================",
		"\n", Banner, "\n",
		"===========================================================\n\n",
	)
}
