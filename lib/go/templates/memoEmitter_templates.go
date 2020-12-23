package templates

import (
	"github.com/onflow/flow-core-contracts/lib/go/templates/internal/assets"
)

const (
	// Admin templates
	memoEmitterDeployFilename = "memoEmitter/admin/deploy.cdc"

	// Simple templates
	emitMemoFilename = "memoEmitter/emit_memo.cdc"

	// Flow token templates
	transferFlowWithMemoFilename = "memoEmitter/transfer_flow_with_memo.cdc"
)

// GenerateMemoEmitterDeploy generates a script that simply
// emits a string as a Memo event.
func GenerateMemoEmitterDeploy(env Environment) []byte {
	code := assets.MustAssetString(memoEmitterDeployFilename)

	return []byte(replaceAddresses(code, env))
}

// GenerateEmitMemoScript generates a script that simply
// emits a string as a Memo event.
func GenerateEmitMemoScript(env Environment) []byte {
	code := assets.MustAssetString(emitMemoFilename)

	return []byte(replaceAddresses(code, env))
}

// GenerateTransferFlowWithMemoScript generates a script
// that transfedrs some FLOW to an account and emits a
// string as a Memo event.
func GenerateTransferFlowWithMemoScript(env Environment) []byte {
	code := assets.MustAssetString(transferFlowWithMemoFilename)

	return []byte(replaceAddresses(code, env))
}
