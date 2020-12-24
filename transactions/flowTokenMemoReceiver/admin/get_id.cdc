import FlowTokenMemoReceiver from 0xFLOWTOKENMEMORECEIVER

/*
    A script to get the unique ID of a FlowMemoReceiver.Wrapper,
    as used in the DepositReceived events that it emits.
 */

pub fun main(address: Address): UInt64 {
    let receiverRef = getAccount(address)
        .getCapability<&FlowTokenMemoReceiver.Wrapper{FungibleTokenMemoReceiver.Receiver}>(
            from: FlowTokenMemoReceiver.receiverPublicPath
        )!.borrow()
        ?? panic("Couldn't borrow Receiver interface")
    
    return receiverRef.uuid
}
