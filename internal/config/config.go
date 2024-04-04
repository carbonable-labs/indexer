package config

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"errors"
	"fmt"
	"os"

	"github.com/carbonable-labs/indexer/internal/storage"
	"gopkg.in/yaml.v3"
)

var (
	ErrParseConfig         = errors.New("unable to parse config")
	ErrUnmarshalYamlFailed = errors.New("unable to unmarshal yaml config")
	ErrUnmarshalJsonFailed = errors.New("unable to unmarshal json config")
)

type ContractRepository interface {
	SaveConfig(Config) error
	GetConfigs() ([]Config, error)
	GetConfig(string) (*Config, error)
}

type NatsContractRepository struct {
	storage storage.Storage
}

func NewPebbleContractRepository(s storage.Storage) *NatsContractRepository {
	return &NatsContractRepository{storage: s}
}

func (r *NatsContractRepository) SaveConfig(c Config) error {
	key := fmt.Sprintf("config.%s", c.AppName)

	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	if err := encoder.Encode(c); err != nil {
		return err
	}
	if err := r.storage.Set([]byte(key), buf.Bytes()); err != nil {
		return err
	}
	return nil
}

func (r *NatsContractRepository) GetConfigs() ([]Config, error) {
	var configs []Config
	data := r.storage.Scan([]byte("config."))
	for _, d := range data {
		var c Config
		decoder := gob.NewDecoder(bytes.NewReader(d))
		if err := decoder.Decode(&c); err != nil {
			return nil, err
		}
		configs = append(configs, c)
	}
	return configs, nil
}

func (r *NatsContractRepository) GetConfig(name string) (*Config, error) {
	var c Config
	data := r.storage.Get([]byte(fmt.Sprintf("config.%s", name)))
	decoder := gob.NewDecoder(bytes.NewReader(data))
	if err := decoder.Decode(&c); err != nil {
		return nil, err
	}
	return &c, nil
}

type Contract struct {
	Events  map[string]string `yaml:"events" json:"events"`
	Address string            `yaml:"address" json:"address"`
	Name    string            `yaml:"name" json:"name"`
}

type Config struct {
	AppName    string `yaml:"app_name" json:"app_name"`
	Hash       string
	Contracts  []Contract `yaml:"contracts" json:"contracts"`
	StartBlock uint64     `yaml:"start_block" json:"start_block"`
}

func (c *Config) GetContract(address string) *Contract {
	for _, contract := range c.Contracts {
		if contract.Address == address {
			return &contract
		}
	}
	return nil
}

func (c *Config) ComputeHash() *Config {
	b, _ := yaml.Marshal(c)
	h := sha256.Sum256(b)
	c.Hash = fmt.Sprintf("%x", h)
	return c
}

func NewCongig() *Config {
	var c Config
	return &c
}

// Create config from yaml input
func FromYamlFile(file string) (*Config, error) {
	cfg := Config{}

	data, err := os.ReadFile(file)
	if err != nil {
		return nil, ErrParseConfig
	}

	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, ErrUnmarshalYamlFailed
	}

	return &cfg, nil
}

// Create config from json input
func FromJson(data []byte) (*Config, error) {
	cfg := Config{}
	err := yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, ErrUnmarshalJsonFailed
	}
	return &cfg, nil
}
