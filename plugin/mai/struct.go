package mai

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"unicode/utf8"

	"github.com/FloatTech/gg"
	"github.com/FloatTech/imgfactory"
	"github.com/MoYoez/Lucy_reibot/utils/toolchain"
	"github.com/MoYoez/Lucy_reibot/utils/transform"
	rei "github.com/fumiama/ReiBot"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/text/width"
)

type BaseUserIDStruct struct {
	UserID int64 `json:"user_id"`
}

type GetMusicStruct struct {
	BaseUserIDStruct
	GetIndex   int64 `json:"get_index"`
	GetCounter int64 `json:"get_counter"`
}

type player struct {
	AdditionalRating int `json:"additional_rating"`
	Charts           struct {
		Dx []struct {
			Achievements float64 `json:"achievements"`
			Ds           float64 `json:"ds"`
			DxScore      int     `json:"dxScore"`
			Fc           string  `json:"fc"`
			Fs           string  `json:"fs"`
			Level        string  `json:"level"`
			LevelIndex   int     `json:"level_index"`
			LevelLabel   string  `json:"level_label"`
			Ra           int     `json:"ra"`
			Rate         string  `json:"rate"`
			SongId       int     `json:"song_id"`
			Title        string  `json:"title"`
			Type         string  `json:"type"`
		} `json:"dx"`
		Sd []struct {
			Achievements float64 `json:"achievements"`
			Ds           float64 `json:"ds"`
			DxScore      int     `json:"dxScore"`
			Fc           string  `json:"fc"`
			Fs           string  `json:"fs"`
			Level        string  `json:"level"`
			LevelIndex   int     `json:"level_index"`
			LevelLabel   string  `json:"level_label"`
			Ra           int     `json:"ra"`
			Rate         string  `json:"rate"`
			SongId       int     `json:"song_id"`
			Title        string  `json:"title"`
			Type         string  `json:"type"`
		} `json:"sd"`
	} `json:"charts"`
	Nickname string      `json:"nickname"`
	Plate    string      `json:"plate"`
	Rating   int         `json:"rating"`
	UserData interface{} `json:"user_data"`
	UserId   interface{} `json:"user_id"`
	Username string      `json:"username"`
}

type playerData struct {
	Achievements float64 `json:"achievements"`
	Ds           float64 `json:"ds"`
	DxScore      int     `json:"dxScore"`
	Fc           string  `json:"fc"`
	Fs           string  `json:"fs"`
	Level        string  `json:"level"`
	LevelIndex   int     `json:"level_index"`
	LevelLabel   string  `json:"level_label"`
	Ra           int     `json:"ra"`
	Rate         string  `json:"rate"`
	SongId       int     `json:"song_id"`
	Title        string  `json:"title"`
	Type         string  `json:"type"`
}

type LxnsUploaderStruct struct {
	Score []LxnsScoreUploader `json:"scores"`
}

type LxnsScoreUploader struct {
	Id           int         `json:"id"`
	Type         string      `json:"type"`
	LevelIndex   int         `json:"level_index"`
	Achievements float64     `json:"achievements"`
	Fc           interface{} `json:"fc"`
	Fs           interface{} `json:"fs"`
	DxScore      int         `json:"dx_score"`
	PlayTime     string      `json:"play_time"`
}

var (
	loadMaiPic        = Root + "pic/"
	defaultCoverLink  = Root + "default_cover.png"
	typeImageDX       = loadMaiPic + "chart_type_dx.png"
	typeImageSD       = loadMaiPic + "chart_type_sd.png"
	titleFontPath     = maifont + "NotoSansSC-Bold.otf"
	UniFontPath       = maifont + "Montserrat-Bold.ttf"
	nameFont          = maifont + "NotoSansSC-Regular.otf"
	maifont           = Root + "font/"
	b50bgOriginal     = loadMaiPic + "b50_bg.png"
	b50bg             = loadMaiPic + "b50_bg.png"
	b50Custom         = loadMaiPic + "b50_bg_custom.png"
	Root              = transform.ReturnLucyMainDataIndex("maidx") + "resources/maimai/"
	userPlate         = engine.DataFolder() + "user/"
	titleFont         font.Face
	scoreFont         font.Face
	rankFont          font.Face
	levelFont         font.Face
	ratingFont        font.Face
	nameTypeFont      font.Face
	diffColor         []color.RGBA
	ratingBgFilenames = []string{
		"rating_white.png",
		"rating_blue.png",
		"rating_green.png",
		"rating_yellow.png",
		"rating_red.png",
		"rating_purple.png",
		"rating_copper.png",
		"rating_silver.png",
		"rating_gold.png",
		"rating_rainbow.png",
	}
)

func init() {
	if _, err := os.Stat(userPlate); os.IsNotExist(err) {
		err := os.MkdirAll(userPlate, 0777)
		if err != nil {
			return
		}
	}
	nameTypeFont = LoadFontFace(nameFont, 36)
	titleFont = LoadFontFace(titleFontPath, 20)
	scoreFont = LoadFontFace(UniFontPath, 32)
	rankFont = LoadFontFace(UniFontPath, 24)
	levelFont = LoadFontFace(UniFontPath, 20)
	ratingFont = LoadFontFace(UniFontPath, 24)
	diffColor = []color.RGBA{
		{69, 193, 36, 255},
		{255, 186, 1, 255},
		{255, 90, 102, 255},
		{134, 49, 200, 255},
		{207, 144, 240, 255},
	}

}

// FullPageRender  Render Full Page
func FullPageRender(data player, ctx *rei.Ctx) (raw image.Image) {
	// muilt-threading.
	getUserID, _ := toolchain.GetChatUserInfoID(ctx)
	var avatarHandler sync.WaitGroup
	avatarHandler.Add(1)
	var getAvatarFormat *gg.Context
	// avatar handler.
	go func() {
		// avatar Round Style
		defer avatarHandler.Done()
		getAvatar := toolchain.GetTargetAvatar(ctx)
		if getAvatar != nil {
			avatarFormat := imgfactory.Size(getAvatar, 180, 180)
			getAvatarFormat = gg.NewContext(180, 180)
			getAvatarFormat.DrawRoundedRectangle(0, 0, 178, 178, 20)
			getAvatarFormat.Clip()
			getAvatarFormat.DrawImage(avatarFormat.Image(), 0, 0)
			getAvatarFormat.Fill()
		}
	}()
	userPlatedCustom := GetUserDefaultBackgroundDataFromDatabase(getUserID)
	// render Header.
	b50Render := gg.NewContext(2090, 1660)
	rawPlateData, errs := gg.LoadImage(userPlate + strconv.Itoa(int(getUserID)) + ".png")
	if errs == nil {
		b50bg = b50Custom
		b50Render.DrawImage(rawPlateData, 595, 30)
		b50Render.Fill()
	} else {
		if userPlatedCustom != "" {
			b50bg = b50Custom
			images, _ := GetDefaultPlate(userPlatedCustom)
			b50Render.DrawImage(images, 595, 30)
			b50Render.Fill()
		} else {
			// show nil
			b50bg = b50bgOriginal
		}
	}
	getContent, _ := gg.LoadImage(b50bg)
	b50Render.DrawImage(getContent, 0, 0)
	b50Render.Fill()
	// render user info
	avatarHandler.Wait()
	if getAvatarFormat != nil {
		b50Render.DrawImage(getAvatarFormat.Image(), 610, 50)
		b50Render.Fill()
	}
	// render Userinfo
	b50Render.SetFontFace(nameTypeFont)
	b50Render.SetColor(color.Black)
	b50Render.DrawStringAnchored(width.Widen.String(data.Nickname), 825, 160, 0, 0)
	b50Render.Fill()
	b50Render.SetFontFace(titleFont)
	setPlateLocalStatus := GetUserPlateInfoFromDatabase(getUserID)
	if setPlateLocalStatus != "" {
		data.Plate = setPlateLocalStatus
	}
	b50Render.DrawStringAnchored(strings.Join(strings.Split(data.Plate, ""), " "), 1050, 207, 0.5, 0.5)
	b50Render.Fill()
	getRating := getRatingBg(data.Rating)
	getRatingBG, err := gg.LoadImage(loadMaiPic + getRating)
	if err != nil {
		return
	}
	b50Render.DrawImage(getRatingBG, 800, 40)
	b50Render.Fill()
	// render Rank
	imgs, err := GetRankPicRaw(data.AdditionalRating)
	if err != nil {
		return
	}
	b50Render.DrawImage(imgs, 1080, 50)
	b50Render.Fill()
	// draw number
	b50Render.SetFontFace(scoreFont)
	b50Render.SetRGBA255(236, 219, 113, 255)
	b50Render.DrawStringAnchored(strconv.Itoa(data.Rating), 1056, 60, 1, 1)
	b50Render.Fill()
	// Render Card Type
	getSDLength := len(data.Charts.Sd)
	getDXLength := len(data.Charts.Dx)
	getDXinitX := 45
	getDXinitY := 1225
	getInitX := 45
	getInitY := 285
	var i int
	for i = 0; i < getSDLength; i++ {
		b50Render.DrawImage(RenderCard(data.Charts.Sd[i], i+1, false), getInitX, getInitY)
		getInitX += 400
		if getInitX == 2045 {
			getInitX = 45
			getInitY += 125
		}
	}

	for dx := 0; dx < getDXLength; dx++ {
		b50Render.DrawImage(RenderCard(data.Charts.Dx[dx], dx+1, false), getDXinitX, getDXinitY)
		getDXinitX += 400
		if getDXinitX == 2045 {
			getDXinitX = 45
			getDXinitY += 125
		}
	}
	return b50Render.Image()
}

// RenderCard Main Lucy Render Page , if isSimpleRender == true, then render count will not show here.
func RenderCard(data playerData, num int, isSimpleRender bool) image.Image {
	getType := data.Type
	var CardBackGround string
	var multiTypeRender sync.WaitGroup
	var CoverDownloader sync.WaitGroup
	CoverDownloader.Add(1)
	multiTypeRender.Add(1)
	// choose Type.
	if getType == "SD" {
		CardBackGround = typeImageSD
	} else {
		CardBackGround = typeImageDX
	}
	charCount := 0.0
	setBreaker := false
	var truncated string
	var charFloatNum float64
	getSongName := data.Title
	var getSongId string
	switch {
	case data.SongId < 1000:
		getSongId = fmt.Sprintf("%05d", data.SongId)
	case data.SongId < 10000:
		getSongId = fmt.Sprintf("1%d", data.SongId)
	default:
		getSongId = strconv.Itoa(data.SongId)
	}
	var Image image.Image
	go func() {
		defer CoverDownloader.Done()
		Image, _ = GetCover(getSongId)
	}()
	// set rune count
	go func() {
		defer multiTypeRender.Done()
		for _, runeValue := range getSongName {
			charWidth := utf8.RuneLen(runeValue)
			if charWidth == 3 {
				charFloatNum = 1.5
			} else {
				charFloatNum = float64(charWidth)
			}
			if charCount+charFloatNum > 19 {
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
	}()
	loadSongType, _ := gg.LoadImage(CardBackGround)
	// draw pic
	drawBackGround := gg.NewContextForImage(GetChartType(data.LevelLabel))
	// draw song pic
	CoverDownloader.Wait()
	drawBackGround.DrawImage(Image, 25, 25)
	// draw name
	drawBackGround.SetColor(color.White)
	drawBackGround.SetFontFace(titleFont)
	multiTypeRender.Wait()
	drawBackGround.DrawStringAnchored(getSongName, 130, 32.5, 0, 0.5)
	drawBackGround.Fill()
	// draw acc
	drawBackGround.SetFontFace(scoreFont)
	drawBackGround.DrawStringAnchored(strconv.FormatFloat(data.Achievements, 'f', 4, 64)+"%", 129, 62.5, 0, 0.5)
	// draw rate
	drawBackGround.DrawImage(GetRateStatusAndRenderToImage(data.Rate), 305, 45)
	drawBackGround.Fill()
	drawBackGround.SetFontFace(rankFont)
	drawBackGround.SetColor(diffColor[data.LevelIndex])
	if !isSimpleRender {
		drawBackGround.DrawString("#"+strconv.Itoa(num), 130, 111)
	}
	drawBackGround.FillPreserve()
	// draw rest of card.
	drawBackGround.SetFontFace(levelFont)
	drawBackGround.DrawString(strconv.FormatFloat(data.Ds, 'f', 1, 64), 195, 111)
	drawBackGround.FillPreserve()
	drawBackGround.SetFontFace(ratingFont)
	drawBackGround.DrawString("▶", 235, 111)
	drawBackGround.FillPreserve()
	drawBackGround.SetFontFace(ratingFont)
	drawBackGround.DrawString(strconv.Itoa(data.Ra), 250, 111)
	drawBackGround.FillPreserve()
	if data.Fc != "" {
		drawBackGround.DrawImage(LoadComboImage(data.Fc), 290, 84)
	}
	if data.Fs != "" {
		drawBackGround.DrawImage(LoadSyncImage(data.Fs), 325, 84)
	}
	drawBackGround.DrawImage(loadSongType, 68, 88)
	return drawBackGround.Image()
}

func GetRankPicRaw(id int) (image.Image, error) {
	var idStr string
	if id < 10 {
		idStr = "0" + strconv.FormatInt(int64(id), 10)
	} else {
		idStr = strconv.FormatInt(int64(id), 10)
	}
	if id == 22 {
		idStr = "21"
	}
	data := Root + "rank/UI_CMN_DaniPlate_" + idStr + ".png"
	imgRaw, err := gg.LoadImage(data)
	if err != nil {
		return nil, err
	}
	return imgRaw, nil
}

func GetDefaultPlate(id string) (image.Image, error) {
	data := Root + "plate/plate_" + id + ".png"
	imgRaw, err := gg.LoadImage(data)
	if err != nil {
		return nil, err
	}
	return imgRaw, nil
}

// GetCover Careful The nil data
func GetCover(id string) (image.Image, error) {
	fileName := id + ".png"
	filePath := Root + "cover/" + fileName
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// Auto download cover from diving fish's site
		downloadURL := "https://www.diving-fish.com/covers/" + fileName
		cover, err := downloadImage(downloadURL)
		if err != nil {
			return LoadPictureWithResize(defaultCoverLink, 90, 90), nil
		}
		saveImage(cover, filePath)
	}
	imageFile, err := os.Open(filePath)
	if err != nil {
		return LoadPictureWithResize(defaultCoverLink, 90, 90), nil
	}
	defer func(imageFile *os.File) {
		err := imageFile.Close()
		if err != nil {
			return
		}
	}(imageFile)
	img, _, err := image.Decode(imageFile)
	if err != nil {
		return LoadPictureWithResize(defaultCoverLink, 90, 90), nil
	}
	return Resize(img, 90, 90), nil
}

// Resize Image width height
func Resize(image image.Image, w int, h int) image.Image {
	return imgfactory.Size(image, w, h).Image()
}

// LoadPictureWithResize Load Picture
func LoadPictureWithResize(link string, w int, h int) image.Image {
	getImage, err := gg.LoadImage(link)
	if err != nil {
		return nil
	}
	return Resize(getImage, w, h)
}

// GetRateStatusAndRenderToImage Get Rate
func GetRateStatusAndRenderToImage(rank string) image.Image {
	// Load rank images
	return LoadPictureWithResize(loadMaiPic+"rate_"+rank+".png", 80, 40)
}

// GetChartType Get Chart Type
func GetChartType(chart string) image.Image {
	data, _ := gg.LoadImage(loadMaiPic + "chart_" + NoHeadLineCase(chart) + ".png")
	return data
}

// LoadComboImage Load combo images
func LoadComboImage(imageName string) image.Image {
	link := loadMaiPic + "combo_" + imageName + ".png"
	return LoadPictureWithResize(link, 60, 40)
}

// LoadSyncImage Load sync images
func LoadSyncImage(imageName string) image.Image {
	link := loadMaiPic + "sync_" + imageName + ".png"
	return LoadPictureWithResize(link, 60, 40)
}

// NoHeadLineCase No HeadLine.
func NoHeadLineCase(word string) string {
	text := strings.ToLower(word)
	textNewer := strings.ReplaceAll(text, ":", "")
	return textNewer
}

// LoadFontFace load font face once before running, to work it quickly and save memory.
func LoadFontFace(filePath string, size float64) font.Face {
	fontFile, _ := os.ReadFile(filePath)
	fontFileParse, _ := opentype.Parse(fontFile)
	fontFace, _ := opentype.NewFace(fontFileParse, &opentype.FaceOptions{Size: size, DPI: 70, Hinting: font.HintingFull})
	return fontFace
}

// Inline Code.
func saveImage(img image.Image, path string) {
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

func downloadImage(url string) (image.Image, error) {
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

func getRatingBg(rating int) string {
	index := 0
	switch {
	case rating >= 15000:
		index++
		fallthrough
	case rating >= 14000:
		index++
		fallthrough
	case rating >= 13000:
		index++
		fallthrough
	case rating >= 12000:
		index++
		fallthrough
	case rating >= 10000:
		index++
		fallthrough
	case rating >= 7000:
		index++
		fallthrough
	case rating >= 4000:
		index++
		fallthrough
	case rating >= 2000:
		index++
		fallthrough
	case rating >= 1000:
		index++
	}
	return ratingBgFilenames[index]
}
