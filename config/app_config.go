package config

import (
	"io"
	"log"
	"mySparkler/pkg/file"
	"mySparkler/pkg/utils/R"
	"os"
	"path"

	"gopkg.in/yaml.v2"
)

var Server *server
var Database *database
var Redis *redis
var Jwt *jwt
var XxlJob *xxlJob
var LogConfig *logConfig

type conf struct {
	Svc         server    `yaml:"server"`
	DB          database  `yaml:"database"`
	RedisConfig redis     `yaml:"redis"`
	Jwt         jwt       `yaml:"jwt"`
	XxlJob      xxlJob    `yaml:"xxl-job"`
	LogConfig   logConfig `yaml:"log"`
}

type server struct {
	Port           int    `yaml:"port"`
	RunMode        string `yaml:"runMode"`
	LogLevel       string `yaml:"logLevel"`
	EnabledSwagger bool   `yaml:"enabledSwagger"`
}

type database struct {
	Type            string `yaml:"type"`
	Host            string `yaml:"host"`
	Port            string `yaml:"port"`
	UserName        string `yaml:"username"`
	Password        string `yaml:"password"`
	DbName          string `yaml:"dbname"`
	DbFileName      string `yaml:"db_file_name"`
	DbPath          string `yaml:"db_path"`
	MaxIdleConn     int    `yaml:"max_idle_conn"`
	MaxOpenConn     int    `yaml:"max_open_conn"`
	ConnMaxLifetime int    `yaml:"conn_max_lifetime"`
}

type redis struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

type jwt struct {
	Secret string `yaml:"secret"`
	JwtTtl int    `yaml:"jwt_ttl"`
}

type xxlJob struct {
	Enabled          bool   `yaml:"enabled"`
	Env              string `yaml:"env"`
	AdminAddress     string `yaml:"admin_address"`
	AccessToken      string `yaml:"access_token"`
	AppName          string `yaml:"app_name"`
	Address          string `yaml:"address"`
	Ip               string `yaml:"ip"`
	Port             int    `yaml:"port"`
	LogPath          string `yaml:"log_path"`
	LogRetentionDays int    `yaml:"log_retention_days"`
	HttpTimeout      int    `yaml:"http_timeout"`
}

type logConfig struct {
	Enabled  bool     `yaml:"enabled"`
	LogMode  string   `yaml:"logMode"`
	FilePath string   `yaml:"filePath"`
	Filtered []string `yaml:"filtered"`
}

func InitAppConfig(dataFile string, ymlDefault string) {

	filePath := path.Join(GetAppPath()+"/", dataFile)
	_, err := os.Stat(filePath)
	if err != nil {
		log.Printf("config file path %s not exist", filePath)
		out, err := os.Create(filePath)
		if err != nil {
			// return "", err
			panic(R.ReturnFailMsg("Marshal:" + err.Error()))
		}

		// 然后将响应流和文件流对接起来
		_, err = io.WriteString(out, ymlDefault)
		if err != nil {
			// return "", err
			panic(R.ReturnFailMsg("Marshal:" + err.Error()))
		}
		defer out.Close()
		c := new(conf)
		data, err := yaml.Marshal(&c)
		if err != nil {
			log.Printf("Marshal: %v", err)
			panic(R.ReturnFailMsg("Marshal:" + err.Error()))
		}
		log.Printf("Marshal: %v", data)
		// panic(R.ReturnFailMsg("config file path " + filePath + " not exist"))
	}
	yamlFile, err := os.ReadFile(filePath)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
		panic(R.ReturnFailMsg("yamlFile.Get err   " + err.Error()))
	}
	c := new(conf)
	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		log.Printf("Unmarshal: %v", err)
		panic(R.ReturnFailMsg("Unmarshal:" + err.Error()))
	}
	log.Printf("load conf success")
	// 绑定到外部可以访问的变量中
	Server = &c.Svc
	Database = &c.DB
	Redis = &c.RedisConfig
	Jwt = &c.Jwt
	XxlJob = &c.XxlJob
	LogConfig = &c.LogConfig
}

// getAppPath 获取应用主目录
func GetAppPath() string {
	//获取我的文档目录
	return file.GetAppPath()
}

// pathExist 判断文件目录是否存在，不存在创建
func PathExist(path string) string {
	_, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		_ = os.MkdirAll(path, os.ModePerm)
	}
	return path
}
