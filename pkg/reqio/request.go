package reqio

import (
	"encoding/json"
	"github.com/evorts/feednomity/pkg/crypt"
	"github.com/evorts/feednomity/pkg/session"
	"io/ioutil"
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
	UnmarshallForm(dst interface{}) error
	UnmarshallBody(dst interface{}) error
	GetFormValue(field string) []string
	GetToken() string
	RenewToken() IRequest
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

func (req *request) RenewToken() IRequest {
	req.Prepare()
	return req
}

func (req *request) IsMethodGet() bool {
	return strings.ToUpper(req.r.Method) == "GET"
}

func (req *request) UnmarshallForm(dst interface{}) error {
	err := req.r.ParseForm()
	if err != nil {
		return err
	}
	result := make(map[string]interface{})
	for k, v := range req.r.Form {
		if len(v) == 1 {
			result[k] = v[0]
		} else {
			result[k] = v
		}
	}
	rs, err2 := json.Marshal(result)
	if err2 != nil {
		return err2
	}
	return json.Unmarshal(rs, dst)
}

func (req *request) UnmarshallBody(dst interface{}) error {
	body, err := ioutil.ReadAll(req.r.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, dst)
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
