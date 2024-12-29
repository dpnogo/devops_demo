package pkg

import (
	"fmt"
	"github.com/spf13/pflag"
	"os"
	"strconv"
)

type versionValue int

// Define some const.
const (
	VersionFalse versionValue = 0 // 默认不打印
	VersionTrue  versionValue = 1 // 正常打印
	VersionRaw   versionValue = 2 // 按照每一行方式打印
)

const (
	versionRaw   = "raw"
	versionTrue  = "true"
	versionFalse = "false"
)

// 打印时候有两种 true / raw
func (v *versionValue) String() string {
	switch *v {
	case VersionRaw:
		return versionRaw
	case VersionTrue:
		return versionTrue
	default:
		return versionFalse
	}
}

// Set 根据 val 赋值给 versionValue
func (v *versionValue) Set(set string) error {
	if set == versionRaw {
		*v = VersionRaw
		return nil
	}
	vBool, err := strconv.ParseBool(set)
	if vBool {
		*v = VersionTrue
	} else {
		*v = VersionFalse
	}
	return err
}

func (v *versionValue) Type() string {

	switch *v {
	case VersionTrue:
		return versionTrue
	case VersionRaw:
		return versionRaw
	default:
		return versionFalse
	}

}

var (
	versionFlagName = "version"
)

// Version 根据命令行得到versionValue的值
func Version(vFlagName string, vVal versionValue, usage string) *versionValue {
	p := new(versionValue)
	*p = vVal
	pflag.Var(p, vFlagName, usage)               // 创建name为"version"的pflag，并添加到CommandLine中(默认*FlagSet)
	pflag.Lookup(vFlagName).NoOptDefVal = "true" // --version 和  --version=true 从pflag的CommandLine找到name为"version"的Flag结构
	return p
}

var versionFlag = Version(versionFlagName, VersionFalse, "")

// PrintAndExitIfRequested will check if the -version flag was passed
// and, if so, print the version and exit.
func PrintAndExitIfRequested() {

	if *versionFlag == VersionRaw {
		fmt.Printf("%#v\n", NewVersion("", "", GitTreeState(gitTreeStateClean), BuildTime()))
		os.Exit(0)
	} else if *versionFlag == VersionTrue {
		fmt.Printf("%s\n", NewVersion("", "", GitTreeState(gitTreeStateClean), BuildTime()))
		os.Exit(0)
	}
}
