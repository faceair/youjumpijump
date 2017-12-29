# 微信跳一跳 AI

思路和原理参考 https://github.com/wangshub/wechat_jump_game [自动跳跃算法细节参考](https://github.com/faceair/wechat_jump_game/blob/master/wechat_jump.py#L50)

用 Golang 重新实现是期望跨平台，方便打包给普通用户使用。代码逻辑精简过，运行速度会有提升。

## 使用须知

1. Android 手机一台，电脑上安装 ADB，连接上电脑后开启调试模式
2. 进入微信打开微信跳一跳，点击开始游戏
2. 运行本 AI

## 跳跃系数

跟系统分辨率有关，但具体换算关系还不太清楚。下面两份数据供参考。
1. 我的手机系统分辨率 720*1280  跳跃系数 2.04
2. 原作者手机分辨率是 1920*1080 跳跃系数是 1.35

## 下载

1. Windows [下载地址](https://github.com/faceair/youjumpijump/releases/latest)
2. Linux [下载地址](https://github.com/faceair/youjumpijump/releases/latest)
3. MacOS [下载地址](https://github.com/faceair/youjumpijump/releases/latest)

## 实验结果

![](http://ww3.sinaimg.cn/large/0060lm7Tly1fmy1dpozipj30k00zkq46.jpg)
