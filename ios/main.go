package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"image"
	"image/png"
	"log"
	"runtime/debug"
	"time"

	jump "github.com/faceair/youjumpijump"
)

var r = jump.NewRequest()

type ScreenshotRes struct {
	Value     string `json:"value"`
	SessionID string `json:"sessionId"`
	Status    int    `json:"status"`
}

func screenshot(ip string) (*ScreenshotRes, image.Image) {
	_, body, err := r.Get(fmt.Sprintf("http://%s/screenshot", ip))
	if err != nil {
		log.Fatal("WebDriverAgentRunner 连接失败，请参考 https://github.com/faceair/youjumpijump/issues/71")
	}

	res := new(ScreenshotRes)
	err = json.Unmarshal(body, res)
	if err != nil {
		log.Fatal("WebDriverAgentRunner 响应数据异常，请检查 WebDriverAgentRunner 运行状态")
	}

	pngValue, err := base64.StdEncoding.DecodeString(res.Value)
	if err != nil {
		log.Fatal("图片解码失败，请参考 https://github.com/faceair/youjumpijump/issues/41")
	}

	src, err := png.Decode(bytes.NewReader(pngValue))
	if err != nil {
		log.Fatal("图片解码失败，请参考 https://github.com/faceair/youjumpijump/issues/41")
	}
	return res, src
}

func main() {
	defer func() {
		jump.Debugger()
		if e := recover(); e != nil {
			log.Printf("%s: %s", e, debug.Stack())
			fmt.Print("程序已崩溃，按任意键退出")
			var c string
			fmt.Scanln(&c)
		}
	}()

	var ip string
	var inputRatio float64
	flag.StringVar(&ip, "ip", "", "WebDriverAgentRunner 监听的 IP 和端口 (例如 192.168.9.94:8100)")
	flag.Float64Var(&inputRatio, "ratio", 0, "跳跃系数(推荐值 2，可适当调整)")
	flag.Parse()

	if ip == "" {
		fmt.Print("请输入 WebDriverAgentRunner 监听的 IP 和端口 (例如 192.168.9.94:8100):")
		_, err := fmt.Scanln(&ip)
		if err != nil {
			log.Fatal("WebDriverAgentRunner 连接失败，请参考 https://github.com/faceair/youjumpijump/issues/71")
		}
	}
	if inputRatio == 0 {
		fmt.Print("请输入跳跃系数(推荐值 2.04，可适当调整):")
		_, err := fmt.Scanln(&inputRatio)
		if err != nil {
			log.Print("未输入跳跃系数，将采用默认跳跃系数 2.04")
			inputRatio = 2.04
		}
	}

	for {
		jump.Debugger()

		res, pic := screenshot(ip)
		go jump.SavePNG("jump.png", pic)

		start, end := jump.Find(pic)
		if start == nil {
			log.Fatal("找不到起点，请把 debugger 目录打包发给开发者检查问题。")
		} else if end == nil {
			log.Fatal("找不到落点，请把 debugger 目录打包发给开发者检查问题。")
		}

		distance := jump.Distance(start, end)
		log.Printf("from:%v to:%v distance:%.2f press:%.2fms ", start, end, distance, distance*inputRatio)

		_, _, err := r.PostJSON(fmt.Sprintf("http://%s/session/%s/wda/touchAndHold", ip, res.SessionID), map[string]interface{}{
			"x":        jump.Random(100, 400),
			"y":        jump.Random(100, 400),
			"duration": distance * inputRatio / 1000,
		})
		if err != nil {
			log.Fatal("WebDriverAgentRunner 连接失败，请参考 https://github.com/faceair/youjumpijump/issues/71")
		}

		time.Sleep(time.Millisecond * 1200)
	}
}
