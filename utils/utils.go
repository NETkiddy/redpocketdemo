/*
Package utils
it holds some util funcs
*/
package utils

import (
	"log"
	"math/rand"
	"os/exec"
	"time"
)

const (
	PWD_BASE = "0123456789abcdefghijkmnopqrstuvwxyzABCDEFGHIJKLMNPQRSTUVWXYZ"
)

/*
This func generates the password

size, the lenth of the password
*/
func GetPassword(size int) (pwd string) {
	if size <= 0{
		return
	}
	tmp := make([]byte, size)
	rand.Seed(time.Now().UnixNano())
	for i := range tmp {
		tmp[i] = PWD_BASE[rand.Intn(len(PWD_BASE))]
	}

	pwd = string(tmp)
	log.Printf("getPassword: Generate pwd: %s", pwd)
	return
}

/*
This func generates a uuid
*/
func GetUuid() (uuid string) {
	tmp, err := exec.Command("uuidgen").Output()
	if err != nil {
		log.Fatal(err)
	}
	uuid = string(tmp)

	return
}
