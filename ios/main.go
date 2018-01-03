package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	"image/png"
	"log"
	"os"
	"runtime/debug"
	"time"

	jump "github.com/faceair/youjumpijump"
)

var similar *jump.Similar

var r = jump.NewRequest()

type ScreenshotRes struct {
	Value     string `json:"value"`
	SessionID string `json:"sessionId"`
	Status    int    `json:"status"`
}

func screenshot(ip string) (*ScreenshotRes, image.Image) {
	_, body, err := r.Get(fmt.Sprintf("http://%s/screenshot", ip))
	if err != nil {
		panic(err)
	}

	res := new(ScreenshotRes)
	err = json.Unmarshal(body, res)
	if err != nil {
		panic(err)
	}

	pngValue, err := base64.StdEncoding.DecodeString(res.Value)
	if err != nil {
		panic(err)
	}

	src, err := png.Decode(bytes.NewReader(pngValue))
	if err != nil {
		panic(err)
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
	fmt.Print("请输入 WebDriverAgentRunner 监听的 IP 和端口 (例如 192.168.9.94:8100):")
	_, err := fmt.Scanln(&ip)
	if err != nil {
		log.Fatal(err)
	}

	var inputRatio float64
	fmt.Print("请输入跳跃系数(推荐值 2，可适当调整):")
	_, err = fmt.Scanln(&inputRatio)
	if err != nil {
		log.Printf("未输入跳跃系数，将采用默认跳跃系数2")
		inputRatio = 2
	}

	similar = jump.NewSimilar(inputRatio)

	for {
		jump.Debugger()

		res, src := screenshot(ip)

		f, _ := os.OpenFile("jump.png", os.O_WRONLY|os.O_CREATE, 0600)
		png.Encode(f, src)
		f.Close()

		start, end := jump.Find(src)
		if start == nil {
			log.Print("找不到起点，请把 debugger 目录打包发给开发者检查问题。")
			break
		} else if end == nil {
			log.Print("找不到落脚点，请把 debugger 目录打包发给开发者检查问题。")
			break
		}

		nowDistance := jump.Distance(start, end)
		similarDistance, nowRatio := similar.Find(nowDistance)

		log.Printf("from:%v to:%v distance:%.2f similar:%.2f ratio:%v press:%.2fms ", start, end, nowDistance, similarDistance, nowRatio, nowDistance*nowRatio)

		_, _, err = r.PostJSON(fmt.Sprintf("http://%s/session/%s/wda/touchAndHold", ip, res.SessionID), map[string]interface{}{
			"x":        200,
			"y":        200,
			"duration": nowDistance * nowRatio / 1000,
		})
		if err != nil {
			panic(err)
		}

		go func() {
			time.Sleep(time.Duration(nowDistance*nowRatio/1000+50) * time.Millisecond)
			_, src := screenshot(ip)
			finally, _ := jump.Find(src)

			f, _ = os.OpenFile("jump.test.png", os.O_WRONLY|os.O_CREATE, 0600)
			png.Encode(f, src)
			f.Close()

			if finally != nil {
				finallyDistance := jump.Distance(start, finally)
				finallyRatio := (nowDistance * nowRatio) / finallyDistance

				if finallyRatio > nowRatio/2 && finallyRatio < nowRatio*2 {
					similar.Add(finallyDistance, finallyRatio)
				}
			}
		}()

		time.Sleep(time.Millisecond * 1500)
	}
}
