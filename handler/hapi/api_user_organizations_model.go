package hapi

import (
	"github.com/evorts/feednomity/domain/users"
	"github.com/evorts/feednomity/pkg/utils"
	"time"
)

type OrganizationRequest struct {
	Id         int64      `json:"id"`
	Name       string     `json:"name"`
	Address    string     `json:"address"`
	Phone      string     `json:"phone"`
	Disabled   bool       `json:"disabled"`
}

type OrganizationResponse struct {
	Id         int64      `json:"id"`
	Name       string     `json:"name"`
	Address    string     `json:"address"`
	Phone      string     `json:"phone"`
	Disabled   bool       `json:"disabled"`
	CreatedAt  *time.Time `json:"created_at"`
	UpdatedAt  *time.Time `json:"updated_at"`
	DisabledAt *time.Time `json:"disabled_at"`
}

func transformOrganizations(f []*OrganizationRequest) (t []*users.Organization) {
	t = make([]*users.Organization, 0)
	for _, fv := range f {
		u := &users.Organization{}
		if err := utils.TransformStruct(u, fv); err == nil {
			t = append(t, u)
		}
	}
	return
}

func transformOrganizationsReverse(f []*users.Organization) (t []*OrganizationResponse) {
	t = make([]*OrganizationResponse, 0)
	for _, fv := range f {
		u := &OrganizationResponse{}
		if err := utils.TransformStruct(u, fv); err == nil {
			t = append(t, u)
		}
	}
	return
}