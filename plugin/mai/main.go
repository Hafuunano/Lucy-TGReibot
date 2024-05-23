package mai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/FloatTech/floatbox/web"
	"github.com/FloatTech/gg"
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/MoYoez/Lucy_reibot/utils/toolchain"
	rei "github.com/fumiama/ReiBot"
	tgba "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"image"
	"os"
	"strconv"
	"strings"
)

var engine = rei.Register("mai", &ctrl.Options[*rei.Ctx]{
	DisableOnDefault:  false,
	Help:              "maimai - bind Username / maimai b50 render",
	PrivateDataFolder: "mai",
})

func init() {
	engine.OnMessageCommand("mai").SetBlock(true).Handle(func(ctx *rei.Ctx) {
		getMsg := ctx.Message.Text
		getSplitLength, getSplitStringList := toolchain.SplitCommandTo(getMsg, 3)
		if getSplitLength >= 2 {
			switch {
			case getSplitStringList[1] == "bind":
				if getSplitLength < 3 {
					ctx.SendPlainMessage(true, "参数提供不足")
					return
				}
				BindUserToMaimai(ctx, getSplitStringList[2])
			case getSplitStringList[1] == "lxbind":
				if getSplitLength < 3 {
					ctx.SendPlainMessage(true, "参数提供不足")
					return
				}
				Toint64, err := strconv.ParseInt(getSplitStringList[2], 10, 64)
				if err != nil {
					ctx.SendPlainMessage(true, "参数的FriendCode为非法")
					return
				}
				BindFriendCode(ctx, Toint64)
			case getSplitStringList[1] == "plate":
				if getSplitLength == 2 {
					SetUserPlateToLocal(ctx, "")
					return
				}
				SetUserPlateToLocal(ctx, getSplitStringList[2])
			case getSplitStringList[1] == "upload":
				// uploadImage
				images := toolchain.RequestImageTo(ctx, "请发送指令同时提供一张图片，图片大小比例适应为6:1 (1260x210) ,如果图片不适应将会自动剪辑到合适大小")
				if images == nil {
					return
				}
				HandlerUserSetsCustomImage(ctx, images)
			case getSplitStringList[1] == "remove":
				RemoveUserLocalCustomImage(ctx)
			case getSplitStringList[1] == "defplate":
				if getSplitLength < 3 {
					SetUserDefaultPlateToDatabase(ctx, "")
					return
				}
				SetUserDefaultPlateToDatabase(ctx, getSplitStringList[2])
			case getSplitStringList[1] == "switch":
				MaimaiSwitcherService(ctx)
			case getSplitStringList[1] == "raw" || getSplitStringList[1] == "file":
				MaimaiRenderBase(ctx, true)
			case getSplitStringList[1] == "query":
				if getSplitLength < 3 {
					ctx.SendPlainMessage(true, "参数提供不足, /mai query [绿黄红紫白][dx|标] <SongAlias> ")
					return
				}
				// CASE: if User Trigger This command, check other settings.
				// getQuery:
				// level_index | song_type
				getLength, getSplitInfo := toolchain.SplitCommandTo(getSplitStringList[2], 2)
				userSettingInterface := map[string]string{}
				var settedSongAlias string
				if getLength > 1 { // prefix judge.
					settedSongAlias = getSplitInfo[1]
					for i, returnLevelValue := range []string{"绿", "黄", "红", "紫", "白"} {
						if strings.Contains(getSplitInfo[0], returnLevelValue) {
							userSettingInterface["level_index"] = strconv.Itoa(i)
							break
						}
					}
					switch {
					case strings.Contains(getSplitInfo[0], "dx"):
						userSettingInterface["song_type"] = "dx"
					case strings.Contains(getSplitInfo[0], "标"):
						userSettingInterface["song_type"] = "standard"
					}
				} else {
					// no other infos. || default setting ==> dx Master | std Master | dx expert | std expert (as the highest score)
					settedSongAlias = getSplitInfo[0]
				}
				// get SongID, render.
				getUserID, _ := toolchain.GetChatUserInfoID(ctx)
				// check the user is Lxns Service | DivingFish Service.
				getBool := GetUserSwitcherInfoFromDatabase(getUserID)
				var isIDChecker bool
				var songIDList []int
				// first read the config.
				getLevelIndex := userSettingInterface["level_index"]
				getSongType := userSettingInterface["song_type"]
				var getReferIndexIsOn bool
				var accStat bool
				if getLevelIndex != "" { // use custom diff
					getReferIndexIsOn = true
				}
				switch {
				case strings.HasPrefix(settedSongAlias, "id"):
					// useID checker.
					isIDChecker = true
					getParse, err := strconv.ParseInt(strings.Replace(settedSongAlias, "id", "", 1), 10, 64)
					if err != nil {
						ctx.SendPlainMessage(true, "ID 查找参数非法")
						return
					}
					songIDList = []int{int(getParse)}
				default:
					isIDChecker = false
					queryStatus, songIDLists, accStats, returnListHere := QueryReferSong(settedSongAlias, getBool)
					songIDList = songIDLists
					accStat = accStats
					if !queryStatus {
						ctx.SendPlainMessage(true, "未找到对应歌曲，可能是数据库未收录（")
						return
					}
					if accStat {
						// Handler Which Song user played.
						var FullList []int
						for _, list := range returnListHere {
							for _, listInsider := range list {
								FullList = append(FullList, listInsider)
							}
						}

						FullList = removeIntDuplicates(FullList)
						// make both song Handler, check this song is from sd | DX pattern.
						var sampleListShown []int
						// sometimes list maybe contain dx | SD, but they are same song.
						if len(FullList) == 2 {
							for _, listSample := range FullList {
								sampleListShown = append(sampleListShown, simpleNumHandler(listSample, false)) //  convert to DX pattern.
							}
							sampleListShown = removeIntDuplicates(sampleListShown)
						}

						if len(sampleListShown) == 1 {
							songIDList = []int{simpleNumHandler(songIDList[0], false)}
						} else {
							// varies handler machine,means it has songs.
							// query them.
							songIDList = FullList
						}

					}
				}

				if getBool { // lxns service.
					getFriendID := GetUserMaiFriendID(getUserID)
					if getFriendID.MaimaiID == 0 {
						ctx.SendPlainMessage(true, "没有绑定哦～ 使用/mai lxbind <friendcode> 以绑定")
						return
					}
					// to convert it, and it can be read by Lxns.
					if !getReferIndexIsOn { // no refer then return the last one.
						var getReport LxnsMaimaiRequestUserReferBestSong
						switch {
						case getSongType == "standard":
							for _, songIdInt := range songIDList {
								getReport = RequestReferSong(getFriendID.MaimaiID, int64(songIdInt), true)
								if getReport.Code == 200 && len(getReport.Data) != 0 {
									break
								}
							}
							if len(getReport.Data) == 0 {
								ctx.SendPlainMessage(true, "没有发现 SD 谱面～ 如不确定可以忽略请求参数, Lucy会自动识别")
								return
							}
						case getSongType == "dx":
							for _, songIdInt := range songIDList {
								getReport = RequestReferSong(getFriendID.MaimaiID, int64(songIdInt), false)
								if getReport.Code == 200 && len(getReport.Data) != 0 {
									break
								}
							}
							if len(getReport.Data) == 0 {
								ctx.SendPlainMessage(true, "没有发现 DX 谱面～ 如不确定可以忽略请求参数, Lucy会自动识别")
								return
							}
						default:
							for _, songIdInt := range songIDList {
								getReport = RequestReferSong(getFriendID.MaimaiID, int64(songIdInt), false)
								fmt.Print(getReport.Code)
								if getReport.Code == 200 && len(getReport.Data) != 0 {
									break
								}

							}

							if getReport.Code != 200 {
								for _, songIdInt := range songIDList {
									getReport = RequestReferSong(getFriendID.MaimaiID, int64(songIdInt), true)
									if getReport.Code == 200 && len(getReport.Data) != 0 {
										break
									}
								}
							}
						}
						getReturnTypeLength := len(getReport.Data)
						if getReturnTypeLength == 0 {
							if !isIDChecker {
								ctx.SendPlainMessage(true, "Lucy 似乎没有查询到你的游玩数据呢（")
							} else {
								ctx.SendPlainMessage(true, "Lucy 查找了对应ID 但是没有发现数据～")
							}
							return
						}
						// DataGet, convert To MaiPlayData Render.
						maiRenderPieces := LxnsMaimaiRequestDataPiece{
							Id:           getReport.Data[len(getReport.Data)-1].Id,
							SongName:     getReport.Data[len(getReport.Data)-1].SongName,
							Level:        getReport.Data[len(getReport.Data)-1].Level,
							LevelIndex:   getReport.Data[len(getReport.Data)-1].LevelIndex,
							Achievements: getReport.Data[len(getReport.Data)-1].Achievements,
							Fc:           getReport.Data[len(getReport.Data)-1].Fc,
							Fs:           getReport.Data[len(getReport.Data)-1].Fs,
							DxScore:      getReport.Data[len(getReport.Data)-1].DxScore,
							DxRating:     getReport.Data[len(getReport.Data)-1].DxRating,
							Rate:         getReport.Data[len(getReport.Data)-1].Rate,
							Type:         getReport.Data[len(getReport.Data)-1].Type,
							UploadTime:   getReport.Data[len(getReport.Data)-1].UploadTime,
						}
						getFinalPic := ReCardRenderBase(maiRenderPieces, 0, true)
						_ = gg.NewContextForImage(getFinalPic).SavePNG(engine.DataFolder() + "save/" + "LXNS_PIC_" + strconv.Itoa(songIDList[0]) + "_" + strconv.Itoa(int(getUserID)) + ".png")
						if accStat {
							ctx.SendPhoto(tgba.FilePath(engine.DataFolder()+"save/"+"LXNS_PIC_"+strconv.Itoa(songIDList[0])+"_"+strconv.Itoa(int(getUserID))+".png"), true, "Lucy 查询到多个别名，此处默认为您返回了 "+getReport.Data[0].SongName+"谱面")
						} else {
							ctx.SendPhoto(tgba.FilePath(engine.DataFolder()+"save/"+"LXNS_PIC_"+strconv.Itoa(songIDList[0])+"_"+strconv.Itoa(int(getUserID))+".png"), true, "")
						}
					} else {
						var getReport LxnsMaimaiRequestUserReferBestSongIndex
						getLevelIndexToint, _ := strconv.ParseInt(getLevelIndex, 10, 64)
						switch {
						case getSongType == "standard":
							for _, p := range songIDList {
								getReport = RequestReferSongIndex(getFriendID.MaimaiID, int64(p), getLevelIndexToint, true)
								if getReport.Code == 200 && getReport.Data.SongName != "" {
									break
								}
							}
							if getReport.Code == 404 {
								ctx.SendPlainMessage(true, "没有发现 SD 谱面～ 如不确定可以忽略请求参数, Lucy会自动识别")
								return
							}
						case getSongType == "dx":
							for _, p := range songIDList {
								getReport = RequestReferSongIndex(getFriendID.MaimaiID, int64(p), getLevelIndexToint, false)
								if getReport.Code == 200 && getReport.Data.SongName != "" {
									break
								}
							}
							if getReport.Code == 404 {
								ctx.SendPlainMessage(true, "没有发现 DX 谱面～ 如不确定可以忽略请求参数, Lucy会自动识别")
								return
							}
						default:
							for _, p := range songIDList {
								getReport = RequestReferSongIndex(getFriendID.MaimaiID, int64(p), getLevelIndexToint, false)
								if getReport.Code == 200 && getReport.Data.SongName != "" {
									break
								}
							}
							if getReport.Code != 200 {
								for _, p := range songIDList {
									getReport = RequestReferSongIndex(getFriendID.MaimaiID, int64(p), getLevelIndexToint, true)
									if getReport.Code == 200 && getReport.Data.SongName != "" {
										break
									}
								}
							}
						}
						if getReport.Data.SongName == "" { // nil pointer.
							if !isIDChecker {
								ctx.SendPlainMessage(true, "Lucy 似乎没有查询到你的游玩数据呢（")
							} else {
								ctx.SendPlainMessage(true, "Lucy 查找了对应ID 但是没有发现数据～")
							}
						}
						maiRenderPieces := LxnsMaimaiRequestDataPiece{
							Id:           getReport.Data.Id,
							SongName:     getReport.Data.SongName,
							Level:        getReport.Data.Level,
							LevelIndex:   getReport.Data.LevelIndex,
							Achievements: getReport.Data.Achievements,
							Fc:           getReport.Data.Fc,
							Fs:           getReport.Data.Fs,
							DxScore:      getReport.Data.DxScore,
							DxRating:     getReport.Data.DxRating,
							Rate:         getReport.Data.Rate,
							Type:         getReport.Data.Type,
							UploadTime:   getReport.Data.UploadTime,
						}
						getFinalPic := ReCardRenderBase(maiRenderPieces, 0, true)
						_ = gg.NewContextForImage(getFinalPic).SavePNG(engine.DataFolder() + "save/" + "LXNS_PIC_" + strconv.Itoa(songIDList[0]) + "_" + strconv.Itoa(int(getUserID)) + ".png")
						if accStat {
							ctx.SendPhoto(tgba.FilePath(engine.DataFolder()+"save/"+"LXNS_PIC_"+strconv.Itoa(songIDList[0])+"_"+strconv.Itoa(int(getUserID))+".png"), true, "Lucy 查询到多个别名，此处默认为您返回了 "+getReport.Data.SongName+"谱面")
						} else {
							ctx.SendPhoto(tgba.FilePath(engine.DataFolder()+"save/"+"LXNS_PIC_"+strconv.Itoa(songIDList[0])+"_"+strconv.Itoa(int(getUserID))+".png"), true, "")
						}
					}

				} else {
					// diving fish checker, data rewrite.
					getUsername := GetUserInfoNameFromDatabase(getUserID)
					if getUsername == "" {
						ctx.SendPlainMessage(true, "你还没有绑定呢！使用/mai bind <UserName> 以绑定")
						return
					}
					fullDevData := QueryDevDataFromDivingFish(getUsername)
					// default setting ==> dx Master | std Master | dx expert | std expert (as the highest score)
					var ReferSongTypeList []int
					switch {
					case getSongType == "standard":
						// roll songIDList first.
						for _, songID := range songIDList {
							if !isIDChecker {
								songID = simpleNumHandler(songID, false)
							}
							for numPosition, index := range fullDevData.Records {
								if index.SongId == songID && index.Type == "SD" {
									ReferSongTypeList = append(ReferSongTypeList, numPosition)
								}
							}
							if len(ReferSongTypeList) != 0 {
								break
							}

						}

						if len(ReferSongTypeList) == 0 {
							ctx.SendPlainMessage(true, "没有发现游玩过的 SD 谱面～ 如不确定可以忽略请求参数, Lucy会自动识别")
							return
						}
					case getSongType == "dx":
						for _, songID := range songIDList {
							if !isIDChecker {
								songID = simpleNumHandler(songID, true)
							}
							for numPosition, index := range fullDevData.Records {
								if index.SongId == songID && index.Type == "DX" {
									ReferSongTypeList = append(ReferSongTypeList, numPosition)
								}
							}
							if len(ReferSongTypeList) != 0 {
								break
							}
						}

						if len(ReferSongTypeList) == 0 {
							ctx.SendPlainMessage(true, "没有发现游玩过的 DX 谱面～ 如不确定可以忽略请求参数, Lucy会自动识别")
							return
						}
					default: // no settings.
						for _, songID := range songIDList {
							if !isIDChecker {
								songID = simpleNumHandler(songID, true)
							}
							for numPosition, index := range fullDevData.Records {
								if index.SongId == songID && index.Type == "DX" {
									ReferSongTypeList = append(ReferSongTypeList, numPosition)
								}
							}
							if len(ReferSongTypeList) != 0 {
								break
							}
						}
						if len(ReferSongTypeList) == 0 {
							for _, songID := range songIDList {
								if !isIDChecker {
									songID = simpleNumHandler(songID, false)
								}
								for numPosition, index := range fullDevData.Records {
									if index.SongId == songID && index.Type == "SD" {
										ReferSongTypeList = append(ReferSongTypeList, numPosition)
									}
								}
								if len(ReferSongTypeList) != 0 {
									break
								}
							}
						}
						if len(ReferSongTypeList) == 0 {
							if !isIDChecker {
								ctx.SendPlainMessage(true, "Lucy 似乎没有查询到你的游玩数据呢（")
							} else {
								ctx.SendPlainMessage(true, "Lucy 查找了对应ID 但是没有发现数据～")
							}
							return
						}
					}

					if !getReferIndexIsOn {
						// index a map =>  level_index = "record_diff"
						levelIndexMap := map[int]string{}
						for _, dataPack := range ReferSongTypeList {
							levelIndexMap[fullDevData.Records[dataPack].LevelIndex] = strconv.Itoa(dataPack)
						}
						var trulyReturnedData string
						for i := 4; i >= 0; i-- { // divingfish is reverse.
							if levelIndexMap[i] != "" {
								trulyReturnedData = levelIndexMap[i]
								break
							}
						}
						getNum, _ := strconv.Atoi(trulyReturnedData)
						// getNum ==> 0
						returnPackage := fullDevData.Records[getNum]
						_ = gg.NewContextForImage(RenderCard(returnPackage, 0, true)).SavePNG(engine.DataFolder() + "save/" + strconv.Itoa(songIDList[0]) + "_" + strconv.Itoa(int(getUserID)) + ".png")
						if accStat {
							ctx.SendPhoto(tgba.FilePath(engine.DataFolder()+"save/"+strconv.Itoa(int(songIDList[0]))+"_"+strconv.Itoa(int(getUserID))+".png"), true, "Lucy 查询到多个别名，此处默认为您返回了 "+returnPackage.Title+" 谱面")
						} else {
							ctx.SendPhoto(tgba.FilePath(engine.DataFolder()+"save/"+strconv.Itoa(int(songIDList[0]))+"_"+strconv.Itoa(int(getUserID))+".png"), true, "")
						}
					} else {
						levelIndexMap := map[int]string{}
						for _, dataPack := range ReferSongTypeList {
							levelIndexMap[fullDevData.Records[dataPack].LevelIndex] = strconv.Itoa(dataPack)
						}
						getDiff, _ := strconv.Atoi(userSettingInterface["level_index"])

						if levelIndexMap[getDiff] != "" {
							getNum, _ := strconv.Atoi(levelIndexMap[getDiff])
							returnPackage := fullDevData.Records[getNum]
							_ = gg.NewContextForImage(RenderCard(returnPackage, 0, true)).SavePNG(engine.DataFolder() + "save/" + strconv.Itoa(songIDList[0]) + "_" + strconv.Itoa(int(getUserID)) + ".png")
							if accStat {
								ctx.SendPhoto(tgba.FilePath(engine.DataFolder()+"save/"+strconv.Itoa(songIDList[0])+"_"+strconv.Itoa(int(getUserID))+".png"), true, "Lucy 查询到多个别名，此处默认为您返回了 "+returnPackage.Title+" 谱面")
							} else {
								ctx.SendPhoto(tgba.FilePath(engine.DataFolder()+"save/"+strconv.Itoa(songIDList[0])+"_"+strconv.Itoa(int(getUserID))+".png"), true, "")
							}
						} else {
							if !isIDChecker {
								ctx.SendPlainMessage(true, "Lucy 貌似你没有玩过这个难度的曲子哦～")
							} else {
								ctx.SendPlainMessage(true, "Lucy 查找了对应ID 在这个难度下没有发现数据～")
							}

						}
					}
				}
			case getSplitStringList[1] == "aliasupdate":
				if rei.SuperUserPermission(ctx) {
					UpdateAliasPackage()
					ctx.SendPlainMessage(true, "更新成功～")
				} else {
					ctx.SendPlainMessage(true, "您似乎没有权限呢(")
				}
			default:
				ctx.SendPlainMessage(true, "未知的指令或者指令出现错误~")
			}
		} else {
			MaimaiRenderBase(ctx, false)
		}
	})
}

// BindFriendCode Bind FriendCode To Users
func BindFriendCode(ctx *rei.Ctx, bindCode int64) {
	getUserID, _ := toolchain.GetChatUserInfoID(ctx)
	FormatMaimaiFriendCode(bindCode, getUserID).BindUserFriendCode()
	ctx.SendPlainMessage(true, "绑定成功~！")
}

// BindUserToMaimai Bind UserMaiMaiID
func BindUserToMaimai(ctx *rei.Ctx, bindName string) {
	getUserID, _ := toolchain.GetChatUserInfoID(ctx)
	FormatUserDataBase(getUserID, GetUserPlateInfoFromDatabase(getUserID), GetUserDefaultBackgroundDataFromDatabase(getUserID), bindName).BindUserDataBase()
	ctx.SendPlainMessage(true, "绑定成功~！")
}

// SetUserPlateToLocal Set Default Plate to Local
func SetUserPlateToLocal(ctx *rei.Ctx, plateID string) {
	getUserID, _ := toolchain.GetChatUserInfoID(ctx)
	FormatUserDataBase(getUserID, plateID, GetUserDefaultBackgroundDataFromDatabase(getUserID), GetUserInfoNameFromDatabase(getUserID)).BindUserDataBase()
	ctx.SendPlainMessage(true, "好哦~ 是个好名称w")
}

// HandlerUserSetsCustomImage  Handle User Custom Image and Send To Local
func HandlerUserSetsCustomImage(ctx *rei.Ctx, ps []tgba.PhotoSize) {
	getUserID, _ := toolchain.GetChatUserInfoID(ctx)
	pic := ps[len(ps)-1]
	picu, _ := ctx.Caller.GetFileDirectURL(pic.FileID)
	imageData, err := web.GetData(picu)
	if err != nil {
		return
	}
	getRaw, _, err := image.Decode(bytes.NewReader(imageData))
	if err != nil {
		return
	}
	// pic Handler
	getRenderPlatePicRaw := gg.NewContext(1260, 210)
	getRenderPlatePicRaw.DrawRoundedRectangle(0, 0, 1260, 210, 10)
	getRenderPlatePicRaw.Clip()
	getHeight := getRaw.Bounds().Dy()
	getLength := getRaw.Bounds().Dx()
	var getHeightHandler, getLengthHandler int
	switch {
	case getHeight < 210 && getLength < 1260:
		getRaw = Resize(getRaw, 1260, 210)
		getHeightHandler = 0
		getLengthHandler = 0
	case getHeight < 210:
		getRaw = Resize(getRaw, getLength, 210)
		getHeightHandler = 0
		getLengthHandler = (getRaw.Bounds().Dx() - 1260) / 3 * -1
	case getLength < 1260:
		getRaw = Resize(getRaw, 1260, getHeight)
		getHeightHandler = (getRaw.Bounds().Dy() - 210) / 3 * -1
		getLengthHandler = 0
	default:
		getLengthHandler = (getRaw.Bounds().Dx() - 1260) / 3 * -1
		getHeightHandler = (getRaw.Bounds().Dy() - 210) / 3 * -1
	}
	getRenderPlatePicRaw.DrawImage(getRaw, getLengthHandler, getHeightHandler)
	getRenderPlatePicRaw.Fill()
	// save.
	_ = getRenderPlatePicRaw.SavePNG(userPlate + strconv.Itoa(int(getUserID)) + ".png")
	ctx.SendPlainMessage(true, "已经存入了哦w~")
}

// RemoveUserLocalCustomImage Remove User Local Image.
func RemoveUserLocalCustomImage(ctx *rei.Ctx) {
	getUserID, _ := toolchain.GetChatUserInfoID(ctx)
	_ = os.Remove(userPlate + strconv.Itoa(int(getUserID)) + ".png")
	ctx.SendPlainMessage(true, "已经移除了~ ")
}

// SetUserDefaultPlateToDatabase Set Default plateID To Database.
func SetUserDefaultPlateToDatabase(ctx *rei.Ctx, plateName string) {
	getUserID, _ := toolchain.GetChatUserInfoID(ctx)
	if plateName == "" {
		FormatUserDataBase(getUserID, GetUserPlateInfoFromDatabase(getUserID), "", GetUserInfoNameFromDatabase(getUserID)).BindUserDataBase()
		ctx.SendPlainMessage(true, "已经删除了预设~")
		return
	}
	getDefaultInfo := plateName
	_, err := GetDefaultPlate(getDefaultInfo)
	if err != nil {
		ctx.SendPlainMessage(true, "设定的预设不正确")
		return
	}
	FormatUserDataBase(getUserID, GetUserPlateInfoFromDatabase(getUserID), getDefaultInfo, GetUserInfoNameFromDatabase(getUserID)).BindUserDataBase()
	ctx.SendPlainMessage(true, "已经设定好了哦w~ ")
}

// MaimaiRenderBase Render Base Maimai B50.
func MaimaiRenderBase(ctx *rei.Ctx, israw bool) {
	// check the user using.
	getUserID, _ := toolchain.GetChatUserInfoID(ctx)
	if GetUserSwitcherInfoFromDatabase(getUserID) {
		// use lxns checker service.
		// check bind first, get user friend id.
		getFriendID := GetUserMaiFriendID(getUserID)
		if getFriendID.MaimaiID == 0 {
			ctx.SendPlainMessage(true, "你还没有绑定呢！使用/mai lxbind <friendcode> 以绑定")
			return
		}
		getUserData := RequestBasicDataFromLxns(getFriendID.MaimaiID)
		if getUserData.Code != 200 {
			ctx.SendPlainMessage(true, "aw 出现了一点小错误~：\n - 请检查你是否有上传过数据\n - 请检查你的设置是否允许了第三方查看")
			return
		}
		getGameUserData := RequestB50DataByFriendCode(getUserData.Data.FriendCode)
		if getGameUserData.Code != 200 {
			ctx.SendPlainMessage(true, "aw 出现了一点小错误~：\n - 请检查你是否有上传过数据\n - 请检查你的设置是否允许了第三方查看")
			return
		}
		getImager, _ := ReFullPageRender(getGameUserData, getUserData, ctx)
		_ = gg.NewContextForImage(getImager).SavePNG(engine.DataFolder() + "save/" + "LXNS_" + strconv.Itoa(int(getUserID)) + ".png")
		if israw {
			getDocumentType := &tgba.DocumentConfig{
				BaseFile: tgba.BaseFile{BaseChat: tgba.BaseChat{
					ChatConfig: tgba.ChatConfig{ChatID: ctx.Message.Chat.ID},
				},
					File: tgba.FilePath(engine.DataFolder() + "save/" + "LXNS_" + strconv.Itoa(int(getUserID)) + ".png")},
				Caption:         "",
				CaptionEntities: nil,
			}
			ctx.Send(true, getDocumentType)
		} else {
			ctx.SendPhoto(tgba.FilePath(engine.DataFolder()+"save/"+"LXNS_"+strconv.Itoa(int(getUserID))+".png"), true, "")
		}
	} else {
		// diving fish checker:
		getUsername := GetUserInfoNameFromDatabase(getUserID)
		if getUsername == "" {
			ctx.SendPlainMessage(true, "你还没有绑定呢！使用/mai bind UserName 以绑定")
			return
		}
		getUserData, err := QueryMaiBotDataFromUserName(getUsername)
		if err != nil {
			ctx.SendPlainMessage(true, err)
			return
		}
		var data player
		_ = json.Unmarshal(getUserData, &data)
		renderImg := FullPageRender(data, ctx)
		_ = gg.NewContextForImage(renderImg).SavePNG(engine.DataFolder() + "save/" + strconv.Itoa(int(getUserID)) + ".png")

		if israw {
			getDocumentType := &tgba.DocumentConfig{
				BaseFile: tgba.BaseFile{BaseChat: tgba.BaseChat{
					ChatConfig: tgba.ChatConfig{ChatID: ctx.Message.Chat.ID},
				},
					File: tgba.FilePath(engine.DataFolder() + "save/" + strconv.Itoa(int(getUserID)) + ".png")},
				Caption:         "",
				CaptionEntities: nil,
			}
			ctx.Send(true, getDocumentType)
		} else {
			ctx.SendPhoto(tgba.FilePath(engine.DataFolder()+"save/"+strconv.Itoa(int(getUserID))+".png"), true, "")
		}
	}
}

// MaimaiSwitcherService True == Lxns Service || False == Diving Fish Service.
func MaimaiSwitcherService(ctx *rei.Ctx) {
	getUserID, _ := toolchain.GetChatUserInfoID(ctx)
	getBool := GetUserSwitcherInfoFromDatabase(getUserID)
	err := FormatUserSwitcher(getUserID, !getBool).ChangeUserSwitchInfoFromDataBase()
	if err != nil {
		panic(err)
	}
	var getEventText string
	// due to it changed, so reverse.
	if !getBool {
		getEventText = "Lxns查分"
	} else {
		getEventText = "Diving Fish查分"
	}
	ctx.SendPlainMessage(true, "已经修改为"+getEventText)
}

func simpleNumHandler(num int, upper bool) int {
	if upper {
		if num < 1000 && num > 100 {
			toint, _ := strconv.Atoi(fmt.Sprintf("10%d", num))
			return toint
		}
		if num > 1000 && num < 10000 {
			toint, _ := strconv.Atoi(fmt.Sprintf("1%d", num))
			return toint
		}
	} else {
		getFmt := fmt.Sprintf("%d", num)
		getFmt = getFmt[2:]
		toint, _ := strconv.Atoi(getFmt)
		return toint
	}
	return num
}
