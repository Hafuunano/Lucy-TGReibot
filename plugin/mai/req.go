package mai

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/FloatTech/floatbox/web"
	"io"
	"net/http"
	"os"
)

type DivingFishB50UserName struct {
	Username string `json:"username"`
	B50      bool   `json:"b50"`
}

type DivingFishDevFullDataRecords struct {
	AdditionalRating int    `json:"additional_rating"`
	Nickname         string `json:"nickname"`
	Plate            string `json:"plate"`
	Rating           int    `json:"rating"`
	Records          []struct {
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
	} `json:"records"`
	Username string `json:"username"`
}

func QueryMaiBotDataFromUserName(username string) (playerdata []byte, err error) {
	// packed json and sent.
	jsonStruct := DivingFishB50UserName{Username: username, B50: true}
	jsonStructData, err := json.Marshal(jsonStruct)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", "https://www.diving-fish.com/api/maimaidxprober/query/player", bytes.NewBuffer(jsonStructData))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		return
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode == 400 {
		return nil, errors.New("- 未找到用户或者用户数据丢失\n\n - 请检查您是否在 https://www.diving-fish.com/maimaidx/prober/ 上 上传过成绩")
	}
	if resp.StatusCode == 403 {
		return nil, errors.New("- 该用户设置禁止查分\n\n - 请检查您是否在 https://www.diving-fish.com/maimaidx/prober/ 上 是否关闭了允许他人查分功能")
	}
	playerDataByte, err := io.ReadAll(resp.Body)
	return playerDataByte, err
}

func QueryDevDataFromDivingFish(username string) DivingFishDevFullDataRecords {
	getData, err := web.RequestDataWithHeaders(web.NewDefaultClient(), "https://www.diving-fish.com/api/maimaidxprober/dev/player/records?username="+username, "GET", func(request *http.Request) error {
		request.Header.Add("Developer-Token", os.Getenv("dvkey"))
		return nil
	}, nil)
	if err != nil {
		return DivingFishDevFullDataRecords{}
	}
	var handlerData DivingFishDevFullDataRecords
	json.Unmarshal(getData, &handlerData)
	return handlerData
}
