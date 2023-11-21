package mai

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type DivingFishB50UserName struct {
	Username string `json:"username"`
	B50      bool   `json:"b50"`
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
