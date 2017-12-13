package test

import (
	"testing"
	"jcqts/redpocketdemo/operator"
)

var values0 = []struct {
	username string
	password string
	initMoney float32
}{
	{"","",-1},
	{"admin","123456",0},
	{"admin","123456",10},
}

func TestNewPerson(t *testing.T) {
	for _, v := range values0 {
		per1 := operator.NewPerson(v.username, v.password, v.initMoney)
		if per1 == nil {
			t.Log("TestNewPerson failed on %v", v)
		}
	}
}

var values1 = []struct {
	initMoney   float32
	initRPMoney float32
	count       int
}{
	{0, -1, 0},
	{0, 0, -1},
	{0, 0, 0},
	{0, 0, 2},
	{0, 10, 0},
	{0, 10, 2},
	{100, 10, 2},
}

func TestSendRedPocket(t *testing.T) {
	for _, v := range values1 {
		per1 := operator.NewPerson("admin","123456",v.initMoney)
		if per1 == nil {
			return
		}
		_, rp := per1.SendRedPocket(v.initRPMoney, v.count) 
		if rp == nil {
			t.Log("TestSendRedPocket failed on %v", v)
		}
	}
}

var values2 = []struct {
	value float32
	count int
}{
	{-1, 1},
	{1, -1},
	{0, 1},
	{1, 0},
	{201, 1},
	{0.02, 3},
	{10, 3},
	{10, 5},
}


func TestOpenRedPocket(t *testing.T) {
	per1 := operator.NewPerson("admin","123456",values1[6].initMoney)
	per2 := operator.NewPerson("admin","123456",values1[6].initMoney)
	if per1 == nil || per2 == nil {
		return
	}

		pwd, rp := operator.NewRedPocket(values2[6].value, values2[6].count, per1.Pid)
		if rp == nil {
			return
		}
		// error pwd
		tmp := "1234qwer"
		if float32(0) != per2.OpenRedPocket(rp, tmp) {
			t.Fatal("TestOpenRedPocket failed on error pwd %v", tmp)
		}

		// correct pwd
		if float32(0) == per2.OpenRedPocket(rp, pwd) {
			t.Fatal("TestOpenRedPocket failed on correct pwd %v", pwd)
		}

		// correct pwd twice, 
		if float32(0) != per2.OpenRedPocket(rp, pwd) {
			t.Fatal("TestOpenRedPocket failed on correct pwd %v twice", pwd)
		}
}

var values5 = []struct {
	initRPMoney float32
	count       int
}{
	{10, 2},
	{10, 5},
	{1, 3},
}
func TestPeekLeftMoney(t *testing.T) {

		per1 := operator.NewPerson("admin","123456",0)
		per2 := operator.NewPerson("admin","123456",1000)

	pwd := make([]string, 0)
	rp := make([]*operator.RedPocket,0)

	for _, v := range values5 {
		if per1 == nil {
			return
		}
		//SendRedPocket won't decrease leftmoney
		t.Log("TestPeekLeftMoney1 per1, %v", per1.PeekLeftMoney())
		t.Log("TestPeekLeftMoney1 per2, %v", per2.PeekLeftMoney())
		pwdtmp, rptmp := per1.SendRedPocket(v.initRPMoney, v.count) 
		if rptmp == nil{
			t.Log("TestPeekLeftMoney failed on %v", v)
		}
		pwd = append(pwd, pwdtmp)
		rp = append(rp, rptmp)
		t.Log("TestPeekLeftMoney2 per1, %v", per1.PeekLeftMoney())
		t.Log("TestPeekLeftMoney2 per2, %v", per2.PeekLeftMoney())
	}
		// open 1
		if float32(0) == per2.OpenRedPocket(rp[0], pwd[0]) {
			t.Fatal("TestPeekLeftMoney failed on %v, %v", rp[0], pwd[0])
		}
		t.Log("TestPeekLeftMoney3 per2, %v", per2.PeekLeftMoney())

		// open 2
		if float32(0) == per2.OpenRedPocket(rp[1], pwd[1]) {
			t.Fatal("TestPeekLeftMoney failed on %v, %v", rp[1], pwd[1])
		}
		t.Log("TestPeekLeftMoney4 per2, %v", per2.PeekLeftMoney())
		if float32(0) == per2.OpenRedPocket(rp[2], pwd[2]) {
			t.Fatal("TestPeekLeftMoney failed on %v, %v", rp[2], pwd[2])
		}
		t.Log("TestPeekLeftMoney5 per2, %v", per2.PeekLeftMoney())
	
	
}

func TestPeekRecvRedPocketList(t *testing.T) {
	per1 := operator.NewPerson("admin","123456",values1[6].initMoney)
	per2 := operator.NewPerson("admin","123456",values1[6].initMoney)
	if per1 == nil || per2 == nil {
		return
	}

		pwd1, rp1 := operator.NewRedPocket(values2[6].value, values2[6].count, per1.Pid)
		if rp1 == nil {
			return
		}
		pwd2, rp2 := operator.NewRedPocket(values2[6].value, values2[6].count, per1.Pid)
		if rp2 == nil {
			return
		}
		pwd3, rp3 := operator.NewRedPocket(values2[6].value, values2[6].count, per1.Pid)
		if rp3 == nil {
			return
		}
		t.Log("TestPeekRecvRedPocketList1, leftMoney %v", per2.PeekLeftMoney())
		t.Log("TestPeekRecvRedPocketList1, len %v", len(per2.PeekRecvRedPocketList()))

		if float32(0) == per2.OpenRedPocket(rp1, pwd1) {
			t.Fatal("TestPeekRecvRedPocketList failed on %v", pwd1)
		}
		t.Log("TestPeekRecvRedPocketList2, leftMoney %v", per2.PeekLeftMoney())
		t.Log("TestPeekRecvRedPocketList2, len %v", len(per2.PeekRecvRedPocketList()))
		if float32(0) == per2.OpenRedPocket(rp2, pwd2) {
			t.Fatal("TestPeekRecvRedPocketList failed on %v", pwd3)
		}
		if float32(0) == per2.OpenRedPocket(rp3, pwd3) {
			t.Fatal("TestPeekRecvRedPocketList failed on %v", pwd3)
		}
		t.Log("TestPeekRecvRedPocketList3, leftMoney %v", per2.PeekLeftMoney())
		t.Log("TestPeekRecvRedPocketList3, len %v", len(per2.PeekRecvRedPocketList()))

}

var values4 = []struct {
	value float32
	count int
	rpid string
}{
	{10, 4,"1111-1"},
	{5, 2,"1111-2"},
	{0.5, 5,"1111-3"},
}

func TestRecycleRedPocket(t *testing.T){
	per1 := operator.NewPerson("admin","123456",values1[6].initMoney)
	if per1 == nil{
		return
	}
	t.Log("TestRecycleRedPocket0, len %v", per1.GetRedPocketListLen())
	t.Log("TestRecycleRedPocket0, leftMoney %v", per1.PeekLeftMoney())
	var rp []*operator.RedPocket
	for _, v := range values5 {
		_, rp1 := per1.SendRedPocket(v.initRPMoney, v.count) 
		if rp1 == nil {
			t.Fatal("TestSendRedPocket failed on %v", v)
		}

		rp = append(rp ,rp1)
	}
	t.Log("TestRecycleRedPocket1, len %v", per1.GetRedPocketListLen())
	t.Log("TestRecycleRedPocket0, leftMoney %v", per1.PeekLeftMoney())
	per1.RecycleRedPocket(rp[0].Rpid)
	t.Log("TestRecycleRedPocket2, len %v", per1.GetRedPocketListLen())
	t.Log("TestRecycleRedPocket0, leftMoney %v", per1.PeekLeftMoney())
	per1.RecycleRedPocket(rp[1].Rpid)
	per1.RecycleRedPocket(rp[2].Rpid)
	t.Log("TestRecycleRedPocket3, len %v", per1.GetRedPocketListLen())
	t.Log("TestRecycleRedPocket0, leftMoney %v", per1.PeekLeftMoney())

}
