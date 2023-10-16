// Package wife Forks from https://github.com/Jiang-Red/ReiBot-Plugin/blob/main/plugin/groupwife/init.go Modfied by MoeMagicMango
package wife

import (
	"math/rand"
	"time"

	"github.com/FloatTech/floatbox/math"
	ctrl "github.com/FloatTech/zbpctrl"
	rei "github.com/fumiama/ReiBot"
	tgba "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	engine = rei.Register("wife", &ctrl.Options[*rei.Ctx]{
		DisableOnDefault:  false,
		Help:              "wife - /wife to choose one.",
		PrivateDataFolder: "wife",
	})
)

func init() {
	go func() {
		db.db.DBPath = engine.DataFolder() + "data.db"
		err := db.db.Open(time.Hour * 24)
		if err != nil {
			panic(err)
		}
		err = db.db.Create("groupinfo", &groupinfo{})
		if err != nil {
			panic(err)
		}
		err = db.db.Create("favorability", &favorability{})
		if err != nil {
			panic(err)
		}
		err = db.db.Create("cooling", &cooling{})
		if err != nil {
			panic(err)
		}
	}()
	engine.OnMessageCommand("wife").SetBlock(true).Handle(func(ctx *rei.Ctx) {
		gid := ctx.Message.Chat.ID
		err := db.checktime(gid)
		if err != nil {
			_, _ = ctx.SendPlainMessage(false, "ERROR: ", err)
			return
		}
		uid := ctx.Message.From.ID
		info, _ := db.findcertificates(gid, uid)
		switch {
		case info != certificates{} && (info.ManID == 0 || info.WomanID == 0):
			_, _ = ctx.SendPlainMessage(false, "今天你选择了保持单身")
			return
		case info.ManID == uid:
			_, _ = ctx.Caller.Send(&tgba.PhotoConfig{
				BaseFile: tgba.BaseFile{
					BaseChat: tgba.BaseChat{
						ChatID: ctx.Message.Chat.ID,
					},
					File: func() tgba.RequestFileData {
						p, err := ctx.Caller.GetUserProfilePhotos(tgba.NewUserProfilePhotos(info.WomanID))
						if err == nil && len(p.Photos) > 0 {
							fp := p.Photos[0]
							return tgba.FileID(fp[len(fp)-1].FileID)
						}
						return nil
					}(),
				},
				Caption: "@" + ctx.Message.From.String() +
					"今天你娶了老婆" + "[" + info.WomanName + "]" +
					"(https://t.me/" + info.WomanName + ")",
				ParseMode: "Markdown",
			})
			return
		case info.WomanID == uid:
			_, _ = ctx.Caller.Send(&tgba.PhotoConfig{
				BaseFile: tgba.BaseFile{
					BaseChat: tgba.BaseChat{
						ChatID: ctx.Message.Chat.ID,
					},
					File: func() tgba.RequestFileData {
						p, err := ctx.Caller.GetUserProfilePhotos(tgba.NewUserProfilePhotos(info.ManID))
						if err == nil && len(p.Photos) > 0 {
							fp := p.Photos[0]
							return tgba.FileID(fp[len(fp)-1].FileID)
						}
						return nil
					}(),
				},
				Caption: "@" + ctx.Message.From.String() +
					"今天你嫁给老公" + "[" + info.ManName + "]" +
					"(https://t.me/" + info.ManName + ")",
				ParseMode: "Markdown",
			})
			return
		}
		groupmemberlist, _ := ctx.Caller.GetChatAdministrators(tgba.ChatAdministratorsConfig{
			ChatConfig: tgba.ChatConfig{
				ChatID: ctx.Message.Chat.ID},
		})
		groupmemberlist = groupmemberlist[math.Max(0, len(groupmemberlist)-30):]
		memberlist := make([]int64, 0, len(groupmemberlist))
		for i := 0; i < len(groupmemberlist); i++ {
			user := groupmemberlist[i].User.ID
			info, _ := db.findcertificates(gid, user)
			if (info != certificates{}) {
				continue
			}
			memberlist = append(memberlist, user)
		}
		if len(memberlist) <= 1 {
			_, _ = ctx.SendPlainMessage(false, "群里没有人是单身了哦~")
			return
		}
		time := nowtime()
		target := memberlist[rand.Intn(len(memberlist))]
		if target == uid {
			switch rand.Intn(10) {
			case 5:
				err := db.updatecertificates(gid, &certificates{
					ManID:      uid,
					WomanID:    0,
					UpdateTime: time,
				})
				if err != nil {
					_, _ = ctx.SendPlainMessage(false, "ERROR: ", err)
					return
				}
				_, _ = ctx.SendPlainMessage(false, "今天你选择了保持单身")
			default:
				_, _ = ctx.SendPlainMessage(false, "谁也没有娶到...")
				return
			}
		}
		uidinfo, _ := ctx.Caller.GetChatMember(tgba.GetChatMemberConfig{
			ChatConfigWithUser: tgba.ChatConfigWithUser{
				ChatID: gid,
				UserID: uid,
			},
		})
		targetinfo, _ := ctx.Caller.GetChatMember(tgba.GetChatMemberConfig{
			ChatConfigWithUser: tgba.ChatConfigWithUser{
				ChatID: gid,
				UserID: target,
			},
		})
		info = certificates{
			ManID:      uid,
			WomanID:    target,
			ManName:    uidinfo.User.UserName,
			WomanName:  targetinfo.User.UserName,
			UpdateTime: time,
		}
		err = db.updatecertificates(gid, &info)
		if err != nil {
			_, _ = ctx.SendPlainMessage(false, "ERROR: ", err)
			return
		}
		_, _ = ctx.Caller.Send(&tgba.PhotoConfig{
			BaseFile: tgba.BaseFile{
				BaseChat: tgba.BaseChat{
					ChatID: gid,
				},
				File: func() tgba.RequestFileData {
					p, err := ctx.Caller.GetUserProfilePhotos(tgba.NewUserProfilePhotos(info.WomanID))
					if err == nil && len(p.Photos) > 0 {
						fp := p.Photos[0]
						return tgba.FileID(fp[len(fp)-1].FileID)
					}
					return nil
				}(),
			},
			Caption: "@" + ctx.Message.From.String() +
				"今天你娶了老婆" + "[" + info.WomanName + "]" +
				"(https://t.me/" + info.WomanName + ")",
			ParseMode: "Markdown",
		})
	})
}
