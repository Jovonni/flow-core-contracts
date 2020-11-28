package test

/*
cd ../contracts && make generate && cd ../templates && make generate && cd ../test && go test -timeout 30s github.com/onflow/lib/go/flow-core-contracts -run ^TestWithdrawalTracker$ -v
*/

import (
	"strconv"
	"strings"
	"testing"

	"github.com/onflow/cadence"
	"github.com/onflow/flow-core-contracts/lib/go/contracts"
	"github.com/onflow/flow-core-contracts/lib/go/templates"
	emulator "github.com/onflow/flow-emulator"
	"github.com/onflow/flow-go-sdk"
	sdk "github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/crypto"
	sdktemplates "github.com/onflow/flow-go-sdk/templates"
	"github.com/onflow/flow-go-sdk/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	trackerInitialBalance = "0.0"
	trackerInitialLimit   = "2000000.0"
)

// Simple error-handling wrapper for Flow account creation.
func createAccount(t *testing.T, b *emulator.Blockchain, accountKeys *test.AccountKeys) (sdk.Address, crypto.Signer, *sdk.AccountKey) {
	accountKey, signer := accountKeys.NewWithSigner()
	address, err := b.CreateAccount([]*sdk.AccountKey{accountKey}, nil)
	require.NoError(t, err)
	return address, signer, accountKey
}

// Create a new Flow account that has the FAT vault installed.
func createWithdrawalTrackerAccount(t *testing.T, b *emulator.Blockchain, accountKeys *test.AccountKeys, withdrawalTrackerAddress sdk.Address) (sdk.Address, crypto.Signer, *sdk.AccountKey) {
	address, signer, accountKey := createAccount(t, b, accountKeys)

	txSetup := flow.NewTransaction().
		SetScript(templates.GenerateSetupAccountWithdrawalTrackerScript(withdrawalTrackerAddress.String())).
		SetGasLimit(100).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address).
		AddAuthorizer(address)

	signAndSubmit(
		t, b, txSetup,
		[]flow.Address{b.ServiceKey().Address, address},
		[]crypto.Signer{b.ServiceKey().Signer(), signer},
		false,
	)

	return address, signer, accountKey
}

func simulateWithdrawal(
	t *testing.T,
	b *emulator.Blockchain,
	withdrawalTrackerAddress sdk.Address,
	withdrawerAddress sdk.Address,
	withdrawerSigner crypto.Signer,
	withdrawalAmount string,
	shouldFail bool,
) {
	cadenceWithdrawalAmount, err := cadence.NewUFix64(withdrawalAmount)
	assert.NoError(t, err)

	tx := flow.NewTransaction().
		SetScript([]byte(templates.ReplaceWithdrawalTrackerAddress(`
	import WithdrawalTracker from 0xWITHDRAWALTRACKER

	transaction(amount: UFix64) {

		prepare(account: AuthAccount) {
			let withdrawalTotalTrackerPath = /storage/withdrawalTotalTracker
	
			let withdrawalTracker = account.borrow<&WithdrawalTracker.WithdrawalTotalTracker>(from: withdrawalTotalTrackerPath)
				?? panic("Could not load withdrawal tracker")
	
			// We have to this here because we need access to the AuthAccount to save the result if it doesn't assert.
			withdrawalTracker!.updateRunningTotal(withdrawalAmount: amount)
		}

	}`, withdrawalTrackerAddress.String()))).
		SetGasLimit(100).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address).
		AddAuthorizer(withdrawerAddress)
	tx.AddArgument(cadenceWithdrawalAmount)

	signAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, withdrawerAddress},
		[]crypto.Signer{b.ServiceKey().Signer(), withdrawerSigner},
		shouldFail,
	)

}

func uFix64(fixme uint64) float64 {
	return float64(fixme) / float64(100000000)
}

// Format the string very carefully to match expected
func uFix64Str(fixme uint64) string {
	str := strconv.FormatFloat(uFix64(fixme), 'f', -1, 64)
	if !strings.Contains(str, ".") {
		str += ".0"
	}
	return str
}

func checkTrackerState(
	t *testing.T,
	b *emulator.Blockchain,
	withdrawalTrackerAddress sdk.Address,
	withdrawerAddress sdk.Address,
	totalWithdrawnShouldBe string,
	limitShouldBe string,
) {
	// Check that the balance and limit were updated correctly
	state, stateError := b.ExecuteScript(
		templates.GenerateCheckWithdrawalTracker(withdrawalTrackerAddress.String(), withdrawerAddress.String()),
		[][]byte{},
	)
	require.NoError(t, stateError)
	require.True(t, state.Succeeded())

	// Check that we got the correct values
	stateValues := state.Value.ToGoValue().([]interface{})
	require.Equal(t, totalWithdrawnShouldBe, uFix64Str(stateValues[0].(uint64)))
	require.Equal(t, limitShouldBe, uFix64Str(stateValues[1].(uint64)))

}

// Get the address of the most recently deployed contract on the emulator blockchain.
func getDeployedContractAddress(t *testing.T, b *emulator.Blockchain) sdk.Address {
	// Get the deployed contract's address.
	var address sdk.Address

	//foundAddress:
	for i := uint64(0); i < 1000; i++ {
		results, _ := b.GetEventsByHeight(i, "flow.AccountCreated")

		for _, event := range results {
			if event.Type == sdk.EventAccountCreated {
				address = sdk.Address(event.Value.Fields[0].(cadence.Address))
				// We want the last created address, and we created one before,
				// so we don't want to break when we find the first address as that
				// will be the wrong one.
				//break foundAddress
			}
		}
	}

	assert.NotEqual(t, address, sdk.EmptyAddress)

	return address
}

func TestWithdrawalTracker(t *testing.T) {
	b := newEmulator()

	accountKeys := test.AccountKeyGenerator()

	// Create new keys for the ID table account
	withdrawalTrackerAccountKey, _ := accountKeys.NewWithSigner()
	withdrawalTrackerCode := contracts.WithdrawalTracker()
	withdrawalTrackerAddress, err := b.CreateAccount(
		[]*sdk.AccountKey{withdrawalTrackerAccountKey},
		[]sdktemplates.Contract{
			{
				Name:   "WithdrawalTracker",
				Source: string(withdrawalTrackerCode),
			},
		})
	assert.NoError(t, err)

	_, err = b.CommitBlock()
	assert.NoError(t, err)

	// Create our test account and place the tracker in it
	// This tests setup_withdrawal_tracker_account.cdc
	withdrawerAddress, withdrawerSigner, _ := createWithdrawalTrackerAccount(t, b, accountKeys, withdrawalTrackerAddress)

	t.Run("Should set tracker state correctly on creation", func(t *testing.T) {
		// Check that the balance and limit were set correctly during resource creation
		checkTrackerState(t, b, withdrawalTrackerAddress, withdrawerAddress, trackerInitialBalance, trackerInitialLimit)
	})

	t.Run("Should set tracker state correctly on deposit", func(t *testing.T) {
		withdrawalAmount := "123.456"

		// Update the tracker but do not withdraw
		simulateWithdrawal(t, b, withdrawalTrackerAddress, withdrawerAddress, withdrawerSigner, withdrawalAmount, false)

		// Check that the balance and limit were updated correctly
		checkTrackerState(t, b, withdrawalTrackerAddress, withdrawerAddress, withdrawalAmount, trackerInitialLimit)
	})

	t.Run("Should fail if withdrawal would exceed limit", func(t *testing.T) {
		withdrawalAmount := "3000000.0"

		// Update the tracker but do not withdraw
		simulateWithdrawal(t, b, withdrawalTrackerAddress, withdrawerAddress, withdrawerSigner, withdrawalAmount, true)
	})
}
