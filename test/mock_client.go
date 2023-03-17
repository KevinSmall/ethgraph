package test

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
	"math/big"
)

var zeroAddress = common.Address{}

func GetMockClient() (*backends.SimulatedBackend, common.Address, error) {
	// See https://goethereumbook.org/en/client-simulated/
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return nil, zeroAddress, err
	}
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(1))
	if err != nil {
		return nil, zeroAddress, err
	}
	balance := new(big.Int)
	balance.SetString("10000000000000000000", 10) // 10 eth in wei
	address := auth.From
	genesisAlloc := map[common.Address]core.GenesisAccount{
		address: {
			Balance: balance,
		},
	}
	blockGasLimit := uint64(4712388)
	client := backends.NewSimulatedBackend(genesisAlloc, blockGasLimit)
	return client, address, nil
}
