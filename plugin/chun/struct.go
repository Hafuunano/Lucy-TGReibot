package chun

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/FloatTech/gg"
	"github.com/FloatTech/imgfactory"
	"github.com/MoYoez/Lucy_reibot/utils/toolchain"
	rei "github.com/fumiama/ReiBot"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/text/width"
	"image"
	"image/color"
	"image/png"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"unicode/utf8"
)

// ChunData Struct.
type ChunData struct {
	Nickname string  `json:"nickname"`
	Rating   float64 `json:"rating"`
	Records  struct {
		B30 []struct {
			Cid        int     `json:"cid"`
			Ds         float64 `json:"ds"`
			Fc         string  `json:"fc"`
			Level      string  `json:"level"`
			LevelIndex int     `json:"level_index"`
			LevelLabel string  `json:"level_label"`
			Mid        int     `json:"mid"`
			Ra         float64 `json:"ra"`
			Score      int     `json:"score"`
			Title      string  `json:"title"`
		} `json:"b30"`
		R10 []struct {
			Cid        int     `json:"cid"`
			Ds         float64 `json:"ds"`
			Fc         string  `json:"fc"`
			Level      string  `json:"level"`
			LevelIndex int     `json:"level_index"`
			LevelLabel string  `json:"level_label"`
			Mid        int     `json:"mid"`
			Ra         float64 `json:"ra"`
			Score      int     `json:"score"`
			Title      string  `json:"title"`
		} `json:"r10"`
	} `json:"records"`
	Username string `json:"username"`
}

// UserDataInner CardBase
type UserDataInner []struct {
	Cid        int     `json:"cid"`
	Ds         float64 `json:"ds"`
	Fc         string  `json:"fc"`
	Level      string  `json:"level"`
	LevelIndex int     `json:"level_index"`
	LevelLabel string  `json:"level_label"`
	Mid        int     `json:"mid"`
	Ra         float64 `json:"ra"`
	Score      int     `json:"score"`
	Title      string  `json:"title"`
}

type DivingFishB50 struct {
	Username string `json:"username"`
	B50      bool   `json:"b50"`
}

var (
	Root    = engine.DataFolder()
	Texture = engine.DataFolder() + "texture/"
)

// QueryChunDataFromUserName Query Chun Data.
func QueryChunDataFromUserName(userName string) (playerdata []byte, err error) {
	// packed json and sent.
	jsonStruct := DivingFishB50{Username: userName, B50: true}
	jsonStructData, err := json.Marshal(jsonStruct)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", "https://www.diving-fish.com/api/chunithmprober/query/player", bytes.NewBuffer(jsonStructData))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		panic(err)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == 400 {
		return nil, errors.New("- 未找到用户或者用户数据丢失\n\n - 请检查您是否在 水鱼查分器 上 上传过成绩并且有绑定QQ号")
	}
	if resp.StatusCode == 403 {
		return nil, errors.New("- 该用户设置禁止查分\n\n - 请检查您是否在 水鱼查分器 上 是否关闭了允许他人查分功能")
	}
	playerData, err := io.ReadAll(resp.Body)
	return playerData, err
}

func RenderCardChun(data UserDataInner, renderCount int) image.Image {
	// get pic
	onloadPic, _ := GetCover(strconv.Itoa(data[renderCount].Mid))
	loadTable, _ := gg.LoadImage(Texture + LevelIndexCount(data[renderCount].LevelIndex) + "_table.png")
	getPic := gg.NewContextForImage(loadTable)
	getPic.DrawImage(onloadPic, 250, 10)
	getPic.Fill()
	// draw Name
	getTitleLoader := LoadFontFace(Texture+"fonts/"+"SourceHanSansCN-Bold.otf", 25)
	getPic.SetFontFace(getTitleLoader)
	getPic.SetColor(color.Black)
	getPic.DrawStringAnchored(BreakWords(data[renderCount].Title, 15), 15, 58, 0, 0)
	getPic.Fill()
	// draw FC/AJ if possible.
	var returnFCAJLink string
	if data[renderCount].Fc != "" {
		returnFCAJLink = Texture + "icon_" + "fullcombo" + ".png"
	} else {
		returnFCAJLink = Texture + "icon_" + "clear" + ".png"
	}
	getLink, _ := gg.LoadImage(returnFCAJLink)
	getPic.DrawImage(getLink, 30, 85)
	getPic.Fill()
	// draw line
	getPic.DrawLine(115, 80, 115, 110)
	getPic.Stroke()
	// draw Upper
	getTitleLoaderS := LoadFontFace(Texture+"fonts/"+"SourceHanSansCN-Bold.otf", 18)
	getPic.SetFontFace(getTitleLoaderS)
	getPic.DrawStringAnchored(strconv.FormatFloat(data[renderCount].Ds, 'f', 1, 64), 125, 105, 0, 0)
	getPic.DrawStringAnchored(">", 165, 105, 0, 0)
	// draw num
	DrawBorderString(getPic, "# "+strconv.Itoa(renderCount+1), 3, 10, 20, 0, 0, color.White, color.Black)
	getPic.SetColor(color.White)
	getPic.DrawString("# "+strconv.Itoa(renderCount+1), 10, 20)
	getPic.SetColor(color.Black)
	getTitleLoaderHeader := LoadFontFace(Texture+"fonts/"+"SourceHanSansCN-Bold.otf", 24)
	getPic.SetFontFace(getTitleLoaderHeader)
	getPic.DrawStringAnchored(strconv.FormatFloat(data[renderCount].Ra, 'f', 1, 64), 180, 105, 0, 0)
	getTitleLoaderScore := LoadFontFace(Texture+"fonts/"+"SourceHanSansCN-Bold.otf", 40)
	getPic.SetFontFace(getTitleLoaderScore)
	getPic.DrawStringAnchored(formatNumber(data[renderCount].Score), 120, 130, 0.5, 0.5)
	return getPic.Image()
}

// Resize Image width height
func Resize(image image.Image, w int, h int) image.Image {
	return imgfactory.Size(image, w, h).Image()
}

func LevelIndexCount(count int) string {
	switch {
	case count == 0:
		return "basic"
	case count == 1:
		return "advance"
	case count == 2:
		return "expert"
	case count == 3:
		return "master"
	case count == 4:
		return "ultra"
	case count == 5:
		return "worldend"
	}
	return ""
}

// Used For format Chunithm Score Game
func formatNumber(number int) string {
	numStr := strconv.Itoa(number)
	length := len(numStr)
	zeroCount := 7 - length
	if zeroCount < 0 {
		zeroCount = 0
	}
	for i := 0; i < zeroCount; i++ {
		numStr = "0" + numStr
	}
	formattedStr := ""
	for i, char := range numStr {
		if i > 0 && (length-i)%3 == 0 {
			formattedStr += ","
		}
		formattedStr += string(char)
	}

	return formattedStr
}

func getColorByRating(value float64) color.Color {
	switch {
	case value >= 0.00 && value <= 3.99:
		return color.NRGBA{G: 255, A: 255}
	case value >= 4.00 && value <= 6.99:
		return color.NRGBA{R: 255, G: 102, A: 255}
	case value >= 7.00 && value <= 9.99:
		return color.NRGBA{R: 255, A: 255}
	case value >= 10.00 && value <= 11.99:
		return color.NRGBA{R: 255, B: 255, A: 255}
	case value >= 12.00 && value <= 13.24:
		return color.NRGBA{R: 153, G: 51, A: 255}
	case value >= 13.25 && value <= 14.49:
		return color.NRGBA{R: 128, G: 128, B: 128, A: 255}
	case value >= 14.50 && value <= 15.24:
		return color.NRGBA{R: 255, G: 204, A: 255}
	case value >= 15.25 && value <= 15.99:
		return color.NRGBA{R: 255, G: 255, A: 255}
	default:
		return color.NRGBA{R: 204, G: 153, B: 255, A: 255}
	}
}

func BaseRender(JsonResultData ChunData, ctx *rei.Ctx) image.Image {
	bgMain, err := gg.LoadImage(Texture + "Background_SUN.png")
	if err != nil {
		panic(err)
	}
	getContent := gg.NewContextForImage(bgMain)
	startCountWidth := 700
	StartCountHeight := 800
	baseCount := 0
	// render B30 + B15
	var sumUserB30 float64
	var SumUserR10 float64
	for renderCount := range JsonResultData.Records.B30 {
		sumUserB30 += JsonResultData.Records.B30[renderCount].Ra
		getRender := RenderCardChun(JsonResultData.Records.B30, renderCount)
		getContent.DrawImage(getRender, startCountWidth, StartCountHeight)
		startCountWidth += 550
		baseCount += 1
		if baseCount == 5 {
			baseCount = 0
			startCountWidth = 700
			StartCountHeight += 230
		}
	}
	sumUserB30Result := sumUserB30 / float64(len(JsonResultData.Records.B30))
	startCountWidthRecent := 3840
	StartCountHeightRecent := 915
	baseCountRecent := 0
	for renderCountBase := range JsonResultData.Records.R10 {
		SumUserR10 += JsonResultData.Records.R10[renderCountBase].Ra
		getRender := RenderCardChun(JsonResultData.Records.R10, renderCountBase)
		getContent.DrawImage(getRender, startCountWidthRecent, StartCountHeightRecent)
		startCountWidthRecent += 440
		baseCountRecent += 1
		if baseCountRecent == 2 {
			baseCountRecent = 0
			startCountWidthRecent = 3840
			StartCountHeightRecent += 230
		}
	}
	sumUserR10Result := SumUserR10 / float64(len(JsonResultData.Records.R10))
	// RENDER USER COUNT
	getRecentUserCount := LoadFontFace(Texture+"fonts/"+"SourceHanSansCN-Bold.otf", 60)
	getContent.SetFontFace(getRecentUserCount)
	getContent.SetColor(color.White)
	getContent.DrawStringAnchored("BEST 30: "+strconv.FormatFloat(sumUserB30Result, 'f', 2, 64), 1500, 730, 0, 0)
	getContent.DrawStringAnchored("RECENT 10: "+strconv.FormatFloat(sumUserR10Result, 'f', 2, 64), 2100, 730, 0, 0)
	getContent.Fill()
	// Draw USERTABLE
	getUserNameFontTitle := LoadFontFace(Texture+"fonts/"+"SourceHanSansCN-Regular.otf", 50)
	getContent.SetFontFace(getUserNameFontTitle)
	getContent.SetColor(color.Black)
	getContent.DrawStringAnchored(width.Widen.String(JsonResultData.Nickname), 630, 330, 0, 0)
	getContent.Fill()
	// Rating
	getUserNameFontTitleSmaller := LoadFontFace(Texture+"fonts/"+"SourceHanSansCN-Regular.otf", 35)
	getContent.SetFontFace(getUserNameFontTitleSmaller)
	DrawBorderString(getContent, "RATING ", 3, 490, 420, 0, 0, getColorByRating(JsonResultData.Rating), color.Black)
	getContent.SetFontFace(getUserNameFontTitle)
	DrawBorderString(getContent, width.Widen.String(strconv.FormatFloat(JsonResultData.Rating, 'f', 2, 64)), 3, 630, 420, 0, 0, getColorByRating(JsonResultData.Rating), color.Black)
	getContent.Fill()
	// draw Avatar
	if toolchain.GetTargetAvatar(ctx) != nil {
		getContent.DrawImage(Resize(toolchain.GetTargetAvatar(ctx), 179, 179), 1030, 258)
		getContent.Fill()
	}
	return getContent.Image()
}

// GetCover Careful The nil data
func GetCover(id string) (image.Image, error) {
	fileName := id + ".png"
	filePath := Root + "cover/" + fileName
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// Auto download cover from diving fish's site
		downloadURL := "https://lx-rec-reproxy.lemonkoi.one/chunithm/jacket/" + id + ".png"
		cover, err := DownloadImage(downloadURL)
		if err != nil {
			return nil, err
		}
		SaveImage(cover, filePath)
	}
	imageFile, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer func(imageFile *os.File) {
		err := imageFile.Close()
		if err != nil {
			return
		}
	}(imageFile)
	img, _, err := image.Decode(imageFile)
	if err != nil {
		return nil, err
	}
	return Resize(img, 155, 154), nil
}

// LoadFontFace load font face once before running, to work it quickly and save memory.
func LoadFontFace(filePath string, size float64) font.Face {
	fontFile, _ := os.ReadFile(filePath)
	fontFileParse, _ := opentype.Parse(fontFile)
	fontFace, _ := opentype.NewFace(fontFileParse, &opentype.FaceOptions{Size: size, DPI: 72, Hinting: font.HintingFull})
	return fontFace
}

// BreakWords Reduce the length of strings, if out of range, use ".." instead.
func BreakWords(getSongName string, breakerCount float64) string {
	charCount := 0.0
	setBreaker := false
	var truncated string
	var charFloatNum float64
	for _, runeValue := range getSongName {
		charWidth := utf8.RuneLen(runeValue)
		if charWidth == 3 {
			charFloatNum = 2
		} else {
			charFloatNum = float64(charWidth)
		}
		if charCount+charFloatNum > breakerCount {
			setBreaker = true
			break
		}
		truncated += string(runeValue)
		charCount += charFloatNum
	}
	if setBreaker {
		getSongName = truncated + ".."
	} else {
		getSongName = truncated
	}
	return getSongName
}

// DrawBorderString GG Package Not support The string render, so I write this (^^)
func DrawBorderString(page *gg.Context, s string, size int, x float64, y float64, ax float64, ay float64, inlineRGB color.Color, outlineRGB color.Color) {
	page.SetColor(outlineRGB)
	n := size
	for dy := -n; dy <= n; dy++ {
		for dx := -n; dx <= n; dx++ {
			if dx*dx+dy*dy >= n*n {
				continue
			}
			renderX := x + float64(dx)
			renderY := y + float64(dy)
			page.DrawStringAnchored(s, renderX, renderY, ax, ay)
		}
	}
	page.SetColor(inlineRGB)
	page.DrawStringAnchored(s, x, y, ax, ay)
}

// SaveImage Save Cover Chun | Maimai
func SaveImage(img image.Image, path string) {
	files, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}
	defer func(files *os.File) {
		err := files.Close()
		if err != nil {
			return
		}
	}(files)
	err = png.Encode(files, img)
	if err != nil {
		log.Fatal(err)
	}
}

// DownloadImage Simple Downloader.
func DownloadImage(url string) (image.Image, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(response.Body)
	img, _, err := image.Decode(response.Body)
	if err != nil {
		return nil, err
	}
	return img, nil
}
