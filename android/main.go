package main

import (
	"bytes"
	"flag"
	"fmt"
	"image"
	"image/png"
	"log"
	"os"
	"os/exec"
	"runtime/debug"
	"strconv"
	"time"

	jump "github.com/faceair/youjumpijump"
)

func screenshot() image.Image {
	out, err := exec.Command("/system/bin/screencap", "-p").Output()
	if err != nil {
		log.Fatal("截图失败，请检查开发者选项中的 USB 调试安全设置是否打开。")
	}
	pic, err := png.Decode(bytes.NewReader(out))
	if err != nil {
		log.Fatal("PNG 截图解码失败。")
	}
	return pic
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

	var inputRatio float64
	flag.Float64Var(&inputRatio, "ratio", 0, "跳跃系数(推荐值 2.04，可适当调整)")
	flag.Parse()

	var err error
	if inputRatio == 0 {
		fmt.Print("请输入跳跃系数(推荐值 2.04，可适当调整):")
		_, err = fmt.Scanln(&inputRatio)
		if err != nil {
			log.Print("未输入跳跃系数，将采用默认跳跃系数 2.04")
			inputRatio = 2.04
		}
	}

	var pic image.Image

	for {
		jump.Debugger()

		if len(os.Getenv("DEBUG")) > 0 {
			pic, err = jump.OpenPNG("jump.png")
			if err != nil {
				log.Fatal("jump.png 图片打开失败")
			}
		} else {
			pic = screenshot()
			go jump.SavePNG("jump.png", pic)
		}

		start, end := jump.Find(pic)
		if start == nil {
			log.Fatal("找不到起点，请使用 adb pull /data/local/tmp/debugger 把 debugger 目录拉下来打包发给开发者检查问题。")
		} else if end == nil {
			log.Fatal("找不到落点，请使用 adb pull /data/local/tmp/debugger 把 debugger 目录拉下来打包发给开发者检查问题。")
		}

		distance := jump.Distance(start, end)
		log.Printf("from:%v to:%v distance:%.2f press:%.2fms ", start, end, distance, distance*inputRatio)

		touchX, touchY := strconv.Itoa(jump.Random(100, 400)), strconv.Itoa(jump.Random(100, 400))
		_, err = exec.Command("/system/bin/sh", "/system/bin/input", "swipe", touchX, touchY, touchX, touchY,
			strconv.Itoa(int(distance*inputRatio))).Output()
		if err != nil {
			log.Fatal("模拟触摸失败，请检查开发者选项中的 USB 调试安全设置是否打开。")
		}

		time.Sleep(time.Millisecond * 1200)
	}
}
