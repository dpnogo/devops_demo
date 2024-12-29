package app

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/moby/term"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	pkg "iam/pkg/app/cli"
	"io"
	"log"
	"os"
	"runtime"
	"strings"
)

var (
	progressMessage = color.GreenString("==>")
)

// 命令行程序：用来启动一个应用。命令行程序需要实现诸如应用描述、help、参数校验等功能。根据需要，还可以实现命令自动补全、打印命令行参数等高级功能。
// 命令行参数解析：用来在启动时指定应用程序的命令行参数，以控制应用的行为。
// 配置文件解析：用来解析不同格式的配置文件。

// Option 应用选项
type Option func(*App)

// App 外部统一构建应用
type App struct {
	name        string               // App 服务名
	basename    string               // 应用的二进制文件名
	description string               // 服务描述
	args        cobra.PositionalArgs // 非选项参数
	rootCmd     *cobra.Command       // 主节点
	command     []*Command           // 子节点
	options     CliOptions           // 命令行/配置参数   ,, 命令行选项 -->
	runFunc     RunFunc              // 回调函数

	silence   bool // 静音模式，不进行启动信息，版本和配置信息
	noVersion bool // 是否显示应用程序版本标志，默认显示
	noConfig  bool // 是否应用程序配置标志，默认显示
}

// Run 执行回调函数
func (a *App) Run() {
	if err := a.rootCmd.Execute(); err != nil {
		fmt.Printf("%v %v\n", color.RedString("Error:"), err)
		os.Exit(1)
	}
}

// NewApp 初始化应用框架
func NewApp(name, basename string, opts ...Option) *App {

	var app = &App{
		name:     name,
		basename: basename,
	}

	for _, o := range opts {
		o(app)
	}

	// 初始化cmd相关
	app.buildCommand()

	return app
}

// 初始化 cmd 相关内容
func (a *App) buildCommand() {
	cmd := &cobra.Command{
		Use:   FormatBaseName(a.basename), // 根据系统
		Short: a.name,
		Long:  a.description,
		// stop printing usage when the command errors
		SilenceUsage:  true, // 不想每次输错命令打印一大堆 usage 信息，你可以通过设置 SilenceUsage: true 来关闭掉 usage
		SilenceErrors: true, // 是否在执行命令时禁止输出错误消息 -->执行命令时不显示错误消息，将其设置为true
		Args:          a.args,
	}

	cmd.SetOut(os.Stdout)
	cmd.SetErr(os.Stderr)
	cmd.Flags().SortFlags = true // 是否对命令的标志进行排序

	// 操作子节点 cmd
	if len(a.command) > 0 {
		for _, subCmd := range a.command {
			cmd.AddCommand(subCmd.cobraCommand())
		}
		// 修改该cmd的help返回
		cmd.SetHelpCommand(helpCommand(FormatBaseName(a.basename))) // 感觉没影响 todo ??
	}

	// 主节点的运行函数
	if a.runFunc != nil {
		cmd.RunE = a.runCommand
	}

	// 将各个返回的flag合并到cmd的flag中
	var namedFlagSets pkg.NamedFlagSets
	if a.options != nil {
		namedFlagSets = a.options.Flags()
		for _, f := range namedFlagSets.FlagSets {
			cmd.Flags().AddFlagSet(f)
		}
	}

	// 添加 --version  根据需要进行修改
	// 需要 version
	if !a.noVersion {
		// 初始化一个 global的pflag并添加到 namedFlagSets 上,并 namedFlagSets.FlagSet("global")的pflag添加全局的version ，全局的version在pkg中versionFlag构建的。
		pkg.AddFlags(namedFlagSets.FlagSet("global")) //  即 名为 global 的 FlagSet 进行addFlag(version 的Flag)
	}

	// 需要配置
	if !a.noConfig {
		addConfigFlag(a.basename, namedFlagSets.FlagSet("global")) // 给 global 的 FlagSet 进行 AddFlag(config的flag) 即配置文件中 global.config 即可
	}

	// 创建 help 到  --> globalflag.AddGlobalFlags(namedFlagSets.FlagSet("global"), cmd.Name())
	namedFlagSets.FlagSet("global").BoolP("help", "h", false, fmt.Sprintf("help for %s", cmd.Name()))
	// 将 global 的 flag 添加到cmd中
	cmd.Flags().AddFlagSet(namedFlagSets.FlagSet("global")) // 整体的 , 相当于将flag分组,

	//  todo -->  给 cmd 设置 命令的使用说明和帮助函数
	addCmdTemplate(cmd, namedFlagSets)

	a.rootCmd = cmd
}

func (a *App) runCommand(cmd *cobra.Command, args []string) error {

	// 打印工作目录
	printWorkingDir()

	// 打印运行相关的版本号
	if !a.noVersion {
		pkg.PrintAndExitIfRequested()
	}

	// 打印运行相关的配置
	if !a.noConfig {
		// 获取cmd中flag参数
		if err := viper.BindPFlags(cmd.Flags()); err != nil {
			return err
		}
		fmt.Println("读取全部配置为", viper.AllSettings())
		// 将viper读取配置文件和命令行值进行保存到传入的 Options 变量中
		if err := viper.Unmarshal(a.options); err != nil {
			return err
		}
	}

	if !a.silence {
		log.Printf("%v Starting %s ...\n", progressMessage, a.name)
		if !a.noVersion {
			// 打印版本结构信息
			log.Printf("%v Version: `%s`\n", progressMessage, pkg.NewVersion("", "", "", pkg.BuildTime()).ToJSON())
		}
		if !a.noConfig {
			log.Printf("%v Config file used: `%s`\n", progressMessage, viper.ConfigFileUsed()) // 打印使用文件
		}
	}

	// todo 搁置>>>>>
	if a.options != nil {
		if err := a.applyOptionRules(); err != nil {
			return err
		}
	}

	return a.runFunc(a.basename)
}

// applyOptionRules  应用选项规则
func (a *App) applyOptionRules() error {

	if completeableOptions, ok := a.options.(CompleteableOptions); ok {
		if err := completeableOptions.Complete(); err != nil {
			return err
		}
	}

	// 验证选项是否存在错误
	if errs := a.options.Validate(); len(errs) != 0 {
		errMsg := ""
		for _, oneErr := range errs {
			errMsg += oneErr.Error()
			errMsg += "::::"
		}
		return fmt.Errorf("%s", errMsg)
	}

	if printableOptions, ok := a.options.(PrintableOptions); ok && !a.silence {
		log.Printf("%v Config: `%s`\n", progressMessage, printableOptions.String())
	}

	return nil
}

// RunFunc 定义应用程序的启动回调函数
type RunFunc func(basename string) error

// WithDescription 用于设置应用程序的描述
func WithDescription(dec string) Option {
	return func(app *App) {
		app.description = dec
	}
}

// WithDefaultValidArgs 非选项参数默认验证方式，其他可见cobra_args中
func WithDefaultValidArgs() Option {
	return func(app *App) {
		app.args = func(cmd *cobra.Command, args []string) error { // 自定义函数
			// 根据需要进行验证
			for _, arg := range args {
				if len(arg) > 0 {
					return fmt.Errorf("%q does not take any arguments", cmd.CommandPath())
				}
			}
			return nil
		}
	}
}

// WithOptions 设置应用程序所获取参数的接口
func WithOptions(ops CliOptions) Option {
	return func(app *App) {
		app.options = ops
	}
}

// WithRunFunc 设置回调函数
func WithRunFunc(run RunFunc) Option {
	return func(a *App) {
		a.runFunc = run
	}
}

func WithSilence(silence bool) Option {
	return func(app *App) {
		app.silence = silence
	}
}

func FormatBaseName(basename string) string {

	// 将执行文件置为小写并去删除可执行后缀
	if runtime.GOOS == "windows" {
		basename = strings.ToLower(basename)
		basename = strings.TrimSuffix(basename, ".exe")
	}
	return basename
}

func printWorkingDir() {
	wd, _ := os.Getwd()
	log.Printf("%v WorkingDir: %s\n", progressMessage, wd)
}

func addCmdTemplate(cmd *cobra.Command, namedFlagSets pkg.NamedFlagSets) {
	usageFmt := "Usage:\n  %s\n"
	cols, _, _ := terminalSize(cmd.OutOrStdout())
	cmd.SetUsageFunc(func(cmd *cobra.Command) error { // 自定义的函数来生成命令的使用说明（usage），而不是使用 Cobra 默认生成的用法信息
		// fmt.Fprintf(cmd.OutOrStderr(), usageFmt, cmd.UseLine())
		pkg.PrintSections(cmd.OutOrStderr(), namedFlagSets, cols)
		return nil
	})
	cmd.SetHelpFunc(func(cmd *cobra.Command, args []string) { // 自定义的帮助函数
		fmt.Fprintf(cmd.OutOrStdout(), "%s\n\n"+usageFmt, cmd.Long, cmd.UseLine())
		pkg.PrintSections(cmd.OutOrStdout(), namedFlagSets, cols)
	})
}

// TerminalSize returns the current width and height of the user's terminal. If it isn't a terminal,
// nil is returned. On error, zero values are returned for width and height.
// Usually w must be the stdout of the process. Stderr won't work.
func terminalSize(w io.Writer) (int, int, error) {
	outFd, isTerminal := term.GetFdInfo(w)
	if !isTerminal {
		return 0, 0, fmt.Errorf("given writer is no terminal")
	}
	winsize, err := term.GetWinsize(outFd)
	if err != nil {
		return 0, 0, err
	}
	return int(winsize.Width), int(winsize.Height), nil
}
