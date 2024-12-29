package errors

import (
	"fmt"
	"sync"
)

var (
	rw    = &sync.RWMutex{}
	codes = map[int]Coder{}
)

type defaultCoder struct {
	// C refers to the integer code of the ErrCode.
	C int

	// HTTP status that should be used for the associated error code.
	HTTP int

	// External (user) facing error text.
	Ext string

	// Ref specify the reference document.
	Ref string
}

var dCoder = defaultCoder{0, 0, "未知code", ""}

func (dc defaultCoder) HTTPStatus() int {
	return dc.HTTP
}

func (dc defaultCoder) Code() int {
	return dc.C
}

func (dc defaultCoder) Msg() string {
	return dc.Ext
}

func (dc defaultCoder) Reference() string {
	return dc.Ref
}

type Coder interface {
	HTTPStatus() int
	Code() int
	Msg() string
	Reference() string
}

// Register 注册 coder，若存在则进行覆盖
func Register(coder Coder) {
	rw.Lock()
	defer rw.Unlock()
	codes[coder.Code()] = coder
}

// MustRegister  判断是否存在，若存在则进行panic，不进行覆盖
func MustRegister(coder Coder) {
	// 判断是否存在，若存在则进行panic
	if _, ok := codes[coder.Code()]; ok {
		panic(fmt.Sprintf("coder code:%d exist", coder.Code()))
	}

	// 加锁
	rw.Lock()
	defer rw.Unlock()

	codes[coder.Code()] = coder
}

// GetCodes 根据 code id 得到对应的信息
func GetCodes(code int) Coder {
	rw.RLock()
	defer rw.RUnlock()

	if coder, ok := codes[code]; ok {
		return coder
	}

	return dCoder
}

func ParseCoder(err error) Coder {
	//if err == nil {
	//	return nil
	//}
	//
	//if v, ok := err.(*withCode); ok {
	//	if coder, ok := codes[v.code]; ok {
	//		return coder
	//	}
	//}
	//
	//return unknownCoder

	return nil
}
