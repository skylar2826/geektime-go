package cookie

import (
	"net/http"
)

type Propagator struct {
	cookieName    string
	cookieOptions func(cookie *http.Cookie)
}

type PropagatorOption func(*Propagator)

func NewPropagator(options ...PropagatorOption) *Propagator {
	res := &Propagator{
		cookieName:    "session_id",
		cookieOptions: func(cookie *http.Cookie) {},
	}

	for _, option := range options {
		option(res)
	}

	return res
}

func PropagatorWithCookieName(cookieName string) PropagatorOption {
	return func(p *Propagator) {
		p.cookieName = cookieName
	}
}

//func NewPropagator() *Propagator {
//	return &Propagator{
//		cookieName:    "session_id",
//		cookieOptions: func(cookie *http.Cookie) {},
//	}
//}

func (p *Propagator) Inject(id string, w http.ResponseWriter) error {
	c := &http.Cookie{Name: p.cookieName, Value: id}
	p.cookieOptions(c)
	http.SetCookie(w, c)
	return nil
}

func (p *Propagator) Extract(r *http.Request) (string, error) {
	c, err := r.Cookie(p.cookieName)
	if err != nil {
		return "", err
	}
	return c.Value, nil
}

func (p *Propagator) Remove(w http.ResponseWriter) error {
	c := &http.Cookie{
		Name: p.cookieName,
		//MaxAge<0 means delete cookie now, equivalently 'Max-Age: 0'
		MaxAge: -1,
	}
	http.SetCookie(w, c)
	return nil
}
