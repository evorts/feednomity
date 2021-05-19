package hcf

import (
	"github.com/evorts/feednomity/pkg/logger"
	"github.com/evorts/feednomity/pkg/reqio"
	"github.com/evorts/feednomity/pkg/view"
	"net/http"
)

func ReviewListing(w http.ResponseWriter, r *http.Request) {
	req := reqio.NewRequest(w, r).Prepare()
	log := req.GetContext().Get("logger").(logger.IManager)
	vm := req.GetContext().Get("view").(view.ITemplateManager)

	log.Log("member_review_listing_handler", "request received")

	renderData := map[string]interface{}{
		"PageTitle":   "Member Review Listing",
	}

	if err := vm.Render(w, http.StatusOK, "member-review-list.html", renderData); err != nil {
		log.Log("member_review_listing_handler", err.Error())
	}
}
