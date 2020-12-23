transaction(publicKey: [UInt8], contractName: String, code: [UInt8],) {

  prepare(signer: AuthAccount) {

    let acct = AuthAccount(payer: signer)
    
    acct.addPublicKey(publicKey)

    acct.contracts.add(name: contractName, code: code)
  }

}
