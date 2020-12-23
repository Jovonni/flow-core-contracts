/*
    MemoEmitter.
    Copyright Dapper Labs 2020.
    License: Unlicense (see: github.com/onflow/flow-core-contracts/blob/master/LICENSE).

    This contract allows anyone to emit informational events from contracts and transactions.
    These can then be watched for and processed off-chain.

    There is a limit to how long memo strings can be, which is curently 255 characters.
    If a memo string longer than this is passed to the emit() function it will revert.
    Do not assume that this limit will remain unchanged, in particular do not store the
    string in a MySQL `TINYTEXT` field or equivalent if you do not control the length of
    the memo strings you are expecting to process of-chain. 
 */

pub contract MemoEmitter {
    // The maximum length of memo string that can be passed to the emit() function
    //
    pub var maximumMemoLength: Int

    // A memo data event.
    // These can be emitted by any party in a transaction, in any order,
    // so be careful with the strategy you use for handling them.
    //
    pub event Memo(memo: String)

    // The function that emits a memo event.
    // If `memo:` is longer than `maximumMemoLength` it will revert.
    // This is to prevent spam.
    //
    pub fun emit(memo: String) {
        pre{
            memo.length <= self.maximumMemoLength: "memo string too long, check Memo.maximumMemoLength"
        }
        emit Memo(memo: memo)
    }

    // An admin resource that allows the user holding it to control memo string length
    pub resource Admin {
        // Set the maximum memo string length.
        // Unless 
        //
        pub fun setMaximumMemoLength(_ newLength: Int) {
            MemoEmitter.maximumMemoLength = newLength
        }
    }

    // Initializer
    init () {
        // Start with a Pascal string's worth of data per memo.
        self.maximumMemoLength = 255

        // Allow the owner of this contract to change that.
        self.account.save(<- create Admin(), to: /storage/memoEmitterAdmin)
    }
}
