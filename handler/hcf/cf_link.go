package hcf

import (
	"fmt"
	"github.com/evorts/feednomity/domain/distribution"
	"github.com/evorts/feednomity/domain/users"
	"github.com/evorts/feednomity/pkg/database"
	"github.com/evorts/feednomity/pkg/logger"
	"github.com/evorts/feednomity/pkg/memory"
	"github.com/evorts/feednomity/pkg/reqio"
	"github.com/evorts/feednomity/pkg/view"
	"net/http"
	"path"
	"time"
)

func Link(w http.ResponseWriter, r *http.Request) {
	req := reqio.NewRequest(w, r).Prepare()
	log := req.GetContext().Get("logger").(logger.IManager)
	vm := req.GetContext().Get("view").(view.ITemplateManager)
	mem := req.GetContext().Get("mem").(memory.IManager)
	ds := req.GetContext().Get("db").(database.IManager)

	log.Log("member_link_handler", "request received")

	if len(req.GetPath()) < 1 {
		_ = vm.Render(w, http.StatusBadRequest, "404.html", map[string]interface{}{
			"PageTitle": "Page Not Found",
		})
		return
	}
	linkHash := path.Base(req.GetPath())
	if len(linkHash) < 1 {
		log.Log(
			"cf_link_hash_invalid",
			fmt.Sprintf("link hash are invalid. url path: %s", req.GetPath()),
		)
		_ = vm.Render(w, http.StatusBadRequest, "404.html", map[string]interface{}{
			"PageTitle": "Page Not Found",
		})
		return
	}
	linkDomain := distribution.NewLinksDomain(ds)
	link, err := linkDomain.FindByHash(req.GetContext().Value(), linkHash)
	if err != nil || link.Id < 1 || link.Disabled {
		log.Log(
			"cf_link_query_error",
			fmt.Sprintf("link error or disabled. error: %v", err),
		)
		_ = vm.Render(w, http.StatusBadRequest, "404.html", map[string]interface{}{
			"PageTitle": "Page Not Found",
		})
		return
	}

	var (
		objects []*distribution.Object
		obj     *distribution.Object
	)

	distDomain := distribution.NewDistributionDomain(ds)
	objects, err = distDomain.FindObjectByLinkIds(req.GetContext().Value(), link.Id)

	if err != nil || len(objects) < 1 {
		log.Log("cf_link_objects_error", fmt.Sprintf("link objects not found. error: %v", err))
		_ = vm.Render(w, http.StatusBadRequest, "404.html", map[string]interface{}{
			"PageTitle": "Page Not Found",
		})
		return
	}

	obj = objects[0]

	if req.IsLoggedIn() {
		http.Redirect(w, r, fmt.Sprintf("/mbr/review/form/%d", obj.Id), http.StatusTemporaryRedirect)
		return
	}

	now := time.Now()
	if link.ExpiredAt != nil && now.After(*link.ExpiredAt) {
		log.Log("cf_link_expired", fmt.Sprintf("link hash %s has expired", linkHash))
		_ = vm.Render(w, http.StatusBadRequest, "404.html", map[string]interface{}{
			"PageTitle": "Page Not Found",
		})
		return
	}

	var usersData []*users.User
	usersDomain := users.NewUserDomain(ds)
	usersData, err = usersDomain.FindByIds(req.GetContext().Value(), obj.RespondentId)
	if err != nil || len(usersData) < 1 {
		log.Log("cf_link_users_error", fmt.Sprintf("link users respondent not found. error: %v", err))
		_ = vm.Render(w, http.StatusBadRequest, "404.html", map[string]interface{}{
			"PageTitle": "Page Not Found",
		})
		return
	}
	_ = linkDomain.RecordLinkVisitor(req.GetContext().Value(), link, obj.RespondentId, usersData[0].Username, req.GetUserAgent(), map[string]interface{}{
		"ref": req.GetUrl().String(),
	})
	if usersData[0].Disabled {
		log.Log("cf_link_users_disabled", fmt.Sprintf("link users respondent are disabled. id: %v", usersData[0].Id))
		_ = vm.Render(w, http.StatusBadRequest, "404.html", map[string]interface{}{
			"PageTitle": "Page Not Found",
		})
		return
	}
	//render login page if necessary
	if len(usersData[0].Password) > 0 {
		http.Redirect(
			w, r,
			fmt.Sprintf("/mbr/login?user=%s", usersData[0].Email),
			http.StatusTemporaryRedirect,
		)
		return
	}
	// @todo: change below logic to api call, for now the fastest way approach
	hash := fmt.Sprintf("fp_%s", linkHash)
	_ = mem.Set(req.GetContext().Value(), hash, usersData[0].Id, 5*60)
	http.Redirect(w, r, fmt.Sprintf("/crp/%s?user=%s", linkHash, usersData[0].Email), http.StatusTemporaryRedirect)
	return
}
