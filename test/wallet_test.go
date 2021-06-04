package test

import (
	"BrunoCoin/pkg"
	"BrunoCoin/pkg/utils"
	"testing"
	"time"
)

func TestHndlTxReqAmtZero(t *testing.T) {
	utils.SetDebug(true)

	genNd := NewGenNd()
	node2 := pkg.New(pkg.DefaultConfig(GetFreePort()))
	genNd.Conf.MnrConf.InitPOWD = utils.CalcPOWD(1)
	genNd.Start()
	genNd.StartMiner()

	genNd.SendTx(0, 0, node2.Id.GetPublicKeyBytes())

	ChkTxSeenLen(t, genNd, 0)
	ChkTxSeenLen(t, node2, 0)
}

func TestHndlTxReqNotEnough(t *testing.T) {
	utils.SetDebug(true)

	genNd := NewGenNd()
	node2 := pkg.New(pkg.DefaultConfig(GetFreePort()))
	genNd.Conf.MnrConf.InitPOWD = utils.CalcPOWD(1)
	genNd.Start()
	genNd.StartMiner()

	genNd.ConnectToPeer(node2.Addr)

	genNd.SendTx(200, 0, node2.Id.GetPublicKeyBytes())

	time.Sleep(5 * time.Second)

	node2.SendTx(999999, 0, genNd.Id.GetPublicKeyBytes())

	ChkTxSeenLen(t, genNd, 0)
	ChkTxSeenLen(t, node2, 0)
}

func TestHndlTxReqChange(t *testing.T) {
	utils.SetDebug(true)
	genNd := NewGenNd()
	node2 := pkg.New(pkg.DefaultConfig(GetFreePort()))
	genNd.Start()
	node2.Start()
	genNd.ConnectToPeer(node2.Addr)
	if peer := genNd.PeerDb.Get(node2.Addr); peer == nil {
		t.Fatal("Seed node did not contain newNode as peer")
	}
	if peer := node2.PeerDb.Get(genNd.Addr); peer == nil {
		t.Fatal("New node did not contain seedNode as peer")
	}
	// Sleep to give time for both nodes to connect
	time.Sleep(1 * time.Second)

	ChkMnChnCons(t, []*pkg.Node{genNd, node2})
	ChkNdPrs(t, genNd, []*pkg.Node{node2})
	ChkNdPrs(t, node2, []*pkg.Node{genNd})

	genNd.SendTx(100, 100, node2.Id.GetPublicKeyBytes())
	node2.SendTx(100, 100, genNd.Id.GetPublicKeyBytes())

	time.Sleep(6 * time.Second)
	node2.StartMiner()

	ChkTxSeenLen(t, genNd, 1)
	ChkTxSeenLen(t, node2, 1)
	time.Sleep(6 * time.Second)

	ChkMnChnCons(t, []*pkg.Node{genNd, node2})

	AsrtBal(t, genNd, 99800)
	AsrtBal(t, node2, 210)
}

func TestHndlTxNoChange(t *testing.T) {
	utils.SetDebug(true)
	genNd := NewGenNd()
	node2 := pkg.New(pkg.DefaultConfig(GetFreePort()))
	genNd.Start()
	node2.Start()
	genNd.ConnectToPeer(node2.Addr)
	if peer := genNd.PeerDb.Get(node2.Addr); peer == nil {
		t.Fatal("Seed node did not contain newNode as peer")
	}
	if peer := node2.PeerDb.Get(genNd.Addr); peer == nil {
		t.Fatal("New node did not contain seedNode as peer")
	}
	// Sleep to give time for both nodes to connect
	time.Sleep(1 * time.Second)

	ChkMnChnCons(t, []*pkg.Node{genNd, node2})
	ChkNdPrs(t, genNd, []*pkg.Node{node2})
	ChkNdPrs(t, node2, []*pkg.Node{genNd})

	genNd.SendTx(100, 0, node2.Id.GetPublicKeyBytes())
	genNd.SendTx(100, 0, genNd.Id.GetPublicKeyBytes())

	time.Sleep(6 * time.Second)
	node2.StartMiner()

	ChkTxSeenLen(t, genNd, 1)
	ChkTxSeenLen(t, node2, 1)
	time.Sleep(6 * time.Second)

	ChkMnChnCons(t, []*pkg.Node{genNd, node2})

	AsrtBal(t, genNd, 100000)
	AsrtBal(t, node2, 0)
}
func TestHndlBlkNilBlck(t *testing.T) {
	utils.SetDebug(true)
	genNd := NewGenNd()
	genNd.Start()

	genNd.Chain.Add(nil)

	ChkMnChnLen(t, genNd, 1)
}
