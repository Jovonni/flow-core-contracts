package test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/onflow/cadence"
	jsoncdc "github.com/onflow/cadence/encoding/json"
	"github.com/onflow/cadence/runtime/interpreter"
	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/crypto"
	"github.com/onflow/flow-go-sdk/test"

	"github.com/onflow/flow-core-contracts/lib/go/templates"
)

const (
	adminID   = "0000000000000000000000000000000000000000000000000000000000000001"
	admin     = 1
	joshID    = "0000000000000000000000000000000000000000000000000000000000000002"
	josh      = 2
	maxID     = "0000000000000000000000000000000000000000000000000000000000000003"
	max       = 3
	bastianID = "0000000000000000000000000000000000000000000000000000000000000004"
	bastian   = 4
	accessID  = "0000000000000000000000000000000000000000000000000000000000000005"
	access    = 5

	nonexistantID = "0000000000000000000000000000000000000000000000000000000000383838383"

	firstDelegatorID        = 1
	firstDelegatorStringID  = "0001"
	secondDelegatorID       = 2
	secondDelegatorStringID = "0002"
)

func TestIDTableDeployment(t *testing.T) {

	t.Parallel()

	b := newBlockchain()

	env := templates.Environment{
		FungibleTokenAddress: emulatorFTAddress,
		FlowTokenAddress:     emulatorFlowTokenAddress,
	}

	accountKeys := test.AccountKeyGenerator()

	// Create new keys for the ID table account
	IDTableAccountKey, IDTableSigner := accountKeys.NewWithSigner()
	var idTableAddress = deployStakingContract(t, b, IDTableAccountKey, env)

	env.IDTableAddress = idTableAddress.Hex()

	t.Run("Should be able to read empty table fields and initialized fields", func(t *testing.T) {

		result := executeScriptAndCheck(t, b, templates.GenerateReturnCurrentTableScript(env), nil)

		idArray := result.(cadence.Array).Values
		assert.Empty(t, idArray)

		result = executeScriptAndCheck(t, b, templates.GenerateReturnProposedTableScript(env), nil)

		idArray = result.(cadence.Array).Values
		assert.Empty(t, idArray)

		// Check that the stake requirements for each node role are initialized correctly

		result = executeScriptAndCheck(t, b, templates.GenerateGetStakeRequirementsScript(env), [][]byte{jsoncdc.MustEncode(cadence.UInt8(1))})
		assertEqual(t, CadenceUFix64("250000.0"), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetStakeRequirementsScript(env), [][]byte{jsoncdc.MustEncode(cadence.UInt8(2))})
		assertEqual(t, CadenceUFix64("500000.0"), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetStakeRequirementsScript(env), [][]byte{jsoncdc.MustEncode(cadence.UInt8(3))})
		assertEqual(t, CadenceUFix64("1250000.0"), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetStakeRequirementsScript(env), [][]byte{jsoncdc.MustEncode(cadence.UInt8(4))})
		assertEqual(t, CadenceUFix64("135000.0"), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetStakeRequirementsScript(env), [][]byte{jsoncdc.MustEncode(cadence.UInt8(5))})
		assertEqual(t, CadenceUFix64("0.0"), result)

		// Check that the total tokens staked for each node role are initialized correctly

		result = executeScriptAndCheck(t, b, templates.GenerateGetTotalTokensStakedByTypeScript(env), [][]byte{jsoncdc.MustEncode(cadence.UInt8(1))})
		assertEqual(t, CadenceUFix64("0.0"), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetTotalTokensStakedByTypeScript(env), [][]byte{jsoncdc.MustEncode(cadence.UInt8(2))})
		assertEqual(t, CadenceUFix64("0.0"), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetTotalTokensStakedByTypeScript(env), [][]byte{jsoncdc.MustEncode(cadence.UInt8(3))})
		assertEqual(t, CadenceUFix64("0.0"), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetTotalTokensStakedByTypeScript(env), [][]byte{jsoncdc.MustEncode(cadence.UInt8(4))})
		assertEqual(t, CadenceUFix64("0.0"), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetTotalTokensStakedByTypeScript(env), [][]byte{jsoncdc.MustEncode(cadence.UInt8(5))})
		assertEqual(t, CadenceUFix64("0.0"), result)

		// Check that the reward ratios were initialized correctly for each node role

		result = executeScriptAndCheck(t, b, templates.GenerateGetRewardRatioScript(env), [][]byte{jsoncdc.MustEncode(cadence.UInt8(1))})
		assertEqual(t, CadenceUFix64("0.168"), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetRewardRatioScript(env), [][]byte{jsoncdc.MustEncode(cadence.UInt8(2))})
		assertEqual(t, CadenceUFix64("0.518"), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetRewardRatioScript(env), [][]byte{jsoncdc.MustEncode(cadence.UInt8(3))})
		assertEqual(t, CadenceUFix64("0.078"), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetRewardRatioScript(env), [][]byte{jsoncdc.MustEncode(cadence.UInt8(4))})
		assertEqual(t, CadenceUFix64("0.236"), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetRewardRatioScript(env), [][]byte{jsoncdc.MustEncode(cadence.UInt8(5))})
		assertEqual(t, CadenceUFix64("0.0"), result)

		// Check that the weekly payout was initialized correctly

		result = executeScriptAndCheck(t, b, templates.GenerateGetWeeklyPayoutScript(env), nil)
		assertEqual(t, CadenceUFix64("1250000.0"), result)

	})

	t.Run("Shouldn't be able to change the cut percentage above 1", func(t *testing.T) {

		tx := createTxWithTemplateAndAuthorizer(b, templates.GenerateChangeCutScript(env), idTableAddress)

		_ = tx.AddArgument(CadenceUFix64("2.10"))

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, idTableAddress},
			[]crypto.Signer{b.ServiceKey().Signer(), IDTableSigner},
			true,
		)
	})

	t.Run("Should be able to change the cut percentage", func(t *testing.T) {

		tx := createTxWithTemplateAndAuthorizer(b, templates.GenerateChangeCutScript(env), idTableAddress)

		_ = tx.AddArgument(CadenceUFix64("0.10"))

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, idTableAddress},
			[]crypto.Signer{b.ServiceKey().Signer(), IDTableSigner},
			false,
		)

		result := executeScriptAndCheck(t, b, templates.GenerateGetCutPercentageScript(env), nil)
		assertEqual(t, CadenceUFix64("0.10"), result)

		tx = createTxWithTemplateAndAuthorizer(b, templates.GenerateChangeCutScript(env), idTableAddress)

		_ = tx.AddArgument(CadenceUFix64("0.08"))

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, idTableAddress},
			[]crypto.Signer{b.ServiceKey().Signer(), IDTableSigner},
			false,
		)

		result = executeScriptAndCheck(t, b, templates.GenerateGetCutPercentageScript(env), nil)
		assertEqual(t, CadenceUFix64("0.08"), result)
	})

	t.Run("Should be able to change the weekly payout", func(t *testing.T) {

		tx := createTxWithTemplateAndAuthorizer(b, templates.GenerateChangePayoutScript(env), idTableAddress)

		_ = tx.AddArgument(CadenceUFix64("5000000.0"))

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, idTableAddress},
			[]crypto.Signer{b.ServiceKey().Signer(), IDTableSigner},
			false,
		)

		result := executeScriptAndCheck(t, b, templates.GenerateGetWeeklyPayoutScript(env), nil)
		assertEqual(t, CadenceUFix64("5000000.0"), result)

		tx = createTxWithTemplateAndAuthorizer(b, templates.GenerateChangePayoutScript(env), idTableAddress)

		_ = tx.AddArgument(CadenceUFix64("1250000.0"))

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, idTableAddress},
			[]crypto.Signer{b.ServiceKey().Signer(), IDTableSigner},
			false,
		)

		result = executeScriptAndCheck(t, b, templates.GenerateGetWeeklyPayoutScript(env), nil)
		assertEqual(t, CadenceUFix64("1250000.0"), result)
	})

}

func TestIDTableStaking(t *testing.T) {

	t.Parallel()

	b := newBlockchain()

	env := templates.Environment{
		FungibleTokenAddress: emulatorFTAddress,
		FlowTokenAddress:     emulatorFlowTokenAddress,
	}

	accountKeys := test.AccountKeyGenerator()

	var totalPayout interpreter.UFix64Value = 125000000000000 // 1.25M
	var cutPercentage interpreter.UFix64Value = 8000000       // 8.0 %

	// Create new keys for the ID table account
	IDTableAccountKey, IDTableSigner := accountKeys.NewWithSigner()
	var idTableAddress = deployStakingContract(t, b, IDTableAccountKey, env)

	env.IDTableAddress = idTableAddress.Hex()

	var totalStaked interpreter.UFix64Value = 0

	committed := make(map[string]interpreter.UFix64Value)
	staked := make(map[string]interpreter.UFix64Value)
	request := make(map[string]interpreter.UFix64Value)
	unstaking := make(map[string]interpreter.UFix64Value)
	unstaked := make(map[string]interpreter.UFix64Value)
	rewards := make(map[string]interpreter.UFix64Value)

	// Create new user accounts
	joshAccountKey, joshSigner := accountKeys.NewWithSigner()
	joshAddress, _ := b.CreateAccount([]*flow.AccountKey{joshAccountKey}, nil)

	// Create a new user account
	maxAccountKey, maxSigner := accountKeys.NewWithSigner()
	maxAddress, _ := b.CreateAccount([]*flow.AccountKey{maxAccountKey}, nil)

	// Create a new user account
	bastianAccountKey, bastianSigner := accountKeys.NewWithSigner()
	bastianAddress, _ := b.CreateAccount([]*flow.AccountKey{bastianAccountKey}, nil)

	// Create a new user account for access node
	accessAccountKey, accessSigner := accountKeys.NewWithSigner()
	accessAddress, _ := b.CreateAccount([]*flow.AccountKey{accessAccountKey}, nil)

	// Create new delegator user accounts
	adminDelegatorAccountKey, adminDelegatorSigner := accountKeys.NewWithSigner()
	adminDelegatorAddress, _ := b.CreateAccount([]*flow.AccountKey{adminDelegatorAccountKey}, nil)

	joshDelegatorOneAccountKey, joshDelegatorOneSigner := accountKeys.NewWithSigner()
	joshDelegatorOneAddress, _ := b.CreateAccount([]*flow.AccountKey{joshDelegatorOneAccountKey}, nil)

	maxDelegatorOneAccountKey, maxDelegatorOneSigner := accountKeys.NewWithSigner()
	maxDelegatorOneAddress, _ := b.CreateAccount([]*flow.AccountKey{maxDelegatorOneAccountKey}, nil)

	maxDelegatorTwoAccountKey, maxDelegatorTwoSigner := accountKeys.NewWithSigner()
	maxDelegatorTwoAddress, _ := b.CreateAccount([]*flow.AccountKey{maxDelegatorTwoAccountKey}, nil)

	t.Run("Should be able to mint tokens for new accounts", func(t *testing.T) {

		mintTokensForAccount(t, b, idTableAddress)

		mintTokensForAccount(t, b, joshAddress)

		mintTokensForAccount(t, b, maxAddress)

		mintTokensForAccount(t, b, accessAddress)

		mintTokensForAccount(t, b, bastianAddress)

		mintTokensForAccount(t, b, maxDelegatorOneAddress)

		mintTokensForAccount(t, b, maxDelegatorTwoAddress)

		mintTokensForAccount(t, b, joshDelegatorOneAddress)

		mintTokensForAccount(t, b, adminDelegatorAddress)

	})

	t.Run("Shouldn't be able to create invalid Node structs", func(t *testing.T) {

		var amountToCommit interpreter.UFix64Value = 25000000000000

		registerNode(t, b, env,
			idTableAddress,
			IDTableSigner,
			// Invalid ID: Too short
			"3039",
			fmt.Sprintf("%0128d", admin),
			fmt.Sprintf("%0128d", admin),
			fmt.Sprintf("%0192d", admin),
			amountToCommit,
			committed[adminID],
			1,
			true)

		registerNode(t, b, env,
			idTableAddress,
			IDTableSigner,
			adminID,
			fmt.Sprintf("%0128d", admin),
			fmt.Sprintf("%0128d", admin),
			fmt.Sprintf("%0192d", admin),
			amountToCommit,
			committed[adminID],
			// Invalid Role: Greater than 5
			6,
			true)

		registerNode(t, b, env,
			idTableAddress,
			IDTableSigner,
			adminID,
			fmt.Sprintf("%0128d", admin),
			fmt.Sprintf("%0128d", admin),
			fmt.Sprintf("%0192d", admin),
			amountToCommit,
			committed[adminID],
			// Invalid Role: Less than 1
			0,
			true)

		registerNode(t, b, env,
			idTableAddress,
			IDTableSigner,
			adminID,
			// Invalid Networking Address: Length cannot be zero
			"",
			fmt.Sprintf("%0128d", admin),
			fmt.Sprintf("%0192d", admin),
			amountToCommit,
			committed[adminID],
			1,
			true)
	})

	t.Run("Should be able to create a valid Node struct and not create duplicates", func(t *testing.T) {

		var amountToCommit interpreter.UFix64Value = 25000000000000

		committed[adminID] = registerNode(t, b, env,
			idTableAddress,
			IDTableSigner,
			adminID,
			// Invalid Networking Address: Length cannot be zero
			fmt.Sprintf("%0128d", admin),
			fmt.Sprintf("%0128d", admin),
			fmt.Sprintf("%0192d", admin),
			amountToCommit,
			committed[adminID],
			1,
			false)

		result := executeScriptAndCheck(t, b, templates.GenerateReturnTableScript(env), nil)

		idArray := result.(cadence.Array).Values
		assert.Len(t, idArray, 1)

		result = executeScriptAndCheck(t, b, templates.GenerateReturnProposedTableScript(env), nil)

		idArray = result.(cadence.Array).Values
		assert.Len(t, idArray, 1)

		result = executeScriptAndCheck(t, b, templates.GenerateGetRoleScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(adminID))})
		assertEqual(t, cadence.NewUInt8(1), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetNetworkingAddressScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(adminID))})
		assertEqual(t, cadence.NewString(fmt.Sprintf("%0128d", admin)), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetNetworkingKeyScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(adminID))})
		assertEqual(t, cadence.NewString(fmt.Sprintf("%0128d", admin)), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetStakingKeyScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(adminID))})
		assertEqual(t, cadence.NewString(fmt.Sprintf("%0192d", admin)), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetInitialWeightScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(adminID))})
		assertEqual(t, cadence.NewUInt64(0), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetStakedBalanceScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(adminID))})
		assertEqual(t, CadenceUFix64(staked[adminID].String()), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetCommittedBalanceScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(adminID))})
		assertEqual(t, CadenceUFix64(committed[adminID].String()), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetUnstakedBalanceScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(adminID))})
		assertEqual(t, CadenceUFix64(unstaked[adminID].String()), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetUnstakingBalanceScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(adminID))})
		assertEqual(t, CadenceUFix64(unstaking[adminID].String()), result)

		registerNode(t, b, env,
			idTableAddress,
			IDTableSigner,
			// Invalid: Admin ID is already in use
			adminID,
			fmt.Sprintf("%0128d", admin),
			fmt.Sprintf("%0128d", admin),
			fmt.Sprintf("%0192d", admin),
			amountToCommit,
			committed[adminID],
			1,
			true)

		registerNode(t, b, env,
			joshAddress,
			joshSigner,
			joshID,
			// Invalid: first admin networking address is already in use
			fmt.Sprintf("%0128d", admin),
			fmt.Sprintf("%0128d", admin),
			fmt.Sprintf("%0192d", admin),
			amountToCommit,
			committed[adminID],
			1,
			true)

		registerNode(t, b, env,
			joshAddress,
			joshSigner,
			joshID,
			fmt.Sprintf("%0128d", josh),
			// Invalid: first admin networking key is already in use
			fmt.Sprintf("%0128d", admin),
			fmt.Sprintf("%0192d", josh),
			amountToCommit,
			committed[adminID],
			1,
			true)

		registerNode(t, b, env,
			joshAddress,
			joshSigner,
			joshID,
			fmt.Sprintf("%0128d", josh),
			fmt.Sprintf("%0128d", josh),
			// Invalid: first admin stake key is already in use
			fmt.Sprintf("%0192d", admin),
			amountToCommit,
			committed[adminID],
			1,
			true)

	})

	t.Run("Shouldn't be able to remove a Node that doesn't exist", func(t *testing.T) {

		tx := createTxWithTemplateAndAuthorizer(b, templates.GenerateRemoveNodeScript(env), idTableAddress)

		_ = tx.AddArgument(cadence.NewString(nonexistantID))

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, idTableAddress},
			[]crypto.Signer{b.ServiceKey().Signer(), IDTableSigner},
			true,
		)
	})

	t.Run("Should be able to create more valid Node structs", func(t *testing.T) {

		var amountToCommit interpreter.UFix64Value = 48000000000000

		committed[joshID] = registerNode(t, b, env,
			joshAddress,
			joshSigner,
			joshID,
			fmt.Sprintf("%0128d", josh),
			fmt.Sprintf("%0128d", josh),
			fmt.Sprintf("%0192d", josh),
			amountToCommit,
			committed[joshID],
			2,
			false)

		result := executeScriptAndCheck(t, b, templates.GenerateGetCommittedBalanceScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(joshID))})
		assertEqual(t, CadenceUFix64(committed[joshID].String()), result)

		amountToCommit = 135000000000000

		committed[maxID] = registerNode(t, b, env,
			maxAddress,
			maxSigner,
			maxID,
			fmt.Sprintf("%0128d", max),
			fmt.Sprintf("%0128d", max),
			fmt.Sprintf("%0192d", max),
			amountToCommit,
			committed[maxID],
			3,
			false)

		amountToCommit = 5000000000000

		committed[accessID] = registerNode(t, b, env,
			accessAddress,
			accessSigner,
			accessID,
			fmt.Sprintf("%0128d", access),
			fmt.Sprintf("%0128d", access),
			fmt.Sprintf("%0192d", access),
			amountToCommit,
			committed[accessID],
			5,
			false)

		result = executeScriptAndCheck(t, b, templates.GenerateReturnCurrentTableScript(env), nil)

		idArray := result.(cadence.Array).Values
		assert.Len(t, idArray, 0)

		result = executeScriptAndCheck(t, b, templates.GenerateReturnProposedTableScript(env), nil)

		idArray = result.(cadence.Array).Values
		assert.Len(t, idArray, 3)
	})

	t.Run("Should be able to remove a Node from the proposed record and add it back", func(t *testing.T) {

		tx := createTxWithTemplateAndAuthorizer(b, templates.GenerateRemoveNodeScript(env), idTableAddress)

		_ = tx.AddArgument(cadence.NewString(joshID))

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, idTableAddress},
			[]crypto.Signer{b.ServiceKey().Signer(), IDTableSigner},
			false,
		)

		result := executeScriptAndCheck(t, b, templates.GenerateReturnCurrentTableScript(env), nil)

		idArray := result.(cadence.Array).Values
		assert.Len(t, idArray, 0)

		result = executeScriptAndCheck(t, b, templates.GenerateReturnProposedTableScript(env), nil)

		idArray = result.(cadence.Array).Values
		assert.Len(t, idArray, 3)

		var amountToCommit interpreter.UFix64Value = 48000000000000

		registerNode(t, b, env,
			idTableAddress,
			IDTableSigner,
			joshID,
			fmt.Sprintf("%0128d", josh),
			fmt.Sprintf("%0128d", josh),
			fmt.Sprintf("%0192d", josh),
			amountToCommit,
			committed[joshID],
			2,
			false)

		result = executeScriptAndCheck(t, b, templates.GenerateReturnProposedTableScript(env), nil)

		idArray = result.(cadence.Array).Values
		assert.Len(t, idArray, 3)

		result = executeScriptAndCheck(t, b, templates.GenerateGetCommittedBalanceScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(joshID))})
		assertEqual(t, CadenceUFix64(committed[joshID].String()), result)
	})

	t.Run("Should be able to commit additional tokens for max's node", func(t *testing.T) {

		var amountToCommit interpreter.UFix64Value = 10000000000000

		committed[maxID] = commitNewTokens(t, b, env,
			maxAddress,
			maxSigner,
			amountToCommit,
			committed[maxID],
			false)

		result := executeScriptAndCheck(t, b, templates.GenerateGetCommittedBalanceScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(maxID))})
		assertEqual(t, CadenceUFix64(committed[maxID].String()), result)
	})

	t.Run("Should not be able request unstaking for more than is available", func(t *testing.T) {

		var amountToRequest interpreter.UFix64Value = 500000000000000

		requestUnstaking(t, b, env,
			maxAddress,
			maxSigner,
			amountToRequest,
			committed[maxID],
			unstaked[maxID],
			request[maxID],
			true,
		)
	})

	t.Run("Should be able to request unstaking which moves from comitted to unstaked", func(t *testing.T) {

		var amountToRequest interpreter.UFix64Value = 10000000000000

		committed[maxID], unstaked[maxID], request[maxID] = requestUnstaking(t, b, env,
			maxAddress,
			maxSigner,
			amountToRequest,
			committed[maxID],
			unstaked[maxID],
			request[maxID],
			false,
		)

		result := executeScriptAndCheck(t, b, templates.GenerateGetCommittedBalanceScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(maxID))})
		assertEqual(t, CadenceUFix64(committed[maxID].String()), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetUnstakedBalanceScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(maxID))})
		assertEqual(t, CadenceUFix64(unstaked[maxID].String()), result)
	})

	t.Run("Should be able to withdraw tokens from unstaked", func(t *testing.T) {

		unstaked[maxID] = unstaked[maxID] - 5000000000000

		tx := createTxWithTemplateAndAuthorizer(b, templates.GenerateWithdrawUnstakedTokensScript(env), maxAddress)

		err := tx.AddArgument(CadenceUFix64("50000.0"))
		require.NoError(t, err)

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, maxAddress},
			[]crypto.Signer{b.ServiceKey().Signer(), maxSigner},
			false,
		)

		result := executeScriptAndCheck(t, b, templates.GenerateGetUnstakedBalanceScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(maxID))})
		assertEqual(t, CadenceUFix64(unstaked[maxID].String()), result)
	})

	t.Run("Should be able to commit unstaked tokens", func(t *testing.T) {

		var amountToCommit interpreter.UFix64Value = 5000000000000

		committed[maxID], unstaked[maxID] = commitUnstaked(t, b, env,
			maxAddress,
			maxSigner,
			amountToCommit,
			committed[maxID],
			unstaked[maxID],
			false)

		result := executeScriptAndCheck(t, b, templates.GenerateGetUnstakedBalanceScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(maxID))})
		assertEqual(t, CadenceUFix64(unstaked[maxID].String()), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetCommittedBalanceScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(maxID))})
		assertEqual(t, CadenceUFix64(committed[maxID].String()), result)
	})

	t.Run("Should be able to end the staking auction, which removes insufficiently staked nodes", func(t *testing.T) {

		unstaked[joshID] = committed[joshID]
		committed[joshID] = 0

		tx := createTxWithTemplateAndAuthorizer(b, templates.GenerateEndStakingScript(env), idTableAddress)

		err := tx.AddArgument(cadence.NewArray([]cadence.Value{cadence.NewString(adminID), cadence.NewString(joshID), cadence.NewString(maxID), cadence.NewString(accessID)}))
		require.NoError(t, err)

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, idTableAddress},
			[]crypto.Signer{b.ServiceKey().Signer(), IDTableSigner},
			false,
		)

		result := executeScriptAndCheck(t, b, templates.GenerateReturnProposedTableScript(env), nil)

		idArray := result.(cadence.Array).Values
		assert.Len(t, idArray, 3)

		result = executeScriptAndCheck(t, b, templates.GenerateGetUnstakedBalanceScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(adminID))})
		assertEqual(t, CadenceUFix64(unstaked[adminID].String()), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetCommittedBalanceScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(adminID))})
		assertEqual(t, CadenceUFix64(committed[adminID].String()), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetUnstakedBalanceScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(joshID))})
		assertEqual(t, CadenceUFix64(unstaked[joshID].String()), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetCommittedBalanceScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(joshID))})
		assertEqual(t, CadenceUFix64(committed[joshID].String()), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetUnstakedBalanceScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(maxID))})
		assertEqual(t, CadenceUFix64(unstaked[maxID].String()), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetCommittedBalanceScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(maxID))})
		assertEqual(t, CadenceUFix64(committed[maxID].String()), result)

	})

	t.Run("Should pay rewards, but no balances are increased because nobody is staked", func(t *testing.T) {

		tx := createTxWithTemplateAndAuthorizer(b, templates.GeneratePayRewardsScript(env), idTableAddress)

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, idTableAddress},
			[]crypto.Signer{b.ServiceKey().Signer(), IDTableSigner},
			false,
		)

		result := executeScriptAndCheck(t, b, templates.GenerateGetUnstakedBalanceScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(adminID))})
		assertEqual(t, CadenceUFix64(unstaked[adminID].String()), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetUnstakedBalanceScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(joshID))})
		assertEqual(t, CadenceUFix64(unstaked[joshID].String()), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetUnstakedBalanceScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(maxID))})
		assertEqual(t, CadenceUFix64(unstaked[maxID].String()), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetRewardBalanceScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(adminID))})
		assertEqual(t, CadenceUFix64(rewards[adminID].String()), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetRewardBalanceScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(joshID))})
		assertEqual(t, CadenceUFix64(rewards[joshID].String()), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetRewardBalanceScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(maxID))})
		assertEqual(t, CadenceUFix64(rewards[maxID].String()), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetRewardBalanceScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(accessID))})
		assertEqual(t, CadenceUFix64(rewards[accessID].String()), result)

	})

	t.Run("Should Move committed tokens to staked buckets", func(t *testing.T) {

		tx := createTxWithTemplateAndAuthorizer(b, templates.GenerateMoveTokensScript(env), idTableAddress)

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, idTableAddress},
			[]crypto.Signer{b.ServiceKey().Signer(), IDTableSigner},
			false,
		)

		committed[adminID], staked[adminID], request[adminID], unstaking[adminID], unstaked[adminID], totalStaked = moveTokens(
			committed[adminID], staked[adminID], request[adminID], unstaking[adminID], unstaked[adminID], totalStaked)

		committed[joshID], staked[joshID], request[joshID], unstaking[joshID], unstaked[joshID], totalStaked = moveTokens(
			committed[joshID], staked[joshID], request[joshID], unstaking[joshID], unstaked[joshID], totalStaked)

		committed[maxID], staked[maxID], request[maxID], unstaking[maxID], unstaked[maxID], totalStaked = moveTokens(
			committed[maxID], staked[maxID], request[maxID], unstaking[maxID], unstaked[maxID], totalStaked)

		committed[accessID], staked[accessID], request[accessID], unstaking[accessID], unstaked[accessID], totalStaked = moveTokens(
			committed[accessID], staked[accessID], request[accessID], unstaking[accessID], unstaked[accessID], totalStaked)

		result := executeScriptAndCheck(t, b, templates.GenerateGetUnstakedBalanceScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(adminID))})
		assertEqual(t, CadenceUFix64(unstaked[adminID].String()), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetUnstakedBalanceScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(joshID))})
		assertEqual(t, CadenceUFix64(unstaked[joshID].String()), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetUnstakedBalanceScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(maxID))})
		assertEqual(t, CadenceUFix64(unstaked[maxID].String()), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetCommittedBalanceScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(adminID))})
		assertEqual(t, CadenceUFix64(committed[adminID].String()), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetCommittedBalanceScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(joshID))})
		assertEqual(t, CadenceUFix64(committed[joshID].String()), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetCommittedBalanceScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(maxID))})
		assertEqual(t, CadenceUFix64(committed[maxID].String()), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetStakedBalanceScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(adminID))})
		assertEqual(t, CadenceUFix64(staked[adminID].String()), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetStakedBalanceScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(joshID))})
		assertEqual(t, CadenceUFix64(staked[joshID].String()), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetStakedBalanceScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(maxID))})
		assertEqual(t, CadenceUFix64(staked[maxID].String()), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetStakedBalanceScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(accessID))})
		assertEqual(t, CadenceUFix64(staked[accessID].String()), result)

	})

	t.Run("Should be able to commit unstaked and new tokens from the node who was not included", func(t *testing.T) {

		var amountToCommit interpreter.UFix64Value = 48000000000000

		committed[joshID], unstaked[joshID] = commitUnstaked(t, b, env,
			joshAddress,
			joshSigner,
			amountToCommit,
			committed[joshID],
			unstaked[joshID],
			false)

		amountToCommit = 10000000000000

		committed[joshID] = commitNewTokens(t, b, env,
			joshAddress,
			joshSigner,
			amountToCommit,
			committed[joshID],
			false)

		result := executeScriptAndCheck(t, b, templates.GenerateReturnProposedTableScript(env), nil)

		idArray := result.(cadence.Array).Values
		assert.Len(t, idArray, 4)

		result = executeScriptAndCheck(t, b, templates.GenerateGetUnstakedBalanceScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(joshID))})
		assertEqual(t, CadenceUFix64(unstaked[joshID].String()), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetCommittedBalanceScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(joshID))})
		assertEqual(t, CadenceUFix64(committed[joshID].String()), result)
	})

	t.Run("Should be able to request unstaking from a staked node", func(t *testing.T) {

		var amountToRequest interpreter.UFix64Value = 5000000000000

		committed[adminID], unstaked[adminID], request[adminID] = requestUnstaking(t, b, env,
			idTableAddress,
			IDTableSigner,
			amountToRequest,
			committed[adminID],
			unstaked[adminID],
			request[adminID],
			false,
		)

		result := executeScriptAndCheck(t, b, templates.GenerateGetUnstakingRequestScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(adminID))})
		assertEqual(t, CadenceUFix64(request[adminID].String()), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetUnstakedBalanceScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(adminID))})
		assertEqual(t, CadenceUFix64(unstaked[adminID].String()), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetStakedBalanceScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(adminID))})
		assertEqual(t, CadenceUFix64(staked[adminID].String()), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetTotalCommitmentBalanceScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(adminID))})
		assertEqual(t, CadenceUFix64((committed[adminID].Plus(staked[adminID].Minus(request[adminID]))).String()), result)

		// josh, max, and access are proposed
		result = executeScriptAndCheck(t, b, templates.GenerateReturnProposedTableScript(env), nil)

		idArray := result.(cadence.Array).Values
		assert.Len(t, idArray, 3)

		// admin, max, and access are staked
		result = executeScriptAndCheck(t, b, templates.GenerateReturnCurrentTableScript(env), nil)

		idArray = result.(cadence.Array).Values
		assert.Len(t, idArray, 3)
	})

	/************* Start of Delegation Tests *******************/

	t.Run("Should be able to register first account to delegate to max", func(t *testing.T) {

		tx := createTxWithTemplateAndAuthorizer(b, templates.GenerateRegisterDelegatorScript(env), maxDelegatorOneAddress)

		err := tx.AddArgument(cadence.String(maxID))
		require.NoError(t, err)

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, maxDelegatorOneAddress},
			[]crypto.Signer{b.ServiceKey().Signer(), maxDelegatorOneSigner},
			false,
		)
	})

	t.Run("Should be able to register second account to delegate to max", func(t *testing.T) {

		tx := createTxWithTemplateAndAuthorizer(b, templates.GenerateRegisterDelegatorScript(env), maxDelegatorTwoAddress)

		err := tx.AddArgument(cadence.String(maxID))
		require.NoError(t, err)

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, maxDelegatorTwoAddress},
			[]crypto.Signer{b.ServiceKey().Signer(), maxDelegatorTwoSigner},
			false,
		)
	})

	t.Run("Should be able to register account to delegate to josh", func(t *testing.T) {

		tx := createTxWithTemplateAndAuthorizer(b, templates.GenerateRegisterDelegatorScript(env), joshDelegatorOneAddress)

		err := tx.AddArgument(cadence.String(joshID))
		require.NoError(t, err)

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, joshDelegatorOneAddress},
			[]crypto.Signer{b.ServiceKey().Signer(), joshDelegatorOneSigner},
			false,
		)
	})

	t.Run("Should not be able to register account to delegate to the admin address, because it has insufficient stake committed", func(t *testing.T) {

		tx := createTxWithTemplateAndAuthorizer(b, templates.GenerateRegisterDelegatorScript(env), adminDelegatorAddress)

		err := tx.AddArgument(cadence.String(adminID))
		require.NoError(t, err)

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, adminDelegatorAddress},
			[]crypto.Signer{b.ServiceKey().Signer(), adminDelegatorSigner},
			true,
		)
	})

	t.Run("Should not be able to register account to delegate to the access node, because access nodes are not allowed to be delegated to", func(t *testing.T) {

		tx := createTxWithTemplateAndAuthorizer(b, templates.GenerateRegisterDelegatorScript(env), adminDelegatorAddress)

		err := tx.AddArgument(cadence.String(accessID))
		require.NoError(t, err)

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, adminDelegatorAddress},
			[]crypto.Signer{b.ServiceKey().Signer(), adminDelegatorSigner},
			true,
		)
	})

	t.Run("Should be able to delegate new tokens to josh", func(t *testing.T) {

		tx := createTxWithTemplateAndAuthorizer(b, templates.GenerateDelegatorStakeNewScript(env), joshDelegatorOneAddress)
		_ = tx.AddArgument(CadenceUFix64("100000.0"))

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, joshDelegatorOneAddress},
			[]crypto.Signer{b.ServiceKey().Signer(), joshDelegatorOneSigner},
			false,
		)

		committed[joshID+firstDelegatorStringID] = 10000000000000

		result := executeScriptAndCheck(t, b, templates.GenerateGetCommittedBalanceScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(joshID))})
		assertEqual(t, CadenceUFix64(committed[joshID].String()), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetTotalCommitmentBalanceScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(joshID))})
		assertEqual(t, CadenceUFix64((committed[joshID].Plus(committed[joshID+firstDelegatorStringID])).(interpreter.UFix64Value).String()), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetDelegatorCommittedScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(joshID)), jsoncdc.MustEncode(cadence.UInt32(firstDelegatorID))})
		assertEqual(t, CadenceUFix64(committed[joshID+firstDelegatorStringID].String()), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetDelegatorStakedScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(joshID)), jsoncdc.MustEncode(cadence.UInt32(firstDelegatorID))})
		assertEqual(t, CadenceUFix64(staked[joshID+firstDelegatorStringID].String()), result)

	})

	t.Run("Should be able to request unstake delegated tokens from Josh, which moves them from committed to unstaked", func(t *testing.T) {

		var amountToUnstake interpreter.UFix64Value = 4000000000000
		committed[joshID+firstDelegatorStringID] = committed[joshID+firstDelegatorStringID].Minus(amountToUnstake).(interpreter.UFix64Value)
		unstaked[joshID+firstDelegatorStringID] = unstaked[joshID+firstDelegatorStringID].Plus(amountToUnstake).(interpreter.UFix64Value)

		tx := createTxWithTemplateAndAuthorizer(b, templates.GenerateDelegatorRequestUnstakeScript(env), joshDelegatorOneAddress)
		_ = tx.AddArgument(CadenceUFix64(amountToUnstake.String()))

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, joshDelegatorOneAddress},
			[]crypto.Signer{b.ServiceKey().Signer(), joshDelegatorOneSigner},
			false,
		)

		result := executeScriptAndCheck(t, b, templates.GenerateGetCommittedBalanceScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(joshID))})
		assertEqual(t, CadenceUFix64(committed[joshID].String()), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetTotalCommitmentBalanceScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(joshID))})
		assertEqual(t, CadenceUFix64((committed[joshID].Plus(committed[joshID+firstDelegatorStringID]).(interpreter.UFix64Value)).String()), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetDelegatorCommittedScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(joshID)), jsoncdc.MustEncode(cadence.UInt32(firstDelegatorID))})
		assertEqual(t, CadenceUFix64(committed[joshID+firstDelegatorStringID].String()), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetDelegatorUnstakedScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(joshID)), jsoncdc.MustEncode(cadence.UInt32(firstDelegatorID))})
		assertEqual(t, CadenceUFix64(unstaked[joshID+firstDelegatorStringID].String()), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetDelegatorRequestScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(joshID)), jsoncdc.MustEncode(cadence.UInt32(firstDelegatorID))})
		assertEqual(t, CadenceUFix64(request[joshID+firstDelegatorStringID].String()), result)

	})

	t.Run("Should be able to withdraw josh delegator's unstaked tokens", func(t *testing.T) {

		var amountToWithdraw interpreter.UFix64Value = 2000000000000

		unstaked[joshID+firstDelegatorStringID] = unstaked[joshID+firstDelegatorStringID].Minus(amountToWithdraw).(interpreter.UFix64Value)

		tx := createTxWithTemplateAndAuthorizer(b, templates.GenerateDelegatorWithdrawUnstakedScript(env), joshDelegatorOneAddress)

		_ = tx.AddArgument(CadenceUFix64(amountToWithdraw.String()))

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, joshDelegatorOneAddress},
			[]crypto.Signer{b.ServiceKey().Signer(), joshDelegatorOneSigner},
			false,
		)

		result := executeScriptAndCheck(t, b, templates.GenerateGetDelegatorCommittedScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(joshID)), jsoncdc.MustEncode(cadence.UInt32(firstDelegatorID))})
		assertEqual(t, CadenceUFix64(committed[joshID+firstDelegatorStringID].String()), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetDelegatorUnstakedScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(joshID)), jsoncdc.MustEncode(cadence.UInt32(firstDelegatorID))})
		assertEqual(t, CadenceUFix64(unstaked[joshID+firstDelegatorStringID].String()), result)

	})

	t.Run("Should be able to delegate unstaked tokens to josh", func(t *testing.T) {

		var amountToCommit interpreter.UFix64Value = 2000000000000

		unstaked[joshID+firstDelegatorStringID] = unstaked[joshID+firstDelegatorStringID].Minus(amountToCommit).(interpreter.UFix64Value)
		committed[joshID+firstDelegatorStringID] = committed[joshID+firstDelegatorStringID].Plus(amountToCommit).(interpreter.UFix64Value)

		tx := createTxWithTemplateAndAuthorizer(b, templates.GenerateDelegatorStakeUnstakedScript(env), joshDelegatorOneAddress)

		_ = tx.AddArgument(CadenceUFix64(amountToCommit.String()))

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, joshDelegatorOneAddress},
			[]crypto.Signer{b.ServiceKey().Signer(), joshDelegatorOneSigner},
			false,
		)

		result := executeScriptAndCheck(t, b, templates.GenerateGetCommittedBalanceScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(joshID))})
		assertEqual(t, CadenceUFix64(committed[joshID].String()), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetTotalCommitmentBalanceScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(joshID))})
		assertEqual(t, CadenceUFix64("660000.0"), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetDelegatorCommittedScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(joshID)), jsoncdc.MustEncode(cadence.UInt32(firstDelegatorID))})
		assertEqual(t, CadenceUFix64(committed[joshID+firstDelegatorStringID].String()), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetDelegatorUnstakedScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(joshID)), jsoncdc.MustEncode(cadence.UInt32(firstDelegatorID))})
		assertEqual(t, CadenceUFix64(unstaked[joshID+firstDelegatorStringID].String()), result)

	})

	t.Run("Should be able to end the staking auction, which marks admin to unstake", func(t *testing.T) {

		tx := createTxWithTemplateAndAuthorizer(b, templates.GenerateEndStakingScript(env), idTableAddress)

		err := tx.AddArgument(cadence.NewArray([]cadence.Value{cadence.NewString(adminID), cadence.NewString(joshID), cadence.NewString(maxID), cadence.NewString(accessID)}))
		require.NoError(t, err)

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, idTableAddress},
			[]crypto.Signer{b.ServiceKey().Signer(), IDTableSigner},
			false,
		)

		request[adminID] = 25000000000000

		result := executeScriptAndCheck(t, b, templates.GenerateReturnProposedTableScript(env), nil)

		idArray := result.(cadence.Array).Values
		assert.Len(t, idArray, 3)

		result = executeScriptAndCheck(t, b, templates.GenerateReturnCurrentTableScript(env), nil)

		idArray = result.(cadence.Array).Values
		assert.Len(t, idArray, 3)

		result = executeScriptAndCheck(t, b, templates.GenerateGetUnstakingRequestScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(adminID))})
		assertEqual(t, CadenceUFix64(request[adminID].String()), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetUnstakedBalanceScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(adminID))})
		assertEqual(t, CadenceUFix64("0.0"), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetCommittedBalanceScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(adminID))})
		assertEqual(t, CadenceUFix64(committed[adminID].String()), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetUnstakedBalanceScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(joshID))})
		assertEqual(t, CadenceUFix64("0.0"), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetCommittedBalanceScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(joshID))})
		assertEqual(t, CadenceUFix64("580000.0"), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetStakedBalanceScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(joshID))})
		assertEqual(t, CadenceUFix64("0.0"), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetTotalTokensStakedScript(env), nil)
		assertEqual(t, CadenceUFix64("1650000.0"), result)
	})

	t.Run("Should pay correct rewards, no delegators are paid because none are staked yet", func(t *testing.T) {

		tx := createTxWithTemplateAndAuthorizer(b, templates.GeneratePayRewardsScript(env), idTableAddress)

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, idTableAddress},
			[]crypto.Signer{b.ServiceKey().Signer(), IDTableSigner},
			false,
		)

		totalStaked = 165000000000000

		result := executeScriptAndCheck(t, b, templates.GenerateGetUnstakedBalanceScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(adminID))})
		assertEqual(t, CadenceUFix64(unstaked[adminID].String()), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetUnstakedBalanceScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(joshID))})
		assertEqual(t, CadenceUFix64(unstaked[joshID].String()), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetUnstakedBalanceScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(maxID))})
		assertEqual(t, CadenceUFix64("0.0"), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetStakedBalanceScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(adminID))})
		assertEqual(t, CadenceUFix64(staked[adminID].String()), result)

		rewardsResult, _ := payRewards(false, totalPayout, totalStaked, cutPercentage, staked[adminID])
		rewards[adminID] = rewards[adminID].Plus(rewardsResult).(interpreter.UFix64Value)

		result = executeScriptAndCheck(t, b, templates.GenerateGetRewardBalanceScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(adminID))})
		assertEqual(t, CadenceUFix64(rewards[adminID].String()), result)

		rewardsResult, _ = payRewards(false, totalPayout, totalStaked, cutPercentage, 0)
		rewards[joshID] = rewards[joshID].Plus(rewardsResult).(interpreter.UFix64Value)

		rewardsResult, delegateeRewardsResult := payRewards(true, totalPayout, totalStaked, cutPercentage, 0)
		rewards[joshID] = rewards[joshID].Plus(delegateeRewardsResult).(interpreter.UFix64Value)
		rewards[joshID+firstDelegatorStringID] = rewards[joshID+firstDelegatorStringID].Plus(rewardsResult).(interpreter.UFix64Value)

		result = executeScriptAndCheck(t, b, templates.GenerateGetRewardBalanceScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(joshID))})
		assertEqual(t, CadenceUFix64(rewards[joshID].String()), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetDelegatorRewardsScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(joshID)), jsoncdc.MustEncode(cadence.UInt32(firstDelegatorID))})
		assertEqual(t, CadenceUFix64(rewards[joshID+firstDelegatorStringID].String()), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetStakedBalanceScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(maxID))})
		assertEqual(t, CadenceUFix64("1400000.0"), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetTotalTokensStakedByTypeScript(env), [][]byte{jsoncdc.MustEncode(cadence.UInt8(3))})
		assertEqual(t, CadenceUFix64("1400000.0"), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetDelegatorStakedScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(maxID)), jsoncdc.MustEncode(cadence.UInt32(firstDelegatorID))})
		assertEqual(t, CadenceUFix64("0.0"), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetDelegatorStakedScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(maxID)), jsoncdc.MustEncode(cadence.UInt32(secondDelegatorID))})
		assertEqual(t, CadenceUFix64("0.0"), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetRewardBalanceScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(maxID))})
		assertEqual(t, CadenceUFix64("1060606.05"), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetDelegatorRewardsScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(maxID)), jsoncdc.MustEncode(cadence.UInt32(firstDelegatorID))})
		assertEqual(t, CadenceUFix64("0.0"), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetRewardBalanceScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(accessID))})
		assertEqual(t, CadenceUFix64("0.0"), result)

	})

	t.Run("Should Move tokens between buckets", func(t *testing.T) {

		tx := createTxWithTemplateAndAuthorizer(b, templates.GenerateMoveTokensScript(env), idTableAddress)

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, idTableAddress},
			[]crypto.Signer{b.ServiceKey().Signer(), IDTableSigner},
			false,
		)

		committed[adminID], staked[adminID], request[adminID], unstaking[adminID], unstaked[adminID], totalStaked = moveTokens(
			committed[adminID], staked[adminID], request[adminID], unstaking[adminID], unstaked[adminID], totalStaked)

		committed[joshID], staked[joshID], request[joshID], unstaking[joshID], unstaked[joshID], totalStaked = moveTokens(
			committed[joshID], staked[joshID], request[joshID], unstaking[joshID], unstaked[joshID], totalStaked)

		committed[joshID+firstDelegatorStringID], staked[joshID+firstDelegatorStringID], request[joshID+firstDelegatorStringID], unstaking[joshID+firstDelegatorStringID], unstaked[joshID+firstDelegatorStringID], totalStaked = moveTokens(
			committed[joshID+firstDelegatorStringID], staked[joshID+firstDelegatorStringID], request[joshID+firstDelegatorStringID], unstaking[joshID+firstDelegatorStringID], unstaked[joshID+firstDelegatorStringID], totalStaked)

		committed[maxID], staked[maxID], request[maxID], unstaking[maxID], unstaked[maxID], totalStaked = moveTokens(
			committed[maxID], staked[maxID], request[maxID], unstaking[maxID], unstaked[maxID], totalStaked)

		committed[accessID], staked[accessID], request[accessID], unstaking[accessID], unstaked[accessID], totalStaked = moveTokens(
			committed[accessID], staked[accessID], request[accessID], unstaking[accessID], unstaked[accessID], totalStaked)

		result := executeScriptAndCheck(t, b, templates.GenerateGetUnstakedBalanceScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(adminID))})
		assertEqual(t, CadenceUFix64(unstaked[adminID].String()), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetCommittedBalanceScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(adminID))})
		assertEqual(t, CadenceUFix64(committed[adminID].String()), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetStakedBalanceScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(adminID))})
		assertEqual(t, CadenceUFix64(staked[adminID].String()), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetUnstakingRequestScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(adminID))})
		assertEqual(t, CadenceUFix64(request[adminID].String()), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetUnstakingBalanceScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(adminID))})
		assertEqual(t, CadenceUFix64(unstaking[adminID].String()), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetRewardBalanceScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(adminID))})
		assertEqual(t, CadenceUFix64(rewards[adminID].String()), result)

		// josh buckets

		result = executeScriptAndCheck(t, b, templates.GenerateGetUnstakedBalanceScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(joshID))})
		assertEqual(t, CadenceUFix64(unstaked[joshID].String()), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetCommittedBalanceScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(joshID))})
		assertEqual(t, CadenceUFix64(committed[joshID].String()), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetStakedBalanceScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(joshID))})
		assertEqual(t, CadenceUFix64(staked[joshID].String()), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetUnstakingRequestScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(joshID))})
		assertEqual(t, CadenceUFix64(request[joshID].String()), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetUnstakingBalanceScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(joshID))})
		assertEqual(t, CadenceUFix64(unstaking[joshID].String()), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetRewardBalanceScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(joshID))})
		assertEqual(t, CadenceUFix64(rewards[joshID].String()), result)

		// Josh Delegator Buckets

		result = executeScriptAndCheck(t, b, templates.GenerateGetDelegatorUnstakedScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(joshID)), jsoncdc.MustEncode(cadence.UInt32(firstDelegatorID))})
		assertEqual(t, CadenceUFix64(unstaked[joshID+firstDelegatorStringID].String()), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetDelegatorCommittedScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(joshID)), jsoncdc.MustEncode(cadence.UInt32(firstDelegatorID))})
		assertEqual(t, CadenceUFix64(committed[joshID+firstDelegatorStringID].String()), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetDelegatorStakedScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(joshID)), jsoncdc.MustEncode(cadence.UInt32(firstDelegatorID))})
		assertEqual(t, CadenceUFix64(staked[joshID+firstDelegatorStringID].String()), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetDelegatorRequestScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(joshID)), jsoncdc.MustEncode(cadence.UInt32(firstDelegatorID))})
		assertEqual(t, CadenceUFix64(request[joshID+firstDelegatorStringID].String()), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetDelegatorUnstakingScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(joshID)), jsoncdc.MustEncode(cadence.UInt32(firstDelegatorID))})
		assertEqual(t, CadenceUFix64(unstaking[joshID+firstDelegatorStringID].String()), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetDelegatorRewardsScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(joshID)), jsoncdc.MustEncode(cadence.UInt32(firstDelegatorID))})
		assertEqual(t, CadenceUFix64(rewards[joshID+firstDelegatorStringID].String()), result)

		// Max buckets

		result = executeScriptAndCheck(t, b, templates.GenerateGetUnstakedBalanceScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(maxID))})
		assertEqual(t, CadenceUFix64("0.0"), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetCommittedBalanceScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(maxID))})
		assertEqual(t, CadenceUFix64("0.0"), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetStakedBalanceScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(maxID))})
		assertEqual(t, CadenceUFix64("1400000.0"), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetUnstakingRequestScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(maxID))})
		assertEqual(t, CadenceUFix64("0.0"), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetUnstakingBalanceScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(maxID))})
		assertEqual(t, CadenceUFix64("0.0"), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetRewardBalanceScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(maxID))})
		assertEqual(t, CadenceUFix64("1060606.05"), result)

		// Max Delegator Buckets

		result = executeScriptAndCheck(t, b, templates.GenerateGetDelegatorUnstakedScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(maxID)), jsoncdc.MustEncode(cadence.UInt32(firstDelegatorID))})
		assertEqual(t, CadenceUFix64("0.0"), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetDelegatorCommittedScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(maxID)), jsoncdc.MustEncode(cadence.UInt32(firstDelegatorID))})
		assertEqual(t, CadenceUFix64("0.0"), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetDelegatorStakedScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(maxID)), jsoncdc.MustEncode(cadence.UInt32(firstDelegatorID))})
		assertEqual(t, CadenceUFix64("0.0"), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetDelegatorRequestScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(maxID)), jsoncdc.MustEncode(cadence.UInt32(firstDelegatorID))})
		assertEqual(t, CadenceUFix64("0.0"), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetDelegatorUnstakingScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(maxID)), jsoncdc.MustEncode(cadence.UInt32(firstDelegatorID))})
		assertEqual(t, CadenceUFix64("0.0"), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetDelegatorRewardsScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(maxID)), jsoncdc.MustEncode(cadence.UInt32(firstDelegatorID))})
		assertEqual(t, CadenceUFix64("0.0"), result)

	})

	t.Run("Should create new execution node", func(t *testing.T) {

		tx := createTxWithTemplateAndAuthorizer(b, templates.GenerateRegisterNodeScript(env), bastianAddress)

		_ = tx.AddArgument(cadence.NewString(bastianID))
		_ = tx.AddArgument(cadence.NewUInt8(3))
		_ = tx.AddArgument(cadence.NewString(fmt.Sprintf("%0128d", bastian)))
		_ = tx.AddArgument(cadence.NewString(fmt.Sprintf("%0128d", bastian)))
		_ = tx.AddArgument(cadence.NewString(fmt.Sprintf("%0192d", bastian)))
		_ = tx.AddArgument(CadenceUFix64("1400000.0"))

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, bastianAddress},
			[]crypto.Signer{b.ServiceKey().Signer(), bastianSigner},
			false,
		)

		result := executeScriptAndCheck(t, b, templates.GenerateGetCommittedBalanceScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(bastianID))})
		assertEqual(t, CadenceUFix64("1400000.0"), result)
	})

	t.Run("Should be able to delegate new tokens to max", func(t *testing.T) {

		tx := createTxWithTemplateAndAuthorizer(b, templates.GenerateDelegatorStakeNewScript(env), maxDelegatorOneAddress)
		_ = tx.AddArgument(CadenceUFix64("100000.0"))

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, maxDelegatorOneAddress},
			[]crypto.Signer{b.ServiceKey().Signer(), maxDelegatorOneSigner},
			false,
		)

		tx = createTxWithTemplateAndAuthorizer(b, templates.GenerateDelegatorStakeNewScript(env), maxDelegatorTwoAddress)
		_ = tx.AddArgument(CadenceUFix64("200000.0"))

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, maxDelegatorTwoAddress},
			[]crypto.Signer{b.ServiceKey().Signer(), maxDelegatorTwoSigner},
			false,
		)

	})

	t.Run("Should not be able request unstaking below the minimum if a node has delegators staked or committed", func(t *testing.T) {

		result := executeScriptAndCheck(t, b, templates.GenerateGetTotalCommitmentBalanceScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(joshID))})
		assertEqual(t, CadenceUFix64("660000.0"), result)

		tx := createTxWithTemplateAndAuthorizer(b, templates.GenerateUnstakeTokensScript(env), joshAddress)

		err := tx.AddArgument(CadenceUFix64("180000.0"))
		require.NoError(t, err)

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, joshAddress},
			[]crypto.Signer{b.ServiceKey().Signer(), joshSigner},
			true,
		)

		tx = createTxWithTemplateAndAuthorizer(b, templates.GenerateUnstakeTokensScript(env), maxAddress)

		err = tx.AddArgument(CadenceUFix64("500000.0"))
		require.NoError(t, err)

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, maxAddress},
			[]crypto.Signer{b.ServiceKey().Signer(), maxSigner},
			true,
		)
	})

	t.Run("Should be able to commit rewarded tokens", func(t *testing.T) {

		var newCommitAmount uint64 = 50000

		tx := createTxWithTemplateAndAuthorizer(b, templates.GenerateStakeRewardedTokensScript(env), idTableAddress)

		err := tx.AddArgument(CadenceUFix64("50000.0"))
		require.NoError(t, err)

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, idTableAddress},
			[]crypto.Signer{b.ServiceKey().Signer(), IDTableSigner},
			false,
		)

		rewards[adminID] = rewards[adminID].Minus(interpreter.NewUFix64ValueWithInteger(newCommitAmount)).(interpreter.UFix64Value)
		committed[adminID] = committed[adminID].Plus(interpreter.NewUFix64ValueWithInteger(newCommitAmount)).(interpreter.UFix64Value)

		result := executeScriptAndCheck(t, b, templates.GenerateGetRewardBalanceScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(adminID))})
		assertEqual(t, CadenceUFix64(rewards[adminID].String()), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetCommittedBalanceScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(adminID))})
		assertEqual(t, CadenceUFix64(committed[adminID].String()), result)
	})

	// Josh Delegator Requests to unstake which marks their request
	t.Run("Should be able to request unstake delegated tokens from Josh, marks as requested", func(t *testing.T) {

		var requestAmount interpreter.UFix64Value = 4000000000000

		request[joshID+firstDelegatorStringID] = request[joshID+firstDelegatorStringID].Plus(requestAmount).(interpreter.UFix64Value)

		tx := createTxWithTemplateAndAuthorizer(b, templates.GenerateDelegatorRequestUnstakeScript(env), joshDelegatorOneAddress)
		_ = tx.AddArgument(CadenceUFix64(requestAmount.String()))

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, joshDelegatorOneAddress},
			[]crypto.Signer{b.ServiceKey().Signer(), joshDelegatorOneSigner},
			false,
		)

		result := executeScriptAndCheck(t, b, templates.GenerateGetCommittedBalanceScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(joshID))})
		assertEqual(t, CadenceUFix64(committed[joshID].String()), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetTotalCommitmentBalanceScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(joshID))})
		assertEqual(t, CadenceUFix64("620000.0"), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetDelegatorCommittedScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(joshID)), jsoncdc.MustEncode(cadence.UInt32(firstDelegatorID))})
		assertEqual(t, CadenceUFix64(committed[joshID+firstDelegatorStringID].String()), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetDelegatorUnstakedScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(joshID)), jsoncdc.MustEncode(cadence.UInt32(firstDelegatorID))})
		assertEqual(t, CadenceUFix64(unstaked[joshID+firstDelegatorStringID].String()), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetDelegatorRequestScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(joshID)), jsoncdc.MustEncode(cadence.UInt32(firstDelegatorID))})
		assertEqual(t, CadenceUFix64(request[joshID+firstDelegatorStringID].String()), result)

	})

	t.Run("Should be able cancel unstake request for delegator", func(t *testing.T) {

		var cancelRequestAmount interpreter.UFix64Value = 2000000000000

		request[joshID+firstDelegatorStringID] = request[joshID+firstDelegatorStringID].Minus(cancelRequestAmount).(interpreter.UFix64Value)

		tx := createTxWithTemplateAndAuthorizer(b, templates.GenerateDelegatorStakeUnstakedScript(env), joshDelegatorOneAddress)

		err := tx.AddArgument(CadenceUFix64(cancelRequestAmount.String()))
		require.NoError(t, err)

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, joshDelegatorOneAddress},
			[]crypto.Signer{b.ServiceKey().Signer(), joshDelegatorOneSigner},
			false,
		)

		result := executeScriptAndCheck(t, b, templates.GenerateGetDelegatorRequestScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(joshID)), jsoncdc.MustEncode(cadence.UInt32(firstDelegatorID))})
		assertEqual(t, CadenceUFix64(request[joshID+firstDelegatorStringID].String()), result)

	})

	// Josh Delegator Requests to unstake which marks their request
	t.Run("Should be able to request unstake delegated tokens from Josh, marks as requested", func(t *testing.T) {

		var requestAmount interpreter.UFix64Value = 2000000000000

		request[joshID+firstDelegatorStringID] = request[joshID+firstDelegatorStringID].Plus(requestAmount).(interpreter.UFix64Value)

		tx := createTxWithTemplateAndAuthorizer(b, templates.GenerateDelegatorRequestUnstakeScript(env), joshDelegatorOneAddress)
		_ = tx.AddArgument(CadenceUFix64(requestAmount.String()))

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, joshDelegatorOneAddress},
			[]crypto.Signer{b.ServiceKey().Signer(), joshDelegatorOneSigner},
			false,
		)

		result := executeScriptAndCheck(t, b, templates.GenerateGetDelegatorRequestScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(joshID)), jsoncdc.MustEncode(cadence.UInt32(firstDelegatorID))})
		assertEqual(t, CadenceUFix64(request[joshID+firstDelegatorStringID].String()), result)
	})

	// End the staking auction
	t.Run("Should be able to end the staking auction", func(t *testing.T) {

		tx := createTxWithTemplateAndAuthorizer(b, templates.GenerateEndStakingScript(env), idTableAddress)

		err := tx.AddArgument(cadence.NewArray([]cadence.Value{cadence.NewString(adminID), cadence.NewString(joshID), cadence.NewString(maxID), cadence.NewString(bastianID)}))
		require.NoError(t, err)

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, idTableAddress},
			[]crypto.Signer{b.ServiceKey().Signer(), IDTableSigner},
			false,
		)

		unstaked[adminID] = unstaked[adminID].Plus(committed[adminID]).(interpreter.UFix64Value)

	})

	// Move tokens between buckets. Josh delegator's should be in the unstaking bucket
	// also make sure that the total totens for the #3 node role is correct
	// Make sure that admin's unstaking were moved into their unstaked
	t.Run("Should Move tokens between buckets", func(t *testing.T) {

		tx := createTxWithTemplateAndAuthorizer(b, templates.GenerateMoveTokensScript(env), idTableAddress)

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, idTableAddress},
			[]crypto.Signer{b.ServiceKey().Signer(), IDTableSigner},
			false,
		)

		committed[adminID], _, request[adminID], unstaking[adminID], unstaked[adminID], totalStaked = moveTokens(
			committed[adminID], staked[adminID], request[adminID], unstaking[adminID], unstaked[adminID], totalStaked)

		committed[joshID], staked[joshID], request[joshID], unstaking[joshID], unstaked[joshID], totalStaked = moveTokens(
			committed[joshID], staked[joshID], request[joshID], unstaking[joshID], unstaked[joshID], totalStaked)

		committed[joshID+firstDelegatorStringID], staked[joshID+firstDelegatorStringID], request[joshID+firstDelegatorStringID], unstaking[joshID+firstDelegatorStringID], unstaked[joshID+firstDelegatorStringID], totalStaked = moveTokens(
			committed[joshID+firstDelegatorStringID], staked[joshID+firstDelegatorStringID], request[joshID+firstDelegatorStringID], unstaking[joshID+firstDelegatorStringID], unstaked[joshID+firstDelegatorStringID], totalStaked)

		committed[maxID], staked[maxID], request[maxID], unstaking[maxID], unstaked[maxID], totalStaked = moveTokens(
			committed[maxID], staked[maxID], request[maxID], unstaking[maxID], unstaked[maxID], totalStaked)

		committed[accessID], staked[accessID], request[accessID], unstaking[accessID], unstaked[accessID], totalStaked = moveTokens(
			committed[accessID], staked[accessID], request[accessID], unstaking[accessID], unstaked[accessID], totalStaked)

		result := executeScriptAndCheck(t, b, templates.GenerateGetStakedBalanceScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(adminID))})
		assertEqual(t, CadenceUFix64(staked[adminID].String()), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetUnstakedBalanceScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(adminID))})
		assertEqual(t, CadenceUFix64(unstaked[adminID].String()), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetStakedBalanceScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(joshID))})
		assertEqual(t, CadenceUFix64(staked[joshID].String()), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetDelegatorStakedScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(joshID)), jsoncdc.MustEncode(cadence.UInt32(firstDelegatorID))})
		assertEqual(t, CadenceUFix64(staked[joshID+firstDelegatorStringID].String()), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetDelegatorUnstakingScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(joshID)), jsoncdc.MustEncode(cadence.UInt32(firstDelegatorID))})
		assertEqual(t, CadenceUFix64(unstaking[joshID+firstDelegatorStringID].String()), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetTotalTokensStakedByTypeScript(env), [][]byte{jsoncdc.MustEncode(cadence.UInt8(2))})
		assertEqual(t, CadenceUFix64("620000.0"), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetStakedBalanceScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(maxID))})
		assertEqual(t, CadenceUFix64("1400000.0"), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetDelegatorStakedScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(maxID)), jsoncdc.MustEncode(cadence.UInt32(firstDelegatorID))})
		assertEqual(t, CadenceUFix64("100000.0"), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetDelegatorStakedScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(maxID)), jsoncdc.MustEncode(cadence.UInt32(secondDelegatorID))})
		assertEqual(t, CadenceUFix64("200000.0"), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetStakedBalanceScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(bastianID))})
		assertEqual(t, CadenceUFix64("1400000.0"), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetTotalTokensStakedByTypeScript(env), [][]byte{jsoncdc.MustEncode(cadence.UInt8(3))})
		assertEqual(t, CadenceUFix64("3100000.0"), result)

	})

	// Pay rewards and make sure josh and josh delegator got paid the right amounts based on the cut
	t.Run("Should pay correct rewards, rewards are split up properly between stakers and delegators", func(t *testing.T) {

		result := executeScriptAndCheck(t, b, templates.GenerateGetRewardBalanceScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(adminID))})
		assertEqual(t, CadenceUFix64(rewards[adminID].String()), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetRewardBalanceScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(joshID))})
		assertEqual(t, CadenceUFix64("0.0"), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetTotalTokensStakedScript(env), nil)
		assertEqual(t, CadenceUFix64("3720000.0"), result)

		totalStaked = 372000000000000

		tx := createTxWithTemplateAndAuthorizer(b, templates.GeneratePayRewardsScript(env), idTableAddress)

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, idTableAddress},
			[]crypto.Signer{b.ServiceKey().Signer(), IDTableSigner},
			false,
		)

		result = executeScriptAndCheck(t, b, templates.GenerateGetRewardBalanceScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(adminID))})
		assertEqual(t, CadenceUFix64(rewards[adminID].String()), result)

		rewardsResult, _ := payRewards(false, totalPayout, totalStaked, cutPercentage, staked[joshID])
		rewards[joshID] = rewards[joshID].Plus(rewardsResult).(interpreter.UFix64Value)

		rewardsResult, delegateeRewardsResult := payRewards(true, totalPayout, totalStaked, cutPercentage, staked[joshID+firstDelegatorStringID])
		rewards[joshID] = rewards[joshID].Plus(delegateeRewardsResult).(interpreter.UFix64Value)
		rewards[joshID+firstDelegatorStringID] = rewards[joshID+firstDelegatorStringID].Plus(rewardsResult).(interpreter.UFix64Value)

		result = executeScriptAndCheck(t, b, templates.GenerateGetRewardBalanceScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(joshID))})
		assertEqual(t, CadenceUFix64(rewards[joshID].String()), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetDelegatorRewardsScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(joshID)), jsoncdc.MustEncode(cadence.UInt32(firstDelegatorID))})
		assertEqual(t, CadenceUFix64(rewards[joshID+firstDelegatorStringID].String()), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetRewardBalanceScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(maxID))})
		assertEqual(t, CadenceUFix64("1539100.666"), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetDelegatorRewardsScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(maxID)), jsoncdc.MustEncode(cadence.UInt32(firstDelegatorID))})
		assertEqual(t, CadenceUFix64("30913.978"), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetDelegatorRewardsScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(maxID)), jsoncdc.MustEncode(cadence.UInt32(secondDelegatorID))})
		assertEqual(t, CadenceUFix64("61827.956"), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetRewardBalanceScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(bastianID))})
		assertEqual(t, CadenceUFix64("470430.1"), result)

	})

	// Move tokens. make sure josh delegators unstaking tokens are moved into their unstaked bucket
	t.Run("Should Move tokens between buckets", func(t *testing.T) {

		tx := createTxWithTemplateAndAuthorizer(b, templates.GenerateMoveTokensScript(env), idTableAddress)

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, idTableAddress},
			[]crypto.Signer{b.ServiceKey().Signer(), IDTableSigner},
			false,
		)

		result := executeScriptAndCheck(t, b, templates.GenerateGetDelegatorStakedScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(joshID)), jsoncdc.MustEncode(cadence.UInt32(firstDelegatorID))})
		assertEqual(t, CadenceUFix64("40000.0"), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetDelegatorUnstakingScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(joshID)), jsoncdc.MustEncode(cadence.UInt32(firstDelegatorID))})
		assertEqual(t, CadenceUFix64("0.0"), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetDelegatorUnstakedScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(joshID)), jsoncdc.MustEncode(cadence.UInt32(firstDelegatorID))})
		assertEqual(t, CadenceUFix64("40000.0"), result)

	})

	// Max Delegator Withdraws rewards
	t.Run("Should be able to withdraw delegator rewards", func(t *testing.T) {

		tx := createTxWithTemplateAndAuthorizer(b, templates.GenerateDelegatorWithdrawRewardsScript(env), maxDelegatorOneAddress)
		_ = tx.AddArgument(CadenceUFix64("2000.0"))

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, maxDelegatorOneAddress},
			[]crypto.Signer{b.ServiceKey().Signer(), maxDelegatorOneSigner},
			false,
		)

		result := executeScriptAndCheck(t, b, templates.GenerateGetDelegatorRewardsScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(maxID)), jsoncdc.MustEncode(cadence.UInt32(firstDelegatorID))})
		assertEqual(t, CadenceUFix64("28913.978"), result)

	})

	t.Run("Should commit more delegator tokens", func(t *testing.T) {

		tx := createTxWithTemplateAndAuthorizer(b, templates.GenerateDelegatorStakeNewScript(env), maxDelegatorOneAddress)
		_ = tx.AddArgument(CadenceUFix64("2000.0"))

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, maxDelegatorOneAddress},
			[]crypto.Signer{b.ServiceKey().Signer(), maxDelegatorOneSigner},
			false,
		)

	})

	t.Run("Should not be able request unstaking for less than the minimum, even if delegators make more than the minimum", func(t *testing.T) {

		result := executeScriptAndCheck(t, b, templates.GenerateGetStakedBalanceScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(maxID))})
		assertEqual(t, CadenceUFix64("1400000.0"), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetDelegatorStakedScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(maxID)), jsoncdc.MustEncode(cadence.UInt32(firstDelegatorID))})
		assertEqual(t, CadenceUFix64("100000.0"), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetDelegatorStakedScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(maxID)), jsoncdc.MustEncode(cadence.UInt32(secondDelegatorID))})
		assertEqual(t, CadenceUFix64("200000.0"), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetCommittedBalanceScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(maxID))})
		assertEqual(t, CadenceUFix64("0.0"), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetDelegatorCommittedScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(maxID)), jsoncdc.MustEncode(cadence.UInt32(firstDelegatorID))})
		assertEqual(t, CadenceUFix64("2000.0"), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetDelegatorCommittedScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(maxID)), jsoncdc.MustEncode(cadence.UInt32(secondDelegatorID))})
		assertEqual(t, CadenceUFix64("0.0"), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetUnstakingRequestScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(maxID))})
		assertEqual(t, CadenceUFix64("0.0"), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetDelegatorUnstakingRequestScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(maxID)), jsoncdc.MustEncode(cadence.UInt32(firstDelegatorID))})
		assertEqual(t, CadenceUFix64("0.0"), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetDelegatorUnstakingRequestScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(maxID)), jsoncdc.MustEncode(cadence.UInt32(secondDelegatorID))})
		assertEqual(t, CadenceUFix64("0.0"), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetTotalCommitmentBalanceScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(maxID))})
		assertEqual(t, CadenceUFix64("1702000.0"), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetTotalCommitmentBalanceWithoutDelegatorsScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(maxID))})
		assertEqual(t, CadenceUFix64("1400000.0"), result)

		tx := createTxWithTemplateAndAuthorizer(b, templates.GenerateUnstakeTokensScript(env), maxAddress)

		err := tx.AddArgument(CadenceUFix64("160000.0"))
		require.NoError(t, err)

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, maxAddress},
			[]crypto.Signer{b.ServiceKey().Signer(), maxSigner},
			true,
		)
	})

	// End the staking auction, saying that Max is not on the approved node list
	t.Run("Should refund delegators when their node is not included in the auction", func(t *testing.T) {

		tx := createTxWithTemplateAndAuthorizer(b, templates.GenerateEndStakingScript(env), idTableAddress)

		err := tx.AddArgument(cadence.NewArray([]cadence.Value{cadence.NewString(adminID), cadence.NewString(joshID)}))
		require.NoError(t, err)

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, idTableAddress},
			[]crypto.Signer{b.ServiceKey().Signer(), IDTableSigner},
			false,
		)

		result := executeScriptAndCheck(t, b, templates.GenerateGetDelegatorUnstakingRequestScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(maxID)), jsoncdc.MustEncode(cadence.UInt32(firstDelegatorID))})
		assertEqual(t, CadenceUFix64("100000.0"), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetDelegatorUnstakedScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(maxID)), jsoncdc.MustEncode(cadence.UInt32(firstDelegatorID))})
		assertEqual(t, CadenceUFix64("2000.0"), result)

	})

	// End staking auction and move tokens in the same transaction
	t.Run("Should end staking auction and move tokens in the same transaction", func(t *testing.T) {

		tx := createTxWithTemplateAndAuthorizer(b, templates.GenerateEndEpochScript(env), idTableAddress)

		err := tx.AddArgument(cadence.NewArray([]cadence.Value{cadence.NewString(adminID), cadence.NewString(joshID)}))
		require.NoError(t, err)

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, idTableAddress},
			[]crypto.Signer{b.ServiceKey().Signer(), IDTableSigner},
			false,
		)

	})

	t.Run("Should be able request unstake all which also requests to unstake all the delegator's tokens", func(t *testing.T) {

		result := executeScriptAndCheck(t, b, templates.GenerateGetCommittedBalanceScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(joshID))})
		assertEqual(t, CadenceUFix64("0.0"), result)

		tx := createTxWithTemplateAndAuthorizer(b, templates.GenerateUnstakeAllScript(env), joshAddress)

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, joshAddress},
			[]crypto.Signer{b.ServiceKey().Signer(), joshSigner},
			false,
		)

		result = executeScriptAndCheck(t, b, templates.GenerateGetUnstakingRequestScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(joshID))})
		assertEqual(t, CadenceUFix64("580000.0"), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetDelegatorRequestScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(joshID)), jsoncdc.MustEncode(cadence.UInt32(firstDelegatorID))})
		assertEqual(t, CadenceUFix64("0.0"), result)

	})

	t.Run("Should be able cancel unstake request for node operator", func(t *testing.T) {

		tx := createTxWithTemplateAndAuthorizer(b, templates.GenerateStakeUnstakedTokensScript(env), joshAddress)
		err := tx.AddArgument(CadenceUFix64("180000.0"))
		require.NoError(t, err)

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, joshAddress},
			[]crypto.Signer{b.ServiceKey().Signer(), joshSigner},
			false,
		)

		result := executeScriptAndCheck(t, b, templates.GenerateGetUnstakingRequestScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(joshID))})
		assertEqual(t, CadenceUFix64("400000.0"), result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetDelegatorRequestScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(joshID)), jsoncdc.MustEncode(cadence.UInt32(firstDelegatorID))})
		assertEqual(t, CadenceUFix64("0.0"), result)

	})

	t.Run("Should be able request unstake all again", func(t *testing.T) {

		tx := createTxWithTemplateAndAuthorizer(b, templates.GenerateUnstakeAllScript(env), joshAddress)

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, joshAddress},
			[]crypto.Signer{b.ServiceKey().Signer(), joshSigner},
			false,
		)
	})

	t.Run("Should end staking auction and move tokens in the same transaction, unstaking unstakeAll delegators' tokens", func(t *testing.T) {

		tx := createTxWithTemplateAndAuthorizer(b, templates.GenerateEndEpochScript(env), idTableAddress)

		err := tx.AddArgument(cadence.NewArray([]cadence.Value{cadence.NewString(adminID), cadence.NewString(joshID)}))
		require.NoError(t, err)

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, idTableAddress},
			[]crypto.Signer{b.ServiceKey().Signer(), IDTableSigner},
			false,
		)

		result := executeScriptAndCheck(t, b, templates.GenerateGetDelegatorUnstakingScript(env), [][]byte{jsoncdc.MustEncode(cadence.String(joshID)), jsoncdc.MustEncode(cadence.UInt32(firstDelegatorID))})
		assertEqual(t, CadenceUFix64("40000.0"), result)
	})

	t.Run("Should end epoch and change payout in the same transaction", func(t *testing.T) {

		tx := createTxWithTemplateAndAuthorizer(b, templates.GenerateEndEpochChangePayoutScript(env), idTableAddress)

		err := tx.AddArgument(cadence.NewArray([]cadence.Value{cadence.NewString(adminID), cadence.NewString(joshID)}))
		require.NoError(t, err)

		err = tx.AddArgument(CadenceUFix64("4000000.0"))
		require.NoError(t, err)

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, idTableAddress},
			[]crypto.Signer{b.ServiceKey().Signer(), IDTableSigner},
			false,
		)

		_ = executeScriptAndCheck(t, b, templates.GenerateGetNodeInfoFromAddressScript(env), [][]byte{jsoncdc.MustEncode(cadence.Address(joshAddress))})

		_ = executeScriptAndCheck(t, b, templates.GenerateGetDelegatorInfoFromAddressScript(env), [][]byte{jsoncdc.MustEncode(cadence.Address(joshDelegatorOneAddress))})

	})

	t.Run("Should be able to change the staking minimums", func(t *testing.T) {

		tx := createTxWithTemplateAndAuthorizer(b, templates.GenerateChangeMinimumsScript(env), idTableAddress)

		colMin := CadenceUFix64("250000.0")
		conMin := CadenceUFix64("250000.0")
		exMin := CadenceUFix64("1250000.0")
		verMin := CadenceUFix64("135000.0")
		accMin := CadenceUFix64("0.0")

		err := tx.AddArgument(cadence.NewArray([]cadence.Value{colMin, conMin, exMin, verMin, accMin}))
		require.NoError(t, err)

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, idTableAddress},
			[]crypto.Signer{b.ServiceKey().Signer(), IDTableSigner},
			false,
		)

		result := executeScriptAndCheck(t, b, templates.GenerateGetStakeRequirementsScript(env), [][]byte{jsoncdc.MustEncode(cadence.UInt8(1))})
		assertEqual(t, colMin, result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetStakeRequirementsScript(env), [][]byte{jsoncdc.MustEncode(cadence.UInt8(2))})
		assertEqual(t, conMin, result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetStakeRequirementsScript(env), [][]byte{jsoncdc.MustEncode(cadence.UInt8(3))})
		assertEqual(t, exMin, result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetStakeRequirementsScript(env), [][]byte{jsoncdc.MustEncode(cadence.UInt8(4))})
		assertEqual(t, verMin, result)

		result = executeScriptAndCheck(t, b, templates.GenerateGetStakeRequirementsScript(env), [][]byte{jsoncdc.MustEncode(cadence.UInt8(5))})
		assertEqual(t, CadenceUFix64("0.0"), result)

	})

	t.Run("Should be able to create public Capability for node", func(t *testing.T) {

		tx := createTxWithTemplateAndAuthorizer(b, templates.GenerateAddPublicNodeCapabilityScript(env), joshAddress)

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, joshAddress},
			[]crypto.Signer{b.ServiceKey().Signer(), joshSigner},
			false,
		)
	})

	t.Run("Should be able to create public Capability for delegator", func(t *testing.T) {

		tx := createTxWithTemplateAndAuthorizer(b, templates.GenerateAddPublicDelegatorCapabilityScript(env), joshDelegatorOneAddress)

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, joshDelegatorOneAddress},
			[]crypto.Signer{b.ServiceKey().Signer(), joshDelegatorOneSigner},
			false,
		)
	})
}
