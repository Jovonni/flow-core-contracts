import FlowToken from 0xFLOWTOKENADDRESS
import FungibleTokenMemoReceiver from 0xFTTOKENMEMORECEIVER
import FlowTokenMemoReceiver from 0xFLOWTOKENMEMORECEIVER

/*
    A transaction to withdraw FLOW tokens from a FlowMemoReceiver.Wrapper
    owned by the calling account.
 */

pub transaction(amount, address): UInt64 {
    let receiverRef = getAccount(address)
        .getCapability<&FlowTokenMemoReceiver.Wrapper{FungibleTokenMemoReceiver.Receiver}>(
            from: FlowTokenMemoReceiver.receiverPublicPath
        )!.borrow()
        ?? panic("Couldn't borrow Receiver interface")
    
    return receiverRef.uuid
}
