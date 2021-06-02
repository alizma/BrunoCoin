package pkg

import (
	"BrunoCoin/pkg/block"
	"BrunoCoin/pkg/block/tx"
	"BrunoCoin/pkg/utils"
)

/*
 *  Brown University, CS1951L, Summer 2021
 *  Designed by: Colby Anderson, John Roy,
 *	Parker Ljung
 *
 */

// ChkBlk (CheckBlock) validates a block based on multiple
// conditions.
// To be valid:
// The block must be syntactically (ChkBlkSyn), semantically
// (ChkBlkSem), and configurally (ChkBlkConf) valid.
// Each transaction on the block must be syntactically (ChkTxSyn),
// semantically (ChkTxSem), and configurally (ChkTxConf) valid.
// Each transaction on the block must reference UTXO on the same
// chain (main or forked chain) and not be a double spend on that
// chain.
// Inputs:
// b *block.Block the block to be checked for validity
// Returns:
// bool True if the block is valid. false
// otherwise
// TODO:
// to be valid

// Each transaction on the block must reference UTXO on the same
// chain (main or forked chain) and not be a double spend on that
// chain.
// The block's size must be less than or equal to the largest
// allowed block size.
// The block hash must be less than the difficulty target.
// The block's first transaction must be of type Coinbase.

// Some helpful functions/methods/fields:
// note: let t be a transaction object
// note: let b be a block object
// t.IsCoinbase()
// b.SatisfiesPOW(...)
// n.Conf.MxBlkSz
// b.Sz()
// n.Chain.ChkChainsUTXO(...)
func (n *Node) ChkBlk(b *block.Block) bool {
	if len(b.Transactions) <= 0 {
		utils.Debug.Printf("block has no transactions")
		return false
	}

	if !b.Transactions[0].IsCoinbase() {
		utils.Debug.Printf("is not coinbase")
	} else if !b.SatisfiesPOW(b.Hdr.DiffTarg) {
		utils.Debug.Printf("does not satisfy difficulty target")
	} else if !(b.Sz() <= n.Conf.MxBlkSz) {
		utils.Debug.Printf("does not satisfy block size")
	} else if !n.Chain.ChkChainsUTXO(b.Transactions[1:], b.Hdr.PrvBlkHsh) {
		utils.Debug.Printf("fails chain check")
	}

	return b.Transactions[0].IsCoinbase() &&
		b.SatisfiesPOW(b.Hdr.DiffTarg) &&
		b.Sz() <= n.Conf.MxBlkSz &&
		n.Chain.ChkChainsUTXO(b.Transactions[1:], b.Hdr.PrvBlkHsh)
}

// ChkTx (CheckTransaction) validates a transaction.
// Inputs:
// t *tx.Transaction the transaction to be checked for validity
// Returns:
// bool True if the transaction is syntactically valid. false
// otherwise
// TODO:
// to be valid:

// The transaction's inputs and outputs must not be empty.
// The transaction's output amounts must be larger than 0.
// The sum of the transaction's inputs must be larger
// than the sum of the transaction's outputs.
// The transaction must not double spend any UTXO.
// The unlocking script on each of the transaction's
// inputs must successfully unlock each of the corresponding
// UTXO.
// The transaction must not be larger than the
// maximum allowed block size.

// Some helpful functions/methods/fields:
// note: let t be a transaction object
// note: let b be a block object
// note: let u be a transaction output object
// n.Conf.MxBlkSz
// t.Sz()
// u.IsUnlckd(...)
// n.Chain.GetUTXO(...)
// n.Chain.IsInvalidInput(...)
// t.SumInputs()
// t.SumOutputs()
func (n *Node) ChkTx(t *tx.Transaction) bool {
	if t.Sz() > n.Conf.MxBlkSz {
		utils.Debug.Printf("failed blocksize")
		return false
	}

	if t.Inputs == nil || t.Outputs == nil || len(t.Outputs) <= 0 {
		utils.Debug.Printf("failed input/output size: len inputs - %v, len outputs: %v", len(t.Inputs), len(t.Outputs))
		return false
	}

	if t.SumOutputs() < uint32(0) {
		utils.Debug.Printf("failed sumoutputs")
		return false
	}

	if t.SumInputs() < t.SumOutputs() {
		utils.Debug.Printf("failed suminputs %v < sumoutputs %v", t.SumInputs(), t.SumOutputs())
		return false
	}

	for _, currInput := range t.Inputs {
		if n.Chain.IsInvalidInput(currInput) {
			utils.Debug.Printf("is invalid input")
			return false
		}

		if !n.Chain.GetUTXO(currInput).IsUnlckd(currInput.UnlockingScript) {
			utils.Debug.Printf("%v failed unlocking", currInput)
			return false
		}
	}

	return true
}
