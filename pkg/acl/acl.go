package acl

import (
	"context"
	"fmt"
	"github.com/evorts/feednomity/domain/users"
	"github.com/evorts/feednomity/pkg/utils"
	"regexp"
	"strings"
)

type AccessScope string

const (
	AccessScopeNone   AccessScope = "none"
	AccessScopeSelf   AccessScope = "self"
	AccessScopeGroup  AccessScope = "group"
	AccessScopeOrg    AccessScope = "org"
	AccessScopeGlobal AccessScope = "global"
)

type IManager interface {
	Populate() error
	IsAllowed(userId int64, method, path string) (allowed bool, scope AccessScope)
}

type access struct {
	Role             string      `json:"role"`
	Path             string      `json:"path"`
	Regex            bool        `json:"regex"`
	MethodAllowed    []string    `json:"method_allowed"`
	MethodDisallowed []string    `json:"method_disallowed"`
	AccessScope      AccessScope `json:"access_scope"`
	Disabled         bool        `json:"disabled"`
	AccessLevel      string      `json:"access_level"`
}

type accessControl struct {
	UserId         int64    `json:"user_id"`
	Username       string   `json:"username"`
	Email          string   `json:"email"`
	OrganizationId int64    `json:"organization_id"`
	GroupId        int64    `json:"group_id"`
	Role           string   `json:"role"`
	RoleAccess     []access `json:"role_access"`
	UserAccess     []access `json:"user_access"`
}

type manager struct {
	u  users.IUsers
	ua users.IUserAccess
	ac map[int64]*accessControl // user_id => accessControl
}

func NewACLManager(u users.IUsers, ua users.IUserAccess) IManager {
	return &manager{u: u, ua: ua}
}

func (m *manager) IsAllowed(userId int64, method, path string) (allowed bool, scope AccessScope) {
	method = strings.ToLower(method)
	allowed = false
	scope = ""
	defer fmt.Println(userId, method, path, allowed, scope)
	v, ok := m.ac[userId]
	if !ok {
		return false, ""
	}
	// check role access
	for _, ra := range v.RoleAccess {
		if !ra.Regex {
			if ra.Path != path {
				continue
			}
		} else {
			match, err := regexp.MatchString(ra.Path, path)
			if err != nil || !match {
				continue
			}
		}
		if utils.InArray(utils.ArrayString(ra.MethodDisallowed).ToArrayInterface(), method) {
			return false, ""
		}
		if !utils.InArray(utils.ArrayString(ra.MethodAllowed).ToArrayInterface(), method) {
			return false, ""
		}
		scope = ra.AccessScope
		allowed = true
		break
	}
	//check user access
	for _, ua := range v.UserAccess {
		if !ua.Regex {
			if ua.Path != path {
				continue
			}
		} else {
			match, err := regexp.MatchString(ua.Path, path)
			if err != nil || !match {
				continue
			}
		}
		if utils.InArray(utils.ArrayString(ua.MethodDisallowed).ToArrayInterface(), method) {
			return false, ""
		}
		if !utils.InArray(utils.ArrayString(ua.MethodAllowed).ToArrayInterface(), method) {
			return false, ""
		}
		scope = ua.AccessScope
	}
	if len(scope) < 1 {
		scope = AccessScopeNone
	}
	return true, scope
}

func (m *manager) Populate() error {
	//read groups data
	g, gErr := m.recursiveFindGroups(context.TODO(), 1, 10)
	if gErr != nil {
		return gErr
	}
	//read users data
	u, uErr := m.recursiveFindUsers(context.TODO(), g, 1, 10)
	if uErr != nil {
		return uErr
	}
	//read user role
	ur, urErr := m.recursiveFindRoleAccess(context.TODO(), 1, 10)
	if urErr != nil {
		return urErr
	}
	//read user access
	ua, uaErr := m.recursiveFindUserAccess(context.TODO(), 1, 10)
	if uaErr != nil {
		return uaErr
	}
	//populate to access control
	if m.ac == nil {
		m.ac = make(map[int64]*accessControl, 0)
	}
	for _, uv := range u {
		if uv.Disabled {
			continue
		}
		m.ac[uv.Id] = &accessControl{
			UserId:         uv.Id,
			Username:       uv.Username,
			Email:          uv.Email,
			GroupId:        uv.GroupId,
			OrganizationId: uv.OrganizationId,
			Role:           uv.AccessRole.String(),
			RoleAccess:     make([]access, 0),
			UserAccess:     make([]access, 0),
		}
		for _, rav := range ur {
			if uv.AccessRole != rav.Role || rav.Disabled {
				continue
			}
			m.ac[uv.Id].RoleAccess = append(m.ac[uv.Id].RoleAccess, access{
				Role:             rav.Role.String(),
				Path:             rav.Path,
				Regex:            rav.Regex,
				MethodAllowed:    users.AccessMethodsToStringArray(rav.AccessAllowed),
				MethodDisallowed: users.AccessMethodsToStringArray(rav.AccessDisallowed),
				AccessScope:      AccessScope(rav.AccessScope),
			})
		}
		for _, uav := range ua {
			if uav.Disabled || uav.Id != uv.Id {
				continue
			}
			m.ac[uv.Id].UserAccess = append(m.ac[uv.Id].UserAccess, access{
				Role:             uv.AccessRole.String(),
				Path:             uav.Path,
				Regex:            uav.Regex,
				MethodAllowed:    users.AccessMethodsToStringArray(uav.AccessAllowed),
				MethodDisallowed: users.AccessMethodsToStringArray(uav.AccessDisallowed),
				AccessScope:      AccessScopeSelf,
			})
		}
	}
	return nil
}

func (m *manager) recursiveFindGroups(ctx context.Context, page, limit int) (map[int64]*users.Group, error) {
	rs := make(map[int64]*users.Group, 0)
	g, gt, gErr := m.u.FindAllGroups(ctx, page, limit)
	if gErr != nil {
		return nil, gErr
	}
	for _, v := range g {
		rs[v.Id] = v
	}
	if (page-1)*limit > gt {
		return rs, nil
	}
	gg, err := m.recursiveFindGroups(ctx, page+1, limit)
	if err != nil {
		return rs, err
	}
	for k, v := range gg {
		rs[k] = v
	}
	return rs, nil
}

func (m *manager) recursiveFindUsers(ctx context.Context, g map[int64]*users.Group, page, limit int) ([]*users.User, error) {
	u, ut, uErr := m.u.FindAll(ctx, page, limit)
	if uErr != nil {
		return nil, uErr
	}
	for k, v := range u {
		if vv, ok := g[v.GroupId]; ok {
			u[k].OrganizationId = vv.OrgId
		}
	}
	if (page-1)*limit > ut {
		return u, nil
	}
	uu, err := m.recursiveFindUsers(ctx, g, page+1, limit)
	if err != nil {
		return u, err
	}
	return append(u, uu...), nil
}

func (m *manager) recursiveFindRoleAccess(ctx context.Context, page, limit int) ([]*users.UserRoleAccess, error) {
	ur, uRoleTotal, uRoleErr := m.ua.FindAllRoleAccess(ctx, page, limit)
	if uRoleErr != nil {
		return nil, uRoleErr
	}
	if (page-1)*limit > uRoleTotal {
		return ur, nil
	}
	uu, err := m.recursiveFindRoleAccess(ctx, page+1, limit)
	if err != nil {
		return ur, err
	}
	return append(ur, uu...), nil
}

func (m *manager) recursiveFindUserAccess(ctx context.Context, page, limit int) ([]*users.UserAccess, error) {
	ua, uat, uaErr := m.ua.FindAllUserAccess(ctx, page, limit)
	if uaErr != nil {
		return nil, uaErr
	}
	if (page-1)*limit > uat {
		return ua, nil
	}
	uu, err := m.recursiveFindUserAccess(ctx, page+1, limit)
	if err != nil {
		return ua, err
	}
	return append(ua, uu...), nil
}
