package test

/*
cd lib/go/test
cd ../contracts && make generate && cd ../templates && make generate && cd ../test && go test -timeout 30s github.com/onflow/lib/go/flow-core-contracts -run ^TestMemoEmitter$ -v
*/

import (
	"fmt"
	"testing"

	jsoncdc "github.com/onflow/cadence/encoding/json"
	emulator "github.com/onflow/flow-emulator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	//"github.com/stretchr/testify/require"

	"github.com/onflow/cadence"
	"github.com/onflow/flow-go-sdk"
	sdk "github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/crypto"
	"github.com/onflow/flow-go-sdk/test"

	//"github.com/onflow/flow-core-contracts/lib/go/contracts"
	"github.com/onflow/flow-core-contracts/lib/go/contracts"
	"github.com/onflow/flow-core-contracts/lib/go/templates"
)

const (
	flowAmount                  = "0.1"
	maxMemoLength               = 255
	memoStringOK                = "Hello Memo!"
	memoStringEmpty             = ""
	memoStringMaxLength         = "0123456789ABCDEF0123456789ABCDEF0123456789ABCDEF0123456789ABCDEF0123456789ABCDEF0123456789ABCDEF0123456789ABCDEF0123456789ABCDEF0123456789ABCDEF0123456789ABCDEF0123456789ABCDEF0123456789ABCDEF0123456789ABCDEF0123456789ABCDEF0123456789ABCDEF0123456789ABCDE"
	memoStringMaxLengthExceeded = memoStringMaxLength + "Z"
)

// Memo event

type MemoEvent interface {
	Memo() string
}

type memoEvent flow.Event

var _ MemoEvent = (*memoEvent)(nil)

// Memo returns the memo string of the Memo event.
func (evt memoEvent) Memo() string {
	return evt.Value.Fields[0].(cadence.String).ToGoValue().(string)
}

// Make sure that a memo has been included in a transaction's events
func checkMemo(t *testing.T, b *emulator.Blockchain, env templates.Environment, tx *flow.Transaction, memoString string) {
	txResult, err := b.GetTransactionResult(tx.ID())
	assert.NoError(t, err)
	assert.Equal(t, flow.TransactionStatusSealed, txResult.Status)

	memoEventType := fmt.Sprintf("A.%s.MemoEmitter.Memo", env.MemoEmitterAddress)

	found := false
	for _, event := range txResult.Events {
		if event.Type == memoEventType {
			memoEvent := memoEvent(event)
			assert.Equal(t, memoString, memoEvent.Memo())
			found = true
			break
		}
	}
	assert.True(t, found)
}

// Simple error-handling wrapper for Flow account creation.
func createAccount(t *testing.T, b *emulator.Blockchain, accountKeys *test.AccountKeys) (sdk.Address, crypto.Signer, *sdk.AccountKey) {
	accountKey, signer := accountKeys.NewWithSigner()
	address, err := b.CreateAccount([]*sdk.AccountKey{accountKey}, nil)
	require.NoError(t, err)
	return address, signer, accountKey
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

func emitMemo(
	t *testing.T,
	b *emulator.Blockchain,
	env templates.Environment,
	senderAddress sdk.Address,
	senderSigner crypto.Signer,
	memoString string,
	shouldRevert bool,
) {
	tx := flow.NewTransaction().
		SetScript(templates.GenerateEmitMemoScript(env)).
		SetGasLimit(100).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address).
		AddAuthorizer(senderAddress)
	tx.AddArgument(cadence.NewString(memoString))

	signAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, senderAddress},
		[]crypto.Signer{b.ServiceKey().Signer(), senderSigner},
		shouldRevert,
	)

	// Only check state if transaction should have gone through
	if !shouldRevert {
		checkMemo(t, b, env, tx, memoString)
	}
}

func transferFlowEmitMemo(
	t *testing.T,
	b *emulator.Blockchain,
	env templates.Environment,
	senderAddress sdk.Address,
	senderSigner crypto.Signer,
	recipientAddress sdk.Address,
	amount string,
	memoString string,
	shouldRevert bool,
) {
	cadenceAmount, err := cadence.NewUFix64(amount)
	assert.NoError(t, err)

	tx := flow.NewTransaction().
		SetScript(templates.GenerateTransferFlowWithMemoScript(env)).
		SetGasLimit(100).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address).
		AddAuthorizer(senderAddress)
	tx.AddArgument(cadenceAmount)
	tx.AddArgument(cadence.NewAddress(recipientAddress))
	tx.AddArgument(cadence.NewString(memoString))

	signAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, senderAddress},
		[]crypto.Signer{b.ServiceKey().Signer(), senderSigner},
		shouldRevert,
	)

	// Only check state if transaction should have gone through
	if !shouldRevert {
		checkMemo(t, b, env, tx, memoString)
	}
}

func TestMemoEmitter(t *testing.T) {
	b := newEmulator()

	env := templates.Environment{
		FungibleTokenAddress: emulatorFTAddress,
		FlowTokenAddress:     emulatorFlowTokenAddress,
	}

	accountKeys := test.AccountKeyGenerator()

	// Deploy the memoEmitter Contract.

	// Create the admin key and signer
	memoEmitterKey, _ /*memoEmitterSigner*/ := accountKeys.NewWithSigner()
	memoEmitterCode := contracts.FlowMemoEmitter()

	cadencePublicKey := bytesToCadenceArray(memoEmitterKey.Encode())
	cadenceCode := bytesToCadenceArray(memoEmitterCode)

	// Deploy the MemoEmitter contract
	createAccountTx := flow.NewTransaction().
		SetScript(templates.GenerateMemoEmitterDeploy(env)).
		SetGasLimit(100).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address).
		AddAuthorizer(b.ServiceKey().Address).
		AddRawArgument(jsoncdc.MustEncode(cadencePublicKey)).
		AddRawArgument(jsoncdc.MustEncode(cadence.NewString("MemoEmitter"))).
		AddRawArgument(jsoncdc.MustEncode(cadenceCode))

	signAndSubmit(
		t, b, createAccountTx,
		[]flow.Address{b.ServiceKey().Address},
		[]crypto.Signer{b.ServiceKey().Signer()},
		false,
	)

	// Get the deployed contract's address.
	env.MemoEmitterAddress = getDeployedContractAddress(t, b).String()

	// Create a user
	userAddress, userSigner, _ := createAccount(t, b, accountKeys)

	t.Run("Should be able to emit Memo event with brief memo string", func(t *testing.T) {
		emitMemo(
			t,
			b,
			env,
			userAddress,
			userSigner,
			memoStringOK,
			false,
		)
	})

	t.Run("Should be able to emit Memo event with empty memo string", func(t *testing.T) {
		emitMemo(
			t,
			b,
			env,
			userAddress,
			userSigner,
			memoStringEmpty,
			false,
		)
	})

	t.Run("Should be able to emit Memo event with memo string of max length", func(t *testing.T) {
		emitMemo(
			t,
			b,
			env,
			userAddress,
			userSigner,
			memoStringMaxLength,
			false,
		)
	})

	t.Run("Should not be able to emit Memo event with memo string that is too long", func(t *testing.T) {
		emitMemo(
			t,
			b,
			env,
			userAddress,
			userSigner,
			memoStringMaxLengthExceeded,
			true,
		)
	})

	t.Run("Should be able to transfer FLOW and emit Memo event with brief memo string", func(t *testing.T) {
		transferFlowEmitMemo(
			t,
			b,
			env,
			b.ServiceKey().Address,
			b.ServiceKey().Signer(),
			userAddress,
			flowAmount,
			memoStringOK,
			false,
		)
	})

	t.Run("Should not be able to transfer FLOW and emit Memo event with memo string that is too long", func(t *testing.T) {
		transferFlowEmitMemo(
			t,
			b,
			env,
			b.ServiceKey().Address,
			b.ServiceKey().Signer(),
			userAddress,
			flowAmount,
			memoStringMaxLengthExceeded,
			true,
		)
	})

}
