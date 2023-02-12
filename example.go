package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/lristar/go-viper/configs"
)

var (
	configFileName = flag.String("c", "./configs.yaml", "设置配置文件地址")
)

func init() {
	//os.Setenv("SETTINGS_APPLICATION_MODE", "test")
	//os.Setenv("SETTINGS_APPLICATION_HOST", "127.0.0.10")
	//os.Setenv("SETTINGS_APPLICATION_NAME", "APP")
	//os.Setenv("SETTINGS_APPLICATION_PORT", "9001")
	//os.Setenv("SETTINGS_APPLICATION_READTIMEOUT", "1")
	//os.Setenv("SETTINGS_APPLICATION_DEMOMSG", "ADAFDSF")
}

var ApplicationConfig = new(Application)

type Application struct {
	ReadTimeout   int    `json:"readtimeout"`
	WriterTimeout int    `json:"writertimeout"`
	Host          string `json:"host"`
	Port          int64  `json:"port"`
	Name          string `json:"name"`
	Mode          string `json:"mode"`
	DemoMsg       string `json:"demomsg"`
}

// Config 配置集合
type Config struct {
	Application *Application `json:"application"`
}

type Settings struct {
	Settings Config `json:"settings"`
}

var _cfg *Settings

func main() {
	_ = configs.Setup(*configFileName, &_cfg)
	marshal, err := json.Marshal(_cfg)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("序列化config: \n", string(marshal))
}
