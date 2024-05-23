module github.com/MoYoez/Lucy_reibot

go 1.22.2

require (
	github.com/FloatTech/floatbox v0.0.0-20230331064925-9af336a84944
	github.com/FloatTech/gg v1.1.2
	github.com/FloatTech/imgfactory v0.2.1
	github.com/FloatTech/sqlite v1.6.2
	github.com/FloatTech/zbpctrl v1.5.3-0.20230130095145-714ad318cd52
	github.com/disintegration/imaging v1.6.2
	github.com/fogleman/gg v1.3.0
	github.com/fumiama/ReiBot v0.0.0-00010101000000-000000000000
	github.com/fumiama/go-base16384 v1.7.0
	github.com/fumiama/gotracemoe v0.0.3
	github.com/go-telegram-bot-api/telegram-bot-api/v5 v5.5.1
	github.com/jinzhu/gorm v1.9.16
	github.com/joho/godotenv v1.5.1
	github.com/mroth/weightedrand v1.0.0
	github.com/shirou/gopsutil v3.21.11+incompatible
	github.com/sirupsen/logrus v1.9.0
	github.com/tidwall/gjson v1.14.4
	github.com/wdvxdr1123/ZeroBot v1.6.11
	github.com/xuri/excelize/v2 v2.8.1
	golang.org/x/image v0.14.0
	golang.org/x/text v0.14.0
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/FloatTech/AnimeAPI v1.7.0 // indirect
	github.com/FloatTech/ttl v0.0.0-20220715042055-15612be72f5b // indirect
	github.com/MoYoez/LucyBaseRouter-Toolchain v0.0.0-20240513105735-fa25f6949091 // indirect
	github.com/RomiChan/syncx v0.0.0-20221202055724-5f842c53020e // indirect
	github.com/ericpauley/go-quantize v0.0.0-20200331213906-ae555eb2afa4 // indirect
	github.com/fumiama/cron v1.3.0 // indirect
	github.com/fumiama/go-registry v0.2.6 // indirect
	github.com/fumiama/go-simple-protobuf v0.1.0 // indirect
	github.com/fumiama/gofastTEA v0.0.10 // indirect
	github.com/fumiama/imgsz v0.0.2 // indirect
	github.com/fumiama/jieba v0.0.0-20221203025406-36c17a10b565 // indirect
	github.com/go-ole/go-ole v1.2.6 // indirect
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/mattn/go-isatty v0.0.16 // indirect
	github.com/mohae/deepcopy v0.0.0-20170929034955-c48cc78d4826 // indirect
	github.com/pbnjay/memory v0.0.0-20210728143218-7b4eea64cf58 // indirect
	github.com/remyoudompheng/bigfft v0.0.0-20200410134404-eec4a21b6bb0 // indirect
	github.com/richardlehane/mscfb v1.0.4 // indirect
	github.com/richardlehane/msoleps v1.0.3 // indirect
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.0 // indirect
	github.com/tklauser/go-sysconf v0.3.12 // indirect
	github.com/tklauser/numcpus v0.6.1 // indirect
	github.com/xuri/efp v0.0.0-20231025114914-d1ff6096ae53 // indirect
	github.com/xuri/nfp v0.0.0-20230919160717-d98342af3f05 // indirect
	github.com/yusufpapurcu/wmi v1.2.3 // indirect
	golang.org/x/crypto v0.21.0 // indirect
	golang.org/x/net v0.22.0 // indirect
	golang.org/x/sys v0.18.0 // indirect
	modernc.org/libc v1.21.5 // indirect
	modernc.org/mathutil v1.5.0 // indirect
	modernc.org/memory v1.4.0 // indirect
	modernc.org/sqlite v1.20.0 // indirect
)

replace modernc.org/sqlite => github.com/fumiama/sqlite3 v1.20.0-with-win386

replace github.com/remyoudompheng/bigfft => github.com/fumiama/bigfft v0.0.0-20211011143303-6e0bfa3c836b

replace github.com/fumiama/ReiBot => github.com/MoYoez/ReiBot v0.0.0-20240311142840-76873b41aa4a

replace github.com/go-telegram-bot-api/telegram-bot-api/v5 => ./telegram-bot-api/
