package config

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"reflect"

	"github.com/caarlos0/env"
	"github.com/ethereum/go-ethereum/crypto"
)

type Config struct {
	IryoAddr        string           `env:"IRYO_ADDR" envDefault:"localhost:8000"`
	EthAddr         string           `env:"ETH_ADDR" envDefault:"localhost:8545"`
	EthPublic       string           `env:"ETH_PUBLIC"`
	EthPrivate      ecdsa.PrivateKey `env:"ETH_PRIVATE,required"`
	EthContractAddr string           `env:"ETH_CONTRACT_ADDR"`
	ClientType      string           `env:"CLIENT_TYPE" envDefault:"Patient"`
	ClientAddr      string           `env:"CLIENT_ADDR" envDefault:"localhost:9000"`
	Debug           bool             `env:"DEBUG"`
	EncryptionKeys  map[string][]byte
	Connections     []string
	Tokens          map[string]string
}

func New() (*Config, error) {
	cfg := &Config{
		EncryptionKeys: make(map[string][]byte),
		Connections:    []string{},
		Tokens:         make(map[string]string),
	}
	parsers := map[reflect.Type]env.ParserFunc{
		reflect.TypeOf(cfg.EthPrivate): parseEthPrivateKey,
	}
	return cfg, env.ParseWithFuncs(cfg, parsers)
}

func parseEthPrivateKey(v string) (interface{}, error) {
	bd, err := hex.DecodeString(v)
	if err != nil {
		return nil, fmt.Errorf("Failed to decode private key from HEX; %v", err)
	}
	k, err := crypto.ToECDSA(bd)
	if err != nil {
		return nil, fmt.Errorf("Failed to convert to private key; %v", err)
	}

	return *k, nil
}

func (c *Config) GetEthPublicAddress() string {
	return crypto.PubkeyToAddress(c.EthPrivate.PublicKey).String()
}
