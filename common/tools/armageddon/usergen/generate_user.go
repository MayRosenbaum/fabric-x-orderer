package main

import (
	"fmt"
	"os"

	"github.com/hyperledger/fabric-x-orderer/common/tools/armageddon"
	"github.com/hyperledger/fabric-x-orderer/common/utils"
	"github.com/hyperledger/fabric-x-orderer/config/generate"
)

func main() {
	// Base paths
	cryptoBasePath := "/Users/maybuzaglo/go/src/github.com/cbdc-platform-deployment/out/control-node/config/cryptogen-artifacts/crypto/ordererOrganizations"
	ordererBasePath := "/Users/maybuzaglo/go/src/github.com/cbdc-platform-deployment/out/local-deployment"

	// Paths for user credentials (using ordererorg1 as the client)
	mspDir := cryptoBasePath + "/ordererorg1.example.com/users/Admin@ordererorg1.example.com/msp"
	userTLSPrivateKeyPath := cryptoBasePath + "/ordererorg1.example.com/users/Admin@ordererorg1.example.com/tls/client.key"
	userTLSCertPath := cryptoBasePath + "/ordererorg1.example.com/users/Admin@ordererorg1.example.com/tls/client.crt"
	outputPath := "/Users/maybuzaglo/go/src/github.com/cbdc-platform-deployment/user-config-generated.yaml"

	// Read TLS CA certs from each orderer's config (these are the CAs the orderers trust for client connections)
	orderers := []string{
		"fabric-orderer-1",
		"fabric-orderer-2",
		"fabric-orderer-3",
		"fabric-orderer-4",
	}

	var tlsCACertsBytesPartiesCollection [][]byte
	for _, orderer := range orderers {
		tlsCACertPath := fmt.Sprintf("%s/%s/config/tls/ca.crt", ordererBasePath, orderer)
		tlsCACertBytes, err := os.ReadFile(tlsCACertPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading TLS CA cert for %s: %s\n", orderer, err)
			os.Exit(1)
		}
		tlsCACertsBytesPartiesCollection = append(tlsCACertsBytesPartiesCollection, tlsCACertBytes)
	}

	// Create network config with parties (router and assembler endpoints are the same)
	networkConfig := &generate.Network{
		Parties: []generate.Party{
			{
				ID:                1,
				RouterEndpoint:    "localhost:7050",
				AssemblerEndpoint: "localhost:7050",
			},
			{
				ID:                2,
				RouterEndpoint:    "localhost:7051",
				AssemblerEndpoint: "localhost:7051",
			},
			{
				ID:                3,
				RouterEndpoint:    "localhost:7052",
				AssemblerEndpoint: "localhost:7052",
			},
			{
				ID:                4,
				RouterEndpoint:    "localhost:7053",
				AssemblerEndpoint: "localhost:7053",
			},
		},
		UseTLSRouter:    "true",
		UseTLSAssembler: "true",
		MaxPartyID:      4,
	}

	// Create user config
	userConfig, err := armageddon.NewUserConfig(mspDir, userTLSPrivateKeyPath, userTLSCertPath, tlsCACertsBytesPartiesCollection, networkConfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating user config: %s\n", err)
		os.Exit(1)
	}

	// Write to YAML
	err = utils.WriteToYAML(userConfig, outputPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing user config yaml: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully generated user config at: %s\n", outputPath)
}

// Made with Bob
