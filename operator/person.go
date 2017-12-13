/*
Package person
It holds all the funcs of person operations
*/
package operator

import (
	"log"
	"math"
	"math/rand"
	"sync"
	"time"

	"jcqts/redpocketdemo/utils"
)

/*
A Person holds the operations of person

pid, the uuid of a person
leftMoney, this money a person current have
redPocketList, it stores all the sent redpockets of a person
*/
type Person struct {
	Pid               string
	Username		string
	Password		string
	LeftMoney         float32
	RedPocketList     map[string]*RedPocket
	RecvRedPocketList map[string]*RedPocket
	PersonLock        *sync.RWMutex
}

/*
This func generates a Person instance
value: init money, normally 0
*/
func NewPerson(username string, password string, value float32) *Person {
	if value < 0 {
		return nil
	}
	person := &Person{}
	person.Pid = utils.GetUuid()
	person.Username =username
	person.Password = password
	person.LeftMoney = value
	person.RedPocketList = make(map[string]*RedPocket, 0)
	person.RecvRedPocketList = make(map[string]*RedPocket, 0)
	person.PersonLock = &sync.RWMutex{}

	return person
}

/*
This func return the pid of a person instance
*/
func (this *Person) GetPid() string {
	return this.Pid
}

/*
This func peeks a person's leftMoney
*/
func (this *Person) PeekLeftMoney() (leftMoney float32) {
	this.PersonLock.Lock()
	defer this.PersonLock.Unlock()

	leftMoney = this.LeftMoney
	return
}

/*
This func peeks a person's recv redpocket list

In this func, we make a deep copy of a person's redpocketList
*/
func (this *Person) PeekRecvRedPocketList() (rpList map[string]*RedPocket) {
	this.PersonLock.Lock()
	defer this.PersonLock.Unlock()

	rpList = make(map[string]*RedPocket, 0)
	for rpid, rp := range this.RecvRedPocketList {
		rpList[rpid] = rp.CopyOne()
	}

	return
}

/*
This func get a person's redpocket list

In this func, we return person's original redpocketList
*/
func (this *Person) GetRedPocketList() map[string]*RedPocket {
	this.PersonLock.Lock()
	defer this.PersonLock.Unlock()

	return this.RedPocketList
}

/*
This func is used to get length of current sent redpocket
*/
func (this *Person) GetRedPocketListLen() int {
	this.PersonLock.Lock()
	defer this.PersonLock.Unlock()

	return len(this.RedPocketList)
}

/*
This func send a redpocket

value, the total money of this redpocket
count, how many persons can get money from this redpocket
password, the token of this redpocket
*/
func (this *Person) SendRedPocket(value float32, count int) (password string, rp *RedPocket) {
	if value <= 0 || count <= 0 {
		log.Printf("SendRedPocket: init value error: %f, %d", value, count)
		return
	}
	password, rp = NewRedPocket(value, count, this.Pid)
	this.PersonLock.Lock()
	defer this.PersonLock.Unlock()

	this.RedPocketList[rp.Rpid] = rp
	log.Printf("SendRedPocket: pid:%s rpid:%s", this.Pid, rp.Rpid)

	return
}

/*
This func get money from redpocket

rp, thie given redpocket
password, the password of current redpocket

This func return the money got from the redpocket,
if the redpocket is empty, it return float32(0)
*/
func (this *Person) OpenRedPocket(rp *RedPocket, password string) (getMoney float32) {
	this.PersonLock.Lock()
	defer this.PersonLock.Unlock()

	var found bool
	rpid := rp.Rpid
	if _, found = this.RecvRedPocketList[rpid]; found { //already opened
		log.Printf("OpenRedPocket: rpid %s already opened", rpid)
		return
	}
	if password != rp.Password {
		log.Printf("OpenRedPocket: rpid %s password failed", rpid)
		return
	}

	log.Printf("OpenRedPocket: rpid %s pass", rpid)
	getMoney = getRandomMoney(rp)
	if getMoney == float32(0) {
		log.Printf("OpenRedPocket: RedPocket empty")
		return
	}

	log.Printf("OpenRedPocket: pid:%s, rpid:%s, getRandomMoney:%f", this.Pid, rp.Rpid, getMoney)
	this.LeftMoney = this.LeftMoney + getMoney
	this.RecvRedPocketList[rp.Rpid] = rp
	log.Printf("OpenRedPocket: pid:%s leftMoney:%f", this.Pid, this.LeftMoney)
	return
}


/*
This func recycles a redpocket by rpid
it saves the money left in the redpocket back to a person's leftMoney

rpid, this uuid of a redpocket
*/
func (this *Person) RecycleRedPocket(rpid string) {
	this.PersonLock.Lock()
	defer this.PersonLock.Unlock()

	var found bool
	var rp *RedPocket
	if rp, found = this.RedPocketList[rpid]; !found {
		return
	}

	this.LeftMoney = this.LeftMoney + rp.LeftMoney
	rp.Release()
	delete(this.RedPocketList, rpid)
	log.Printf("RecycleRedPocket: Person leftMoney: %f", this.LeftMoney)
}

/*
This func calcs the money a person can get from a given redpocket

notice: the minimum value is 0.01 and the unit is forbidden to Fen
If the redpocket is empty, it returns float32(0)
*/
func getRandomMoney(rp *RedPocket) (getMoney float32) {
	rp.RpLock.Lock()
	defer rp.RpLock.Unlock()

	if rp.LeftCount == 0 {
		return
	}
	if rp.LeftCount == 1 {
		rp.LeftCount--
		getMoney = (float32)(math.Floor(float64(rp.LeftMoney*100+0.5)) / (float64)(100))
		return
	}

	min := float32(0.01) //minimum money
	max := rp.LeftMoney / (float32)(rp.LeftCount) * 2
	rand.Seed(time.Now().Unix())
	getMoney = rand.Float32() * max
	if getMoney <= min {
		getMoney = min
	}
	getMoney = (float32)(math.Floor(float64(getMoney*100)) / 100)
	rp.LeftCount--
	rp.LeftMoney = rp.LeftMoney - getMoney

	return
}
