package kafka

import (
	"crypto/sha256"
	"crypto/sha512"
	"hash"

	"github.com/IBM/sarama"
)

type ClientOption func(*sarama.Config)

var (
	SHA256 = func() hash.Hash { return sha256.New() }
	SHA512 = func() hash.Hash { return sha512.New() }
)

type XDGSCRAMClient struct {
	HashGeneratorFcn func() hash.Hash
}

func (x *XDGSCRAMClient) Begin(userName, password, authzID string) error {
	return nil
}

func (x *XDGSCRAMClient) Step(challenge string) (string, error) {
	return "", nil
}

func (x *XDGSCRAMClient) Done() bool {
	return true
}

func WithSASL(username, password, mechanism string) ClientOption {
	return func(config *sarama.Config) {
		config.Net.SASL.Enable = true
		config.Net.SASL.User = username
		config.Net.SASL.Password = password
		switch mechanism {
		case "SCRAM-SHA-512":
			config.Net.SASL.Mechanism = sarama.SASLTypeSCRAMSHA512
			config.Net.SASL.SCRAMClientGeneratorFunc = func() sarama.SCRAMClient {
				return &XDGSCRAMClient{HashGeneratorFcn: SHA512}
			}
		case "SCRAM-SHA-256":
			config.Net.SASL.Mechanism = sarama.SASLTypeSCRAMSHA256
			config.Net.SASL.SCRAMClientGeneratorFunc = func() sarama.SCRAMClient {
				return &XDGSCRAMClient{HashGeneratorFcn: SHA256}
			}
		}
	}
}

func WithTLS(enable bool) ClientOption {
	return func(config *sarama.Config) {
		config.Net.TLS.Enable = enable
	}
}
