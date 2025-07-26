package config

import (
	"bytes"
	"crypto/rand"
	_ "embed"
	"encoding/hex"
	"log"
	"os"

	model "github.com/lin-snow/ech0/internal/model/common"
	"github.com/spf13/viper"
)

// Config 全局配置变量
var Config AppConfig

// JWT_SECRET 用于JWT签名的密钥
var JWT_SECRET []byte

// AppConfig 应用程序配置结构体
type AppConfig struct {
	Server struct {
		Port string `yaml:"port"`
		Host string `yaml:"host"`
		Mode string `yaml:"mode"`
	} `yaml:"server"`
	Database struct {
		Type string `yaml:"type"`
		Path string `yaml:"path"`
		// Pragma string `yaml:"pragma"`
	} `yaml:"database"`
	Auth struct {
		Jwt struct {
			Expires  int    `yaml:"expires"`
			Issuer   string `yaml:"issuer"`
			Audience string `yaml:"audience"`
		} `yaml:"jwt"`
	} `yaml:"auth"`
	Upload struct {
		ImageMaxSize int      `yaml:"imagemaxsize"`
		AudioMaxSize int      `yaml:"audiomaxsize"`
		AllowedTypes []string `yaml:"allowedtypes"`
		ImagePath    string   `yaml:"imagepath"`
		AudioPath    string   `yaml:"audiopath"`
	} `yaml:"upload"`
	Setting struct {
		SiteTitle     string `yaml:"sitetitle"`
		Servername    string `yaml:"servername"`
		Serverurl     string `yaml:"serverurl"`
		AllowRegister bool   `yaml:"allowregister"`
		Icpnumber     string `yaml:"icpnumber"`
		MetingAPI     string `yaml:"metingapi"`
		CustomCSS     string `yaml:"customcss"`
		CustomJS      string `yaml:"customjs"`
	} `yaml:"setting"`
	Comment struct {
		EnableComment bool   `yaml:"enablecomment"`
		Provider      string `yaml:"provider"`   // 评论提供者
		CommentAPI    string `yaml:"commentapi"` // 评论 API 地址
	} `yaml:"comment"`
	SSH struct {
		Port string `yaml:"port"` // SSH 端口
		Host string `yaml:"host"` // SSH 主机地址
		Key  string `yaml:"key"`  // SSH 私钥路径
	} `yaml:"ssh"`
}

//go:embed config.yaml
var configData []byte

// LoadAppConfig 加载应用程序配置
func LoadAppConfig() {
	// viper.SetConfigFile("config/config.yaml")
	viper.SetConfigType("yaml")
	// 使用嵌入的配置数据而不是从文件系统读取
	err := viper.ReadConfig(bytes.NewReader(configData))
	if err != nil {
		panic(model.READ_CONFIG_PANIC + ":" + err.Error())
	}

	// 将配置文件内容反序列化到结构体 Config 中
	err = viper.Unmarshal(&Config)
	if err != nil {
		panic(model.READ_CONFIG_PANIC + ":" + err.Error())
	}

	JWT_SECRET = GetJWTSecret()
}

// GetJWTSecret 加载JWT密钥
func GetJWTSecret() []byte {
	// 从环境变量中获取JWT密钥
	secret := os.Getenv("JWT_SECRET")
	if secret == "" { // 如果没有设置环境变量，则使用UUID生成默认密钥
		b := make([]byte, 16)
		_, err := rand.Read(b)
		if err != nil {
			log.Fatal("failed to generate random JWT secret:", err)
		}
		secret = hex.EncodeToString(b)
	}

	return []byte(secret)
}
