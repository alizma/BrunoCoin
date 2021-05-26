package miner

import (
	"BrunoCoin/pkg/block"
	"BrunoCoin/pkg/block/tx"
	"BrunoCoin/pkg/proto"
	"BrunoCoin/pkg/utils"
	"context"
	"encoding/hex"
)

/*
 *  Brown University, CS1951L, Summer 2021
 *  Designed by: Colby Anderson, Parker Ljung
 */

// Mine waits to be told to mine a block
// or to kill it's thread. If it is asked
// to mine, it selects the transactions
// with the highest priority to add to the
// mining pool. The nonce is then attempted
// to be found unless the miner is stopped.
func (m *Miner) Mine() {
	ctx, cancel := context.WithCancel(context.Background())
	for {
		<-m.PoolUpdated
		cancel()
		if !m.Active.Load() {
			continue
		}
		ctx, cancel = context.WithCancel(context.Background())
		go func(ctx context.Context) {
			if !m.TxP.PriMet() {
				//utils.Debug.Printf("current priority: %v, PriLim: %v", m.TxP.CurPri.Load(), m.TxP.PriLim)
				return
			}
			m.Mining.Store(true)
			m.MiningPool = m.NewMiningPool()
			txs := append([]*tx.Transaction{m.GenCBTx(m.MiningPool)}, m.MiningPool...)
			b := block.New(m.PrvHsh, txs, m.DifTrg())
			//utils.Debug.Printf("new block %v", b.NameTag())
			result := m.CalcNonce(ctx, b)
			m.Mining.Store(false)
			if result {
				utils.Debug.Printf("%v mined %v %v", utils.FmtAddr(m.Addr), b.NameTag(), b.Summarize())
				m.SendBlk <- b
				m.HndlBlk(b)
			}
		}(ctx)
	}
	cancel()
}

// Returns boolean to indicate success
func (m *Miner) CalcNonce(ctx context.Context, b *block.Block) bool {
	for i := uint32(0); i < m.Conf.NncLim; i++ {
		select {
		case <-ctx.Done():
			return false
		default:
			b.Hdr.Nonce = i
			if b.SatisfiesPOW(m.DifTrg()) {
				return true
			}
		}
	}
	return false
}

// DifTrg (DifficultyTarget) calculates the
// current difficulty target.
// Returns:
// string the difficulty target as a hex
// string
func (m *Miner) DifTrg() string {
	return m.Conf.InitPOWD
}

// GenCBTx (GenerateCoinbaseTransaction) generates a coinbase
// transaction based off the transactions in the mining pool.
// It does this by adding the fee reward to the minting reward.
// Inputs:
// txs	[]*tx.Transaction the transactions (besides the
// coinbase tx) that the miner is mining to a block
// Returns:
// the coinbase transaction that pays the miner the reward
// for mining the block
// TODO
// 1. calculate how many fees you can collect
// 2. calculate the minting reward
// for the minting reward: it starts as c.InitSubsdy
// and gets cut in half every c.SubsdyHlvRt until
// c.MxHlvgs
// 3. Then a transaction is made with zero inputs
// and with an output paying the miner fees + mint

// some functions/fields/methods that might be helpful
// tx.Deserialize(...)
// proto.NewTx(...)
// m.Conf.Ver
// m.Id.GetPublicKeyBytes()
// proto.NewTxOutpt(...)
// hex.EncodeToString(...)
// m.Conf.DefLckTm
// m.ChnLen.Load()
// c.SubsdyHlvRt
// c.MxHlvgs
// c.InitSubsdy
// t.SumInputs()
// t.SumOutputs()
func (m *Miner) GenCBTx(txs []*tx.Transaction) *tx.Transaction {
	// ignore the fact that we need to check for the TimeLock
	// calculate the total fees from all the transactions
	var total_fees uint32 = 0
	// fees are the difference between output and input values on Txs

	for _, curr_transaction := range txs {
		total_fees += curr_transaction.SumInputs() - curr_transaction.SumOutputs()
	}

	mintingReward := uint32(0)

	if chnLen := m.ChnLen.Load(); chnLen < m.Conf.SubsdyHlvRt*m.Conf.MxHlvgs {
		mintingReward += m.Conf.InitSubsdy / (chnLen/m.Conf.SubsdyHlvRt + 1)
	}

	/*
		// compute minting reward
		numHlvngs := uint32(math.Floor(math.Pow(0.5, float64(m.ChnLen.Load()/m.Conf.SubsdyHlvRt))))
		if numHlvngs > m.Conf.MxHlvgs {
			numHlvngs = m.Conf.MxHlvgs
		}
	*/

	CBTOutput := []*proto.TransactionOutput{proto.NewTxOutpt(total_fees+mintingReward, hex.EncodeToString(m.Id.GetPublicKeyBytes()))}

	ptrCBT := proto.NewTx(m.Conf.Ver, nil, CBTOutput, m.Conf.DefLckTm)

	//utils.Debug.Printf("new transaction %v", ptrCBT.Version)

	return tx.Deserialize(ptrCBT)
}
