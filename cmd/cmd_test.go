package cmd

import (
	"testing"
	"fmt"
	"encoding/hex"
	"github.com/FactomProject/go-bip39"
	"github.com/FactomProject/go-bip32"
	"github.com/FactomProject/go-bip44"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"io/ioutil"
)

func TestAddrule(t *testing.T) {

	entropy, err := bip39.NewEntropy(128)
	if err != nil {
		panic(err.Error())
	}
	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("Mnemonic: ", mnemonic)

	// Generate a Bip32 HD wallet for the mnemonic and a user supplied password
	seed := bip39.NewSeed(mnemonic, "")
	fmt.Println("BIP39 Seed", hex.EncodeToString(seed))

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

	ks := keystore.NewKeyStore("test123", keystore.LightScryptN, keystore.LightScryptP)
	aa, err := ks.ImportECDSA(p1, "123456")
	if err != nil {
		panic(err.Error())
	}

	if err := ioutil.WriteFile(fmt.Sprintf("mnemonic.%s", aa.Address.String()), []byte(mnemonic), 0755); err != nil {
		panic(err.Error())
	}

	fmt.Println(aa.Address.String())
}
