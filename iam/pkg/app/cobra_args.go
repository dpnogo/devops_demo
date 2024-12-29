package app

import "github.com/spf13/cobra"

// 使用 cobra 自带args的验证方法
// NoArgs：如果存在任何非选项参数，该命令将报错。
// ArbitraryArgs：该命令将接受任何非选项参数。
// OnlyValidArgs：如果有任何非选项参数不在 Command 的 ValidArgs 字段中，该命令将报错。
// MinimumNArgs(int)：如果没有至少 N 个非选项参数，该命令将报错。
// MaximumNArgs(int)：如果有多于 N 个非选项参数，该命令将报错。
// ExactArgs(int)：如果非选项参数个数不为 N，该命令将报错。
// ExactValidArgs(int)：如果非选项参数的个数不为 N，或者非选项参数不在 Command 的 ValidArgs 字段中，该命令将报错。
// RangeArgs(min, max)：如果非选项参数的个数不在 min 和 max 之间，该命令将报错。

// WithValidArgs 设置验证非选项参数方法，可使用下面cobra实现好的，或者自己进行实现
func WithValidArgs(postArgs cobra.PositionalArgs) Option {
	return func(app *App) {
		app.args = postArgs
	}
}

// CobraNoArgs 如果存在任何非选项参数，该命令将报错
func CobraNoArgs() func(cmd *cobra.Command, args []string) error {
	return cobra.NoArgs
}

// CobraArbitraryArgs 该命令将接受任何非选项参数
func CobraArbitraryArgs() func(cmd *cobra.Command, args []string) error {
	return cobra.ArbitraryArgs
}

// CobraOnlyValidArgs 如果有任何非选项参数不在 Command 的 ValidArgs 字段中，该命令将报错
func CobraOnlyValidArgs() func(cmd *cobra.Command, args []string) error {
	return cobra.OnlyValidArgs
}

// CobraMinimumNArgs 如果没有至少 N 个非选项参数，该命令将报错
func CobraMinimumNArgs(n int) func(cmd *cobra.Command, args []string) error {
	return cobra.MinimumNArgs(n)
}

// CobraMaximumNArgs 如果有多于 N 个非选项参数，该命令将报错
func CobraMaximumNArgs(n int) func(cmd *cobra.Command, args []string) error {
	return cobra.MaximumNArgs(n)
}

// CobraExactArgs 如果非选项参数的个数不为 N，或者非选项参数不在 Command 的 ValidArgs 字段中，该命令将报错
func CobraExactArgs(n int) func(cmd *cobra.Command, args []string) error {
	return cobra.ExactArgs(n)
}

// CobraRangeArgs 如果非选项参数的个数不在 min 和 max 之间，该命令将报错
func CobraRangeArgs(min, max int) func(cmd *cobra.Command, args []string) error {
	return cobra.RangeArgs(min, max)
}

// todo 要自己实现时候，要调用cobra
