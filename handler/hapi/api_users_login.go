package hapi

import (
	"github.com/evorts/feednomity/domain/users"
	"github.com/evorts/feednomity/pkg/api"
	"github.com/evorts/feednomity/pkg/crypt"
	"github.com/evorts/feednomity/pkg/database"
	"github.com/evorts/feednomity/pkg/jwe"
	"github.com/evorts/feednomity/pkg/logger"
	"github.com/evorts/feednomity/pkg/reqio"
	"github.com/evorts/feednomity/pkg/validate"
	"github.com/evorts/feednomity/pkg/view"
	"gopkg.in/square/go-jose.v2/jwt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func ApiLogin(w http.ResponseWriter, r *http.Request) {
	req := reqio.NewRequest(w, r).PrepareRestful()
	log := req.GetContext().Get("logger").(logger.IManager)
	vm := req.GetContext().Get("view").(view.IManager)
	hash := req.GetContext().Get("hash").(crypt.ICryptHash)
	datasource := req.GetContext().Get("db").(database.IManager)

	log.Log("login_api_handler", "request received")

	var payload struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Remember string `json:"remember"`
		Csrf     string `json:"csrf"`
	}

	err := req.UnmarshallBody(&payload)
	if err != nil {
		_ = vm.RenderJson(w, http.StatusBadRequest, api.Response{
			Status:  http.StatusBadRequest,
			Content: make(map[string]interface{}, 0),
			Error: &api.ResponseError{
				Code:    "LOG:ERR:BND",
				Message: "Bad Request! Something wrong with the payload of your request.",
				Reasons: make(map[string]string, 0),
				Details: make([]interface{}, 0),
			},
		})
		return
	}
	// Validate request
	errs := make(map[string]string, 0)
	if !validate.ValidUsername(payload.Username) || !validate.ValidEmail(payload.Username) {
		errs["username"] = "Not a valid username or email!"
	}
	if !validate.ValidPassword(payload.Password) {
		errs["password"] = "Not a valid password!"
	}
	if len(errs) > 0 {
		_ = vm.RenderJson(w, http.StatusBadRequest, api.Response{
			Status:  http.StatusBadRequest,
			Content: make(map[string]interface{}, 0),
			Error: &api.ResponseError{
				Code:    "LOG:ERR:VAL",
				Message: "Bad Request! Validation error.",
				Reasons: errs,
				Details: make([]interface{}, 0),
			},
		})
		return
	}
	var user *users.User
	usersDomain := users.NewUserDomain(datasource)

	if validate.ValidEmail(payload.Username) {
		user, err = usersDomain.FindByUserEmail(req.GetContext().Value(), payload.Username)
	} else {
		user, err = usersDomain.FindByUsername(req.GetContext().Value(), payload.Username)
	}

	if err != nil {
		_ = vm.RenderJson(w, http.StatusBadRequest, api.Response{
			Status:  http.StatusBadRequest,
			Content: make(map[string]interface{}, 0),
			Error: &api.ResponseError{
				Code:    "LOG:ERR:USR",
				Message: "Bad Request! User not found.",
				Reasons: map[string]string{"err": err.Error()},
				Details: make([]interface{}, 0),
			},
		})
		return
	}
	// ensure the user and password are correct
	passCrypt := hash.RenewHash().HashWithoutSalt(payload.Password)
	if len(strings.Trim(user.Password, " ")) < 1 || strings.ToLower(passCrypt) != strings.ToLower(strings.TrimLeft(user.Password, "\\x")) {
		_ = vm.RenderJson(w, http.StatusBadRequest, api.Response{
			Status:  http.StatusBadRequest,
			Content: make(map[string]interface{}, 0),
			Error: &api.ResponseError{
				Code:    "LOG:ERR:ATH",
				Message: "Bad Request! authentication failed.",
				Reasons: make(map[string]string, 0),
				Details: make([]interface{}, 0),
			},
		})
		return
	}

	//@todo: check multi-device login limitation

	// find out the organizations
	var (
		g        []*users.Group
		gid, oid int64
	)
	g, err = usersDomain.FindGroupByIds(req.GetContext().Value(), user.GroupId)
	if err != nil {
		_ = vm.RenderJson(w, http.StatusBadRequest, api.Response{
			Status:  http.StatusBadRequest,
			Content: make(map[string]interface{}, 0),
			Error: &api.ResponseError{
				Code:    "LOG:ERR:GRP",
				Message: "Bad Request! Invalid group.",
				Reasons: map[string]string{"err": err.Error()},
				Details: make([]interface{}, 0),
			},
		})
		return
	}
	gid = g[0].Id
	oid = g[0].OrgId
	g, err = usersDomain.FindGroupByOrgId(req.GetContext().Value(), oid)
	if err != nil {
		_ = vm.RenderJson(w, http.StatusBadRequest, api.Response{
			Status:  http.StatusBadRequest,
			Content: make(map[string]interface{}, 0),
			Error: &api.ResponseError{
				Code:    "LOG:ERR:ORG",
				Message: "Bad Request! Invalid organization.",
				Reasons: map[string]string{"err": err.Error()},
				Details: make([]interface{}, 0),
			},
		})
		return
	}

	//by default expiration only in 6 hour
	expiration := 6 * time.Hour
	if len(payload.Remember) > 0 {
		expiration = 3 * 24 * time.Hour
	}
	userData := reqio.UserData{
		Id:          user.Id,
		Username:    user.Username,
		DisplayName: user.DisplayName,
		Attributes:  user.Attributes,
		Email:       user.Email,
		Phone:       user.Phone,
		AccessRole:  string(user.AccessRole),
		JobRole:     user.JobRole,
		Assignment:  user.Assignment,
		GroupId:     gid,
		OrgId:       oid,
		OrgGroupIds: make([]int64, 0),
	}
	for _, gv := range g {
		userData.OrgGroupIds = append(userData.OrgGroupIds, gv.Id)
	}
	jwxToken, _ := req.GetJwx().Encode(
		jwt.Claims{
			Issuer:   jwe.ISSUER,
			Subject:  "basic",
			Expiry:   jwt.NewNumericDate(time.Now().Add(expiration)),
			IssuedAt: jwt.NewNumericDate(time.Now()),
			ID:       strconv.FormatInt(userData.Id, 10),
		},
		jwe.PrivateClaims{
			ClientId:    req.GetClientId(),
			Id:          userData.Id,
			Username:    userData.Username,
			DisplayName: userData.DisplayName,
			Attributes:  userData.Attributes,
			Email:       userData.Email,
			Phone:       userData.Phone,
			AccessRole:  userData.AccessRole,
			JobRole:     userData.JobRole,
			Assignment:  userData.Assignment,
			GroupId:     userData.GroupId,
			OrgId:       userData.OrgId,
			OrgGroupIds: userData.OrgGroupIds,
		},
	)
	_ = vm.RenderJson(w, http.StatusOK, api.Response{
		Status: http.StatusOK,
		Content: map[string]interface{}{
			"token": jwxToken,
		},
	})
}
