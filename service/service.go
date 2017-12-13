/*
Package service
*/
package service

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"jcqts/redpocketdemo/operator"
)

const (
	RELEASE_INTERVAL = 24
)

/*
A Service is the body of this server

population, the user list of persons
lock, locking while operating on population
*/
type Service struct {
	Population  map[string]*operator.Person //user list
	ServiceLock *sync.RWMutex
}

/*
This func generates a Service instance
*/
func NewService() *Service {
	service := &Service{}
	service.Population = make(map[string]*operator.Person, 0)
	service.ServiceLock = &sync.RWMutex{}

	return service
}

/*
This func handling the http request
*/
func (this *Service) Start() {
	go this.updateRedPocket()
	http.HandleFunc("/", this.Handler)
	http.HandleFunc("/register", this.Register)
	http.HandleFunc("/login", this.Login)
	http.ListenAndServeTLS(":8282", "server.crt", "server.key", nil)
}

/*
The root handler
we just make a simple text here
*/
func (this *Service) Handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Https connected!")
}

/*
The register handler
we just make a simple text here
*/
func (this *Service) Register(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Register succeed!")
}

/*
The login handler
we just make a simple text here
*/
func (this *Service) Login(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Login succeed!")
}

/*
This func runs in a goroutine independently,
it makes a timer and calls the doUpdateRedPocket() repeatly
*/
func (this *Service) updateRedPocket() {
	ticker := time.NewTicker(time.Hour * RELEASE_INTERVAL)
	for {
		<-ticker.C
		log.Printf("updateRedPocket: Try release")
		this.doUpdateRedPocket()
	}
}

/*
This func recycle all the expired redpockets
it ranges over all the redpockets for every person in the population
*/
func (this *Service) doUpdateRedPocket() {
	this.ServiceLock.Lock()
	defer this.ServiceLock.Unlock()

	for _, p := range this.Population {
		for _, rp := range p.GetRedPocketList() {
			if time.Now().Unix()-rp.GetTimestamp() >= RELEASE_INTERVAL {
				p.RecycleRedPocket(rp.Rpid)
			}
		}
	}
}
