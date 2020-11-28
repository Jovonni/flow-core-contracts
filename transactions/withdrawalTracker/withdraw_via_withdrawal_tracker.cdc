import LockedTokens from 0xLOCKEDTOKENADDRESS
import StakingProxy from 0xSTAKINGPROXYADDRESS
import WithdrawalTracker from 0xWITHDRAWALTRACKER

transaction(amount: UFix64) {

    let holderRef: &LockedTokens.TokenHolder

    prepare(account: AuthAccount) {
        let withdrawalTotalTrackerPath = /storage/withdrawalTotalTracker

        self.holderRef = account.borrow<&LockedTokens.TokenHolder>(from: LockedTokens.TokenHolderStoragePath)
            ?? panic("Could not borrow reference to TokenHolder")

        let withdrawalTracker = account.load<WithdrawalTracker.WithdrawalTotalTracker>(from: withdrawalLimitPath)
            ?? panic("Could not load withdrawal tracker")

        // We have to this here because we need access to the AuthAccount to save the result if it doesn't assert.
        withdrawalTracker.updateRunningTotal(withdrawalAmount: amount)
        account.save<WithdrawalTracker.WithdrawalTotalTracker>(tracker, to: withdrawalTotalTrackerPath)
    }

    execute {
        let stakerProxy = self.holderRef.borrowStaker()
        stakerProxy.withdrawRewardedTokens(amount: amount)
    }
}
