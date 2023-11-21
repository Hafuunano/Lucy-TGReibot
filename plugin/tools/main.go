package tools

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/FloatTech/ReiBot-Plugin/utils/CoreFactory"
	"github.com/FloatTech/ReiBot-Plugin/utils/toolchain"
	ctrl "github.com/FloatTech/zbpctrl"
	rei "github.com/fumiama/ReiBot"
	tgba "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
)

var engine = rei.Register("tools", &ctrl.Options[*rei.Ctx]{
	DisableOnDefault: false,
	Help:             "tools for Lucy",
})

func init() {
	engine.OnMessageCommand("leave", rei.SuperUserPermission).SetBlock(true).Handle(func(ctx *rei.Ctx) {
		arg := strings.TrimSpace(ctx.State["args"].(string))
		var gid int64
		var err error
		if arg != "" {
			gid, err = strconv.ParseInt(arg, 10, 64)
			if err != nil {
				_, _ = ctx.SendPlainMessage(false, "ERROR: ", err)
				return
			}
		} else {
			gid = ctx.Message.Chat.ID
		}
		_, _ = ctx.Caller.Send(&tgba.LeaveChatConfig{ChatID: gid})
	})
	engine.OnMessageCommand("status").SetBlock(true).Handle(func(ctx *rei.Ctx) {
		ctx.SendPlainMessage(false, "* Hosted On Azure JP Cloud.\n",
			"* CPU Usage: ", cpuPercent(), "%\n",
			"* RAM Usage: ", memPercent(), "%\n",
			"* DiskInfo Usage Check: ", diskPercent(), "\n",
			"  Lucyは、高性能ですから！")
	})
	engine.OnMessage().SetBlock(false).Handle(func(ctx *rei.Ctx) {
		toolchain.FastSaveUserStatus(ctx)
	})
	engine.OnMessageCommand("dataupdate").SetBlock(true).Handle(func(ctx *rei.Ctx) {
		if !toolchain.GetTheTargetIsNormalUser(ctx) {
			return
		}
		getUserName := ctx.Message.From.UserName
		getUserID := ctx.Message.From.ID
		newUserName := CoreFactory.GetUserSampleUserinfobyid(getUserID).UserName
		if newUserName == getUserName {
			ctx.SendPlainMessage(true, "不需要更新的~用户名为最新w")
			return
		}
		CoreFactory.StoreUserDataBase(getUserID, newUserName)
	})
}

func cpuPercent() float64 {
	percent, err := cpu.Percent(time.Second, false)
	if err != nil {
		return -1
	}
	return math.Round(percent[0])
}

func memPercent() float64 {
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		return -1
	}
	return math.Round(memInfo.UsedPercent)
}

func diskPercent() string {
	parts, err := disk.Partitions(true)
	if err != nil {
		return err.Error()
	}
	msg := ""
	for _, p := range parts {
		diskInfo, err := disk.Usage(p.Mountpoint)
		if err != nil {
			msg += "\n  - " + err.Error()
			continue
		}
		pc := uint(math.Round(diskInfo.UsedPercent))
		if pc > 0 {
			msg += fmt.Sprintf("\n  - %s(%dM) %d%%", p.Mountpoint, diskInfo.Total/1024/1024, pc)
		}
	}
	return msg
}
