import MemoEmitter from 0xMEMOEMITTER

transaction(memo: String) {

    prepare(account: AuthAccount) {
        MemoEmitter.emit(memo: memo)
    }

}