# 微信跳一跳 AI

思路和原理参考 https://github.com/wangshub/wechat_jump_game [自动跳跃算法细节参考](https://github.com/faceair/wechat_jump_game/blob/master/wechat_jump.py#L50)

用 Golang 重新实现是期望跨平台，方便打包给普通用户使用。代码逻辑精简过，运行速度会有提升。

## 使用须知

1. Android 手机一台，电脑上安装 ADB，连接上电脑后开启调试模式
2. 进入微信打开微信跳一跳，点击开始游戏
2. 运行本 AI

## 跳跃系数

目前推荐设为 2.04，截图后会先 resize 成 720p 的图片然后再找点和跳跃。

## 下载

1. Windows [下载地址](https://github.com/faceair/youjumpijump/releases/latest)
2. Linux [下载地址](https://github.com/faceair/youjumpijump/releases/latest)
3. MacOS [下载地址](https://github.com/faceair/youjumpijump/releases/latest)

## 实验结果

![](http://ww3.sinaimg.cn/large/0060lm7Tly1fmy1dpozipj30k00zkq46.jpg)

## FAQ

1. MacOS 和 Linux 下怎么运行这个程序？

在终端下运行 `chmod +x ./jumpAI-darwin-amd64 && ./jumpAI-darwin-amd64`

2. 为什么我的手机没有反应？

请手动执行 `adb shell input swipe 320 410 320 410 300` 看是否有反应，否则可能是手机设置问题，比如开发者选项中模拟点击是否有打开

3. 为什么有时候跳很远 & 跳不远？

不要动数据线，接触不良可能导致 ADB 命令执行失败。另外可以将这一跳的截图通过下面方式反馈。

3. 其他疑难杂症

可以开新 issue 或加 QQ 群反馈 684623076，最好能附上截图。
