package templates

import (
	"strings"

	"github.com/onflow/flow-core-contracts/lib/go/templates/internal/assets"
)

const (
	// admin templates
	deployWithdrawalTrackerFilename = "withdrawalTracker/admin_deploy_contract.cdc"

	// User templates
	setupWithdrawalTrackerAccountFilename = "withdrawalTracker/setup_withdrawal_tracker_account.cdc"
	checkWithdrawalTrackerFilename        = "withdrawalTracker/check_withdrawal_tracker.cdc"

	// Placehokders
	placeholderWithdrawalTrackerAddress     = "0xWITHDRAWALTRACKER"
	placeholderWithdrawalTrackerUserAddress = "0xPLACEHOLDERTRACKERUSERADDRESS"
)

/************ WithdrawalTracker Utility code ****************/

func ReplaceWithdrawalTrackerAddress(code string, wtAddr string) string {
	return strings.ReplaceAll(
		code,
		placeholderWithdrawalTrackerAddress,
		withHexPrefix(wtAddr),
	)
}

/************ WithdrawalTracker Admin Transactions ****************/

func GenerateDeployWithdrawalTracker() []byte {
	return assets.MustAsset(filePath + deployWithdrawalTrackerFilename)
}

/************ WithdrawalTracker User Transactions and scripts ****************/

func GenerateSetupAccountWithdrawalTrackerScript(withdrawalTrackerAddress string) []byte {
	code := assets.MustAssetString(filePath + setupWithdrawalTrackerAccountFilename)

	code = ReplaceWithdrawalTrackerAddress(code, withdrawalTrackerAddress)

	return []byte(code)
}

func GenerateSetupWithdrawalTrackerAccountScript(withdrawalTrackerAddress string) []byte {
	code := assets.MustAssetString(filePath + setupWithdrawalTrackerAccountFilename)

	code = ReplaceWithdrawalTrackerAddress(code, withdrawalTrackerAddress)

	return []byte(code)
}
func GenerateCheckWithdrawalTracker(withdrawalTrackerAddress, userAddress string) []byte {
	code := assets.MustAssetString(filePath + checkWithdrawalTrackerFilename)

	code = ReplaceWithdrawalTrackerAddress(code, withdrawalTrackerAddress)

	code = strings.ReplaceAll(
		code,
		placeholderWithdrawalTrackerUserAddress,
		withHexPrefix(userAddress),
	)

	return []byte(code)
}
