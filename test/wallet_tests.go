package test

import (
	"BrunoCoin/pkg"
	"BrunoCoin/pkg/utils"
	"testing"
	"time"
)

func HndlTxReqAmtZero(t *testing.T) {
	utils.SetDebug(true)

	genNd := NewGenNd()
	node2 := pkg.New(pkg.DefaultConfig(GetFreePort()))
	genNd.Conf.MnrConf.InitPOWD = utils.CalcPOWD(1)
	genNd.Start()
	genNd.StartMiner()

	genNd.SendTx(0, 0, node2.Id.GetPublicKeyBytes())

	ChkMnChnCons(t, []*pkg.Node{genNd, node2})
	ChkNdPrs(t, genNd, []*pkg.Node{node2})
	ChkNdPrs(t, node2, []*pkg.Node{genNd})

	ChkTxSeenLen(t, genNd, 0)
	ChkTxSeenLen(t, node2, 0)
}

func HndlTxReqNotEnough(t *testing.T) {
	utils.SetDebug(true)

	genNd := NewGenNd()
	node2 := pkg.New(pkg.DefaultConfig(GetFreePort()))
	genNd.Conf.MnrConf.InitPOWD = utils.CalcPOWD(1)
	genNd.Start()
	genNd.StartMiner()

	genNd.SendTx(200, 0, node2.Id.GetPublicKeyBytes())

	time.Sleep(5 * time.Second)

	node2.SendTx(999999, 0, genNd.Id.GetPublicKeyBytes())

	ChkMnChnCons(t, []*pkg.Node{genNd, node2})
	ChkNdPrs(t, genNd, []*pkg.Node{node2})
	ChkNdPrs(t, node2, []*pkg.Node{genNd})

	ChkTxSeenLen(t, genNd, 0)
	ChkTxSeenLen(t, node2, 0)
}
