package test

import (
	"BrunoCoin/pkg/block/tx"
	"BrunoCoin/pkg/block/tx/txi"
	"BrunoCoin/pkg/block/tx/txo"
	"BrunoCoin/pkg/utils"
	"encoding/hex"
	"testing"
)

func TestGenCBTXNoTxs(t *testing.T) {
	utils.SetDebug(true)

	genNd := NewGenNd()
	genNd.Conf.MnrConf.InitPOWD = utils.CalcPOWD(1)
	genNd.Start()
	genNd.StartMiner()

	resultingtx := genNd.Mnr.GenCBTx([]*tx.Transaction{})

	if resultingtx != nil {
		t.Errorf("expected error since no transactions provided")
	}
}

func TestGenCBTxNilTxPtr(t *testing.T) {
	utils.SetDebug(true)

	genNd := NewGenNd()
	genNd.Conf.MnrConf.InitPOWD = utils.CalcPOWD(1)
	genNd.Start()
	genNd.StartMiner()

	resultingtx := genNd.Mnr.GenCBTx(nil)

	if resultingtx != nil {
		t.Errorf("expected error since nil provided")
	}
}

func TestGenCBTxListofTxWithNil(t *testing.T) {
	utils.SetDebug(true)

	genNd := NewGenNd()
	genNd.Conf.MnrConf.InitPOWD = utils.CalcPOWD(1)
	genNd.Start()
	genNd.StartMiner()

	tx1 := MakeSingleTx(genNd, genNd.Id.GetPublicKeyBytes(), 2)

	txarr := genNd.Mnr.GenCBTx([]*tx.Transaction{nil, tx1})

	if txarr != nil {
		t.Errorf("expected error since nil provided")
	}
}

func TestGenCBTxLowFees(t *testing.T) {
	utils.SetDebug(true)

	genNd := NewGenNd()
	genNd.Conf.MnrConf.InitPOWD = utils.CalcPOWD(1)
	genNd.Start()
	genNd.StartMiner()

	tx1 := MakeSingleTx(genNd, genNd.Id.GetPublicKeyBytes(), 0)

	cbtx := genNd.Mnr.GenCBTx([]*tx.Transaction{tx1})

	if cbtx.SumOutputs() > genNd.Mnr.Conf.InitSubsdy {
		t.Errorf("EXPECTED OUTPUTS TO BE LESS THAN SUBSIDY")
	}
}

func TestGenCBTxNoFees(t *testing.T) {
	utils.SetDebug(true)

	genNd := NewGenNd()
	genNd.Conf.MnrConf.InitPOWD = utils.CalcPOWD(1)
	genNd.Start()
	genNd.StartMiner()

	inputUtxo1 := &txo.TransactionOutput{
		Amount:        0,
		LockingScript: hex.EncodeToString(genNd.Id.GetPublicKeyBytes()),
	}
	unlck1, _ := inputUtxo1.MkSig(genNd.Id)
	txii1 := &txi.TransactionInput{
		TransactionHash: inputUtxo1.Hash(),
		OutputIndex:     0,
		UnlockingScript: unlck1,
		Amount:          inputUtxo1.Amount,
	}
	inputUtxo2 := &txo.TransactionOutput{
		Amount:        0,
		LockingScript: hex.EncodeToString(genNd.Id.GetPublicKeyBytes()),
	}
	unlck2, _ := inputUtxo2.MkSig(genNd.Id)
	txii2 := &txi.TransactionInput{
		TransactionHash: inputUtxo2.Hash(),
		OutputIndex:     1,
		UnlockingScript: unlck2,
		Amount:          inputUtxo2.Amount,
	}

	txoo1 := &txo.TransactionOutput{
		Amount:        0,
		LockingScript: hex.EncodeToString(genNd.Id.GetPublicKeyBytes()),
	}
	txoo2 := &txo.TransactionOutput{
		Amount:        0,
		LockingScript: hex.EncodeToString(genNd.Id.GetPublicKeyBytes()),
	}

	transaction := &tx.Transaction{
		Version:  genNd.Wallet.Conf.TxVer,
		Inputs:   []*txi.TransactionInput{txii1, txii2},
		Outputs:  []*txo.TransactionOutput{txoo1, txoo2},
		LockTime: genNd.Wallet.Conf.DefLckTm,
	}

	cbtx := genNd.Mnr.GenCBTx([]*tx.Transaction{transaction})

	if cbtx.SumOutputs()-genNd.Mnr.Conf.InitSubsdy != 0 {
		t.Errorf("FEE WUZ ZERO, DIFF MUST BE ZERO")
	}
}
