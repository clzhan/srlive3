package conf

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/sdming/gosnow"
)

var AppConf struct {
	RtmpPort       int `json:"RtmpPort"`
	HttpPort       int `json:"HttpPort"`
	LogPath        string `json:"LogPath"`
	LogLvl         int    `json:"LogLvl"`
	Srvid          int    `json:"Srvid"`
	PProf          bool   `json:"PProf"`
}

var (
	Snow *gosnow.SnowFlake
)

func Init() {

	data, err := ioutil.ReadFile("./conf/server.json")
	if err != nil {
		fmt.Printf("%v", err)
		os.Exit(1)
	}
	err = json.Unmarshal(data, &AppConf)
	if err != nil {
		fmt.Printf("%v", err)
		os.Exit(2)
	}

	Snow, err = gosnow.NewSnowFlake(uint32(AppConf.Srvid))
	if err != nil {
		fmt.Printf("%v", err)
		os.Exit(3)
	}

	fmt.Println("conf : ", AppConf)
}
