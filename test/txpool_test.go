package test

import (
	"BrunoCoin/pkg/miner"
	"BrunoCoin/pkg/utils"
	"testing"
)

func TestAddTxWthZeroFee(t *testing.T) {
	utils.SetDebug(true)
	genNd := NewGenNd()
	genNd.Start()

	txpool := miner.NewTxPool(miner.DefaultConfig(-1))

	tx := MakeSingleTx(genNd, genNd.Id.GetPublicKeyBytes(), 10)
	txpool.Add(tx)

	if txpool.Ct.Load() != 2 {
		t.Errorf("expected cap: 2, actual cap: %d", txpool.Ct.Load())
	}

	qooo := (*txpool.TxQ)[0]
	if qooo.P != 1 {
		t.Errorf("exp pri: 1, actual pri: %d", qooo.P)
	}
}

func TestAddMassiveFee(t *testing.T) {
	utils.SetDebug(true)
	genNd := NewGenNd()
	genNd.Start()

	txpool := miner.NewTxPool(miner.DefaultConfig(-1))

	tx := MakeSingleTx(genNd, genNd.Id.GetPublicKeyBytes(), 9999999)
	txpool.Add(tx)

	if txpool.Ct.Load() != 2 {
		t.Errorf("expected cap: 2, actual cap: %d", txpool.Ct.Load())
	}

	qooo := (*txpool.TxQ)[0]
	if qooo.P != miner.CalcPri(tx) {
		t.Errorf("exp pri: 1, actual pri: %d", qooo.P)
	}

}
