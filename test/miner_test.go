package test

import (
	"BrunoCoin/pkg/miner"
	"BrunoCoin/pkg/utils"
	"testing"
)

func TestHndlMnrTxInactive(t *testing.T) {
	utils.SetDebug(true)
	genNd := NewGenNd()
	genNd.Start()

	mnr := miner.New(genNd.Mnr.Conf, genNd.Id)
	mnr.Active.Store(false)

	tx := MakeSingleTx(genNd, genNd.Id.GetPublicKeyBytes(), 10)

	mnr.HndlTx(tx)

	if mnr.TxP.Ct.Load() != 1 {
		t.Errorf("failed, count is: %v, actual count is: %v", 0, m.TxP.Ct.Load())
	}

}
