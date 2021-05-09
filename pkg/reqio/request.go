package reqio

import (
	"encoding/json"
	"github.com/evorts/feednomity/pkg/acl"
	"github.com/evorts/feednomity/pkg/crypt"
	"github.com/evorts/feednomity/pkg/jwe"
	"github.com/evorts/feednomity/pkg/session"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const (
	HeaderAuth         = "X-Authorization"
	HeaderClientId     = "X-Client-Id"
	UserContextKey     = "user"
	UserAccessScopeKey = "access_scope"
)

type UserData struct {
	Id          int64                  `json:"id"`
	Username    string                 `json:"username"`
	DisplayName string                 `json:"display_name"`
	Attributes  map[string]interface{} `json:"attributes"`
	Email       string                 `json:"email"`
	Phone       string                 `json:"phone"`
	AccessRole  string                 `json:"access_role"`
	JobRole     string                 `json:"job_role"`
	Assignment  string                 `json:"assignment"`
	GroupId     int64                  `json:"group_id"`
	OrgId       int64                  `json:"org_id"`
	OrgGroupIds []int64                `json:"org_group_ids"`
}

type request struct {
	w               http.ResponseWriter
	r               *http.Request
	ctx             IContext
	csrfToken       string
	jweToken        string
	clientId        string
	hash            crypt.ICryptHash
	session         session.IManager
	userData        *UserData
	expireAt        *time.Time
	userAccessScope acl.AccessScope
	jwx             jwe.IManager
}

type IRequest interface {
	IsMethodGet() bool
	IsMethodPost() bool
	IsMethodPut() bool
	IsMethodDelete() bool
	IsMethodOptions() bool
	IsLoggedIn() bool

	PrepareRestful() IRequest
	Prepare() IRequest
	UnmarshallForm(dst interface{}) error
	UnmarshallBody(dst interface{}) error
	GetFormValue(field string) []string
	GetCsrfToken() string
	GetToken() string
	RenewSessionToken() IRequest
	GetContext() IContext
	GetUserData() *UserData
	GetUserAccessScope() acl.AccessScope
	GetJweToken() string
	GetJwx() jwe.IManager
	GetClientId() string

	getUserAccessScopeFromContext() acl.AccessScope
	getUserAccessScopeFromSession() acl.AccessScope
	getUserDataFromSession() *UserData
	getUserDataFromContext() *UserData
}

func NewRequest(w http.ResponseWriter, r *http.Request) IRequest {
	ctx := NewContext(r.Context())
	req := &request{
		w:   w,
		r:   r,
		ctx: ctx,
	}
	if sm := ctx.Get("sm"); sm != nil {
		req.session = sm.(session.IManager)
	}
	if c := ctx.Get("hash"); c != nil {
		req.hash = c.(crypt.ICryptHash)
	}
	if j := ctx.Get("jwx"); j != nil {
		req.jwx = j.(jwe.IManager)
	}
	return req
}

func (req *request) GetContext() IContext {
	return req.ctx
}

func (req *request) GetCsrfToken() string {
	return req.csrfToken
}

func (req *request) GetToken() string {
	return req.jweToken
}

func (req *request) GetJwx() jwe.IManager {
	return req.jwx
}

func (req *request) RenewSessionToken() IRequest {
	req.Prepare()
	return req
}

func (req *request) GetFormValue(field string) []string {
	return req.r.Form[field]
}

func (req *request) PrepareRestful() IRequest {
	req.userData = req.getUserDataFromContext()
	//req.csrfToken = req.hash.HashWithSalt(time.Now().String())
	req.userAccessScope = req.getUserAccessScopeFromContext()
	return req
}

func (req *request) Prepare() IRequest {
	if req.hash == nil {
		return req
	}
	req.csrfToken = req.hash.HashWithSalt(time.Now().String())
	req.userData = req.getUserDataFromSession()
	req.userAccessScope = req.getUserAccessScopeFromSession()
	return req
}

func (req *request) IsLoggedIn() bool {
	if req.userData == nil || req.userData.Id < 1 {
		return false
	}
	return true
}

func (req *request) GetUserData() *UserData {
	return req.userData
}

func (req *request) getUserDataFromSession() *UserData {
	if req.session.GetJson(req.GetContext().Value(), UserContextKey, req.userData) != nil {
		return nil
	}
	return req.userData
}

func (req *request) getUserDataFromContext() *UserData {
	u := req.GetContext().Get(UserContextKey)
	if u == nil {
		return nil
	}
	return u.(*UserData)
}

func (req *request) GetJweToken() string {
	return strings.Trim(req.r.Header.Get(HeaderAuth), " ")
}

func (req *request) GetClientId() string {
	return strings.Trim(req.r.Header.Get(HeaderClientId), " ")
}

func (req *request) GetUserAccessScope() acl.AccessScope {
	return req.userAccessScope
}

func (req *request) getUserAccessScopeFromContext() acl.AccessScope {
	acc := req.GetContext().Get(UserAccessScopeKey)
	if acc == nil {
		return ""
	}
	return acc.(acl.AccessScope)
}

func (req *request) getUserAccessScopeFromSession() acl.AccessScope {
	return acl.AccessScope(req.session.GetString(req.GetContext().Value(), UserAccessScopeKey))
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

func (req *request) IsMethodGet() bool {
	return strings.ToUpper(req.r.Method) == "GET"
}

func (req *request) IsMethodPost() bool {
	return strings.ToUpper(req.r.Method) == "POST"
}

func (req *request) IsMethodPut() bool {
	return strings.ToUpper(req.r.Method) == "PUT"
}

func (req *request) IsMethodDelete() bool {
	return strings.ToUpper(req.r.Method) == "DELETE"
}

func (req *request) IsMethodOptions() bool {
	return strings.ToUpper(req.r.Method) == "OPTIONS"
}
