package batchCmd

import (
	"fmt"
	"github.com/cihub/seelog"
	"github.com/fanyang1988/eos-go"
	"github.com/fanyang1988/eos-go/eosforce"
	"github.com/fanyang1988/eosc/eosc/cmd"
	"github.com/fanyang1988/eosforce-batch-tool/config"
	"github.com/fanyang1988/eosforce-batch-tool/utils"
	"github.com/spf13/cobra"
)

func Transfer(idx int, api *eos.API, from, to, quantityStr, memo string) error {
	quantity, err := eos.NewEOSAssetFromString(quantityStr)
	utils.ErrorCheck(err, "invalid amount %s", quantityStr)

	action := eosforce.NewTransfer(
		toAccount(from, "from"),
		toAccount(to, "to"),
		quantity, fmt.Sprintf("id:%v,from:%v,memo:%v", idx, from, memo))

	// in eosforce the sys token is use `eosio.transfer` in System to transfer coin
	action.Account = eos.AN("eosio")
	// action.Account = toAccount(viper.GetString("transfer-cmd-contract"), "--contract")

	seelog.Infof("Send Transfer %s, %s, %s, %s to %s", from, to, quantityStr, memo, api.BaseURL)

	return pushEOSCActions(api, action)
}

var batchTransferCmd = &cobra.Command{
	Use:   "batch transfer [json_file]",
	Short: "Batch transfer from tokens from an account to another",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		var transferCmds = struct {
			Cmds []struct {
				From     string `json:"from"`
				To       string `json:"to"`
				Quantity string `json:"quantity"`
				Memo     string `json:"memo"`
			} `json:"cmds"`
		}{}

		err := utils.LoadJsonFile(args[1], &transferCmds)
		utils.ErrorCheck(err, "load file %s err by", args[1])

		apis := getAPIs(config.CommonCfg.ApiUrls)

		for idx, data := range transferCmds.Cmds {
			err := Transfer(idx, apis[idx%len(apis)], data.From, data.To, data.Quantity, data.Memo)
			if err != nil {
				seelog.Errorf("transfer err by %v", err.Error())
				return
			}
			//time.Sleep(time.Duration(1) * time.Millisecond)
		}
	},
}

func init() {
	cmd.RootCmd.AddCommand(batchTransferCmd)
}
