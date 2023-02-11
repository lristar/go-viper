package main

import (
	"encoding/json"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"os"
	"strings"
)

func init() {
	os.Setenv("SETTINGS_APPLICATION_MODE", "test")
	os.Setenv("SETTINGS_APPLICATION_HOST", "127.0.0.10")
	os.Setenv("SETTINGS_APPLICATION_NAME", "APP")
	os.Setenv("SETTINGS_APPLICATION_PORT", "9001")
	os.Setenv("SETTINGS_APPLICATION_READTIMEOUT", "1")
	os.Setenv("SETTINGS_APPLICATION_DEMOMSG", "ADAFDSF")
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

func setTag(m *mapstructure.DecoderConfig) {
	m.TagName = "json"
}

// Setup 载入配置文件
func Setup(configFile string, fs ...func()) {
	_cfg = &Settings{
		//Settings: Config{
		//	Application: ApplicationConfig,
		//},
	}
	v := viper.New()
	//自动获取全部的env加入到viper中。（如果环境变量多就全部加进来）默认别名和环境变量名一致
	v.AutomaticEnv()
	//将加入的环境变量*_*_格式替换成 *.*格式
	//（因为从环境变量读是按"a.b.c.d"的格式读取，所以要给在viper维护一个别名对象，给环境变量一个别名）
	v.SetEnvKeyReplacer(strings.NewReplacer("_", "."))
	SetEnvToViper(v)
	//配置文件位置
	v.SetConfigFile(configFile)
	// 获取远程配置文件
	v.AddRemoteProvider("etcd3", "http://127.0.0.1:4001", "/config/hugo.json")
	v.ReadRemoteConfig()
	//支持 "yaml", "yml", "json", "toml", "hcl", "tfvars", "ini", "properties", "props", "prop", "dotenv", "env":
	v.SetConfigType("yaml")

	//读文件到viper配置中
	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	// 系列化成config对象
	if err := v.Unmarshal(&_cfg, setTag); err != nil {
		fmt.Println(err)
	}

	fmt.Println("所有的config: \n", v.AllSettings())
	fmt.Println(v.Get("SETTINGS_APPLICATION_READTIMEOUT"), "------fdafdaf")
}

func main() {
	Setup("configs.yaml")
	marshal, err := json.Marshal(_cfg)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("序列化config: \n", string(marshal))
}

func SetEnvToViper(v *viper.Viper) {
	keys := os.Environ()
	for i := range keys {
		cache := strings.Split(keys[i], "=")
		if strings.Contains(cache[0], "PATH") {
			continue
		}
		if len(cache) > 1 {
			v.Set(strings.ReplaceAll(strings.ToLower(cache[0]), "_", "."), cache[1])
		}
	}
}
