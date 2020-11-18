package handler

import (
	"encoding/json"
	"fmt"
	"github.com/evorts/feednomity/domain/feedbacks"
	"github.com/evorts/feednomity/pkg/api"
	"github.com/evorts/feednomity/pkg/config"
	"github.com/evorts/feednomity/pkg/crypt"
	"github.com/evorts/feednomity/pkg/database"
	"github.com/evorts/feednomity/pkg/logger"
	"github.com/evorts/feednomity/pkg/reqio"
	"github.com/evorts/feednomity/pkg/session"
	"github.com/evorts/feednomity/pkg/template"
	"github.com/evorts/feednomity/pkg/validate"
	"github.com/pkg/errors"
	"github.com/segmentio/ksuid"
	"net/http"
	"time"
)

type HashData struct {
	ExpireAt   time.Time   `json:"expire_at"`
	RealHash   string      `json:"real_hash"`
	Attributes interface{} `json:"attributes"`
}

type hashHelper struct {
	aes crypt.ICryptAES
}

type IHashHelper interface {
	Generate(expireAt time.Time, realHash string, attributes interface{}) string
	Decode(value string) (*HashData, error)
}

func NewHashHelper(aes crypt.ICryptAES) IHashHelper {
	return &hashHelper{aes: aes}
}

func (h *hashHelper) Generate(expireAt time.Time, realHash string, attributes interface{}) string {
	data := HashData{
		ExpireAt:   expireAt,
		RealHash:   realHash,
		Attributes: attributes,
	}
	jData, err := json.Marshal(data)
	if err != nil {
		return ""
	}
	if hash, err2 := h.aes.Encrypt(string(jData)); err2 == nil {
		return hash
	}
	return ""
}

func (h *hashHelper) Decode(value string) (*HashData, error) {
	jData, err := h.aes.Decrypt(value)
	if err != nil {
		return nil, err
	}
	var data HashData
	err = json.Unmarshal([]byte(jData), &data)
	return &data, err
}

func LinksAPI(w http.ResponseWriter, r *http.Request) {
	req := reqio.NewRequest(w, r).Prepare()
	log := req.GetContext().Get("logger").(logger.IManager)
	view := req.GetContext().Get("view").(template.IManager)

	log.Log("links_create_api_handler", "request received")

	var payload struct {
		Page  Page  `json:"page"`
		Limit Limit `json:"limit"`
	}

	_ = req.UnmarshallBody(&payload)

	datasource := req.GetContext().Get("db").(database.IManager)
	linkDomain := feedbacks.NewLinksDomain(datasource)

	links, total, err := linkDomain.FindLinks(req.GetContext().Value(), payload.Page.Value(), payload.Limit.Value())
	if err != nil {
		_ = view.RenderJson(w, http.StatusBadRequest, api.Response{
			Status:  http.StatusBadRequest,
			Content: make(map[string]interface{}, 0),
			Error: &api.ResponseError{
				Code:    "LNK:ERR:FND",
				Message: "Bad Request! Some problems occurred when searching the data.",
				Reasons: make(map[string]string, 0),
				Details: make([]interface{}, 0),
			},
		})
		return
	}
	_ = view.RenderJson(w, http.StatusOK, api.Response{
		Status: http.StatusOK,
		Content: map[string]interface{}{
			"total": total,
			"links": links,
		},
		Error: nil,
	})
}

func LinksCreateAPI(w http.ResponseWriter, r *http.Request) {
	req := reqio.NewRequest(w, r).Prepare()
	log := req.GetContext().Get("logger").(logger.IManager)
	view := req.GetContext().Get("view").(template.IManager)
	aes := req.GetContext().Get("aes").(crypt.ICryptAES)
	cfg := req.GetContext().Get("cfg").(config.IManager)

	log.Log("links_create_api_handler", "request received")

	var payload struct {
		Csrf  string           `json:"csrf"`
		Links []feedbacks.Link `json:"links"`

		DisableAutoGenerateHash bool `json:"disable_auto_generate_hash"`
	}
	err := req.UnmarshallBody(&payload)
	if err != nil || len(payload.Links) < 1 {
		_ = view.RenderJson(w, http.StatusBadRequest, api.Response{
			Status:  http.StatusBadRequest,
			Content: make(map[string]interface{}, 0),
			Error: &api.ResponseError{
				Code:    "LNK:ERR:BND",
				Message: "Bad Request! Something wrong with the payload of your request.",
				Reasons: make(map[string]string, 0),
				Details: make([]interface{}, 0),
			},
		})
		return
	}
	// let's do validation
	errs := make(map[string]string, 0)
	// csrf check
	sm := req.GetContext().Get("sm").(session.IManager)
	sessionCsrf := sm.Get(r.Context(), "token")
	if validate.IsEmpty(payload.Csrf) || sessionCsrf == nil || payload.Csrf != sessionCsrf.(string) {
		errs["session"] = "Not a valid request session!"
	}
	hh := NewHashHelper(aes)
	expireAt := time.Now().Add(time.Duration(cfg.GetConfig().App.HashExpire) * time.Hour)
	for li, link := range payload.Links {
		hash := link.Hash.Value()
		if !payload.DisableAutoGenerateHash {
			hash = ksuid.New().String()
		}
		if len(hash) > 0 {
			payload.Links[li].Hash = feedbacks.Hash(hh.Generate(expireAt, hash, map[string]interface{}{
				"usage_limit": link.UsageLimit,
				"group_id":    link.GroupId,
				"pin":         link.PIN,
			}))
			link.Hash = feedbacks.Hash(hash)
		}
		if !link.Hash.Valid() {
			errs[fmt.Sprintf("%d_hash", li)] = "invalid hash"
		}
		if link.GroupId < 1 {
			errs[fmt.Sprintf("%d_group_id", li)] = "invalid group"
		}
		if !link.PIN.Valid() {
			errs[fmt.Sprintf("%d_pin", li)] = "pin must be 6 character length"
		}
	}
	if len(errs) > 0 {
		_ = view.RenderJson(w, http.StatusBadRequest, api.Response{
			Status:  http.StatusBadRequest,
			Content: make(map[string]interface{}, 0),
			Error: &api.ResponseError{
				Code:    "LNK:ERR:VAL",
				Message: "Bad Request! Your request resulting validation error.",
				Reasons: errs,
				Details: make([]interface{}, 0),
			},
		})
		return
	}
	datasource := req.GetContext().Get("db").(database.IManager)
	linkDomain := feedbacks.NewLinksDomain(datasource)
	if err = linkDomain.SaveLinks(req.GetContext().Value(), payload.Links); err != nil {
		_ = view.RenderJson(w, http.StatusExpectationFailed, api.Response{
			Status:  http.StatusExpectationFailed,
			Content: make(map[string]interface{}, 0),
			Error: &api.ResponseError{
				Code:    "LNK:ERR:SAV",
				Message: "Fail to save your request. Please check your data and try again.",
				Reasons: map[string]string{
					"save_error": errors.Wrap(err, "something wrong with the execution syntax").Error(),
				},
				Details: make([]interface{}, 0),
			},
		})
		return
	}
	_ = view.RenderJson(w, http.StatusOK, api.Response{
		Status:  http.StatusOK,
		Content: make(map[string]interface{}, 0),
	})
}

func LinkUpdateAPI(w http.ResponseWriter, r *http.Request) {
	req := reqio.NewRequest(w, r).Prepare()
	log := req.GetContext().Get("logger").(logger.IManager)
	view := req.GetContext().Get("view").(template.IManager)
	aes := req.GetContext().Get("aes").(crypt.ICryptAES)
	cfg := req.GetContext().Get("cfg").(config.IManager)

	log.Log("links_update_api_handler", "request received")

	var payload struct {
		Csrf           string         `json:"csrf"`
		RegenerateHash bool           `json:"regenerate_hash"`
		Link           feedbacks.Link `json:"link"`
	}
	err := req.UnmarshallBody(&payload)
	if err != nil {
		_ = view.RenderJson(w, http.StatusBadRequest, api.Response{
			Status:  http.StatusBadRequest,
			Content: make(map[string]interface{}, 0),
			Error: &api.ResponseError{
				Code:    "LNK:ERR:BND",
				Message: "Bad Request! Something wrong with the payload of your request.",
				Reasons: make(map[string]string, 0),
				Details: make([]interface{}, 0),
			},
		})
		return
	}
	// let's do validation
	errs := make(map[string]string, 0)
	// csrf check
	sm := req.GetContext().Get("sm").(session.IManager)
	sessionCsrf := sm.Get(r.Context(), "token")

	if validate.IsEmpty(payload.Csrf) || sessionCsrf == nil || payload.Csrf != sessionCsrf.(string) {
		errs["session"] = "Not a valid request session!"
	}
	if payload.Link.Id < 1 {
		errs["id"] = "not a valid identifier"
	}
	if !payload.RegenerateHash && !payload.Link.Hash.Valid() {
		errs["hash"] = "not a valid hash code"
	}
	if payload.Link.GroupId < 1 {
		errs["group_id"] = "not a valid group"
	}
	if !payload.Link.PIN.Valid() {
		errs["pin"] = "pin must be 6 character length"
	}

	if len(errs) > 0 {
		_ = view.RenderJson(w, http.StatusBadRequest, api.Response{
			Status:  http.StatusBadRequest,
			Content: make(map[string]interface{}, 0),
			Error: &api.ResponseError{
				Code:    "LNK:ERR:VAL",
				Message: "Bad Request! Your request resulting validation error.",
				Reasons: errs,
				Details: make([]interface{}, 0),
			},
		})
		return
	}
	hh := NewHashHelper(aes)
	expireAt := time.Now().Add(time.Duration(cfg.GetConfig().App.HashExpire) * time.Hour)
	if payload.RegenerateHash {
		payload.Link.Hash = feedbacks.Hash(hh.Generate(expireAt, ksuid.New().String(), map[string]interface{}{
			"usage_limit": payload.Link.UsageLimit,
			"group_id":    payload.Link.GroupId,
			"pin":         payload.Link.PIN,
		}))
	}
	datasource := req.GetContext().Get("db").(database.IManager)
	linkDomain := feedbacks.NewLinksDomain(datasource)
	if err = linkDomain.UpdateLink(req.GetContext().Value(), payload.Link); err != nil {
		_ = view.RenderJson(w, http.StatusExpectationFailed, api.Response{
			Status:  http.StatusExpectationFailed,
			Content: make(map[string]interface{}, 0),
			Error: &api.ResponseError{
				Code:    "LNK:ERR:UPD",
				Message: "Fail to update your request. Please check your data and try again.",
				Reasons: map[string]string{"update_error": err.Error()},
				Details: make([]interface{}, 0),
			},
		})
		return
	}
	_ = view.RenderJson(w, http.StatusOK, api.Response{
		Status:  http.StatusOK,
		Content: make(map[string]interface{}, 0),
	})
}

func LinksRemoveAPI(w http.ResponseWriter, r *http.Request) {
	req := reqio.NewRequest(w, r).Prepare()
	log := req.GetContext().Get("logger").(logger.IManager)
	view := req.GetContext().Get("view").(template.IManager)

	log.Log("links_update_api_handler", "request received")

	var payload struct {
		Csrf   string `json:"csrf"`
		LinkId int64  `json:"link_id"`
	}

	err := req.UnmarshallBody(&payload)
	if err != nil {
		_ = view.RenderJson(w, http.StatusBadRequest, api.Response{
			Status:  http.StatusBadRequest,
			Content: make(map[string]interface{}, 0),
			Error: &api.ResponseError{
				Code:    "LNK:ERR:BND",
				Message: "Bad Request! Something wrong with the payload of your request.",
				Reasons: make(map[string]string, 0),
				Details: make([]interface{}, 0),
			},
		})
		return
	}
	// let's do validation
	errs := make(map[string]string, 0)
	// csrf check
	sm := req.GetContext().Get("sm").(session.IManager)
	sessionCsrf := sm.Get(r.Context(), "token")

	if validate.IsEmpty(payload.Csrf) || sessionCsrf == nil || payload.Csrf != sessionCsrf.(string) {
		errs["session"] = "Not a valid request session!"
	}
	if payload.LinkId < 1 {
		errs["id"] = "not a valid identifier"
	}
	if len(errs) > 0 {
		_ = view.RenderJson(w, http.StatusBadRequest, api.Response{
			Status:  http.StatusBadRequest,
			Content: make(map[string]interface{}, 0),
			Error: &api.ResponseError{
				Code:    "LNK:ERR:VAL",
				Message: "Bad Request! Your request resulting validation error.",
				Reasons: errs,
				Details: make([]interface{}, 0),
			},
		})
		return
	}
	datasource := req.GetContext().Get("db").(database.IManager)
	linkDomain := feedbacks.NewLinksDomain(datasource)
	if err = linkDomain.DisableLinksByIds(req.GetContext().Value(), payload.LinkId); err != nil {
		_ = view.RenderJson(w, http.StatusExpectationFailed, api.Response{
			Status:  http.StatusExpectationFailed,
			Content: make(map[string]interface{}, 0),
			Error: &api.ResponseError{
				Code:    "LNK:ERR:UPD",
				Message: "Fail to update your request. Please check your data and try again.",
				Reasons: map[string]string{"save_error": err.Error()},
				Details: make([]interface{}, 0),
			},
		})
		return
	}
	_ = view.RenderJson(w, http.StatusOK, api.Response{
		Status:  http.StatusOK,
		Content: make(map[string]interface{}, 0),
	})
}
