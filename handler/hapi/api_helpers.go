package hapi

import (
	"github.com/evorts/feednomity/pkg/acl"
	"github.com/evorts/feednomity/pkg/reqio"
	"github.com/evorts/feednomity/pkg/utils"
)

type Page int

func (p Page) Value() int {
	if p < 1 {
		return 1
	}
	return int(p)
}

type Limit int

func (l Limit) Value() int {
	if l < 1 {
		return 10
	}
	return int(l)
}

func Eligible(u reqio.UserData, as acl.AccessScope, uid, gid int64) bool {
	switch as {
	case acl.AccessScopeSelf:
		if uid != u.Id {
			return false
		}
	case acl.AccessScopeGroup:
		if gid != u.GroupId {
			return false
		}
	case acl.AccessScopeOrg:
		if len(u.OrgGroupIds) < 1 {
			return false
		}
		return utils.InArray(utils.ArrayInt64(u.OrgGroupIds).ToArrayInterface(), gid)
	case acl.AccessScopeGlobal:
		return true
	default:
		return false
	}
	return true
}
