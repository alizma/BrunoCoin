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

	if mnr.TxP.Ct.Load() != 0 {
		t.Errorf("expected failed, count is: %v, actual count is: %v", 0, mnr.TxP.Ct.Load())
	}
}

func TestHndlActiveMnr(t *testing.T) {
	utils.SetDebug(true)
	genNd := NewGenNd()
	genNd.Start()

	mnr := miner.New(genNd.Mnr.Conf, genNd.Id)
	mnr.Active.Store(false)

	tx := MakeSingleTx(genNd, genNd.Id.GetPublicKeyBytes(), 10)
	mnr.HndlTx(tx)

	if mnr.TxP.Ct.Load() != 1 {
		t.Errorf("expected count is: %v, actual count is: %v", 1, mnr.TxP.Ct.Load())
	}
}

func TestNilTxHndl(t *testing.T) {
	utils.SetDebug(true)
	genNd := NewGenNd()
	genNd.Start()

	mnr := miner.New(genNd.Mnr.Conf, genNd.Id)
	mnr.Active.Store(false)

	mnr.HndlTx(nil)

	if mnr.TxP.Ct.Load() != 1 {
		t.Errorf("expected count is: %v, actual count is: %v", 0, mnr.TxP.Ct.Load())
	}
}
