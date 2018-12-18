package master

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	ApiPort         int      `json:"apiPort"`
	ApiReadTimeout  int      `json:"apiReadTimeout"`
	ApiWriteTimeout int      `json:"apiWriteTimeout"`
	EtcdEndpoints   []string `json:"etcdEndpoints"`
	EtcdDialTimeout int      `json:"etcdDialTimeout"`
	WebRoot         string   `json:"webroot"`
	// MongodbUri            string   `json:"mongodbUri"`
	// MongodbConnectTimeout int      `json:"mongodbConnectTimeout"`
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
