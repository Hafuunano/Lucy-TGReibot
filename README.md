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

  `import _ "github.com/FloatTech/ZeroBot-Plugin/plugin/b14"`

  - [x] 加密xxx

  - [x] 解密xxx

  - [x] 用yyy加密xxx

  - [x] 用yyy解密xxx

</details>

<details>
  <summary>lolicon</summary>

  `import _ "github.com/FloatTech/ReiBot-Plugin/plugin/lolicon"`

  - [x] 来份萝莉

</details>

## 特别感谢

- [ZeroBot](https://github.com/wdvxdr1123/ZeroBot)
