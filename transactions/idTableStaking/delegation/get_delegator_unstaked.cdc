import FlowIDTableStaking from 0xIDENTITYTABLEADDRESS

// This script returns the balance of unstaked tokens of a delegator

pub fun main(nodeID: String, delegatorID: UInt32): UFix64 {
    return FlowIDTableStaking.getDelegatorUnStakedBalance(nodeID, delegatorID: delegatorID)!
}