# go-mev-geth

go mev-geth example. signTx,send transaction to relay.

Updated mev v0.2

## sample smart contract

`/eth/contract/mevtransfer/mevtransfer.sol`

```solidity
pragma solidity 0.6.6;
interface IChiToken {
    function freeFromUpTo(address from, uint256 value) external;
    function balanceOf(address owner) external view returns (uint);
    function allowance(address owner, address spender) external view returns (uint);
}

contract MevTransfer{
    IChiToken constant public chi = IChiToken(0x0000000000004946c0e9F43F4Dee607b0eF1fA1c);

    modifier discountCHI {
        uint gasStart = gasleft();
        _;
        uint gasSpent = 21000 + gasStart - gasleft() + 16 * msg.data.length;
        uint chiToUse = (gasSpent + 14154) / 41947;
        if (chiToUse > 0 && chi.balanceOf(msg.sender) >= chiToUse && chi.allowance(msg.sender,address(this)) >= chiToUse){
            chi.freeFromUpTo(msg.sender, (gasSpent + 14154) / 41947);
        }
    }

    function transfer(uint bribe)  external payable discountCHI {
        block.coinbase.transfer(bribe);
        msg.sender.transfer(address(this).balance);
    }
}
```

a sample contract to 1. give bribe 2. use chi token 3. do a sample transfer.

use abigen to get go binding code.

## sign & send tx

`/bundle/bundle.go`

## build tx

`/bundle/bundle_test.go`

## how to use

change [sender](https://github.com/huahuayu/github.com/phanletrunghieu/go-mev-geth/blob/785eb0a2a4fff313653fa67e9dac29d467bb75b9/bundle/bundle_test.go#L16) & [contract](https://github.com/huahuayu/github.com/phanletrunghieu/go-mev-geth/blob/785eb0a2a4fff313653fa67e9dac29d467bb75b9/bundle/bundle_test.go#L17) var in `TestBundle_Simulate` and run test.

the `sender` is the EOA to send the tx.

the `contract` is the `mevTransfer.sol` address after deploy.
