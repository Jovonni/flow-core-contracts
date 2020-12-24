import FlowToken from 0xFLOWTOKENADDRESS
import FungibleTokenMemoReceiver from 0xFTTOKENMEMORECEIVER
import FlowTokenMemoReceiver from 0xFLOWTOKENMEMORECEIVER

pub transaction() {
    setup(acct: AuthAccount) {
        // Check for Flow Vault
        assert(
            acct.borrow<&FlowToken.Vault>(from: /storage/flowTokenVault) != nil,
            message: "Account has no FLOW token vault resource"
        )

        // Check for existing wrapper
        assert(
            acct.borrow<&FlowTokenMemoReceiver.Wrapper>(from: FlowTokenMemoReceiver.wrapperStoragePath) == nil,
            message: "Account already has no FLOW token vault memo wrapper resource"
        )

        // Wrap Flow Vault
        let flowVault = acct.load<FlowToken.Vault>(from: /storage/flowTokenVault)
            ?? panic("Cannot load FLOW vault resource")
        
        let wrapper <- FlowTokenMemoReceiver.createNewWrapper(vault: <-flowVault)
 
        // Install wrapper
        acct.save(<-wrapper, to: FlowTokenMemoReceiver.wrapperStoragePath)
 
        // Install wrapper public capabilities
        acct.link<&FlowTokenMemoReceiver.Wrapper{FungibleTokenMemoReceiver.Balance}>(FlowTokenMemoReceiver.balancePublicPath, target: FlowTokenMemoReceiver.wrapperStoragePath)
        acct.link<&FlowTokenMemoReceiver.Wrapper{FungibleTokenMemoReceiver.Receiver}>(FlowTokenMemoReceiver.receiverPublicPath, target: FlowTokenMemoReceiver.wrapperStoragePath)

        // Remove Flow Vault public capabilities, if present
        acct.unlink(/public/flowTokenReceiver)
        acct.unlink(/public/flowTokenBalance)
    }
}