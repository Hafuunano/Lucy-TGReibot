package score

import (
	"encoding/json"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/FloatTech/ReiBot-Plugin/utils/CoreFactory"
	coins "github.com/FloatTech/ReiBot-Plugin/utils/coins"
	"github.com/FloatTech/ReiBot-Plugin/utils/ctxext"
	"github.com/FloatTech/ReiBot-Plugin/utils/toolchain"
	rei "github.com/fumiama/ReiBot"
	tgba "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/wdvxdr1123/ZeroBot/extension/rate"
)

type partygame struct {
	Name  string `json:"name"`
	Desc  string `json:"desc"`
	Coins string `json:"coins"`
}

var (
	pgs            = make(pg, 256)
	RobTimeManager = rate.NewManager[int64](time.Minute*70, 163)
	checkLimit     = rate.NewManager[int64](time.Minute*1, 5) // time setup
	catchLimit     = rate.NewManager[int64](time.Hour*1, 9)   // time setup
	processLimit   = rate.NewManager[int64](time.Hour*1, 5)   // time setup
	wagerData      map[string]int
)

type pg = map[string]partygame

func init() {
	wagerData = make(map[string]int)
	wagerData["data"] = rand.Intn(2000)
	sdb := coins.Initialize("./data/score/score.db")
	data, err := os.ReadFile(engine.DataFolder() + "loads.json")
	err = json.Unmarshal(data, &pgs)
	if err != nil {
		return
	}

	engine.OnMessageCommand("coinroll").SetBlock(true).Limit(ctxext.LimitByUser).Handle(func(ctx *rei.Ctx) {
		userID, _ := toolchain.GetChatUserInfoID(ctx)
		getMsgID := ctx.Message.MessageID // caller

		if !checkLimit.Load(userID).Acquire() {
			ctx.SendPlainMessage(true, "太贪心了哦~过会试试吧")
			return
		}
		getProtectStatus := CheckUserIsEnabledProtectMode(userID, sdb)
		if getProtectStatus {
			ctx.SendPlainMessage(true, "已经启动保护模式，不允许参与任何抽奖性质类互动")
			return
		}
		var mutex sync.RWMutex // 添加读写锁以保证稳定性
		mutex.Lock()
		uid := userID
		si := coins.GetSignInByUID(sdb, uid) // 获取用户目前状况信息
		userCurrentCoins := si.Coins         // loading coins status
		if userCurrentCoins < 0 {
			_ = coins.InsertUserCoins(sdb, uid, 0)
			ctx.SendPlainMessage(true, "本次参与的柠檬片不够哦~请多多打卡w")
			return
		} // fix unexpected bug during the code error
		checkEnoughCoins := coins.CheckUserCoins(userCurrentCoins)
		if !checkEnoughCoins {
			ctx.SendPlainMessage(true, "本次参与的柠檬片不够哦~请多多打卡w")
			return
		}
		all := rand.Intn(43) // 一共44种可能性
		referpg := pgs[(strconv.Itoa(all))]
		getName := referpg.Name
		getCoinsStr := referpg.Coins
		getCoinsInt, _ := strconv.Atoi(getCoinsStr)
		getDesc := referpg.Desc
		addNewCoins := si.Coins + getCoinsInt - 60
		_ = coins.InsertUserCoins(sdb, uid, addNewCoins)
		msgUnique, _ := ctx.SendPlainMessage(true, "呼~让咱看看你抽到了什么东西ww\n"+
			"你抽到的是~ "+getName+"\n"+"获得了柠檬片 "+strconv.Itoa(getCoinsInt)+"\n"+getDesc+"\n目前的柠檬片总数为："+strconv.Itoa(addNewCoins))
		time.Sleep(time.Second * 20)
		getReplyMsgID := msgUnique.MessageID
		ctx.Caller.Request(tgba.NewDeleteMessage(ctx.Message.Chat.ID, getMsgID))
		ctx.Caller.Request(tgba.NewDeleteMessage(ctx.Message.Chat.ID, getReplyMsgID))
		mutex.Unlock()
	})
	engine.OnMessageCommand("cointhrow").SetBlock(true).SetBlock(true).Limit(ctxext.LimitByUser).Handle(func(ctx *rei.Ctx) {
		_, splitText := toolchain.SplitCommandTo(ctx.Message.Text, 2)
		modifyCoins := splitText[1]
		modifyCoinsToint, _ := strconv.ParseInt(modifyCoins, 10, 64)
		userID, _ := toolchain.GetChatUserInfoID(ctx)
		handleUser := coins.GetSignInByUID(sdb, userID)
		currentUserCoins := handleUser.Coins
		if currentUserCoins-int(modifyCoinsToint) < 0 {
			ctx.SendPlainMessage(true, "貌似你的柠檬片不够处理呢(")
			return
		}
		hadModifyCoins := currentUserCoins - int(modifyCoinsToint)
		_ = coins.InsertUserCoins(sdb, handleUser.UID, hadModifyCoins)
		ctx.SendPlainMessage(true, "已经帮你扔掉了哦")
	})
	engine.OnMessageCommand("coinwager").SetBlock(true).Limit(ctxext.LimitByUser).Handle(func(ctx *rei.Ctx) {
		// 得到本身奖池大小，如果没有或者被get的情况下获胜
		// this method should deal when we have less starter.
		userid, _ := toolchain.GetChatUserInfoID(ctx)
		_, splitText := toolchain.SplitCommandTo(ctx.Message.Text, 2)
		var rawNumber string
		if len(splitText) == 1 {
			rawNumber = "50"
		} else {
			rawNumber = splitText[1]
		}
		getProtectStatus := CheckUserIsEnabledProtectMode(userid, sdb)
		if getProtectStatus {
			ctx.SendPlainMessage(true, "已经启动保护模式，不允许参与任何抽奖性质类互动")
			return
		}
		modifyCoins, _ := strconv.Atoi(rawNumber)
		if modifyCoins > 1000 {
			ctx.SendPlainMessage(true, "一次性最大投入为1k")
			return
		}
		handleUser := coins.GetSignInByUID(sdb, userid)
		currentUserCoins := handleUser.Coins
		if currentUserCoins-modifyCoins < 0 {
			ctx.SendPlainMessage(true, "貌似你的柠檬片不够处理呢(")
			return
		}
		// first of all , check the user status
		handlerWagerUser := coins.GetWagerUserStatus(sdb, userid)
		if handlerWagerUser.UserExistedStoppedTime > time.Now().Add(-time.Hour*12).Unix() {
			// then not pass | in the freeze time.
			ctx.SendPlainMessage(true, "目前在冷却时间，距离下个可用时间为: ", time.Unix(handlerWagerUser.UserExistedStoppedTime, 0).Add(time.Hour*12).Format(time.DateTime))
			return
		}
		// passed,delete this one and continue || before max is 3500.
		checkUserWagerCoins := handlerWagerUser.InputCountNumber
		if int64(modifyCoins)+checkUserWagerCoins > 3500 {
			ctx.SendPlainMessage(true, "达到冷却最大值，您目前可投入："+strconv.Itoa(int(3500-checkUserWagerCoins)))
			return
		}
		// get wager
		getWager := coins.GetWagerStatus(sdb)
		if getWager.Expected == 0 {
			// it shows that no condition happened.
			// if not maxzine
			// in the wager mode. || start to load
			getGenOne := toolchain.RandSenderPerDayN(time.Now().Unix(), 16500)
			getRandNumber := getGenOne + toolchain.RandSenderPerDayN(time.Now().Unix()+userid, 5000) + 3000
			_ = coins.WagerCoinsInsert(sdb, modifyCoins+wagerData["data"], 0, getRandNumber)
			if int64(modifyCoins)+checkUserWagerCoins == 3500 {
				_ = coins.UpdateWagerUserStatus(sdb, userid, time.Now().Unix(), 0)
			} else {
				_ = coins.UpdateWagerUserStatus(sdb, userid, 0, int64(modifyCoins)+checkUserWagerCoins)
			}
			if getRandNumber <= modifyCoins {
				// winner, he | she is so lucky.^^
				// Lucy will cost 10 percent Coins.
				willRunCoins := math.Round(float64(modifyCoins+getWager.Wagercount) * 0.9)
				_ = coins.InsertUserCoins(sdb, userid, handleUser.Coins+int(willRunCoins)-modifyCoins)
				_ = coins.WagerCoinsInsert(sdb, 0, int(userid), 0)
				wagerData["data"] = int(math.Round(float64(modifyCoins+getWager.Wagercount)*0.1)) - 200
				ctx.SendPlainMessage(true, "w！恭喜哦，奖池中奖了ww，一共获得 ", willRunCoins, " 个柠檬片，当前有 ", handleUser.Coins+int(willRunCoins)-modifyCoins, " 个柠檬片 (获胜者得到奖池 x0.9的柠檬片总数)")
				return
			}
			// not winner
			_ = coins.InsertUserCoins(sdb, handleUser.UID, handleUser.Coins-modifyCoins)
			ctx.SendPlainMessage(true, "没有中奖哦~，当前奖池为："+strconv.Itoa(modifyCoins))
			return
		}
		// not init,start to add.
		getExpected := getWager.Expected
		if int64(modifyCoins)+checkUserWagerCoins == 3500 {
			_ = coins.UpdateWagerUserStatus(sdb, userid, time.Now().Unix(), 0)
		} else {
			_ = coins.UpdateWagerUserStatus(sdb, userid, 0, int64(modifyCoins)+checkUserWagerCoins)
		}
		if getWager.Wagercount+modifyCoins >= getExpected {
			// you are winner!
			willRunCoins := math.Round(float64(modifyCoins+getWager.Wagercount) * 0.9)
			_ = coins.InsertUserCoins(sdb, userid, handleUser.Coins+int(willRunCoins)-modifyCoins)
			_ = coins.WagerCoinsInsert(sdb, 0, int(userid), 0)
			wagerData["data"] = int(math.Round(float64(modifyCoins+getWager.Wagercount)*0.1)) - 200
			ctx.SendPlainMessage(true, "w！恭喜哦，奖池中奖了ww，一共获得 ", willRunCoins, " 个柠檬片，当前有 ", handleUser.Coins+int(willRunCoins)-modifyCoins, " 个柠檬片 (获胜者得到奖池 x0.9的柠檬片总数)")
			return
		}
		_ = coins.WagerCoinsInsert(sdb, getWager.Wagercount+modifyCoins, 0, getExpected)
		_ = coins.InsertUserCoins(sdb, userid, handleUser.Coins-modifyCoins)
		if rand.Intn(10) == 8 {
			ctx.SendPlainMessage(true, "呐～，不会还有大哥哥到现在 "+strconv.Itoa(getWager.Wagercount+modifyCoins)+" 个柠檬片了都没中奖吧？杂鱼～❤，杂鱼～❤")
		} else {
			ctx.SendPlainMessage(true, "没有中奖哦~，当前奖池为: ", getWager.Wagercount+modifyCoins)
		}
		
	})
	engine.OnMessageCommand("coinfull").SetBlock(true).Limit(ctxext.LimitByUser).Handle(func(ctx *rei.Ctx) {
		getList := toolchain.ListEntitiesMention(ctx)
		if len(getList) > 0 {
			TargetInt := toolchain.GetUserIDFromUserName(ctx, getList[0])
			siTargetUser := coins.GetSignInByUID(sdb, TargetInt)
			getTargetName := toolchain.GetNickNameFromUsername(getList[0])
			ctx.SendPlainMessage(true, "这位 ( ", getTargetName, " ) 的柠檬片为", siTargetUser.Coins, "个")
		} else {
			uid, _ := toolchain.GetChatUserInfoID(ctx)
			si := coins.GetSignInByUID(sdb, uid)
			ctx.SendPlainMessage(true, "你的柠檬片数量一共是: "+strconv.Itoa(si.Coins))
		}
	})
	engine.OnMessageCommand("coinrob").SetBlock(true).Limit(ctxext.LimitByUser).Handle(func(ctx *rei.Ctx) {
		_, getCommandSplitInt := toolchain.SplitCommandTo(ctx.Message.Text, 3)
		// handler
		getUserID, _ := toolchain.GetChatUserInfoID(ctx)
		if !catchLimit.Load(getUserID).Acquire() {
			ctx.SendPlainMessage(true, "太贪心了哦~一小时后再来试试吧")
			return
		}
		// check Lucy's Permission
		if len(getCommandSplitInt) < 2 {
			ctx.SendPlainMessage(true, "使用方法为 : /command [@user]")
			return
		}
		if !toolchain.GetTheTargetIsNormalUser(ctx) {
			ctx.SendPlainMessage(true, "暂时不支持匿名身份哦~")
			return
		}
		// getUserInfo
		acquireEntity := toolchain.ListEntitiesMention(ctx)
		if len(acquireEntity) == 0 {
			return
		}
		getUserData := CoreFactory.GetUserSampleUserinfo(strings.Replace(acquireEntity[0], "@", "", 1))
		if getUserData.UserID == 0 {
			ctx.SendPlainMessage(true, "为非正常用户组或者该用户没有在Lucy处记录x~")
			return
		}
		// check the user is in group?
		if !toolchain.CheckIfthisUserInThisGroup(getUserData.UserID, ctx) {
			ctx.SendPlainMessage(true, "该用户不在这个群哦x")
			return
		}
		// the targetUser
		getTargetID := getUserData.UserID
		// get UserID.

		if getTargetID == getUserID {
			ctx.SendPlainMessage(true, "哈? 干嘛骗自己的?坏蛋哦")
			return
		}

		// get coin status, origin from Lucy_zerobot.
		siEventUser := coins.GetSignInByUID(sdb, getUserID)    // 获取主用户目前状况信息
		siTargetUser := coins.GetSignInByUID(sdb, getTargetID) // 获得被抢用户目前情况信息
		switch {
		case siEventUser.Coins < 400:
			ctx.SendPlainMessage(true, "貌似没有足够的柠檬片去准备哦~请多多打卡w")
			return
		case siTargetUser.Coins < 400:
			ctx.SendPlainMessage(true, "太坏了~试图的对象貌似没有足够多的柠檬片~")
			return
		}
		eventUserName := toolchain.GetNickNameFromUsername(ctx.Message.From.UserName)
		eventTargetName := toolchain.GetNickNameFromUsername(acquireEntity[0])
		// token chance.
		// add more possibility to get the chance (0-200)
		getTicket := RobOrCatchLimitManager(getUserID) // full is 1 , least 3. level 1,2,3,
		// however, the total is still 0-400.
		fullChanceToken := rand.Intn(10)
		var modifyCoins int
		if fullChanceToken > 7 { // use it to reduce the chance to lower coins.
			modifyCoins = rand.Intn(200) + 200
		} else {
			modifyCoins = rand.Intn(200)
		}
		getRandomNum := rand.Intn(10)
		PossibilityNum := 6 / getTicket
		setIsTrue := getRandomNum/PossibilityNum != 0
		var remindTicket string
		if getTicket == 3 {
			remindTicket = "目前已经达到疲倦状态，成功率下调到15%，或许考虑一下不要做一个坏人呢～ ^^ "
		}
		if setIsTrue {
			_ = coins.InsertUserCoins(sdb, siEventUser.UID, siEventUser.Coins-modifyCoins)
			_ = coins.InsertUserCoins(sdb, siTargetUser.UID, siTargetUser.Coins+modifyCoins)
			ctx.SendPlainMessage(true, "试着去拿走 ", eventTargetName, " 的柠檬片时,被发现了.\n所以 ", eventUserName, " 失去了 ", modifyCoins, " 个柠檬片\n\n同时 ", eventTargetName, " 得到了 ", modifyCoins, " 个柠檬片\n", remindTicket)
			return
		}
		_ = coins.InsertUserCoins(sdb, siEventUser.UID, siEventUser.Coins+modifyCoins)
		_ = coins.InsertUserCoins(sdb, siTargetUser.UID, siTargetUser.Coins-modifyCoins)
		ctx.SendPlainMessage(true, "试着去拿走 ", eventTargetName, " 的柠檬片时,成功了.\n所以 ", eventUserName, " 得到了 ", modifyCoins, " 个柠檬片\n\n同时 ", eventTargetName, " 失去了 ", modifyCoins, " 个柠檬片\n", remindTicket)

	})
	engine.OnMessageCommand("coincheat").SetBlock(true).Limit(ctxext.LimitByUser).Handle(func(ctx *rei.Ctx) {
		_, info := toolchain.SplitCommandTo(ctx.Message.Text, 3)
		if len(info) < 3 {
			ctx.SendPlainMessage(true, "缺少参数~ 使用方法为 : /command [@user] (number)")
			return
		}
		if !catchLimit.Load(ctx.Message.From.ID).Acquire() {
			ctx.SendPlainMessage(true, "太贪心了哦~一小时后再来试试吧")
			return
		}
		getProtectStatus := CheckUserIsEnabledProtectMode(ctx.Message.From.ID, sdb)
		if getProtectStatus {
			ctx.SendPlainMessage(true, "已经启动保护模式，不允许参与任何抽奖性质类互动")
			return
		}
		getEntitiy := toolchain.ListEntitiesMention(ctx)
		if len(getEntitiy) == 0 {
			ctx.SendPlainMessage(true, "未指定用户~ 使用方法为 : /command [@user] (number)")
			return
		}
		getID := toolchain.GetUserIDFromUserName(ctx, getEntitiy[0])
		if getID == 0 {
			ctx.SendPlainMessage(true, "未记录此用户id或者此用户不存在于本组")
			return
		}
		if !toolchain.GetTheTargetIsNormalUser(ctx) {
			ctx.SendPlainMessage(true, "暂时不支持匿名身份哦~")
			return
		}
		TargetInt := getID

		getProtectTargetStatus := CheckUserIsEnabledProtectMode(TargetInt, sdb)
		if getProtectTargetStatus {
			ctx.SendPlainMessage(true, "对方已经启动保护模式，不允许此处理操作")
			return
		}
		modifyCoins := int(toolchain.ExtractNumbers(info[2]))
		if TargetInt == ctx.Message.From.ID {
			ctx.SendPlainMessage(true, "哈? 干嘛骗自己的?坏蛋哦")
			return
		}
		switch {
		case modifyCoins <= 100:
			ctx.SendPlainMessage(true, "貌似你是想倒贴别人来着嘛?可以试试多骗一点哦，既然都骗了那就多点吧x")
			return
		case modifyCoins > 2000:
			ctx.SendPlainMessage(true, "不要太贪心了啦！太坏了 ")
			return
		}
		siEventUser := coins.GetSignInByUID(sdb, ctx.Message.From.ID) // 获取主用户目前状况信息
		siTargetUser := coins.GetSignInByUID(sdb, TargetInt)          // 获得被抢用户目前情况信息
		switch {
		case siTargetUser.Coins < modifyCoins:
			ctx.SendPlainMessage(true, "太坏了~试图的对象貌似没有足够多的柠檬片~")
			return
		case siEventUser.Coins < modifyCoins:
			ctx.SendPlainMessage(true, "貌似你需要有那么多数量的柠檬片哦w")
			return
		}
		eventUserName := toolchain.GetNickNameFromUsername(ctx.Message.From.UserName)
		eventTargetName := toolchain.GetNickNameFromUsername(getEntitiy[0])
		// get random numbers.
		getTargetChanceToDealRaw := math.Round(float64(modifyCoins / 20)) // the total is 0-100，however I don't allow getting chance 0. lmao. max is 100 if modify is 2000
		getTicket := RobOrCatchLimitManager(ctx.Message.From.ID)
		var remindTicket string
		if getTicket == 3 {
			remindTicket = "目前已经达到疲倦状态，成功率下调本身概率的15%，或许考虑一下不要做一个坏人呢～ ^^ "
		}
		getTargetChanceToDealPossibilityKey := rand.Intn(102 / getTicket)
		if getTargetChanceToDealPossibilityKey < int(getTargetChanceToDealRaw) { // failed
			doubledModifyNum := modifyCoins * 2
			if doubledModifyNum > siEventUser.Coins {
				doubledModifyNum = siEventUser.Coins
				_ = coins.InsertUserCoins(sdb, siEventUser.UID, siEventUser.Coins-doubledModifyNum)
				_ = coins.InsertUserCoins(sdb, siTargetUser.UID, siTargetUser.Coins+doubledModifyNum)
				ctx.SendPlainMessage(true, "试着去骗走 ", eventTargetName, " 的柠檬片时,被 ", eventTargetName, " 发现了.\n本该失去 ", modifyCoins*2, "\n但因为 ", eventUserName, " 的柠檬片过少，所以 ", eventUserName, " 失去了 ", doubledModifyNum, " 个柠檬片\n\n同时 ", eventTargetName, " 得到了 ", doubledModifyNum, " 个柠檬片\n", remindTicket)
				return
			}
			_ = coins.InsertUserCoins(sdb, siEventUser.UID, siEventUser.Coins-doubledModifyNum)
			_ = coins.InsertUserCoins(sdb, siTargetUser.UID, siTargetUser.Coins+doubledModifyNum)
			ctx.SendPlainMessage(true, "试着去骗走 ", eventTargetName, " 的柠檬片时,被 ", eventTargetName, " 发现了.\n所以 ", eventUserName, " 失去了 ", doubledModifyNum, " 个柠檬片\n\n同时 ", eventTargetName, " 得到了 ", doubledModifyNum, " 个柠檬片\n", remindTicket)
			return
		}
		_ = coins.InsertUserCoins(sdb, siEventUser.UID, siEventUser.Coins+modifyCoins)
		_ = coins.InsertUserCoins(sdb, siTargetUser.UID, siTargetUser.Coins-modifyCoins)
		ctx.SendPlainMessage(true, "试着去拿走 ", eventTargetName, " 的柠檬片时,成功了.\n所以 ", eventUserName, " 得到了 ", modifyCoins, " 个柠檬片\n\n同时 ", eventTargetName, " 失去了 ", modifyCoins, " 个柠檬片\n", remindTicket)

	})
	engine.OnMessageCommand("coinhand").SetBlock(true).Limit(ctxext.LimitByUser).Handle(func(ctx *rei.Ctx) {
		if !toolchain.GetTheTargetIsNormalUser(ctx) {
			ctx.SendPlainMessage(true, "暂时不支持匿名身份哦~")
			return
		}
		if !processLimit.Load(ctx.Message.From.ID).Acquire() {
			ctx.SendPlainMessage(true, "请等一会再转账哦w")
			return
		}
		_, info := toolchain.SplitCommandTo(ctx.Message.Text, 3)
		if len(info) < 3 {
			ctx.SendPlainMessage(true, "缺少参数~ 使用方法为 : /command [@user] (number)")
			return
		}
		getProtectStatus := CheckUserIsEnabledProtectMode(ctx.Message.From.ID, sdb)
		if getProtectStatus {
			ctx.SendPlainMessage(true, "已经启动保护模式，不允许参与任何抽奖性质类互动")
			return
		}
		getEntitiy := toolchain.ListEntitiesMention(ctx)
		if len(getEntitiy) == 0 {
			ctx.SendPlainMessage(true, "未指定用户~ 使用方法为 : /command [@user] (number)")
			return
		}
		getID := toolchain.GetUserIDFromUserName(ctx, getEntitiy[0])
		if getID == 0 {
			ctx.SendPlainMessage(true, "未记录此用户id或者此用户不存在于本组")
			return
		}
		TargetInt := getID

		getProtectTargetStatus := CheckUserIsEnabledProtectMode(TargetInt, sdb)
		if getProtectTargetStatus {
			ctx.SendPlainMessage(true, "对方已经启动保护模式，不允许此处理操作")
			return
		}
		modifyCoins := int(toolchain.ExtractNumbers(info[2]))
		if modifyCoins < 1 {
			ctx.SendPlainMessage(true, "然而你不能转账低于0个柠檬片哦w～ 敲")
			return
		}
		if TargetInt == ctx.Message.From.ID {
			ctx.SendPlainMessage(true, "不可以给自己转账哦w（敲）")
			return
		}
		uid := ctx.Message.From.ID
		siEventUser := coins.GetSignInByUID(sdb, uid)        // 获取主用户目前状况信息
		siTargetUser := coins.GetSignInByUID(sdb, TargetInt) // 获得被转账用户目前情况信息
		if modifyCoins > siEventUser.Coins {
			ctx.SendPlainMessage(true, "貌似你的柠檬片数量不够哦~")
			return
		}
		siEventUserName := toolchain.GetNickNameFromUsername(ctx.Message.From.UserName)
		siTargetUserName := toolchain.GetNickNameFromUsername(getEntitiy[0])
		ctx.SendPlainMessage(true, "转账成功了哦~\n", siEventUserName, " 变化为 ", siEventUser.Coins, " - ", modifyCoins, "= ", siEventUser.Coins-modifyCoins, "\n", siTargetUserName, " 变化为: ", siTargetUser.Coins, " + ", modifyCoins, "= ", siTargetUser.Coins+modifyCoins)
		_ = coins.InsertUserCoins(sdb, siEventUser.UID, siEventUser.Coins-modifyCoins)
		_ = coins.InsertUserCoins(sdb, siTargetUser.UID, siTargetUser.Coins+modifyCoins)
	})
	engine.OnMessageCommand("coinsuhand", rei.SuperUserPermission).SetBlock(true).Handle(func(ctx *rei.Ctx) {
		getEntitiy := toolchain.ListEntitiesMention(ctx)
		if len(getEntitiy) == 0 {
			ctx.SendPlainMessage(true, "未指定用户~ 使用方法为 : /command [@user] (number)")
			return
		}
		getID := toolchain.GetUserIDFromUserName(ctx, getEntitiy[0])
		if getID == 0 {
			ctx.SendPlainMessage(true, "未记录此用户id或者此用户不存在于本组")
			return
		}
		TargetInt := getID
		siTargetUser := coins.GetSignInByUID(sdb, TargetInt) // get user info
		unModifyCoins := siTargetUser.Coins
		_, modifyCoins := toolchain.SplitCommandTo(ctx.Message.Text, 3)
		coins.InsertUserCoins(sdb, TargetInt, unModifyCoins+int(toolchain.ExtractNumbers(modifyCoins[2])))
		ctx.SendPlainMessage(true, "Handle Coins Successfully.\n")

	})
	engine.OnMessageCommand("coinstaus").SetBlock(true).Handle(func(ctx *rei.Ctx) {
		getLength, getCodeRaw := toolchain.SplitCommandTo(ctx.Message.Text, 2)
		if getLength < 2 {
			ctx.SendPlainMessage(true, "缺少指令 , 应当为 /command (disable|禁用|enable|启用)")
			return
		}
		getCode := getCodeRaw[1]
		if getCode != "禁用" && getCode != "disable" && getCode != "enable" && getCode != "启用" {
			ctx.SendPlainMessage(true, "指令错误 , 应当为 /command (disable|禁用|enable|启用)")
			return
		}
		uid := ctx.Message.From.ID
		if getCode == "禁用" || getCode == "disable" {
			getStatus := CheckUserIsEnabledProtectMode(uid, sdb)
			if getStatus {
				ctx.SendPlainMessage(true, "你已经关闭了~")
				return
			}
			// start to handle
			suser := coins.GetProtectModeStatus(sdb, uid)
			boolStatus := suser.Time+60*60*24 < time.Now().Unix()
			if !boolStatus {
				ctx.SendPlainMessage(true, "仅允许24小时修改一次")
				return
			} // not the time
			// handle it.
			_ = coins.ChangeProtectStatus(sdb, uid, 1)
			ctx.SendPlainMessage(true, "修改完成~")
			return
		}
		getStatus := CheckUserIsEnabledProtectMode(uid, sdb)
		if !getStatus {
			ctx.SendPlainMessage(true, "你已经启用了~")
			return
		}
		// start to handle
		suser := coins.GetProtectModeStatus(sdb, uid)
		boolStatus := suser.Time+60*60*24 < time.Now().Unix()
		if !boolStatus {
			ctx.SendPlainMessage(true, "仅允许24小时修改一次")
			return
		} // not the time
		// handle it.
		_ = coins.ChangeProtectStatus(sdb, uid, 0)
		ctx.SendPlainMessage(true, "修改完成~")
	})
}

func RobOrCatchLimitManager(id int64) (ticket int) {
	// use limitManager to reduce the chance of true.
	// 33 * 4 + 6 * 5 + null key (4 time tired.)
	/*
		first time to get full chance to win.
		second time reduce it to 50 % chance to win
		last time is null , you are tired and reduce it to 33% chance to win.
	*/
	switch {
	case RobTimeManager.Load(id).AcquireN(33):
		return 1
	case RobTimeManager.Load(id).AcquireN(6):
		return 2
	case RobTimeManager.Load(id).Acquire():
		return 3
	default:
		return 3
	}
}

// CheckUserIsEnabledProtectMode 1 is enabled protect mode.
func CheckUserIsEnabledProtectMode(uid int64, sdb *coins.Scoredb) bool {
	s := coins.GetProtectModeStatus(sdb, uid)
	getCode := s.Status
	if getCode == 0 {
		return false
	}
	return true

}
