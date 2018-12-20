package worker

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	EtcdEndpoints   []string `json:"etcdEndpoints"`
	EtcdDialTimeout int      `json:"etcdDialTimeout"`
}

var (
	G_config *Config
)

func InitConfig(filename string) (err error) {
	var (
		content []byte
		conf    Config
	)
	//读取配置
	if content, err = ioutil.ReadFile(filename); err != nil {
		return
	}

	//json 返序列
	if err = json.Unmarshal(content, &conf); err != nil {
		return
	}
	//单利
	G_config = &conf
	// fmt.Println(G_config)
	return
}
