package jump

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

var basePath string

func init() {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	basePath = filepath.Dir(ex)

	if ok, _ := Exists(basePath + "/debugger"); !ok {
		os.MkdirAll(basePath+"/debugger", os.ModePerm)
	}

	os.Remove(basePath + "/debugger/debug.log")
	logFile, _ := os.OpenFile(basePath+"/debugger/debug.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	log.SetOutput(io.MultiWriter(os.Stdout, logFile))
}

func Debugger() {
	if ok, _ := Exists(basePath + "/jump.png"); ok {
		os.Rename(basePath+"/jump.png", basePath+"/debugger/"+strconv.Itoa(TimeStamp())+".png")

		files, err := ioutil.ReadDir(basePath + "/debugger/")
		if err != nil {
			panic(err)
		}

		for _, f := range files {
			fname := f.Name()
			ext := filepath.Ext(fname)
			name := fname[0 : len(fname)-len(ext)]
			if ts, err := strconv.Atoi(name); err == nil {
				if TimeStamp()-ts > 20 {
					os.Remove(basePath + "/debugger/" + fname)
				}
			}
		}
	}
}

func Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func TimeStamp() int {
	return int(time.Now().UnixNano() / int64(time.Second))
}
