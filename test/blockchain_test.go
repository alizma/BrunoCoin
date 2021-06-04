package test

import (
	"BrunoCoin/pkg/utils"
	"testing"
)

func TestAddNilBlck(t *testing.T) {
	utils.SetDebug(true)
	genNd := NewGenNd()
	genNd.Conf.MnrConf.InitPOWD = utils.CalcPOWD(1)
	genNd.Start()
	genNd.StartMiner()

	genNd.Chain.Add(nil)

	ChkMnChnLen(t, genNd, 1)
}

/*
func TestAddEmptyTxBlk(t *testing.T) {
	utils.SetDebug(true)
	genNd := NewGenNd()
	genNd.Conf.MnrConf.InitPOWD = utils.CalcPOWD(1)
	genNd.Start()
	genNd.StartMiner()

	newBlk := block.New(genNd.Chain.LastBlock.Block.Hash(), []*tx.Transaction{}, genNd.Chain.LastBlock.Block.Hdr.DiffTarg)

	genNd.Chain.Add(newBlk)

	ChkMnChnLen(t, genNd, 1)
}
*/
