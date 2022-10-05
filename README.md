<div align="center">
  <a href="https://crypko.ai/crypko/GtWYDpVMx5GYm/">
  <img src=".github/Misaki.png" alt="看板娘" width = "256">
  </a><br>

  <h1>ReiBot-Plugin</h1>
  基于 ReiBot 的 Telegram 插件<br><br>

  <img src="http://cmoe.azurewebsites.net/cmoe?name=ReiBot&theme=r34" /><br>

</div>

## 命令行参数
> `[]`代表是可选参数
```bash
reibot [-tbdoTh] ID1 ID2 ...

-T int
        timeout (default 60)
  -b int
        message sequence length (default 256)
  -d    enable debug-level log output
  -h    print this help
  -o int
        the last Update ID to include
  -t string
        telegram api token
```

## 功能
> 在编译时，以下功能均可通过注释`main.go`中的相应`import`而物理禁用，减小插件体积。

<details>
  <summary>base16384加解密</summary>

  `import _ "github.com/FloatTech/ReiBot-Plugin/plugin/b14"`

  - [x] 加密xxx

  - [x] 解密xxx

  - [x] 用yyy加密xxx

  - [x] 用yyy解密xxx

</details>

<details>
  <summary>base64卦加解密</summary>

  `import _ "github.com/FloatTech/ReiBot-Plugin/plugin/base64gua"`

  - [x] 六十四卦加密xxx

  - [x] 六十四卦解密xxx

  - [x] 六十四卦用yyy加密xxx

  - [x] 六十四卦用yyy解密xxx

</details>

<details>
  <summary>base天城文加解密</summary>

  `import _ "github.com/FloatTech/ReiBot-Plugin/plugin/baseamasiro"`

  - [x] 天城文加密xxx

  - [x] 天城文解密xxx

  - [x] 天城文用yyy加密xxx

  - [x] 天城文用yyy解密xxx

</details>

<details>
  <summary>b站视频链接解析</summary>

  `import _ "github.com/FloatTech/ReiBot-Plugin/plugin/bilibili_parse"`

  - [x] https://www.bilibili.com/video/BV1xx411c7BF | https://www.bilibili.com/video/av1605 | https://b23.tv/I8uzWCA | https://www.bilibili.com/video/bv1xx411c7BF

</details>

<details>
  <summary>英文字符翻转</summary>

  `import _ "github.com/FloatTech/ReiBot-Plugin/plugin/chrev"`

  - [x] 翻转 I love you

</details>

<details>
  <summary>合成emoji</summary>

  `import _ "github.com/FloatTech/ReiBot-Plugin/plugin/emojimix"`

  - [x] [emoji][emoji]

</details>

<details>
  <summary>每日运势</summary>

  `import _ "github.com/FloatTech/ReiBot-Plugin/plugin/fortune"`

  - [x] 运势 | 抽签

  - [x] 设置底图[车万 DC4 爱因斯坦 星空列车 樱云之恋 富婆妹 李清歌 公主连结 原神 明日方舟 碧蓝航线 碧蓝幻想 战双 阴阳师 赛马娘 东方归言录 奇异恩典 夏日口袋 ASoul]

</details>

<details>
  <summary>原神抽卡</summary>

  `import _ "github.com/FloatTech/ReiBot-Plugin/plugin/genshin"`

  - [x] 切换原神卡池

  - [x] 原神十连

</details>

<details>
  <summary>百人一首</summary>

  `import _ "github.com/FloatTech/ReiBot-Plugin/plugin/hyaku"`

  - [x] 百人一首

  - [x] 百人一首之n

</details>

<details>
  <summary>lolicon</summary>

  `import _ "github.com/FloatTech/ReiBot-Plugin/plugin/lolicon"`

  - [x] 来份萝莉

</details>

<details>
  <summary>bot管理相关</summary>

  `import _ "github.com/FloatTech/ReiBot-Plugin/plugin/manager"`

  - [x] /离开 (ChatID)

</details>

<details>
  <summary>日韩 VITS 模型拟声</summary>

  `import _ "github.com/FloatTech/ReiBot-Plugin/plugin/moegoe"`

  - [x] 让[宁宁|爱瑠|芳乃|茉子|丛雨|小春|七海]说(日语)

  - [x] 让[수아|미미르|아린|연화|유화|선배]说(韩语)

</details>

<details>
  <summary>NovelAI作画</summary>

  `import _ "github.com/FloatTech/ReiBot-Plugin/plugin/novelai"`

  - [x] novelai作图 tag1 tag2...

	- [x] 设置 novelai key [key]

</details>

<details>
  <summary>在线代码运行</summary>

  `import _ "github.com/FloatTech/ReiBot-Plugin/plugin/runcode"`

  - [x] >runcode [language] help

  - [x] >runcode [language] [code block]

  - [x] >runcoderaw [language] [code block]

</details>

<details>
  <summary>搜图</summary>

  `import _ "github.com/FloatTech/ReiBot-Plugin/plugin/saucenao"`

  - [x] 以图搜图 | 搜索图片 | 以图识图[图片]

  - [x] 搜图[P站图片ID]

  - [x] 设置 saucenao api key [apikey]

</details>

<details>
  <summary>搜番</summary>

  `import _ "github.com/FloatTech/ReiBot-Plugin/plugin/tracemoe"`

  - [x] 搜番 | 搜索番剧[图片]

</details>

## 特别感谢

- [ZeroBot](https://github.com/wdvxdr1123/ZeroBot)
