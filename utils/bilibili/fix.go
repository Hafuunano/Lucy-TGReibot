package bilibili

import (
	"github.com/FloatTech/AnimeAPI/bilibili"
	"regexp"
)

func BilibiliFixedLink(link string) string {
	// query for the last link here.
	getRealLink, err := bilibili.GetRealURL("https://" + link)
	if err != nil {
		return ""
	}

	getReq, err := regexp.Compile("bilibili.com\\\\?/video\\\\?/(?:av(\\d+)|([bB][vV][0-9a-zA-Z]+))")
	if err != nil {
		return ""
	}
	return getReq.FindAllString(getRealLink, 1)[0]
}
