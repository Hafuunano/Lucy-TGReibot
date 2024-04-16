package tools

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/FloatTech/floatbox/web"
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/MoYoez/Lucy_reibot/utils/toolchain"
	CoreFactory "github.com/MoYoez/Lucy_reibot/utils/userpackage"
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
		_, _ = ctx.Caller.Send(&tgba.LeaveChatConfig{ChatConfig: tgba.ChatConfig{ChatID: gid}})
	})
	engine.OnMessageCommand("status").SetBlock(true).Handle(func(ctx *rei.Ctx) {
		ctx.SendPlainMessage(false, "* Hosted On Azure JP Cloud.\n",
			"* CPU Usage: ", cpuPercent(), "%\n",
			"* RAM Usage: ", memPercent(), "%\n",
			"* DiskInfo Usage Check: ", diskPercent(), "\n",
			"  Lucyは、高性能ですから！")
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
	engine.OnMessage().SetBlock(false).Handle(func(ctx *rei.Ctx) {
		toolchain.FastSaveUserStatus(ctx)

	})
	engine.OnMessage().SetBlock(false).Handle(func(ctx *rei.Ctx) {
		toolchain.FastSaveUserGroupList(ctx) // error
	})
	engine.OnMessageCommand("runpanic", rei.SuperUserPermission).Handle(func(ctx *rei.Ctx) {
		ctx.SendPlainMessage(true, "run panic , check debug.")
		panic("Test Value ERR")
	})

	engine.OnMessageCommand("qpic").SetBlock(true).Handle(func(ctx *rei.Ctx) {
		getLength, List := rei.SplitCommandTo(ctx.Message.Text, 2)
		if getLength == 2 {
			getDataRaw, err := web.GetData("https://gchat.qpic.cn/gchatpic_new/0/0-0-" + List[1] + "/0")
			if err != nil {
				ctx.SendPlainMessage(true, "获取对应图片错误,或许是图片已过期")
				return
			}
			ctx.SendPhoto(tgba.FileBytes{Name: List[1], Bytes: getDataRaw}, true, "Link: "+"https://gchat.qpic.cn/gchatpic_new/0/0-0-"+List[1]+"/0")
		} else {
			ctx.SendPlainMessage(true, "缺少参数/ 应当是 /qpic <md5> ")
		}
	})

	engine.OnMessageCommand("title").SetBlock(true).Handle(func(ctx *rei.Ctx) {
		getCutterLength, cutterTypeList := rei.SplitCommandTo(ctx.Message.Text, 2)
		// Check if Lucy has Permission to modify
		// get user permissions.
		_, err := ctx.Caller.Request(tgba.PromoteChatMemberConfig{
			ChatMemberConfig: tgba.ChatMemberConfig{UserID: ctx.Message.From.ID, ChatConfig: tgba.ChatConfig{
				ChatID: ctx.Message.Chat.ID,
			}},
			CanManageChat: true,
		})
		if err != nil {
			ctx.SendPlainMessage(true, " 发生了一点错误: 将对方提升管理员失效 , Err: ", err)
			return
		}

		if getCutterLength == 1 {
			getendpoint, errs := ctx.Caller.Request(tgba.SetChatAdministratorCustomTitle{
				ChatMemberConfig: tgba.ChatMemberConfig{
					ChatConfig: tgba.ChatConfig{
						ChatID: ctx.Message.Chat.ID,
					},
					UserID: ctx.Message.From.ID,
				},
				CustomTitle: ctx.Message.From.UserName,
			})
			if getendpoint.Ok {
				ctx.SendPlainMessage(true, "是个不错的头衔呢w~")
			} else {
				ctx.SendPlainMessage(true, "貌似出错了( | ", errs)
			}
			return
		}

		getendpoint, err := ctx.Caller.Request(tgba.SetChatAdministratorCustomTitle{
			ChatMemberConfig: tgba.ChatMemberConfig{
				ChatConfig: tgba.ChatConfig{
					ChatID: ctx.Message.Chat.ID,
				},
				UserID: ctx.Message.From.ID,
			},
			CustomTitle: cutterTypeList[1],
		})

		if getendpoint.Ok {
			ctx.SendPlainMessage(true, "返回正常, 帮你贴上去了w 现在的头衔是 ", cutterTypeList[1], " 了")
		} else {
			ctx.SendPlainMessage(true, "貌似出错了( | ", err)
		}

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
