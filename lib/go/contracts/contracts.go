package contracts

//go:generate go run github.com/kevinburke/go-bindata/go-bindata -prefix ../../../contracts/... -o internal/assets/assets.go -pkg assets -nometadata -nomemcopy ../../../contracts/...

import (
	"strings"

	"github.com/onflow/flow-core-contracts/lib/go/contracts/internal/assets"
)

const (
	// contractPrefix				= "../../../contracts/"
	contractPrefix				= ""
	flowFeesFilename           = "FlowFees.cdc"
	flowServiceAccountFilename = "FlowServiceAccount.cdc"
	flowTokenFilename          = "FlowToken.cdc"
	flowIdentityTableFilename  = "epochs/FlowIDTableStaking.cdc"
	flowQCFilename             = "epochs/FlowQuorumCertificate.cdc"
	flowDKGFilename            = "epochs/FlowDKG.cdc"
	flowEpochFilename          = "epochs/FlowEpoch.cdc"
	flowStakingHelper          = "FlowStakingHelper.cdc"
	flowStakingScaffoldFilename = "FlowStakingScaffold.cdc"

	hexPrefix                = "0x"
	defaultFungibleTokenAddr = "FUNGIBLETOKENADDRESS"
	defaultFlowTokenAddr     = "FLOWTOKENADDRESS"
	defaultIDTableAddr       = "FLOWIDTABLESTAKINGADDRESS"
	defaultQCAddr            = "QCADDRESS"
	defaultDKGAddr           = "DKGADDRESS"
)

// FlowToken returns the FlowToken contract. importing the
//
// The returned contract will import the FungibleToken contract from the specified address.
func FlowToken() []byte {
	code := assets.MustAssetString(contractPrefix + flowTokenFilename)
	return []byte(code)
}

// FlowFees returns the FlowFees contract.
//
// The returned contract imports the FungibleToken and FlowToken
// contracts from the default addresses.
func FlowFees() []byte {
	code := assets.MustAssetString(contractPrefix + flowFeesFilename)

	return []byte(code)
}

// FlowServiceAccount returns the FlowServiceAccount contract.
//
// The returned contract imports the FungibleToken, FlowToken and FlowFees
// contracts from the default addresses.
func FlowServiceAccount() []byte {
	code := assets.MustAssetString(contractPrefix + flowServiceAccountFilename)

	return []byte(code)
}

// FlowIDTableStaking returns the FlowIDTableStaking contract
func FlowIDTableStaking(ftAddr, flowTokenAddr string) []byte {
	code := assets.MustAssetString(contractPrefix + flowIdentityTableFilename)

	code = strings.ReplaceAll(code, defaultFungibleTokenAddr, ftAddr)
	code = strings.ReplaceAll(code, defaultFlowTokenAddr, flowTokenAddr)

	return []byte(code)
}

// FlowFees returns the FlowFees contract.
//
// The returned contract imports FlowIDTableStaking, FungibleToken and FlowToken
// contracts from the default addresses.
func FlowStakingHelper(ftAddr, flowTokenAddr, flowIdTableAddr string) []byte {
	code := assets.MustAssetString(contractPrefix + flowStakingHelper)

	code = strings.ReplaceAll(code, defaultFungibleTokenAddr, ftAddr)
	code = strings.ReplaceAll(code, defaultFlowTokenAddr, flowTokenAddr)
	code = strings.ReplaceAll(code, defaultIDTableAddr, flowIdTableAddr)

	return []byte(code)
}

func FlowStakingScaffold(ftAddr, flowTokenAddr, flowIdTableAddr string) []byte {
	println(ftAddr, flowTokenAddr, flowIdTableAddr)
	code := assets.MustAssetString(contractPrefix + flowStakingScaffoldFilename)

	code = strings.ReplaceAll(code, defaultFungibleTokenAddr, ftAddr)
	code = strings.ReplaceAll(code, defaultFlowTokenAddr, flowTokenAddr)
	code = strings.ReplaceAll(code, defaultIDTableAddr, flowIdTableAddr)

	println(code)

	return []byte(code)
}
