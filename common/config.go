package common

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v1"
)

type Config struct {
	Common struct {
		Version  string
		IsDebug  bool `yaml:"debug"`
		LogPath  string
		LogLevel string
		Service  string
		RealIp   string
	}

	Api struct {
		GatewayPort   string
		ServerID      int
		ApiUpdatePort string
	}

	Admin struct {
		ManagerPort string
	}

	Mysql struct {
		Addr     string
		Port     string
		Database string
		Acc      string
		Pw       string
	}

	Etcd struct {
		Addrs     []string
		ServerKey string
	}
}

var Conf = &Config{}

func InitConfig() {
	data, err := ioutil.ReadFile("openapi.yaml")
	if err != nil {
		log.Fatal("read config error :", err)
	}

	err = yaml.Unmarshal(data, &Conf)
	if err != nil {
		log.Fatal("yaml decode error :", err)
	}

	log.Println(Conf)
}
