package config

import (
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"os"
	log "github.com/inconshreveable/log15"
	"path"
	"github.com/kooksee/pstoff/cmn"
	"time"
	"context"
	"sync"
	"github.com/ethereum/go-ethereum/common"
)

var once1 sync.Once

func (c *Config) GetNonce() uint64 {
	once1.Do(func() {

		nonce := uint64(0)
		var err error

		ctx, _ := context.WithTimeout(context.Background(), time.Minute)
		pon := os.Getenv("PON")

		c.l.Info("primas onwer", "pon", pon)

		if pon == "o1" {
			nonce, err = c.GetEthClient().NonceAt(ctx, common.HexToAddress("cd593e2fabd6ff935ba2d44070e599fce242ca09"), nil)
			//nonce, err = c.GetEthClient().NonceAt(ctx, common.HexToAddress("ef924f46ff4d6315867f1580a15a9617ea68463f"), nil)
		}

		if pon == "o2" {
			//nonce, err = c.GetEthClient().NonceAt(ctx, common.HexToAddress("93ee3eef1c32c63aefc15c36d468e5770f358d39"), nil)
			nonce, err = c.GetEthClient().NonceAt(ctx, common.HexToAddress("66ff46896da45915993fe9a785defe5c49144963"), nil)
		}

		if pon == "o3" {
			nonce, err = c.GetEthClient().NonceAt(ctx, common.HexToAddress("8f930297bcd24d9567afb9ad8631411145711a58"), nil)
			//nonce, err = c.GetEthClient().NonceAt(ctx, common.HexToAddress("87d234c23dc04efc56153fb7cd58053560be04b2"), nil)
		}

		if pon != "o1" && pon != "o2" && pon != "o3" {
			nonce, err = c.GetEthClient().NonceAt(ctx, c.GetNodeAccount().Address, nil)
		}

		if err != nil {
			panic(err.Error())
		}

		c.Nonce = nonce
		c.isNonce = true
	})

	if c.isNonce {
		Log().Info("nonce", "nonce", c.Nonce, "isnonce", c.isNonce)
		c.isNonce = false
		return c.Nonce
	}

	c.Nonce += 1
	Log().Info("nonce", "nonce", c.Nonce)
	return c.Nonce
}

func (c *Config) GetNodeAccount() *accounts.Account {
	if c.nodeAccount == nil {
		panic("please init node account")
	}
	return c.nodeAccount
}

func (c *Config) GetNodeKeyStore() *keystore.KeyStore {
	if c.nodeKeystore == nil {
		panic("please init nodeKeystore")
	}
	return c.nodeKeystore
}

func GetCfg() *Config {
	if instance == nil {
		panic("please init config")
	}
	return instance
}

func GetHomeDir(defaultHome string) string {
	if len(os.Args) > 2 && os.Args[len(os.Args)-2] == "--home" {
		defaultHome = os.Args[len(os.Args)-1]
		os.Args = os.Args[:len(os.Args)-2]
	}
	return defaultHome
}

func Log() log.Logger {
	cfg := GetCfg()
	if cfg.l == nil {
		panic("please init log")
	}
	return cfg.l
}

func NewCfg(defaultHomeDir string) *Config {
	defaultHomeDir = GetHomeDir(defaultHomeDir)
	instance = &Config{}

	instance.home = defaultHomeDir
	instance.LogLevel = "debug"
	instance.cfgFile = path.Join(defaultHomeDir, "kdata.yaml")
	instance.LogLevel = "debug"
	instance.LogPath = path.Join(defaultHomeDir, "log")
	instance.KeystoreDir = path.Join(defaultHomeDir, "keystore")
	instance.IFile = "input.json"
	instance.OFile = "output.json"
	instance.ChainId = 0
	instance.Passphrase = "Test123:::"
	instance.GasLimit = 750000
	instance.Gasprice = 100000000000
	instance.Nonce = 1
	instance.EthAddr = "http://localhost:8545"

	cmn.EnsureDir(instance.home, os.FileMode(0755))
	cmn.EnsureDir(instance.LogPath, os.FileMode(0755))
	cmn.EnsureDir(instance.KeystoreDir, os.FileMode(0755))

	return instance
}
