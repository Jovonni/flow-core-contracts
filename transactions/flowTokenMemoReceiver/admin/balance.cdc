import FungibleTokenMemoReceiver from 0xFTTOKENMEMORECEIVER
import FlowTokenMemoReceiver from 0xFLOWTOKENMEMORECEIVER

/*
    A script to get the balance of a wrapped FlowToken.Vault
    from its enclosing FlowMemoReceiver.Wrapper .
 */

pub fun main(address: Address): UFix64 {
    let balanceRef = getAccount(address)
        .getCapability<&FlowTokenMemoReceiver.Wrapper{FungibleTokenMemoReceiver.Balance}>(
            from: FlowTokenMemoReceiver.balancePublicPath
        )!.borrow()
        ?? panic("Couldn't borrow Balance interface")
    
    return balanceRef.balance()
}