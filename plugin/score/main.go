package score

import (
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"sync"
	"time"

	coins "github.com/MoYoez/Lucy_reibot/utils/coins"
	"github.com/MoYoez/Lucy_reibot/utils/toolchain"
	"github.com/MoYoez/Lucy_reibot/utils/transform"
	"github.com/FloatTech/floatbox/file"
	"github.com/FloatTech/gg"
	ctrl "github.com/FloatTech/zbpctrl"
	rei "github.com/fumiama/ReiBot"
	tgba "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	engine = rei.Register("score", &ctrl.Options[*rei.Ctx]{
		DisableOnDefault:  false,
		Help:              "Hi NekoPachi!\n说明书: https://lucy-sider.lemonkoi.one",
		PrivateDataFolder: "score",
	})
)

func init() {
	cachePath := engine.DataFolder() + "scorecache/"
	sdb := coins.Initialize("./data/score/score.db")
	engine.OnMessageCommand("sign").SetBlock(true).Handle(func(ctx *rei.Ctx) {
		var mutex sync.Mutex // 添加读写锁以保证稳定性
		uid, username := toolchain.GetChatUserInfoID(ctx)
		// remove emoji.
		emojiRegex := regexp.MustCompile(`[\x{1F600}-\x{1F64F}|[\x{1F300}-\x{1F5FF}]|[\x{1F680}-\x{1F6FF}]|[\x{1F700}-\x{1F77F}]|[\x{1F780}-\x{1F7FF}]|[\x{1F800}-\x{1F8FF}]|[\x{1F900}-\x{1F9FF}]|[\x{1FA00}-\x{1FA6F}]|[\x{1FA70}-\x{1FAFF}]|[\x{1FB00}-\x{1FBFF}]|[\x{1F170}-\x{1F251}]|[\x{1F300}-\x{1F5FF}]|[\x{1F600}-\x{1F64F}]|[\x{1FC00}-\x{1FCFF}]|[\x{1F004}-\x{1F0CF}]|[\x{1F170}-\x{1F251}]]+`)
		username = emojiRegex.ReplaceAllString(username, "")
		// save time data by add 30mins (database save it, not to handle it when it gets ready.)
		// just handle data time when it on,make sure to decrease 30 mins when render the page(

		// not sure what happened
		getNowUnixFormatElevenThirten := time.Now().Add(time.Minute * 30).Format("20060102")

		mutex.Lock()
		si := coins.GetSignInByUID(sdb, uid)
		mutex.Unlock()
		// in case
		drawedFile := cachePath + strconv.FormatInt(uid, 10) + getNowUnixFormatElevenThirten + "signin.png"
		if si.UpdatedAt.Add(time.Minute*30).Format("20060102") == getNowUnixFormatElevenThirten && si.Count != 0 {
			fmt.Print("DEBUGGER: " + si.UpdatedAt.Add(time.Minute*30).Format("20060102"))
			ctx.SendPlainMessage(true, "w~ 你今天已经签到过了哦w")
			if file.IsExist(drawedFile) {
				ctx.SendPhoto(tgba.FilePath(drawedFile), true, "~")
			}
			return
		}
		coinsGet := 300 + rand.Intn(200)
		mutex.Lock()

		_ = coins.InsertUserCoins(sdb, uid, si.Coins+coinsGet)
		_ = coins.InsertOrUpdateSignInCountByUID(sdb, uid, si.Count+1) // 柠檬片获取
		score := coins.GetScoreByUID(sdb, uid).Score
		score++ //  每日+1
		_ = coins.InsertOrUpdateScoreByUID(sdb, uid, score)
		CurrentCountTable := coins.GetCurrentCount(sdb, getNowUnixFormatElevenThirten)
		handledTodayNum := CurrentCountTable.Counttime + 1
		_ = coins.UpdateUserTime(sdb, handledTodayNum, getNowUnixFormatElevenThirten)
		mutex.Unlock()
		if time.Now().Hour() > 6 && time.Now().Hour() < 19 {
			// package for test draw.
			getTimeReplyMsg := coins.GetHourWord(time.Now()) // get time and msg
			currentTime := time.Now().Format("2006-01-02 15:04:05")
			// time day.
			dayTimeImg, _ := gg.LoadImage(transform.ReturnLucyMainDataIndex("score") + "BetaScoreDay.png")
			dayGround := gg.NewContext(1920, 1080)
			dayGround.DrawImage(dayTimeImg, 0, 0)
			_ = dayGround.LoadFontFace(transform.ReturnLucyMainDataIndex("score")+"dyh.ttf", 60)
			dayGround.SetRGB(0, 0, 0)
			// draw something with cautions Only (
			dayGround.DrawString(currentTime, 1270, 950)   // draw time
			dayGround.DrawString(getTimeReplyMsg, 50, 930) // draw text.
			dayGround.DrawString(username, 310, 110)       // draw name :p why I should do this???
			_ = dayGround.LoadFontFace(transform.ReturnLucyMainDataIndex("score")+"dyh.ttf", 60)
			dayGround.DrawStringWrapped(strconv.Itoa(handledTodayNum), 350, 255, 1, 1, 0, 1.3, gg.AlignCenter)   // draw first part
			dayGround.DrawStringWrapped(strconv.Itoa(si.Count+1), 1000, 255, 1, 1, 0, 1.3, gg.AlignCenter)       // draw second part
			dayGround.DrawStringWrapped(strconv.Itoa(coinsGet), 220, 370, 1, 1, 0, 1.3, gg.AlignCenter)          // draw third part
			dayGround.DrawStringWrapped(strconv.Itoa(si.Coins+coinsGet), 720, 370, 1, 1, 0, 1.3, gg.AlignCenter) // draw forth part
			// level array with rectangle work.
			rankNum := coins.GetLevel(score)
			RankGoal := rankNum + 1
			achieveNextGoal := coins.LevelArray[RankGoal]
			achievedGoal := coins.LevelArray[rankNum]
			currentNextGoalMeasure := achieveNextGoal - score  // measure rest of the num. like 20 - currentLink(TestRank 15)
			measureGoalsLens := achieveNextGoal - achievedGoal // like 20 - 10
			currentResult := float64(currentNextGoalMeasure) / float64(measureGoalsLens)
			// draw this part
			dayGround.SetRGB255(180, 255, 254)        // aqua color
			dayGround.DrawRectangle(70, 570, 600, 50) // draw rectangle part1
			dayGround.Fill()
			dayGround.SetRGB255(130, 255, 254)
			dayGround.DrawRectangle(70, 570, 600*currentResult, 50) // draw rectangle part2
			dayGround.Fill()
			dayGround.SetRGB255(0, 0, 0)
			dayGround.DrawString("Lv. "+strconv.Itoa(rankNum)+" 签到天数 + 1", 80, 490)
			_ = dayGround.LoadFontFace(transform.ReturnLucyMainDataIndex("score")+"dyh.ttf", 40)
			dayGround.DrawString(strconv.Itoa(currentNextGoalMeasure)+"/"+strconv.Itoa(measureGoalsLens), 710, 610)
			_ = dayGround.SavePNG(drawedFile)
			ctx.SendPhoto(tgba.FilePath(drawedFile), true, "[sign]签到完毕~")
		} else {
			// nightVision
			// package for test draw.
			getTimeReplyMsg := coins.GetHourWord(time.Now()) // get time and msg
			currentTime := time.Now().Format("2006-01-02 15:04:05")
			nightTimeImg, _ := gg.LoadImage(transform.ReturnLucyMainDataIndex("score") + "BetaScoreNight.png")
			nightGround := gg.NewContext(1886, 1060)
			nightGround.DrawImage(nightTimeImg, 0, 0)
			_ = nightGround.LoadFontFace(transform.ReturnLucyMainDataIndex("score")+"dyh.ttf", 60)
			nightGround.SetRGB255(255, 255, 255)
			// draw something with cautions Only (
			nightGround.DrawString(currentTime, 1360, 910)   // draw time
			nightGround.DrawString(getTimeReplyMsg, 60, 930) // draw text.
			nightGround.DrawString(username, 350, 140)       // draw name :p why I should do this???
			_ = nightGround.LoadFontFace(transform.ReturnLucyMainDataIndex("score")+"dyh.ttf", 60)
			nightGround.DrawStringWrapped(strconv.Itoa(handledTodayNum), 345, 275, 1, 1, 0, 1.3, gg.AlignCenter)   // draw first part
			nightGround.DrawStringWrapped(strconv.Itoa(si.Count+1), 990, 275, 1, 1, 0, 1.3, gg.AlignCenter)        // draw second part
			nightGround.DrawStringWrapped(strconv.Itoa(coinsGet), 225, 360, 1, 1, 0, 1.3, gg.AlignCenter)          // draw third part
			nightGround.DrawStringWrapped(strconv.Itoa(si.Coins+coinsGet), 720, 360, 1, 1, 0, 1.3, gg.AlignCenter) // draw forth part
			// level array with rectangle work.
			rankNum := coins.GetLevel(score)
			RankGoal := rankNum + 1
			achieveNextGoal := coins.LevelArray[RankGoal]
			achievedGoal := coins.LevelArray[rankNum]
			currentNextGoalMeasure := achieveNextGoal - score  // measure rest of the num. like 20 - currentLink(TestRank 15)
			measureGoalsLens := achieveNextGoal - achievedGoal // like 20 - 10
			currentResult := float64(currentNextGoalMeasure) / float64(measureGoalsLens)
			// draw this part
			nightGround.SetRGB255(49, 86, 157)          // aqua color
			nightGround.DrawRectangle(70, 570, 600, 50) // draw rectangle part1
			nightGround.Fill()
			nightGround.SetRGB255(255, 255, 255)
			nightGround.DrawRectangle(70, 570, 600*currentResult, 50) // draw rectangle part2
			nightGround.Fill()
			nightGround.SetRGB255(255, 255, 255)
			nightGround.DrawString("Lv. "+strconv.Itoa(rankNum)+" 签到天数 + 1", 80, 490)
			_ = nightGround.LoadFontFace(transform.ReturnLucyMainDataIndex("score")+"dyh.ttf", 40)
			nightGround.DrawString(strconv.Itoa(currentNextGoalMeasure)+"/"+strconv.Itoa(measureGoalsLens), 710, 610)
			_ = nightGround.SavePNG(drawedFile)
			ctx.SendPhoto(tgba.FilePath(drawedFile), true, "[sign]签到完成~")
		}
	})

}
