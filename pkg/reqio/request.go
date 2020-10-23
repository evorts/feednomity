package reqio

import (
	"github.com/evorts/godash/pkg/crypt"
	"github.com/evorts/godash/pkg/session"
	"net/http"
	"strings"
	"time"
)

type request struct {
	w    http.ResponseWriter
	r    *http.Request
	ctx  IContext
	token   string
	crypt   crypt.ICrypt
	session session.IManager
}

type IRequest interface {
	IsMethodGet() bool
	IsLoggedIn() bool
	Prepare() IRequest
	ParseForm() error
	GetFormValue(field string) []string
	GetToken() string
	GetContext() IContext
}

func NewRequest(w http.ResponseWriter, r *http.Request) IRequest {
	ctx := NewContext(r.Context())
	req := &request{
		w:       w,
		r:       r,
		ctx:     ctx,
	}
	if sm := ctx.Get("sm"); sm != nil {
		req.session = sm.(session.IManager)
	}
	if c := ctx.Get("crypt"); c != nil {
		req.crypt = c.(crypt.ICrypt)
	}
	return req
}

func (req *request) GetContext() IContext {
	return req.ctx
}

func (req *request) GetToken() string {
	return req.token
}

func (req *request) IsMethodGet() bool {
	return strings.ToUpper(req.r.Method) == "GET"
}

func (req *request) ParseForm() error {
	return req.r.ParseForm()
}

func (req *request) GetFormValue(field string) []string {
	return req.r.Form[field]
}

func (req *request) Prepare() IRequest {
	if req.crypt == nil {
		return req
	}
	req.token = req.crypt.CryptWithSalt(time.Now().String())
	return req
}

func (req *request) IsLoggedIn() bool {
	user := req.session.Get(req.GetContext().Value(), "user")
	if user == nil || user == "" {
		return false
	}
	return true
}
