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
	payLimit       = rate.NewManager[int64](time.Hour*1, 10)  // time setup
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
		panic(err)
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
		} else {
			_ = coins.WagerCoinsInsert(sdb, getWager.Wagercount+modifyCoins, 0, getExpected)
			_ = coins.InsertUserCoins(sdb, userid, handleUser.Coins-modifyCoins)
			if rand.Intn(10) == 8 {
				ctx.SendPlainMessage(true, "呐～，不会还有大哥哥到现在 "+strconv.Itoa(getWager.Wagercount+modifyCoins)+" 个柠檬片了都没中奖吧？杂鱼～❤，杂鱼～❤")
			} else {
				ctx.SendPlainMessage(true, "没有中奖哦~，当前奖池为: ", getWager.Wagercount+modifyCoins)
			}
		}
	})
	engine.OnMessageCommand("coinfull").SetBlock(true).Limit(ctxext.LimitByUser).Handle(func(ctx *rei.Ctx) {
		uid, _ := toolchain.GetChatUserInfoID(ctx)
		si := coins.GetSignInByUID(sdb, uid)
		ctx.SendPlainMessage(true, "你的柠檬片数量一共是: "+strconv.Itoa(si.Coins))
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
			ctx.SendPlainMessage(true, "参数不足")
			return
		}

		// getUserInfo
		acquireEntity := toolchain.ListEntitiesMention(ctx)
		if len(acquireEntity) == 0 {
			return
		}
		getUserData := CoreFactory.GetUserSampleUserinfo(strings.Replace(acquireEntity[0], "@", "", 1))
		if getUserData.UserID == 0 {
			ctx.SendPlainMessage(true, "为非正常用户组或者该用户为在Lucy处记录x")
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
	} else {
		return true
	}
}
