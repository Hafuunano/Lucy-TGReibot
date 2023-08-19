module github.com/FloatTech/ReiBot-Plugin

go 1.19

require (
	github.com/FloatTech/AnimeAPI v1.6.1-0.20230207081411-573533b18194
	github.com/FloatTech/floatbox v0.0.0-20230207080446-026a2f086c74
	github.com/FloatTech/zbpctrl v1.5.3-0.20230130095145-714ad318cd52
	github.com/fumiama/ReiBot v0.0.0-20230215122748-dab38cf6b71b
	github.com/fumiama/go-base16384 v1.6.4
	github.com/fumiama/gotracemoe v0.0.3
	github.com/go-telegram-bot-api/telegram-bot-api/v5 v5.5.1
	github.com/jozsefsallai/gophersauce v1.0.1
	github.com/sirupsen/logrus v1.9.0
	github.com/wdvxdr1123/ZeroBot v1.6.9
)

require (
	github.com/FloatTech/sqlite v1.5.7 // indirect
	github.com/FloatTech/ttl v0.0.0-20220715042055-15612be72f5b // indirect
	github.com/RomiChan/syncx v0.0.0-20221202055724-5f842c53020e // indirect
	github.com/fumiama/cron v1.3.0 // indirect
	github.com/fumiama/go-registry v0.2.5 // indirect
	github.com/fumiama/go-simple-protobuf v0.1.0 // indirect
	github.com/fumiama/gofastTEA v0.0.10 // indirect
	github.com/gabriel-vasile/mimetype v1.0.4 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/mattn/go-isatty v0.0.16 // indirect
	github.com/remyoudompheng/bigfft v0.0.0-20200410134404-eec4a21b6bb0 // indirect
	github.com/tidwall/gjson v1.14.4 // indirect
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.0 // indirect
	golang.org/x/sys v0.1.1-0.20221102194838-fc697a31fa06 // indirect
	golang.org/x/text v0.6.0 // indirect
	modernc.org/libc v1.21.5 // indirect
	modernc.org/mathutil v1.5.0 // indirect
	modernc.org/memory v1.4.0 // indirect
	modernc.org/sqlite v1.20.0 // indirect
)

replace modernc.org/sqlite => github.com/fumiama/sqlite3 v1.20.0-with-win386

replace github.com/remyoudompheng/bigfft => github.com/fumiama/bigfft v0.0.0-20211011143303-6e0bfa3c836b
