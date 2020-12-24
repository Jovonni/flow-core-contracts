import FlowToken from 0xFLOWTOKENADDRESS
import FungibleTokenMemoReceiver from 0xFTTOKENMEMORECEIVER
import FlowTokenMemoReceiver from 0xFLOWTOKENMEMORECEIVER

transaction(amount: UFix64, to: Address, memo: String) {

    // The Vault resource that holds the tokens that are being transferred
    let sentVault: @FungibleToken.Vault

    prepare(signer: AuthAccount) {

        // Get a reference to the signer's stored vault
        let vaultRef = signer.borrow<&FlowToken.Vault>(from: /storage/flowTokenVault)
			?? panic("Could not borrow reference to the owner's Vault!")

        // Withdraw tokens from the signer's stored vault
        self.sentVault <- vaultRef.withdraw(amount: amount)
    }

    execute {

        // Get the recipient's public account object
        let recipient = getAccount(to)

        // Get a reference to the recipient's Receiver
        let receiverRef = recipient.getCapability(
            FlowTokenMemoReceiver.receiverPublicPath
        )!.borrow<&FlowTokenMemoReceiver.Wrapper{FungibleTokenMemoReceiver.Receiver}>()
			?? panic("Could not borrow receiver reference to the recipient's Vault")

        // Deposit the withdrawn tokens in the recipient's receiver, sending the memo
        receiverRef.depositWithMemo(from: <-self.sentVault, memo: memo)
    }
}
 