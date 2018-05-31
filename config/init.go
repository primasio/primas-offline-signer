package config

import (
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/accounts"
	log "github.com/inconshreveable/log15"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"sync"
	"path"
	"github.com/ethereum/go-ethereum/ethclient"
)

var (
	once     sync.Once
	instance *Config
)

type Contract struct {
	Name    string `yaml:"name"`
	Address string `yaml:"address"`
	Abi     string `yaml:"abi"`
	Byc     string `yaml:"byc"`
}

type Config struct {
	ethClient    *ethclient.Client
	l            log.Logger
	cfgFile      string
	nonce        chan int
	nodeAccount  *accounts.Account
	nodeKeystore *keystore.KeyStore
	home         string

	IFile  string `yaml:"-"`
	OFile  string `yaml:"-"`
	PassWD string `yaml:"-"`

	EthAddr        string     `yaml:"eth_addr"`
	Nonce          uint64     `yaml:"nonce"`
	KeystoreDir    string     `yaml:"keystore"`
	Passphrase     string     `yaml:"passphrase"`
	AccountAddress string     `yaml:"account_address"`
	GasLimit       int        `yaml:"gas_limit"`
	Gasprice       int        `yaml:"gas_price"`
	ChainId        int        `yaml:"chain_id"`
	Contracts      []Contract `yaml:"contracts"`
}

func (c *Config) LoadConfig() {
	c.l.Info("cfgFile", "file", c.cfgFile)

	d, err := ioutil.ReadFile(c.cfgFile)
	if err != nil {
		panic(fmt.Sprintf("配置文件读取错误\n%s", err.Error()))
	}

	if err := yaml.Unmarshal(d, c); err != nil {
		panic(fmt.Sprintf("配置文件加载错误\n%s", err.Error()))
	}

	c.cfgFile = path.Join(c.home, "kdata.yaml")
	instance.KeystoreDir = path.Join(c.home, "keystore")
}

func (c *Config) Dumps() {
	d, err := yaml.Marshal(c)
	if err != nil {
		panic(err.Error())
	}

	if err := ioutil.WriteFile(c.cfgFile, d, 0755); err != nil {
		panic(fmt.Sprintf("写入配置文件\n%s", err.Error()))
	}
}

func (c *Config) InitNode() {
	c.nodeKeystore = keystore.NewKeyStore(c.KeystoreDir, keystore.LightScryptN, keystore.LightScryptP)

	if len(c.nodeKeystore.Accounts()) == 0 {
		c.l.Warn("keystore file does not exist")
		return
	}

	c.nodeAccount = &c.nodeKeystore.Accounts()[0]
	if err := c.nodeKeystore.Unlock(*c.nodeAccount, c.Passphrase); err != nil {
		panic(fmt.Sprintf("账号解锁失败\n%s", err.Error()))
	}
}

func (t *Config) InitLog() {
	t.l = log.New()
	ll, err := log.LvlFromString("debug")
	if err != nil {
		panic(err.Error())
	}
	t.l.SetHandler(log.LvlFilterHandler(ll, log.StreamHandler(os.Stdout, log.TerminalFormat())))
}

func (c *Config) InitEthClient() {
	client, err := ethclient.Dial(c.EthAddr)
	if err != nil {
		c.l.Error("以太坊连接失败")
		panic(err.Error())
	}

	c.ethClient = client
}

func (c *Config) GetEthClient() *ethclient.Client {
	if c.ethClient == nil {
		panic("请初始化以太坊客户端")
	}
	return c.ethClient
}
