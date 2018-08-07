package config

import (
	"github.com/cihub/seelog"
	"github.com/fanyang1988/eosforce-batch-tool/utils"
)

// CommonCfg cfg value from config.json
type CommonCfgTyp struct {
	ApiUrls   []string `json:"api-urls"`
	WalletURL string   `json:"walletURL"`
}

var CommonCfg CommonCfgTyp

func init() {
	err := utils.LoadJsonFile("./config.json", &CommonCfg)
	if err != nil {
		seelog.Errorf("load cfg err by %s", err.Error())
		panic(err)
	}
}
