package mai

import (
	"encoding/json"
	"github.com/FloatTech/floatbox/web"
	"net/http"
	"os"
)

type LxnsAliases struct {
	Aliases []struct {
		SongId  int      `json:"song_id"`
		Aliases []string `json:"aliases"`
	} `json:"aliases"`
}

// only support LXNS because => DivingFish need Token.

// RequestAliasFromLxns Get Alias From LXNS Network.
func RequestAliasFromLxns() LxnsAliases {
	getData, err := web.RequestDataWithHeaders(web.NewDefaultClient(), "https://maimai.lxns.net/api/v0/maimai/alias/list", "GET", func(request *http.Request) error {
		request.Header.Add("Authorization", os.Getenv("lxnskey"))
		return nil
	}, nil)
	if err != nil {
		return LxnsAliases{}
	}
	var handlerData LxnsAliases
	json.Unmarshal(getData, &handlerData)
	return handlerData
}

// QueryReferSong CASE: In DivingFish Mode this SongID will be added "00" AHEAD IF SONGID IS LOWER THAN 1000 (if lower than 100 then it will be added "000" ) , Otherwise it will be added "1" AHEAD. || LXNS don't need to do anything. (DEFAULT RETURN LXNS SONGDATA)
func (requester *LxnsAliases) QueryReferSong(songAlias string) (status bool, songID int64) {
	for i, dataInnner := range requester.Aliases {
		for _, v := range dataInnner.Aliases {
			if songAlias == v {
				return true, int64(requester.Aliases[i].SongId)
			}
		}
	}
	return false, 0
}
