## Overview 
BrunoCoin is a cryptocurrency built from scratch. It contains almost all aspects of a cryptocurrency: making transactions, broadcasting transactions/blocks to the network, mining transactions to a block, validation, maintaining the blockchain, etc. It is merely intended as a sample implementation of a cryptocurrency for educational purposes. It is only meant to run on one machine in a multithreaded environment in order to simulate a real network. For reference, BrunoCoin is based on the implementation of Bitcoin, so there will be many crossovers in terms of implementation-specific details.

## BrunoCoin Layout

The design of BrunoCoin was attempted to be designed in the most modular way possible. There are four main parts of the project: node, blockchain, miner, wallet. The node represents a single computer within the network. It acts as an interface between the person using the computer and all functionality. For example, if you wanted to make a transaction, the wallet will ultimately do this; however, you would tell the node that you wanted to make the transaction. The node, acting as an interface, will then tell the wallet to make the transaction. The nodes (represented as a struct) stores references to the other pieces of functionality (also represented as structs), namely, miner, blockchain, and wallet.

## BrunoCoin Node
The node struct represents a node on the BrunoCoin network. It contains networking functionality to send/receive blocks, transactions, and other types of requests. It also has a reference to a local copy of the blockchain, a miner who actively mines transactions to blocks, and a wallet to make transactions. A user of BrunoCoin should only interact through the node struct.

## BrunoCoin Blockchain
	
The blockchain struct stores all valid blocks mined from the node it is attached to or the network. This means that it will be storing multiple forked chains. It keeps track of all these chains using a tree data structure. This tree data structure is stored in a map where each value in the map has a backpointer to a previous node (this allows backwards traversal from leaf nodes up to genesis block). Each element in this tree is actually a blockchain node, not a block. A blockchain node just stores extra relevant information (such as the current depth). It also stores a map of UTXO. The UTXO is structured in a dynamic programming fashion. Every node stores a map of UTXO up until that point on the respective chain. Every new node added to a specific chain will have the same UTXO as the node before it except for any new UTXO and UTXO spent within the new block.

## BrunoCoin Wallet
	
The wallet is in charge of making transactions. In order to do this, it has to interface with the Node’s copy of the blockchain in order to request for UTXO. It also keeps track of transactions that it has made, but that have not been added “deep enough” in the blockchain. After a certain amount of time, if the wallet doesn’t see a previous transaction it made being added “deep enough” in the blockchain, then it will resend the transaction out.

## BrunoCoin Miner

BrunoCoin’s miner is responsible for mining transactions into a block (when appropriate to do so), as well as storing all related information to make this process happen. The miner keeps track of one pool of transactions (AdvancedCoin keeps track of two). This pool of transactions represent a sort of waiting list that the miner will pull from to start mining. When a person wants to start mining, they will tell the node.

The node will then call the mining functionality in a separate go routine. This separate go routine waits on a channel to start mining. This channel will become active when another go routine sends something using this channel. For example, another go routine might realize that we have enough transactions to mine, so it sends a message to the channel. On the receiving end, the channel lets the other go routine knowing that it should start mining now.
