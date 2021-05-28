package hcf

import (
	"fmt"
	"github.com/evorts/feednomity/pkg/logger"
	"github.com/evorts/feednomity/pkg/reqio"
	"github.com/evorts/feednomity/pkg/utils"
	"github.com/evorts/feednomity/pkg/view"
	"net/http"
)

func CreatePassword(w http.ResponseWriter, r *http.Request) {
	req := reqio.NewRequest(w, r).Prepare()
	log := req.GetContext().Get("logger").(logger.IManager)
	vm := req.GetContext().Get("vm").(view.ITemplateManager)

	log.Log("create_password_handler", "request received")

	ref := req.GetQueryParam("ref")
	renderData := map[string]interface{}{
		"PageTitle":   "Create Password Page",
		"RedirectUrl": utils.IIf(len(ref) > 0, ref, fmt.Sprintf("/mbr/login?user=%s",req.GetQueryParam("user"))),
		"Hash": req.GetPathLastValue(),
	}

	if err := vm.InjectData("Csrf", req.GetToken()).Render(w, http.StatusOK, "_create-password.html", renderData); err != nil {
		log.Log("create_password_handler", err.Error())
	}
}
