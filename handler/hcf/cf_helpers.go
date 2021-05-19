package hcf

import (
	"context"
	"github.com/evorts/feednomity/domain/distribution"
	"github.com/evorts/feednomity/domain/users"
	"github.com/evorts/feednomity/pkg/database"
	"github.com/pkg/errors"
)

func QueryAndValidate(
	ctx context.Context,
	ds database.IManager,
	linkHash string,
) (
	link distribution.Link,
	linkDomain distribution.ILinksManager,
	linkUsageCount int,
	dist *distribution.Distribution,
	distObject *distribution.Object,
	recipient *users.User,
	respondent *users.User,
	group *users.Group,
	user *users.User,
	errs map[string]error,
) {
	var (
		err error
		d   []*distribution.Distribution
		do  []*distribution.Object
		g   []*users.Group
		u   []*users.User
	)
	errs = make(map[string]error, 0)
	linkDomain = distribution.NewLinksDomain(ds)
	link, err = linkDomain.FindByHash(ctx, linkHash)
	if err != nil || !link.Published || link.Disabled {
		errs["SUB:ERR:ENA0"] = err
		return
	}
	distDomain := distribution.NewDistributionDomain(ds)
	do, err = distDomain.FindObjectByLinkIds(ctx, link.Id)
	if err != nil || len(do) < 1 {
		errs["SUB:ERR:DIO4"] = errors.New("Could not find respective information about distribution!")
		return
	}
	distObject = do[0]
	d, err = distDomain.FindByIds(ctx, distObject.DistributionId)
	if err != nil || len(d) < 1 {
		errs["SUB:ERR:DIS0"] = errors.New("Could not find respective information about distribution!")
		return
	}
	dist = d[0]
	usersDomain := users.NewUserDomain(ds)
	u, err = usersDomain.FindByIds(ctx, distObject.RecipientId, distObject.RespondentId)
	if err != nil || len(u) < 2 {
		errs["SUB:ERR:OBJ4"] = errors.New("Could not find respective information on objects!")
		return
	}
	recipient = u[0]
	respondent = u[1]
	g, err = usersDomain.FindGroupByIds(ctx, recipient.GroupId)
	if err != nil || len(g) < 1 {
		errs["SUB:ERR:USG4"] = errors.New("Could not find respective group of objects!")
		return
	}
	group = g[0]
	u, err = usersDomain.FindByIds(ctx, dist.CreatedBy)
	if err != nil || len(u) < 1 {
		errs["SUB:ERR:USR4"] = errors.New("Could not find respective users owner!")
		return
	}
	user = u[0]
	return
}
