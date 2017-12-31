package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image/png"
	"log"
	"math"
	"os"
	"runtime/debug"
	"time"

	jump "github.com/faceair/youjumpijump"
)

type ScreenshotRes struct {
	Value     string `json:"value"`
	SessionID string `json:"sessionId"`
	Status    int    `json:"status"`
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

	var ratio float64
	fmt.Print("请输入跳跃系数(推荐值 2，可适当调整):")
	_, err = fmt.Scanln(&ratio)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("现在跳跃系数是 %f", ratio)

	r := jump.NewRequest()

	for {
		jump.Debugger()

		_, body, err := r.Get(fmt.Sprintf("http://%s/screenshot", ip))
		if err != nil {
			panic(err)
		}

		scr := new(ScreenshotRes)
		err = json.Unmarshal(body, scr)
		if err != nil {
			panic(err)
		}

		pngValue, err := base64.StdEncoding.DecodeString(scr.Value)
		if err != nil {
			panic(err)
		}

		src, err := png.Decode(bytes.NewReader(pngValue))
		if err != nil {
			panic(err)
		}

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

		ms := float64(math.Pow(math.Pow(float64(start[0]-end[0]), 2)+math.Pow(float64(start[1]-end[1]), 2), 0.5) * ratio)
		log.Printf("from:%v to:%v press:%vms", start, end, ms)

		_, _, err = r.PostJSON(fmt.Sprintf("http://%s/session/%s/wda/touchAndHold", ip, scr.SessionID), map[string]interface{}{
			"x":        200,
			"y":        200,
			"duration": ms / 1000,
		})
		if err != nil {
			panic(err)
		}

		time.Sleep(time.Millisecond * 1500)
	}
}
