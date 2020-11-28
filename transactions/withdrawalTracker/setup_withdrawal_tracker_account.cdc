import WithdrawalTracker from 0xWITHDRAWALTRACKER

transaction() {

    prepare(account: AuthAccount) {
        let withdrawalLimitAmount = UFix64(2000000.0)
        let startingRunningTotal = UFix64(0.0)
        let withdrawalTotalTrackerPath = /storage/withdrawalTotalTracker
        let withdrawalTotalCheckerPath = /public/withdrawalTotalChecker

        // Make and save the tracker
        let tracker <- WithdrawalTracker.createWithdrawalTotalTracker(initialLimit: withdrawalLimitAmount, initialRunningTotal: startingRunningTotal)
        account.save<@WithdrawalTracker.WithdrawalTotalTracker>(<-tracker, to: withdrawalTotalTrackerPath)

        // Link the tracker checker to the account's public area
        let checkerCapability = account.link<&{WithdrawalTracker.WithdrawalTotalChecker}>(
            withdrawalTotalCheckerPath,
            target: withdrawalTotalTrackerPath
        )
    }

}
