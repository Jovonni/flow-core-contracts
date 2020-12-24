import FungibleToken from 0xFUNGIBLETOKEN

/*
    FungibleTokenMemoReceiver
    Copyright 2020 Dapper Labs
    License: Unlicense

    An interface for resources that wrap FungibleToken.Vaults
    and provide an alternative to FungibleToken.Vault.deposit()
    that emit events containing informational messages.
 */

pub contract FungibleTokenMemoReceiver {
    // DepositReceived
    // Publish the memo, along with the ID of the receiver resource
    // that published it and the amount.
    //
    pub event DepositReceived(receiverId: UInt64, amount: UFix64, memo: String)

    // Receiver
    // A public interface to a Vault wrapper that allows deposits with memos.
    //
    pub resource interface Receiver {
        // depositWithMemo
        // Allows anyone to deposit a Vault, accompanied by a memo string.
        //
        pub fun depositWithMemo(from: @FungibleToken.Vault, memo: String)
    }

    // Balance
    // An interface for a Vault wrapper to provide the balance of the wrapped Vault.
    //
    pub resource interface Balance {
        // balance
        // Allows the owner of the object to get the balance of the
        // wrapped Vault.
        //
        pub fun balance(): UFix64
    }

    // Wrapper
    // An interface for a Vault wrapper that allows deposits with memos.
    //
    pub resource interface Wrapper {
        // withdraw
        // Allows the owner of the object to withdraw from the
        // wrapped Vault.
        //
        pub fun withdraw(amount: UFix64): @FungibleToken.Vault
    }
}