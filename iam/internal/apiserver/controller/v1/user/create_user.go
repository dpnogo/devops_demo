package user

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"iam/pkg/api/user"
	"iam/pkg/core"
	"net/http"
	"time"
)

func (ctl *UserController) Create(c *gin.Context) {

	var (
		err   error
		uInfo = new(user.User)
	)

	err = c.ShouldBind(uInfo)
	if err != nil {
		logrus.Errorf("should bind user err:%v", err)
		core.WriteResponse(c, http.StatusBadRequest, nil, fmt.Sprintf("err:%v", err))
		return
	}

	if uInfo.CreatedAt.IsZero() {
		uInfo.CreatedAt = time.Now()
	}

	// todo 验证用户的参数
	fields := uInfo.Validate()
	if len(fields) > 0 {
		// 说明此时存在问题,例如密码长度不符合等操作
		logrus.Errorf("hash pwd err:%v", err)
		core.WriteResponse(c, http.StatusBadRequest, nil, fmt.Sprintf("用户参数存在问题:%v", fields))
		return
	}

	uInfo.Password, err = user.GenerateHashPwd(uInfo.Password)
	if err != nil {
		logrus.Errorf("hash pwd err:%v", err)
		core.WriteResponse(c, http.StatusInternalServerError, nil, "创建失败")
		return
	}

	logrus.Debugln("uInfo", uInfo)

	timeCtx, cFunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cFunc()

	err = ctl.svc.User().CreateUser(timeCtx, uInfo)
	if err != nil {
		core.WriteResponse(c, http.StatusInternalServerError, nil, fmt.Sprintf("operate db err:%v", err))
		return
	}

	core.WriteResponse(c, http.StatusOK, nil, "ok")

	return
}
