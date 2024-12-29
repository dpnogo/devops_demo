package core

import (
	"github.com/gin-gonic/gin"
)

type ErrResponse struct {
	Code      int
	Message   string
	Reference string
}

func WriteResponse(c *gin.Context, code int, err error, data interface{}) {

	//if err != nil  || {
	//
	//	coder := errors.GetCodes(code)
	//
	//	//coder := errors.ParseCoder(err)
	//	c.JSON(coder.HTTPStatus(), ErrResponse{
	//		Code:      coder.Code(),
	//		Message:   coder.Msg(),
	//		Reference: coder.Reference(),
	//	})
	//
	//	return
	//}

	c.JSON(code, data)
	return
}
