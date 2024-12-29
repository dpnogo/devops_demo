package errors

// stack represents a stack of program counters.
type stack []uintptr

type withCode struct {
	err   error
	code  int
	cause error
	*stack
}

// 支持，通过 withCode 的方式构建成 err ,然后通过code得到coder
