package mai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/xuri/excelize/v2"
	"io"
	"net/http"
	"os"
	"strconv"
)

type MappedListStruct struct {
	DingFishId int      `json:"dv_id"`
	SongName   string   `json:"song_name"`
	SongId     []int    `json:"song_id"`
	Aliases    []string `json:"aliases"`
}

type LxnsSongListInfo struct {
	Songs []struct {
		Id           int    `json:"id"`
		Title        string `json:"title"`
		Artist       string `json:"artist"`
		Genre        string `json:"genre"`
		Bpm          int    `json:"bpm"`
		Version      int    `json:"version"`
		Difficulties struct {
			Standard []struct {
				Type         string  `json:"type"`
				Difficulty   int     `json:"difficulty"`
				Level        string  `json:"level"`
				LevelValue   float64 `json:"level_value"`
				NoteDesigner string  `json:"note_designer"`
				Version      int     `json:"version"`
			} `json:"standard"`
			Dx []struct {
				Type         string  `json:"type"`
				Difficulty   int     `json:"difficulty"`
				Level        string  `json:"level"`
				LevelValue   float64 `json:"level_value"`
				NoteDesigner string  `json:"note_designer"`
				Version      int     `json:"version"`
			} `json:"dx"`
		} `json:"difficulties"`
	} `json:"songs"`
	Genres []struct {
		Id    int    `json:"id"`
		Title string `json:"title"`
		Genre string `json:"genre"`
	} `json:"genres"`
	Versions []struct {
		Id      int    `json:"id"`
		Title   string `json:"title"`
		Version int    `json:"version"`
	} `json:"versions"`
}

type LxnsAliases struct {
	Aliases []struct {
		SongId  int      `json:"song_id"`
		Aliases []string `json:"aliases"`
	} `json:"aliases"`
}

type AliasesReturnValue struct {
	Aliases []struct {
		DvId     int      `json:"dv_id"`
		SongName string   `json:"song_name"`
		SongId   []int    `json:"song_id"`
		Aliases  []string `json:"aliases"`
	} `json:"aliases"`
}

// only support LXNS because => if DivingFish then need Token.

// QueryReferSong use LocalStorageData.
func QueryReferSong(Alias string, isLxnet bool) (status bool, id []int, needAcc bool) {
	// unpackedData
	getData, err := os.ReadFile(engine.DataFolder() + "alias.json")
	if err != nil {
		panic(err)
	}
	var DataHandler AliasesReturnValue
	json.Unmarshal(getData, &DataHandler)
	var onloadList [][]int
	for _, dataSearcher := range DataHandler.Aliases {
		for _, aliasSearcher := range dataSearcher.Aliases {
			if aliasSearcher == Alias {
				onloadList = append(onloadList, dataSearcher.SongId) // write in memory
			}
		}
	}
	// if list is 2,query them is from the same song? | if above 2(3 or more ,means this song need acc.)
	switch {
	case len(onloadList) == 1: // only one query.
		if isLxnet {
			for _, listhere := range onloadList[0] {
				if listhere < 10000 {
					return true, []int{listhere}, false
				}
			}
		} else {
			return true, onloadList[0], false
		}
	// query length is 2,it means this maybe same name but diff id ==> (E.G: Oshama Scramble!)
	case len(onloadList) == 2:
		for _, listHere := range onloadList[0] {
			for _, listAliasTwo := range onloadList[1] {
				if listHere == listAliasTwo {
					// same list here.
					var returnIntList []int
					returnIntList = append(returnIntList, onloadList[0]...)
					returnIntList = append(returnIntList, onloadList[1]...)
					returnIntList = removeIntDuplicates(returnIntList)
					if isLxnet {
						for _, listhere := range returnIntList {
							if listhere < 10000 {
								return true, []int{listhere}, false
							}
						}
					} else {
						return true, returnIntList, false
					}
				}
			}
		}
		// if query is none, means it need moreacc
		return true, nil, true
	case len(onloadList) >= 3:
		return true, nil, true
	}
	// no found.
	return false, nil, false
}

// UpdateAliasPackage Use simple action to update alias.
func UpdateAliasPackage() {
	// get Lxns Data
	respls, err := http.Get("https://maimai.lxns.net/api/v0/maimai/alias/list")
	defer respls.Body.Close()
	getData, err := io.ReadAll(respls.Body)
	var lxnsAliasData LxnsAliases
	json.Unmarshal(getData, &lxnsAliasData)

	// get Lxns Data SongListInfo
	resplsSongList, err := http.Get("https://maimai.lxns.net/api/v0/maimai/song/list")
	defer respls.Body.Close()
	getDataSongList, err := io.ReadAll(resplsSongList.Body)
	var lxnsSongListData LxnsSongListInfo
	json.Unmarshal(getDataSongList, &lxnsSongListData)

	// get AkiraBot AliasData
	url := "https://docs.google.com/spreadsheets/d/e/2PACX-1vRwHptWLUyMG9ASCgk9MhI693jmAA1_CJrPfTxjX9J8f3wGHlR09Ja_h5i3InPbFhK1BjJp5cO_kugM/pub?output=xlsx"
	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error fetching URL:", err)
		return
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	reader := bytes.NewReader(body)
	f, err := excelize.OpenReader(reader)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	savedListMap := map[string][]string{}
	getRows, err := f.GetRows("主表")
	var titleStart bool
	for _, rows := range getRows {
		if rows[0] == "ID" {
			titleStart = true
			continue
		}
		if titleStart {
			var mappedList []string
			for _, rowList := range rows {
				mappedList = append(mappedList, rowList)
			}
			savedListMap[mappedList[0]] = mappedList[1:]
		}
	}
	// generate a json file here.
	var tempList []interface{}
	for i, listData := range savedListMap {
		var vartiesList []int
		getInt, _ := strconv.Atoi(i)
		vartiesList = append(vartiesList, getInt)
		var referListData []string
		// check this alias in lxns network pattern, it maybe slowly(
		for _, listLxns := range lxnsSongListData.Songs {
			if listLxns.Title == listData[0] && listLxns.Id != getInt {
				vartiesList = append(vartiesList, listLxns.Id)
			}
		}
		// due to AkihaBot use two packed id, so make the id together.
		referListData = append(referListData, listData[1:]...)
		// add same alias to lxns
		for _, lxnsAliasRefer := range lxnsAliasData.Aliases {
			for _, listLocation := range vartiesList {
				if listLocation == lxnsAliasRefer.SongId {
					// prefix, add alias to it.
					referListData = append(referListData, lxnsAliasRefer.Aliases...)
					referListData = removeDuplicates(referListData)
				}
			}
		}
		tempList = append(tempList, &MappedListStruct{DingFishId: getInt, SongName: listData[0], Aliases: referListData, SongId: vartiesList})
	}
	GeneratedList := map[string]interface{}{
		"aliases": tempList,
	}
	getBytes, err := json.Marshal(GeneratedList)
	if err != nil {
		panic(err)
	}
	os.WriteFile(engine.DataFolder()+"alias.json", getBytes, 0777)
	GeneratedList = nil // revoke
}

func removeDuplicates(list []string) []string {
	seen := make(map[string]bool)
	var result []string
	for _, item := range list {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}
	return result
}

func removeIntDuplicates(list []int) []int {
	seen := make(map[int]bool)
	var result []int
	for _, item := range list {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}
	return result
}
