package test

import (
	"BrunoCoin/pkg/miner"
	"BrunoCoin/pkg/utils"
	"testing"
)

func TestNilPriCalc(t *testing.T) {
	bullshitpriority := miner.CalcPri(nil)
	if bullshitpriority != 0 {
		t.Errorf("expected 0, actual; %d", bullshitpriority)
	}
}

func TestPriCalcTxiTxo(t *testing.T) {
	utils.SetDebug(true)
	genNd := NewGenNd()
	genNd.Start()

	minitx := MakeSingleTx(genNd, genNd.Id.GetPublicKeyBytes(), 2)

	calculatedpri := miner.CalcPri(minitx)

	if calculatedpri != (minitx.SumInputs()-minitx.SumOutputs())*100.0/minitx.Sz() {
		t.Errorf("expected %d, actual; %d", (minitx.SumInputs()-minitx.SumOutputs())*100.0/minitx.Sz(), calculatedpri)
	}
}
