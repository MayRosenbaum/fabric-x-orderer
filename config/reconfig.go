package config

import (
	"bytes"
	"errors"
	"net"
	"strconv"

	"github.com/hyperledger/fabric-x-orderer/common/types"
	config_protos "github.com/hyperledger/fabric-x-orderer/config/protos"
)

type NodeConfig interface {
	GetHost() string
	GetPort() uint32
	GetTlsCert() []byte
}

type NodeConfigWthSign interface {
	NodeConfig
	GetSignCert() []byte
}

// IsPartyEvicted checks if a party exists in the configuration and returns true if it is evicted or not.
func IsPartyEvicted(partyID types.PartyID, newConfig *Configuration) bool {
	newSharedPartyConfig := FindParty(partyID, newConfig)
	if newSharedPartyConfig == nil {
		return true
	}
	return false
}

// FindParty returns the PartyConfig associated with the given partyID from the shared configuration.
// It returns nil if the party is not found.
func FindParty(partyID types.PartyID, config *Configuration) *config_protos.PartyConfig {
	for _, party := range config.SharedConfig.PartiesConfig {
		if types.PartyID(party.PartyID) == partyID {
			return party
		}
	}
	return nil
}

// IsNodeConfigChangeRestartRequired reports whether a restart is required due to changed between current configuration and new configuration.
// A restart is required if any of the following fields differ:
//   - host or port
//   - TLS certificate
//   - sign certificate (if both configs implement NodeConfigWithSignCert)
//
// Both arguments must represent the same node type and be non-nil.
func IsNodeConfigChangeRestartRequired(currentConfig, newConfig NodeConfig) (bool, error) {
	if currentConfig == nil {
		return false, errors.New("current config is nil")
	}

	if newConfig == nil {
		return false, errors.New("new config is nil")
	}

	currAddr := net.JoinHostPort(currentConfig.GetHost(), strconv.Itoa(int(currentConfig.GetPort())))
	newAddr := net.JoinHostPort(newConfig.GetHost(), strconv.Itoa(int(newConfig.GetPort())))

	if currAddr != newAddr || !bytes.Equal(currentConfig.GetTlsCert(), newConfig.GetTlsCert()) {
		return true, nil
	}

	if extendedCurrConfig, ok := currentConfig.(NodeConfigWthSign); ok {
		if extendedNewConfig, ok := newConfig.(NodeConfigWthSign); ok {
			if !bytes.Equal(extendedCurrConfig.GetSignCert(), extendedNewConfig.GetSignCert()) {
				return true, nil
			}
		}
	}

	return false, nil
}
