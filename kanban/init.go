// Package kanban 打印版本信息
package kanban

import "fmt"

//go:generate go run github.com/FloatTech/ReiBot-Plugin/kanban/gen

func init() {
	PrintBanner()
}

func PrintBanner() {
	fmt.Print(
		"\n======================[ReiBot-Plugin]======================",
		"\n", Banner, "\n",
		"===========================================================\n\n",
	)
}
