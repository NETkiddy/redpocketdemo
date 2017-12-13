package test

import (
	"testing"
	"jcqts/redpocketdemo/operator"
)


func TestNewRedPocket(t *testing.T) {
	per1 := operator.NewPerson("admin","123456",values1[6].initMoney)
	if per1 == nil {
		return
	}

	for _, v := range values2 {
		_, rp := operator.NewRedPocket(v.value, v.count, per1.Pid)
		if rp == nil {
			t.Log("TestSendRedPocket failed on %v", v)
		}
	}
}


func TestRelease(t *testing.T){
    per1 := operator.NewPerson("admin","123456",values1[6].initMoney)
	if per1 == nil {
		return
	}

		_, rp1 := operator.NewRedPocket(values2[6].value, values2[6].count, per1.Pid)
		if rp1 == nil {
			t.Log("TestSendRedPocket failed on %v", values2[6])
		}

		t.Log("TestSendRedPocket1, LeftMoney %v", rp1.LeftMoney)
		t.Log("TestSendRedPocket1, LeftCount %v", rp1.LeftCount)

		rp1.Release()

		t.Log("TestSendRedPocket2, LeftMoney %v", rp1.LeftMoney)
		t.Log("TestSendRedPocket2, LeftCount %v", rp1.LeftCount)


}


