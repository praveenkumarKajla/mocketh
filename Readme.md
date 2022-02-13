# Indexer for Transfer events of a given list of tokens.


## Motivation

Tokens can be transferred by calling transfer or transferFrom functions and each time the tokens are transferred, Transfer event is emitted. This event can be used to track all token transfers by offchain systems.

This Service works two ways

1. Store new transfer events as soon as they are fired.
ERC20 token Smart contracts generate Transfer logs by firing events whenever transfer method is called.

2. Look for  transfer events in the past from the current block and store them in the db. We store full history of the transfer events for a specific ERC20 token over time as it is taken care by the periodic task which runs every x seconds.


