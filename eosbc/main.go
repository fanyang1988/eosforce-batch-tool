package main

import (
	// Load all contracts here, so we can always read and decode
	// transactions with those contracts.
	_ "github.com/fanyang1988/eos-go/msig"
	_ "github.com/fanyang1988/eos-go/system"
	_ "github.com/fanyang1988/eos-go/token"

	_ "github.com/fanyang1988/eosforce-batch-tool/batch_cmd"

	"github.com/cihub/seelog"
	"github.com/fanyang1988/eosc/eosc/cmd"
)

var version = "dev"

func init() {
	cmd.Version = version
}

func main() {
	defer seelog.Flush()
	cmd.Execute()
}
