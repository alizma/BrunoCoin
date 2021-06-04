package test

import (
	"BrunoCoin/pkg/utils"
	"testing"
)

func TestAddNilBlck(t *testing.T) {
	utils.SetDebug(true)
	genNd := NewGenNd()
	genNd.Start()

	genNd.Chain.Add(nil)

	ChkMnChnLen(t, genNd, 1)
}
