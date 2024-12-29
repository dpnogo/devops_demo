package authorize

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ory/ladon"
	"iam/internal/authzserver/authorization"
	"iam/internal/authzserver/authorization/authorizer"
)

// controller -> authorization (授权人) -> 是否通过

type AuthorizeCtl struct {
	store authorizer.PolicyGetter // 根据 key 得到其拥有的所有权限
}

func NewAuthorizeCtl(store authorizer.PolicyGetter) AuthorizeCtl {
	return AuthorizeCtl{store: store}
}

func (ctl *AuthorizeCtl) Authorize(c *gin.Context) {
	var r ladon.Request
	if err := c.ShouldBind(&r); err != nil {
		// todo
		// core.WriteResponse(c, errors.WithCode(code.ErrBind, err.Error()), nil)

		return
	}

	auth := authorization.NewAuthorizer(authorizer.NewAuthorization(ctl.store))
	if r.Context == nil {
		r.Context = ladon.Context{}
	}

	r.Context["username"] = c.GetString("username")
	rsp := auth.Authorize(&r)

	fmt.Println("rsp", rsp)

	// core.WriteResponse(c, nil, rsp)
}
