/*
Package service
*/
package service

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"net/url"
	"time"
	"html/template"
	"strconv"

	"jcqts/redpocketdemo/utils"
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

type UserList struct{
	v  map[string]*operator.Person
	sync.RWMutex
}

type Service struct {
	Population  UserList //user list
	providers map[string]*operator.Provider
	ssManager *operator.Manager
	mgrLock *sync.RWMutex
}

/*
This func generates a Service instance
*/
func NewService() *Service {
	service := &Service{}
	service.Population.v = make(map[string]*operator.Person, 0)
    service.ssManager = operator.NewManager("memory","gosessionid",3600)
	service.mgrLock = &sync.RWMutex{}

	return service
}

/*
This func handling the http request
*/
func (this *Service) Start() {
	go this.updateRedPocket()
	http.HandleFunc("/", this.home)
	http.HandleFunc("/register", this.register)
	http.HandleFunc("/login", this.login)
	http.HandleFunc("/logout", this.logout)
	http.HandleFunc("/send", this.sendRedPocket)
	http.HandleFunc("/open", this.openRedPocket)
	http.HandleFunc("/peek", this.peekRedPocketList)
	http.HandleFunc("/get", this.getRedPocketList)
	http.ListenAndServeTLS(":8282", "server.crt", "server.key", nil)
}

/*
The root handler
we just make a simple text here
*/
func (this *Service) home(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method, ", r.Method)
	if r.Method == "GET"{
		t, _ := template.ParseFiles("view/home.html")
        log.Println(t.Execute(w, nil))
	}
}

/*
The register handler
we just make a simple text here
*/
func (this *Service) register(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method, ", r.Method)
	if r.Method == "GET"{
		t, _ := template.ParseFiles("view/register.html")
        log.Println(t.Execute(w, nil))
	}else{
		r.ParseForm()
		username:=r.Form["username"][0]
		password1:=r.Form["password1"][0]
		password2:=r.Form["password2"][0]
        fmt.Println("username: ", username)
        fmt.Println("password1: ", password1)
        fmt.Println("password2: ", password2)

		if username == "" || password1 == "" || password2 == ""{
			fmt.Fprintf(w, "Register failed! Empty param")
			return
		}

		this.Population.Lock()
		defer this.Population.Unlock()
		if _, found:=this.Population.v[username]; found{
			fmt.Fprintf(w, "Register failed! Username[%s] occupied", username)
			return
		}
		p := operator.NewPerson(username, password1, 0)
		this.Population.v[username] = p
		http.Redirect(w, r, "/login", 302)
	}
}

/*
The login handler
we just make a simple text here
*/
func (this *Service) login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method, ", r.Method)
	if r.Method == "GET"{
		t, _ := template.ParseFiles("view/login.html")
        log.Println(t.Execute(w, nil))
	}else{

		s := this.startSesson(w,r)
		if s == nil{
			fmt.Fprintf(w, "Login falied!, session is nil")
			return
		}

		r.ParseForm()
		username:=r.Form["username"][0]
		password:=r.Form["password"][0]
        fmt.Println("username: ", username)
        fmt.Println("password: ", password)

		this.Population.RLock()
		defer this.Population.RUnlock()
		p, found := this.Population.v[username]
		if !found{
			fmt.Fprintf(w, "Login failed! Username[%s] Not found", username)
			return
		}
		if p.Password != password{
			fmt.Fprintf(w, "Login failed! Username[%s] Username mismatch", username)
			return
		}
		s.Set("username", username)
		s.SetPerson(p)
		fmt.Printf("Login succeed!, username[%s], pid[%s], sid[%s]", username, p.Pid, s.SId)
		http.Redirect(w, r, "/send", 302)
	}
}

/*
The login handler
we just make a simple text here
*/
func (this *Service) logout(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method, ", r.Method)


	cookie, err := r.Cookie(this.ssManager.CookieName)
	if err !=nil || cookie.Value == ""{
		http.Redirect(w, r, "/login", 302)
		return
	}

	if cookie.Value != ""{
		cookie.MaxAge = -1
		http.SetCookie(w, cookie)
	}
	
	http.Redirect(w, r, "/", 302)
}

/*
This func get a session from manager
if not found, then create a new session
*/
func (this *Service) startSesson(w http.ResponseWriter, r *http.Request) (session *operator.Session){
	this.mgrLock.Lock()
	defer this.mgrLock.Unlock()

	cookie, err := r.Cookie(this.ssManager.CookieName)
	if err !=nil || cookie.Value == ""{
		sessionId := utils.GetUuid()
		session = this.ssManager.InitSession(sessionId)

		cookie := &http.Cookie{Name: this.ssManager.CookieName, Value: url.QueryEscape(sessionId), Path: "/", HttpOnly: true, MaxAge: int(this.ssManager.Period)}
        http.SetCookie(w, cookie)
	}else{
		fmt.Println("cookie.Value:",cookie.Value)
		sessionId, _ := url.QueryUnescape(cookie.Value)
		fmt.Println("sessionId:", sessionId)
        session = this.ssManager.ReadSession(sessionId)
	}

	return 
}

/*
This func hanles the send request
it will call Person.SendResPocket() to send a redpocket
*/
func (this *Service) sendRedPocket(w http.ResponseWriter, r *http.Request) {
	s := this.startSesson(w,r)
	if s == nil{
		fmt.Fprintf(w, "Send falied!, session is nil")
		return
	}

	fmt.Println("method, ", r.Method)
	if r.Method == "GET"{
		t, _ := template.ParseFiles("view/send.html")
        log.Println(t.Execute(w, nil))
	}else{
		r.ParseForm()
		money, err:=strconv.ParseFloat(r.Form["money"][0], 32)
		if err !=nil{
        	fmt.Println("money is not digital: ", r.Form["money"][0])
			fmt.Fprintf(w, "Send falied!, money is not digital")
			return
		}
		count, err:=strconv.Atoi(r.Form["count"][0])
		if err !=nil{
        	fmt.Println("count is not interger: ", r.Form["count"][0])
			fmt.Fprintf(w, "Send falied!, count is not digital")
			return
		}
        fmt.Println("money: ", float32(money))
        fmt.Println("count: ", count)

		pwd, rp := s.P.SendRedPocket(float32(money), count)
		fmt.Fprintf(w, "Send succeed!, redpocket id[%s] password[%s]", rp.Rpid, pwd)
	}

}

/*
This func hanles the open request
it will call Person.OpenResPocket() to open a redpocket
*/
func (this *Service) openRedPocket(w http.ResponseWriter, r *http.Request) {
	s := this.startSesson(w,r)
	if s == nil{
		fmt.Fprintf(w, "Open falied!, session is nil")
		return
	}

	fmt.Println("method, ", r.Method)
	if r.Method == "GET"{
		t, _ := template.ParseFiles("view/open.html")
        log.Println(t.Execute(w, nil))
	}else{
		r.ParseForm()
		rpid:=r.Form["rpid"][0]
		password:=r.Form["password"][0]
        fmt.Println("rpid: ", rpid)
        fmt.Println("password: ", password)

		rp := this.findRedPocketByRpid(rpid)
		if rp == nil{
			fmt.Fprintf(w, "Open falied!, findRedPocketByRpid[%s] Not Found", rpid)
			return
		}
		money := s.P.OpenRedPocket(rp, password)
		if money == float32(0){
			fmt.Fprintf(w, "Open falied!, redpocket[%s] already opened", rpid)
			return
		}
		fmt.Fprintf(w, "Open succeed!, redpocket[%s] password[%s] money[%s]", rpid, password, money)
	}

}

/*
This func hanles the get request
it will call Person.PeekResPocketList() to open a redpocket
*/
func (this *Service) peekRedPocketList(w http.ResponseWriter, r *http.Request) {
	s := this.startSesson(w,r)
	if s == nil{
		fmt.Fprintf(w, "Peek falied!, session is nil")
		return
	}

	rpList := s.P.PeekRecvRedPocketList()
	var rpids string
	for _, v := range rpList{
		rpids = rpids +"\n"+v.Rpid
	}
	fmt.Fprintf(w, "Peek succeed!, redpockets[%s]", rpids)
}

/*
This func hanles the get request
it will call Person.GetResPocketList() to open a redpocket
*/
func (this *Service) getRedPocketList(w http.ResponseWriter, r *http.Request) {
	s := this.startSesson(w,r)
	if s == nil{
		fmt.Fprintf(w, "Get falied!, session is nil")
		return
	}

	rpList := s.P.GetRedPocketList()
	var rpids string
	for _, v := range rpList{
		rpids = rpids +"\n"+v.Rpid
	}
	fmt.Fprintf(w, "Get succeed!, redpockets[%s]", rpids)
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
	this.Population.Lock()
	defer this.Population.Unlock()
	for _, p := range this.Population.v {
		for _, rp := range p.GetRedPocketList() {
			if time.Now().Unix()-rp.GetTimestamp() >= RELEASE_INTERVAL {
				p.RecycleRedPocket(rp.Rpid)
			}
		}
	}
}


func (this *Service) findRedPocketByRpid(rpid string) (rp * operator.RedPocket){
	this.Population.Lock()
	defer this.Population.Unlock()
	for _, p := range this.Population.v {
		log.Printf("findRedPocketByRpid: pid[%s]\n", p.Pid)
		for _, rp = range p.GetRedPocketList() {
			log.Printf("findRedPocketByRpid: rpid[%s]\n", rp.Rpid)
			if rp.Rpid == rpid{
				return 
			}
		}
	}

	return
}