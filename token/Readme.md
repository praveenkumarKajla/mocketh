# How to generate erc20.go

```
solc --abi erc20.sol |  awk '/JSON ABI/{x=1;next}x' > ERC20.abi  

abigen --abi=ERC20.abi --pkg=token --out=erc20.go

```