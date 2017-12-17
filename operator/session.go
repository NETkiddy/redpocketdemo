/*
Package person
It holds all the funcs of person operations
*/
package operator

import(
    "sync"
    "fmt"
    "reflect"

)

type ISession interface {
    Set(key, value interface{}) error //set session value
    Get(key interface{}) interface{}  //get session value
    Delete(key interface{}) error     //delete session value
    SessionID() string                //back current sessionID
}

type Session struct {
    SId string
    Username string
    P *Person
}

/*
This interface abstracts the operations of a session
then the session can be coded easily to save to DB/redis/cache etc.
*/
type Provider interface {
    SessionInit(sid string) (Session, error)
    SessionRead(sid string) (Session, error)
    SessionDestroy(sid string) error
    SessionGC(maxLifeTime int64)
}

/*
This struct manages all the sessions
*/
type Manager struct {
    CookieName  string     //private cookiename
    Period int64    // max lifetime
    Sessions    map[string]*Session 
    Lock        *sync.RWMutex // protects session
}

func NewManager(providerName string, cookieName string, period int64) (m *Manager){
    m = &Manager{}
    m.CookieName = cookieName
    m.Period = period
    m.Sessions = make(map[string]*Session, 0)
    m.Lock = &sync.RWMutex{}

    return
}

func (this *Manager) InitSession(sessionId string) (s *Session){
    this.Lock.Lock()
    defer this.Lock.Unlock()
    
    s = &Session{}
    s.SId = sessionId
    this.Sessions[sessionId] = s
    fmt.Printf("sessionId[%s] Created", sessionId)
    return s
}

func (this *Manager) ReadSession(sessionId string) (s *Session){
    this.Lock.RLock()
    defer this.Lock.RUnlock()
    
    var found bool
    if s, found =this.Sessions[sessionId]; !found{
        fmt.Printf("sessionId[%s] Not Found", sessionId)
        return nil
    }
    return
}

func (this *Session)Set(key string, value interface{}){
    switch key{
        case "username":
            this.Username = reflect.ValueOf(value).String()
    }

}

func (this *Session)SetPerson(p *Person){
    this.P = p
}