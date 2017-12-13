/*
Package redpocket
It holds all the funcs of redpocket operations
*/
package operator

import (
	"log"
	"sync"
	"time"

	"jcqts/redpocketdemo/utils"
)

/*
A RedPocket holds the body of the redpocket

rpid, this uuid of a redpocket
leftMoney, the money a redpocket current left
leftCount, the count a redpocket current left
password, the token of a redpocket
timestamp, the unix time when a redpocket is sent
owner, the person who send this redpocket
*/
type RedPocket struct {
	Rpid      string
	LeftMoney float32
	LeftCount int
	Password  string
	Timestamp int64
	RpLock    *sync.RWMutex
}

/*
This func generates a redpocket instance

value, total money
count, total count
password, the redpocket token
owner, who send this redpocket
*/
func NewRedPocket(value float32, count int, pid string) (pwd string, rp *RedPocket) {
	if value <= 0 || count <= 0 || value > 200 {
		log.Printf("NewRedPocket: failed, init value not available")
		return 
	}
	if value/float32(count) < 0.01 {
		log.Printf("NewRedPocket: failed, value/count<0.01")
		return 
	}
	rp = &RedPocket{}
	rp.Rpid = pid + "_" + utils.GetUuid()
	rp.LeftMoney = value
	rp.LeftCount = count
	rp.Password = utils.GetPassword(8)
	rp.Timestamp = time.Now().Unix()
	rp.RpLock = &sync.RWMutex{}

	pwd = rp.Password
	return
}

/*
This func recycles the redpocket when is expired
*/
func (this *RedPocket) Release() {
	this.RpLock.Lock()
	defer this.RpLock.Unlock()

	if this.LeftCount == 0 || this.LeftMoney == float32(0){
		return
	}

	this.LeftCount = 0
	this.LeftMoney = float32(0)
}

/*
This func makes a deep copy of a redpocket
*/
func (this *RedPocket) CopyOne() *RedPocket {
	this.RpLock.Lock()
	defer this.RpLock.Unlock()

	rp := &RedPocket{}
	rp.Rpid = this.Rpid
	rp.LeftMoney = this.LeftMoney
	rp.Password = this.Password
	rp.Timestamp = this.Timestamp
	rp.RpLock = &sync.RWMutex{}

	return rp
}

/*
This func return timestamp
*/
func (this *RedPocket) GetTimestamp() int64 {
	return this.Timestamp
}
