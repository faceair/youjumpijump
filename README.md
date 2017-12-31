# 微信跳一跳 AI

思路和原理参考 https://github.com/wangshub/wechat_jump_game [自动跳跃算法细节参考](https://github.com/faceair/wechat_jump_game/blob/master/wechat_jump.py#L50)

用 Golang 重新实现是期望跨平台，方便打包给普通用户使用。代码逻辑精简过，运行速度会有提升。

## 为什么要将文件 Push 到 Android 手机内执行？

我们发现会有偶尔的情况下定位的关键点都是准的，ADB 命令也执行了，但是就是没跳过去；有 MIUI 用户的截图被自动重命名了；也有朋友跑着跑着突然 ADB 报错了，无法继续执行了；同时也有朋友使用的电脑上的模拟器，非常的稳定。所以我们推测 ADB 命令在某些情况下会有问题，可能也跟连接的线材有关。所以我们决定将程序移植到 Android 上，直接在 Android 上运行，也可以避开手机 USB 调试的一些安全设置。实验证明，确实稳定多了。

## 下载地址

Android [下载地址](https://github.com/faceair/youjumpijump/releases/latest)

## 使用须知

1. Android 手机一台，电脑上安装 ADB，连接上电脑后开启 USB 调试模式
2. 进入微信打开微信跳一跳，点击开始游戏
3. 将下载的文件 Push 到手机上 `adb push ./youjumpijump /data/local/tmp/ && adb shell`
4. 跑起来 `cd /data/local/tmp/ && chmod +x ./youjumpijump && ./youjumpijump`

## 跳跃系数

目前推荐设为 2.04，截图后会先 resize 成 720p 的图片然后再找点和跳跃。

## FAQ

1. 怎么安装 ADB？怎么开启手机的 USB 调试？为什么电脑连不上手机？为什么我 Push 时说找不到文件？这些命令怎么执行？

请自行搜索解决，实在连不上可以[用模拟器](https://github.com/wangshub/wechat_jump_game/tree/master/%E6%96%B0%E6%89%8B%E5%B0%8F%E7%99%BD%E8%AF%B7%E4%BD%BF%E7%94%A8%E8%BF%99%E4%B8%AA%E4%BB%A3%E7%A0%81%20%20%E4%B8%8D%E9%9C%80%E8%A6%81%E4%BD%BF%E7%94%A8%E7%9C%9F%E6%9C%BA%E7%9A%84%E4%B8%93%E7%94%A8%E4%BB%A3%E7%A0%81)（不保证安全性，自行承担风险）。

2. 锤子手机运行异常？

关大爆炸功能。

3. 怎么编译 Android 版本？

`CGO_ENABLED=1 GOARCH=arm GOOS=linux go build .`

4. 其他疑难杂症？

开新 issue 请并附上日志 `adb pull /data/local/tmp/debugger` 供我们排查。

## 实验结果

![](http://ww3.sinaimg.cn/large/0060lm7Tly1fmy1dpozipj30k00zkq46.jpg)
