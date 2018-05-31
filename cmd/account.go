package cmd

import (
	"github.com/urfave/cli"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/FactomProject/go-bip44"
	"io/ioutil"
	"fmt"
	"github.com/FactomProject/go-bip39"
	"github.com/FactomProject/go-bip32"
	"github.com/ethereum/go-ethereum/crypto"
	"path"
)

func AccountCmd() cli.Command {
	return cli.Command{
		Name:    "newAccount",
		Aliases: []string{"nc"},
		Usage:   "create account",
		Flags: []cli.Flag{
			passwdFlag(),
		},
		Action: func(c *cli.Context) error {
			logger.Info("keystore file", "file", cfg.KeystoreDir)

			if cfg.PassWD == "" {
				logger.Error("请输入密码参数")
				panic("请输入密码参数")
			}

			entropy, err := bip39.NewEntropy(128)
			if err != nil {
				panic(err.Error())
			}

			mnemonic, err := bip39.NewMnemonic(entropy)
			if err != nil {
				panic(err.Error())
			}

			logger.Info("Mnemonic", "msg", mnemonic)

			seed := bip39.NewSeed(mnemonic, "")

			masterKey, err := bip32.NewMasterKey(seed)
			if err != nil {
				panic(err)
			}

			fKey, err := bip44.NewKeyFromMasterKey(masterKey, bip44.TypeEther, bip32.FirstHardenedChild, 0, 0)
			if err != nil {
				panic(err)
			}

			p1, err := crypto.ToECDSA(fKey.Key)
			if err != nil {
				panic(err.Error())
			}

			ks := keystore.NewKeyStore(cfg.KeystoreDir, keystore.LightScryptN, keystore.LightScryptP)
			aa, err := ks.ImportECDSA(p1, cfg.PassWD)
			if err != nil {
				panic(err.Error())
			}

			p := path.Join(cfg.KeystoreDir, fmt.Sprintf("mnemonic.%s", aa.Address.String()))
			if err := ioutil.WriteFile(p, []byte(mnemonic), 0755); err != nil {
				panic(err.Error())
			}

			logger.Info("account", "address", aa.Address.String())
			return nil
		},
	}
}
