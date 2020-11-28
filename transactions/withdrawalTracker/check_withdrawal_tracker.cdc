
import WithdrawalTracker from 0xWITHDRAWALTRACKER

pub fun main(): [UFix64] {
    let withdrawalTotalCheckerPath = /public/withdrawalTotalChecker

    let account = getAccount(0xPLACEHOLDERTRACKERUSERADDRESS)

    // Link the tracker checker to the account's public area
    let checkerCapability = account.getCapability<&{WithdrawalTracker.WithdrawalTotalChecker}>(
        withdrawalTotalCheckerPath
    ) ?? panic("Cannot get withdrawal tracker checker capability from account")

    let checker = checkerCapability.borrow()!
    return [checker.getCurrentRunningTotal(), checker.getCurrentWithdrawalLimit()]
}
