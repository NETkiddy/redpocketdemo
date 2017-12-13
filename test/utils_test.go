package test

import (
	"testing"
	"jcqts/redpocketdemo/utils"
)

var values6 = []struct {
	size int
}{
	{-1},
	{0},
	{1},
	{8},
	{128},
}

func TestGetPassword(t *testing.T) {
	for _,v:=range values6{
		pwd := utils.GetPassword(v.size)
		t.Log("TestGetPassword got size: ",v.size, "-->",pwd)

	}
}

func TestGetUuid(t *testing.T){
	for i:=0;i<10;i++{
   		uuid := utils.GetUuid()
		t.Log("TestGetUuid got uuid:",uuid)
	}
}


