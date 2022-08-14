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
  <summary>b站视频链接解析</summary>

  `import _ "github.com/FloatTech/ReiBot-Plugin/plugin/bilibili_parse"`

  - [x] https://www.bilibili.com/video/BV1xx411c7BF | https://www.bilibili.com/video/av1605 | https://b23.tv/I8uzWCA | https://www.bilibili.com/video/bv1xx411c7BF

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
  <summary>日韩 VITS 模型拟声</summary>

  `import _ "github.com/FloatTech/ReiBot-Plugin/plugin/moegoe"`

  - [x] 让[宁宁|爱瑠|芳乃|茉子|丛雨|小春|七海]说(日语)

  - [x] 让[수아|미미르|아린|연화|유화|선배]说(韩语)

</details>

## 特别感谢

- [ZeroBot](https://github.com/wdvxdr1123/ZeroBot)
