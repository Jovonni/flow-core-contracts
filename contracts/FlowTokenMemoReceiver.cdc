import FungibleToken from 0xFUNGIBLETOKEN
import FlowToken from 0xFLOWTOKEN
import FungibleTokenMemoReceiver from 0xFTMEMORECEIVER

/*
    FlowTokenMemoReceiver
    Copyright 2020 Dapper Labs
    License: Unlicense

    A contract providing a resourc that wraps a FlowToken.Vaults
    and provide an alternative to FlowToken.Vault.deposit()
    that emits event containing informational messages.
 */

pub contract FlowTokenMemoReceiver {
    // DepositReceived
    // Publish the memo, along with the ID of the receiver resource
    // that published it and the amount.
    //
    pub event DepositReceived(receiverId: UInt64, amount: UFix64, memo: String)

    // wrapperStoragePath
    // The path to the full receiver resource in storage.
    //
    pub let wrapperStoragePath: Path

    // receiverPublicPath
    // The path to the public deposit interface for the receiver resource.
    //
    pub let receiverPublicPath: Path

    // balancePublicPath
    // The path to the public balance interface for the receiver resource.
    //
    pub let balancePublicPath: Path

    // Wrapper
    // A resource that wraps a FlowToken.Vault and allows anyone to
    // deposit FLOW to it along with a memo, and allows the owner to
    // withdraw FLOW from it.
    //
    pub resource Wrapper:
        FungibleTokenMemoReceiver.Wrapper,
        FungibleTokenMemoReceiver.Receiver,
        FungibleTokenMemoReceiver.Balance
    {
        // flowVault
        // The (permanently) wrapped FLOW Vault.
        access(self) let flowVault: @FlowToken.Vault

        // depositWithMemo
        // Allows anyone to deposit FLOW tokens, accompanied by a memo string.
        //
        pub fun depositWithMemo(from: @FungibleToken.Vault, memo: String) {
            emit DepositReceived(receiverId: self.uuid, amount: from.balance, memo: memo)
            self.flowVault.deposit(from: <-from)
        }

        // withdraw
        // Allows the owner of the object to withdraw FLOW from the
        // wrapped Vault.
        //
        pub fun withdraw(amount: UFix64): @FungibleToken.Vault {
            return <-self.flowVault.withdraw(amount: amount)
        }

        // balance
        // Allows the owner of the object to get the balance of the
        // wrapped Vault.
        //
        pub fun balance(): UFix64 {
            return self.flowVault.balance
        }

        // destructor
        // Make sure to withdraw any remaining FLOW Vault balance
        // before calling this, otherwise it will fail.
        //
        destroy() {
            assert(self.flowVault.balance == 0.0, message: "Wrapped FLOW vault is not empty")
            destroy self.flowVault
        }

        // initializer
        // The vault is wrapped permanently, and cannot be removed.
        //
        init(vault: @FlowToken.Vault) {
            self.flowVault <- vault
        }
    }

    // createNewReceiver creates a new Wrapper resource with the provided Vault
    //
    pub fun createNewWrapper(vault: @FlowToken.Vault): @Wrapper {
        return <-create Wrapper(vault: <-vault)
    }

    // initializer
    //
    init () {
        // Set our named paths
        self.wrapperStoragePath = /storage/FlowTokenMemoReciever
        self.receiverPublicPath = /public/FlowTokenMemoRecieverDeposit
        self.balancePublicPath = /public/FlowTokenMemoRecieverBalance
    }
}
