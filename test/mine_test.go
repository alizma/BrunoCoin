package test

import (
	"BrunoCoin/pkg/block/tx"
	"BrunoCoin/pkg/utils"
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
