package test

import (
	"BrunoCoin/pkg"
	"BrunoCoin/pkg/utils"
	"testing"
)

func HndlTxReqAmtZero(t *testing.T) {
	utils.SetDebug(true)

	genNd := NewGenNd()
	node2 := pkg.New(pkg.DefaultConfig(GetFreePort()))
	genNd.Conf.MnrConf.InitPOWD = utils.CalcPOWD(1)
	genNd.Start()
	genNd.StartMiner()

	genNd.SendTx(0, 0, node2.Id.GetPublicKeyBytes())
}
