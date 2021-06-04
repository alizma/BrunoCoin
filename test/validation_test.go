package test

import (
	"BrunoCoin/pkg"
	"BrunoCoin/pkg/block/tx"
	"BrunoCoin/pkg/utils"
	"testing"
)

func TestChkTxsDuplicates(t *testing.T) {
	utils.SetDebug(true)

	genNd := NewGenNd()
	node2 := pkg.New(pkg.DefaultConfig(GetFreePort()))
	genNd.Conf.MnrConf.InitPOWD = utils.CalcPOWD(1)
	genNd.Start()
	genNd.StartMiner()

	genNd.ConnectToPeer(node2.Addr)

	tx1 := MakeSingleTx(genNd, genNd.Id.GetPublicKeyBytes(), 2)
	tx2 := MakeSingleTx(genNd, genNd.Id.GetPublicKeyBytes(), 2)
	tx3 := MakeSingleTx(genNd, genNd.Id.GetPublicKeyBytes(), 2)

	txarr := []*tx.Transaction{tx1, tx2, tx3}

	genNd.Mnr.TxP.ChkTxs(txarr)

	if genNd.Mnr.TxP.Ct.Load() != 0 {
		t.Errorf("expected %d, actual; %d", 0, genNd.Mnr.TxP.Ct.Load())
	}
}

func TestNoDuplicatesChkTxs(t *testing.T) {
	utils.SetDebug(true)

	genNd := NewGenNd()
	node2 := pkg.New(pkg.DefaultConfig(GetFreePort()))
	genNd.Conf.MnrConf.InitPOWD = utils.CalcPOWD(1)
	genNd.Start()
	genNd.StartMiner()

	genNd.ConnectToPeer(node2.Addr)

	tx1 := MakeSingleTx(genNd, genNd.Id.GetPublicKeyBytes(), 2)
	tx2 := MakeSingleTx(genNd, genNd.Id.GetPublicKeyBytes(), 3)
	tx3 := MakeSingleTx(genNd, genNd.Id.GetPublicKeyBytes(), 4)

	genNd.Mnr.TxP.Add(tx1)
	genNd.Mnr.TxP.Add(tx2)
	genNd.Mnr.TxP.Add(tx3)

	tx4 := MakeSingleTx(genNd, genNd.Id.GetPublicKeyBytes(), 5)
	tx5 := MakeSingleTx(genNd, genNd.Id.GetPublicKeyBytes(), 6)
	tx6 := MakeSingleTx(genNd, genNd.Id.GetPublicKeyBytes(), 7)

	txarr := []*tx.Transaction{tx4, tx5, tx6}

	genNd.Mnr.TxP.ChkTxs(txarr)

	if genNd.Mnr.TxP.Ct.Load() != 3 {
		t.Errorf("expected %d, actual; %d", 3, genNd.Mnr.TxP.Ct.Load())
	}
}
