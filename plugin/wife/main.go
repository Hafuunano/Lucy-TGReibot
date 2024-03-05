// Package wife From "github.com/MoYoez/Lucy_zerobot"
package wife

import (
	"bytes"
	"image/png"
	"math/rand"
	"regexp"
	"strconv"
	"time"

	text "github.com/FloatTech/imgfactory"
	coins "github.com/MoYoez/Lucy_reibot/utils/coins"
	"github.com/MoYoez/Lucy_reibot/utils/toolchain"
	"github.com/MoYoez/Lucy_reibot/utils/transform"

	ctrl "github.com/FloatTech/zbpctrl"
	rei "github.com/fumiama/ReiBot"
	tgba "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	engine = rei.Register("wife", &ctrl.Options[*rei.Ctx]{
		DisableOnDefault:  false,
		Help:              "Hi NekoPachi!",
		PrivateDataFolder: "wife",
	})
)

/*
StatusID:

Type 1: Normal Mode, nothing happened.

Type 2: Cannot be the Target, Target became initiative, so reverse.

(However the target and the initiative should be in their position, DO NOT CHANGE. )

Type 3: Something is wrong, you are Target == initiative Person. (Drop The Person Before.)

Type 4: Removed.
(When User get others person. || IF REMARRIED, CHANGE IT TO TYPE1.) || (Be check more Time to reduce to err.)

Type 5: NTR Mode
(Tips: NTR means changed their pairkey & TargetID || UserID, need to do some changes. ) ||
(Attempt to do once more every person.)

Type 6: No wife Mod?
Fake - Invisible person here.
(Lucy Hides this and shows it in the next Time if a person uses NTR,
shows nothing, and Lucy will make it for joke. LMAO)

Type 7: NTRED BY SOMEONE.
*/

func init() {
	sdb := coins.Initialize("./data/score/score.db")
	dict := make(map[string][]string) // this dict is used to reply
	// dict path.
	dict["block"] = []string{"嗯哼？貌似没有找到哦w", "再试试哦w，或许有帮助w", "运气不太好哦，想一下办法呢x"}
	dict["success"] = []string{"Lucky For You~", "恭喜哦ww~ ", "这边来恭喜一下哦w～", "貌似很成功的一次尝试呢w~"}
	dict["failed"] = []string{"今天的运气有一点背哦~这一次没有成功呢x", "_(:з」∠)_下次还有机会 抱抱w", "没关系哦，虽然失败了但还有机会呢x"}
	dict["ntr"] = []string{"嗯哼～这位还是成功了呢x", "aaa 好怪 不过还是让你通过了 ^^ "}
	dict["lost_failed"] = []string{"为什么要分呢? 让咱捏捏w", "太坏了啦！不许！"}
	dict["lost_success"] = []string{"好呢w 就这样呢(", "已经成功了哦w"}
	dict["hide_mode"] = []string{"哼哼～ 哼唧", "喵喵喵？！"}

	engine.OnMessageCommand("marry", rei.OnlyGroupOrSuperGroup).SetBlock(true).Handle(func(ctx *rei.Ctx) {
		// 结婚
		// command patterns
		// marry @user
		// in telegram, we should consider user more. || marry to someone (you(ctx) are the main.)
		getEntities := toolchain.ListEntitiesMention(ctx)
		uid := ctx.Message.From.ID
		if len(getEntities) == 0 {
			ctx.SendPlainMessage(true, "没有找到用户QAQ, help: /command @User")
			return
		}
		fiancee := toolchain.GetUserIDFromUserName(ctx, getEntities[0])
		if !CheckDisabledListIsExistedInThisGroup(marryList, uid, ctx.Message.Chat.ID) {
			ctx.SendPlainMessage(true, "你已经禁用了被随机，所以不可以参与娶群友哦w")
			return
		}
		// fast check
		if !CheckTheUserStatusAndDoRepeat(ctx) {
			return
		}
		if !CheckTheTargetUserStatusAndDoRepeat(ctx, fiancee) {
			return
		}
		// check the target status.
		getStatusIfBannned := CheckTheUserIsInBlackListOrGroupList(fiancee, uid, ctx.Message.Chat.ID)
		/*
			disabled_Target
			blacklist_Target
		*/
		if getStatusIfBannned {
			// blocked.
			GlobalCDModelCost(ctx)
			getReply := dict["block"][rand.Intn(len(dict["block"]))]
			ctx.SendPlainMessage(true, getReply)
			return
		}
		if GlobalCDModelCostLeastReply(ctx) == 0 {
			ctx.SendPlainMessage(true, "今天的机会已经使用完了哦～12小时后再来试试吧")
			return
		}
		if uid == fiancee {
			switch rand.Intn(5) {
			case 1:
				GlobalCDModelCost(ctx)
				ReplyMeantMode("貌似Lucy故意添加了 --force 的命令，成功了(笑 ", uid, 1, ctx)
				generatePairKey := GenerateMD5(uid, uid, ctx.Message.Chat.ID)
				err := InsertUserGlobalMarryList(marryList, ctx.Message.Chat.ID, uid, uid, 3, generatePairKey)
				if err != nil {
					panic(err)
				}
			default:
				GlobalCDModelCost(ctx)
				ctx.SendPlainMessage(true, "笨蛋！娶你自己干什么a")
			}
			return
		}
		// However Lucy is only available to be married. LOL.
		if fiancee == ctx.Caller.Self.ID {
			// not work yet, so just the next path.
			if rand.Intn(100) > 90 {
				ctx.SendPlainMessage(true, "笨蛋！不准娶~ ama")
				GlobalCDModelCost(ctx)
				return
			}
			// do it.
			GlobalCDModelCost(ctx)
			getSuccessMsg := dict["success"][rand.Intn(len(dict["success"]))]
			// normal mode. nothing happened.
			ReplyMeantMode(getSuccessMsg, fiancee, 1, ctx)
			generatePairKey := GenerateMD5(uid, fiancee, ctx.Message.Chat.ID)
			_ = InsertUserGlobalMarryList(marryList, ctx.Message.Chat.ID, uid, fiancee, 1, generatePairKey)
			return

		}
		ResuitTheReferUserAndMakeIt(ctx, dict, uid, fiancee)
	})
	engine.OnMessageCommand("wife", rei.OnlyGroupOrSuperGroup).SetBlock(true).Handle(func(ctx *rei.Ctx) {
		/*
			Work:
			- Check the User Status, if the user is 1 or 0 || 10 ,then pause and do this handler.
			- Choose a person, do something acciednt. (if the person had, pause and give one more chance.)
			- Check the banned or Disabled Status (To Target,if had,then stoppped it and give no chance. Others has checked itself too. )
			- add this key.
			- add more feature.
		*/
		/*
			TODO: HIDE MODE TYPE 6
		*/
		uid := ctx.Message.From.ID
		gid := ctx.Message.Chat.ID
		// Check if Disabled this group.
		if !CheckDisabledListIsExistedInThisGroup(marryList, uid, gid) {
			ctx.SendPlainMessage(true, "你已经禁用了被随机，所以不可以参与娶群友哦w")
			return
		}
		// fast check (Check User Itself.)
		if !CheckTheUserStatusAndDoRepeat(ctx) {
			return
		}
		ChooseAPerson := GetUserListAndChooseOne(ctx)
		if ChooseAPerson == 0 {
			ctx.SendPlainMessage(true, "貌似你需要等一会试试呢~ Lucy正在确认群里的人数w")
			return
		}
		// ok , go next. || before that we should check this person is in the lucky list?
		// Luck Path. (Only available in marry action.)
		getLuckyChance, getLuckyPeople, getLuckyTime := CheckTheOrderListAndBackDetailed(uid, gid)
		getCurrentTime := time.Now().Unix()
		getLuckyTimeToInt64, _ := strconv.ParseInt(getLuckyTime, 10, 64)
		if getLuckyChance > 10 && getLuckyTimeToInt64 < getCurrentTime {
			if getLuckyTimeToInt64+(24*60*60) < getCurrentTime {
				ctx.SendPlainMessage(true, "貌似时间过去的太久了 这一次算是撤销掉了哦x,不过返回消耗的柠檬片，并且机会不变")
				_ = RemoveOrderToList(marryList, uid, gid)
				getUserID := coins.GetSignInByUID(sdb, uid)
				_ = coins.InsertUserCoins(sdb, uid, getUserID.Coins+1000)
				return
			}
			getTargetStatusCode, _ := CheckTheUserIsTargetOrUser(marryList, ctx, getLuckyPeople) // 判断这个target是否已经和别人在一起了，同时判断Type3
			if getTargetStatusCode == -1 {
				// do this?
				getExistedToken := GlobalCDModelCostLeastReply(ctx)
				if getExistedToken == 0 {
					ctx.SendPlainMessage(true, "今天的机会已经使用完了哦～12小时后再来试试吧，不过这边可以透露一下～许愿池已经实现了哦w～不过给你保留这次机会w")
					return
				}
				// check the target status.
				getStatusIfBannned := CheckTheUserIsInBlackListOrGroupList(getLuckyPeople, uid, gid)
				if getStatusIfBannned {
					// blocked.
					ctx.SendPlainMessage(true, "看起来挺倒霉的～貌似对方在许愿的过程中加入了黑名单x,或者对方已经禁用了对于本群的功能，只能无情删掉了哦x,不过这一次机会不会被浪费掉，并且会返回相应的柠檬片～")
					getUserID := coins.GetSignInByUID(sdb, uid)
					_ = coins.InsertUserCoins(sdb, uid, getUserID.Coins+1000)
					_ = RemoveOrderToList(marryList, uid, gid)
					return
				}
				// success .
				GlobalCDModelCost(ctx)
				getSuccessMsg := dict["success"][rand.Intn(len(dict["success"]))]
				// normal mode. nothing happened.
				ReplyMeantMode("许愿池生效～\n"+getSuccessMsg, getLuckyPeople, 1, ctx)
				generatePairKey := GenerateMD5(uid, getLuckyPeople, gid)
				_ = InsertUserGlobalMarryList(marryList, gid, uid, getLuckyPeople, 1, generatePairKey)
				_ = RemoveOrderToList(marryList, uid, gid)
				return
			}
			// target not -1 (has others.)
			// didn't do it.
			GlobalCDModelCost(ctx)
			_ = RemoveOrderToList(marryList, uid, gid)
			ctx.SendPlainMessage(true, "抱歉哦～虽然已经使用了愿望池，不过仍然没有成功呢awa～")
			// handle this chance but no cares
			return

		}

		// Luck Path end.

		if !CheckTheTargetUserStatusAndDoRepeat(ctx, ChooseAPerson) {
			return
		}
		// check the target status.
		getStatusIfBannned := CheckTheUserIsInBlackListOrGroupList(ChooseAPerson, uid, gid)
		/*
			disabled_Target
			blacklist_Target
		*/
		if getStatusIfBannned {
			// blocked.
			GlobalCDModelCost(ctx)
			getReply := dict["block"][rand.Intn(len(dict["block"]))]
			ctx.SendPlainMessage(true, getReply)
			return
		}
		// go next. do something colorful, pls cost something.
		getExistedToken := GlobalCDModelCostLeastReply(ctx)
		if getExistedToken == 0 {
			ctx.SendPlainMessage(true, "今天的机会已经使用完了哦～12小时后再来试试吧")
			return
		}
		// one chance to get himself | herself
		if ChooseAPerson == uid {
			// status code 3
			GlobalCDModelCost(ctx)
			// drop target pls.
			ReplyMeantMode("嗯哼哼～抽到了自己，然而 Lucy 还是将双方写成一个人哦w （笑w ", uid, 1, ctx)
			generatePairKey := GenerateMD5(uid, uid, gid)
			_ = InsertUserGlobalMarryList(marryList, gid, uid, uid, 3, generatePairKey)

		}
		returnNumber := GetSomeRanDomChoiceProps(ctx)
		switch {
		case returnNumber == 1:
			GlobalCDModelCost(ctx)
			getSuccessMsg := dict["success"][rand.Intn(len(dict["success"]))]
			// normal mode. nothing happened.
			ReplyMeantMode(getSuccessMsg, ChooseAPerson, 1, ctx)
			generatePairKey := GenerateMD5(uid, ChooseAPerson, gid)
			_ = InsertUserGlobalMarryList(marryList, gid, uid, ChooseAPerson, 1, generatePairKey)
		case returnNumber == 2:
			GlobalCDModelCost(ctx)
			ReplyMeantMode("貌似很奇怪哦～因为某种奇怪的原因～1变成了0,0变成了1", ChooseAPerson, 0, ctx)
			generatePairKey := GenerateMD5(ChooseAPerson, uid, gid)
			_ = InsertUserGlobalMarryList(marryList, gid, ChooseAPerson, uid, 2, generatePairKey)
		// reverse Target Mode
		case returnNumber == 3:
			GlobalCDModelCost(ctx)
			// drop target pls.
			ReplyMeantMode("嗯哼哼～发生了一些错误～本来应当抽到别人的变成了自己～所以", uid, 1, ctx)
			generatePairKey := GenerateMD5(uid, uid, gid)
			_ = InsertUserGlobalMarryList(marryList, gid, uid, uid, 3, generatePairKey)
		// you became your own target
		case returnNumber == 6:
			GlobalCDModelCost(ctx)
			// now no wife mode.
			getHideMsg := dict["hide_mode"][rand.Intn(len(dict["hide_mode"]))]
			ctx.SendPlainMessage(true, getHideMsg, "\n貌似没有任何反馈～")
			generatePairKey := GenerateMD5(uid, ChooseAPerson, gid)
			_ = InsertUserGlobalMarryList(marryList, gid, uid, uid, 6, generatePairKey)
		}
	})
	engine.OnMessageCommand("divorce", rei.OnlyGroupOrSuperGroup).SetBlock(true).Handle(func(ctx *rei.Ctx) {
		getStatusCode, _ := CheckTheUserIsTargetOrUser(marryList, ctx, ctx.Message.From.ID)
		if getStatusCode == -1 {
			ctx.SendPlainMessage(true, "貌似？没有对象的样子x")
			return
		}
		if LeaveCDModelCostLeastReply(ctx) == 0 {
			ctx.SendPlainMessage(true, "今天的次数已经用完了哦～或许可以试一下别的方式？")
			return
		}
		reverseCheckTheUserIsDisabled := CheckDisabledListIsExistedInThisGroup(marryList, ctx.Message.From.ID, ctx.Message.Chat.ID)
		if !reverseCheckTheUserIsDisabled {
			ctx.SendPlainMessage(true, "你已经禁用了被随机，所以不可以参与娶群友哦w")
			return
		}
		getlostSuccessedMsg := dict["lost_success"][rand.Intn(len(dict["lost_success"]))]
		getLostFailedMsg := dict["lost_failed"][rand.Intn(len(dict["lost_failed"]))]
		if rand.Intn(4) >= 2 {
			LeaveCDModelCost(ctx)
			ctx.SendPlainMessage(true, getLostFailedMsg)
		} else {
			LeaveCDModelCost(ctx)
			getPairKey := CheckThePairKey(marryList, ctx.Message.From.ID, ctx.Message.Chat.ID)
			RemoveUserGlobalMarryList(marryList, getPairKey, ctx.Message.Chat.ID)
			ctx.SendPlainMessage(true, getlostSuccessedMsg)
		}
	})
	engine.OnMessageCommand("chwaifu", rei.OnlyGroupOrSuperGroup).SetBlock(true).Handle(func(ctx *rei.Ctx) {
		// command patterns
		// marry @user
		// in telegram, we should consider user more.
		getEntities := toolchain.ListEntitiesMention(ctx)
		if len(getEntities) == 0 {
			ctx.SendPlainMessage(true, "没有找到用户QAQ, help: /command @User")
			return
		}
		fiancee := toolchain.GetUserIDFromUserName(ctx, getEntities[0])
		uid := ctx.Message.From.ID
		groupID := ctx.Message.Chat.ID
		reverseCheckTheUserIsDisabled := CheckDisabledListIsExistedInThisGroup(marryList, uid, groupID)
		if !reverseCheckTheUserIsDisabled {
			ctx.SendPlainMessage(true, "你已经禁用了被随机，所以不可以参与娶群友哦w")
			return
		}
		if fiancee == uid {
			ctx.SendPlainMessage(true, "要骗别人哦~为什么要骗自己呢x")
			return
		}
		// this case should other people existed.
		getStatusCode, _ := CheckTheUserIsTargetOrUser(marryList, ctx, uid)
		if getStatusCode != -1 {
			ctx.SendPlainMessage(true, "貌似你已经有了哦～？难不成时要找 ^^ 别人嘛（恼")
			return
		}
		getTargetStatusCode, _ := CheckTheUserIsTargetOrUser(marryList, ctx, fiancee)
		if getTargetStatusCode == -1 {
			ctx.SendPlainMessage(true, "嗯哼～这位还是一个人哦w～可以不用这个的哦w")
			return
		}
		// low possibility to get this chance.
		if LeaveCDModelCostLeastReply(ctx) <= 0 {
			ctx.SendPlainMessage(true, "今日机会不够哦w，过段时间再来试试吧w")
			return
		}
		LeaveCDModelCost(ctx)
		if rand.Intn(100) < 30 {
			// win this goal
			getNTRMsg := dict["ntr"][rand.Intn(len(dict["ntr"]))]
			ReplyMeantMode(getNTRMsg, fiancee, 5, ctx)
			CustomRemoveUserGlobalMarryList(marryList, CheckThePairKey(marryList, fiancee, groupID), groupID, 7)
			pairKey := GenerateMD5(uid, fiancee, groupID)
			err := InsertUserGlobalMarryList(marryList, groupID, uid, fiancee, 5, pairKey)
			if err != nil {
				ctx.SendPlainMessage(true, "ERR: ", err)
				return
			}
		} else {
			getFailed := dict["failed"][rand.Intn(len(dict["failed"]))]
			ctx.SendPlainMessage(true, getFailed)
			return
		}
	})
	engine.OnMessageCommand("waifulist", rei.OnlyGroupOrSuperGroup).SetBlock(true).Handle(func(ctx *rei.Ctx) {
		gid := ctx.Message.Chat.ID
		getList, num := GetTheGroupList(gid)
		var RawMsg string
		if num == 0 {
			ctx.SendPlainMessage(true, "本群貌似还没有人结婚来着（")
			return
		}
		emojiRegex := regexp.MustCompile(`[\x{1F600}-\x{1F64F}|[\x{1F300}-\x{1F5FF}]|[\x{1F680}-\x{1F6FF}]|[\x{1F700}-\x{1F77F}]|[\x{1F780}-\x{1F7FF}]|[\x{1F800}-\x{1F8FF}]|[\x{1F900}-\x{1F9FF}]|[\x{1FA00}-\x{1FA6F}]|[\x{1FA70}-\x{1FAFF}]|[\x{1FB00}-\x{1FBFF}]|[\x{1F170}-\x{1F251}]|[\x{1F300}-\x{1F5FF}]|[\x{1F600}-\x{1F64F}]|[\x{1FC00}-\x{1FCFF}]|[\x{1F004}-\x{1F0CF}]|[\x{1F170}-\x{1F251}]]+`)
		for i := 0; i < num; i++ {
			getUserInt64, _ := strconv.ParseInt(getList[i][0], 10, 64)
			getTargetInt64, _ := strconv.ParseInt(getList[i][1], 10, 64)
			RawMsg += strconv.FormatInt(int64(i+1), 10) + ". " + emojiRegex.ReplaceAllString(toolchain.GetUserNickNameByIDInGroup(ctx, getUserInt64), "") + "( " + getList[i][0] + " )" + "  -->  " + emojiRegex.ReplaceAllString(toolchain.GetUserNickNameByIDInGroup(ctx, getTargetInt64), "") + "( " + getList[i][1] + " )" + "\n"
		}
		headerMsg := "群老婆列表～ For Group( " + strconv.FormatInt(gid, 10) + " )" + " [ " + ctx.Message.Chat.Title + " ]\n\n"
		base64Font, err := text.RenderText(headerMsg+RawMsg+"\n\n Tips: 此列表将会在 23：00 PM (GMT+8) 重置", transform.ReturnLucyMainDataIndex("Font")+"regular-bold.ttf", 1920, 45)
		var buf bytes.Buffer
		png.Encode(&buf, base64Font)
		if err != nil {
			ctx.SendPlainMessage(true, "ERR: ", err)
			return
		}
		ctx.SendPhoto(tgba.FileReader{Name: "waifu_" + strconv.FormatInt(gid, 10), Reader: bytes.NewReader(buf.Bytes())}, true, "")
	})
	engine.OnMessageCommand("waifuwish", rei.OnlyGroupOrSuperGroup).SetBlock(true).Handle(func(ctx *rei.Ctx) {
		getEntities := toolchain.ListEntitiesMention(ctx)
		uid := ctx.Message.From.ID
		gid := ctx.Message.Chat.ID
		si := coins.GetSignInByUID(sdb, uid)
		if len(getEntities) == 0 {
			ctx.SendPlainMessage(true, "没有找到用户QAQ, help: /command @User")
			return
		}
		fiancee := toolchain.GetUserIDFromUserName(ctx, getEntities[0])
		reverseCheckTheUserIsDisabled := CheckDisabledListIsExistedInThisGroup(marryList, uid, gid)
		if !reverseCheckTheUserIsDisabled {
			ctx.SendPlainMessage(true, "你已经禁用了被随机，所以不可以参与娶群友哦w")
			return
		}
		if si.Coins < 1000 {
			ctx.SendPlainMessage(true, "本次许愿的柠檬片不足哦～需要1000个柠檬片才可以哦w")
			return
		}
		if !CheckTheBlackListIsExistedToThisPerson(marryList, fiancee, uid) || !CheckTheBlackListIsExistedToThisPerson(marryList, uid, fiancee) {
			ctx.SendPlainMessage(true, "已经被Ban掉了，愿望无法实现～")
			return
		}
		_, getTargetID, _ := CheckTheOrderListAndBackDetailed(uid, gid)
		if getTargetID != fiancee && getTargetID != 0 {
			ctx.SendPlainMessage(true, "每次仅可以许愿一个人w 不允许第二个人")
			return
		}
		if fiancee == uid {
			ctx.SendPlainMessage(true, "坏哦！为什么要许自己的x")
			return
		}
		if getTargetID == fiancee {
			ctx.SendPlainMessage(true, "已经许过一次了哦～不需要第二次")
			return
		}
		// Handler
		_ = coins.InsertUserCoins(sdb, uid, si.Coins-1000)
		timeStamp := time.Now().Unix() + (6 * 60 * 60)
		_ = AddOrderToList(marryList, uid, fiancee, strconv.FormatInt(timeStamp, 10), gid)
		ctx.SendPlainMessage(true, "已经许愿成功了哦～w 给", toolchain.GetUserNickNameByIDInGroup(ctx, fiancee), " 的许愿已经生效，将会在6小时后增加70%可能性实现w")

	})
}
