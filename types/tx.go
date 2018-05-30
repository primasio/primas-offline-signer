package types

import (
	"encoding/json"
)

type Tx struct {
	IsCreateContract bool   `json: "isCreateContract"`
	Nonce            uint64 `json: "nonce"`
	To               string `json: "to"`
	Amount           int64  `json: "amount"`
	GasLimit         int64  `json: "gasLimit"`
	GasPrice         int64  `json: "gasPrice"`
	Data             []byte `json: "data"`
}

func (t *Tx) Encode() []byte {
	dd, err := json.Marshal(t)
	if err != nil {
		panic(err.Error())
	}
	return dd
}

func (t *Tx) Decode(d []byte) error {
	return json.Unmarshal(d, t)
}
