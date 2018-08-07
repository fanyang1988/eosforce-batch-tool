package batchCmd

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/bronze1man/go-yaml2json"
	"github.com/cihub/seelog"
	"github.com/fanyang1988/eos-go"
	"github.com/fanyang1988/eosc/cli"
	"github.com/fanyang1988/eosc/eosc/fee"
	eosvault "github.com/fanyang1988/eosc/vault"
	"github.com/fanyang1988/eosforce-batch-tool/utils"
	"github.com/spf13/viper"
	"time"
)

func mustGetWallet() *eosvault.Vault {
	vault, err := setupWallet()
	utils.ErrorCheck(err, "wallet setup")
	return vault
}

func setupWallet() (*eosvault.Vault, error) {
	walletFile := viper.GetString("global-vault-file")
	if _, err := os.Stat(walletFile); err != nil {
		return nil, fmt.Errorf("Wallet file %q missing, ", walletFile)
	}

	vault, err := eosvault.NewVaultFromWalletFile(walletFile)
	if err != nil {
		return nil, fmt.Errorf("loading vault, %s", err)
	}

	boxer, err := eosvault.SecretBoxerForType(vault.SecretBoxWrap, viper.GetString("global-kms-gcp-keypath"))
	if err != nil {
		return nil, fmt.Errorf("secret boxer, %s", err)
	}

	if err := vault.Open(boxer); err != nil {
		return nil, err
	}

	return vault, nil
}

func getAPI(url string) *eos.API {
	res := eos.New(url)
	return res
}

func getAPIs(urls []string) []*eos.API {
	vault, err := setupWallet()
	utils.ErrorCheck(err, "setting up wallet")

	res := make([]*eos.API, 0, len(urls))
	for _, u := range urls {
		api := getAPI(u)
		api.SetSigner(vault.KeyBag)
		res = append(res, api)
	}
	return res
}

func permissionToPermissionLevel(in string) (out eos.PermissionLevel, err error) {
	return eos.NewPermissionLevel(in)
}

func permissionsToPermissionLevels(in []string) (out []eos.PermissionLevel, err error) {
	// loop all parameters
	for _, singleArg := range in {

		// if they specified "account@active,account2", handle that too..
		for _, val := range strings.Split(singleArg, ",") {
			level, err := permissionToPermissionLevel(strings.TrimSpace(val))
			if err != nil {
				return out, err
			}

			out = append(out, level)
		}
	}

	return
}
func yamlUnmarshal(cnt []byte, v interface{}) error {
	jsonCnt, err := yaml2json.Convert(cnt)
	if err != nil {
		return err
	}

	return json.Unmarshal(jsonCnt, v)
}

func loadYAMLOrJSONFile(filename string, v interface{}) error {
	cnt, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	if strings.HasSuffix(strings.ToLower(filename), ".json") {
		return json.Unmarshal(cnt, v)
	}
	return yamlUnmarshal(cnt, v)
}

func toAccount(in, field string) eos.AccountName {
	acct, err := cli.ToAccountName(in)
	if err != nil {
		utils.ErrorCheck(err, "invalid account format for %q", field)
	}

	return acct
}

func toName(in, field string) eos.Name {
	name, err := cli.ToName(in)
	if err != nil {
		utils.ErrorCheck(err, "invalid name format for %q", field)
	}

	return name
}

func toSHA256Bytes(in, field string) eos.SHA256Bytes {
	if len(in) != 64 {
		utils.ErrorCheck(errors.New("should be 64 hexadecimal characters"), "%q invalid", field)
	}

	bytes, err := hex.DecodeString(in)
	utils.ErrorCheck(err, "invalid hex in %q", field)

	return bytes
}

func pushEOSCActions(api *eos.API, actions ...*eos.Action) error {
	opts := &eos.TxOptions{}

	if err := opts.FillFromChain(api); err != nil {
		return seelog.Errorf(
			"Error fetching tapos + chain_id from the chain: %v",
			err.Error())
	}

	tx := eos.NewTransaction(actions, opts)
	tx.SetExpiration(time.Duration(viper.GetInt("global-expiration")) * time.Second)
	tx.Fee = fee.GetFeeByActions(actions)

	_, packedTx, err := api.SignTransaction(tx, opts.ChainID, eos.CompressionNone)
	if err != nil {
		return seelog.Errorf("signing transaction err by %v", err.Error())
	}

	if packedTx == nil {
		return seelog.Errorf("pack trx err by nil")
	}

	resp, err := api.PushTransaction(packedTx)
	if err != nil {
		return seelog.Errorf("pushing trx err by %v", err.Error())
	}

	seelog.Infof("Transaction submitted to the network. Transaction ID: %s", resp.TransactionID)

	return nil
}
