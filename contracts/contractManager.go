package contracts

import "github.com/ethereum/go-ethereum/common"

var eventNameHashMap map[string]*Contract

func GetContract(name string) *Contract {
	return eventNameHashMap[name]
}

func initContracts() error {
	eventNameHashMap = make(map[string]*Contract)
	for _, c := range cfg.Contracts {
		cfb := common.FromHex(c.Byc)
		if cfb == nil {
			panic("bytecode decode失败")
		}

		nContract, err := NewContract(c.Address, c.Abi, common.FromHex(c.Byc))
		if err != nil {
			return err
		}
		eventNameHashMap[c.Name] = nContract
	}
	return nil
}
