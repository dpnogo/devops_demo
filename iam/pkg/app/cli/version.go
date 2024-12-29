package pkg

import (
	"encoding/json"
	"fmt"
	"runtime"
	"time"
)

var (
	// GitVersion is semantic version.
	GitVersion = "v0.0.0-master+$Format:%h$"
	// BuildDate in ISO8601 format, output of $(date -u +'%Y-%m-%dT%H:%M:%SZ').
	BuildDate = "1970-01-01T00:00:00Z"
	// GitCommit sha1 from git, output of $(git rev-parse HEAD).
	GitCommit = "$Format:%H$"
	// GitTreeStateStr state of git tree, either "clean" or "dirty".
	GitTreeStateStr = ""
)

const (
	gitTreeStateDirty   = iota // 源代码在编译时处于 Git 树的脏状态
	gitTreeStateClean          // 编译时的 Git 树状态为干净状态 (无提交的更改)
	gitTreeStateUnknown        // 无法确定源代码Git状态
)

type VersionInfo struct {
	GitVersion   string `json:"git_version"`
	GitCommit    string `json:"git_commit"`
	GitTreeState string `json:"git_tree_state"`

	BuildDate string `json:"build_date"`
	GoVersion string `json:"go_version"`
	Compiler  string `json:"compiler"`
	Platform  string `json:"platform"`
}

func (info VersionInfo) ToJSON() string {
	s, _ := json.Marshal(info)

	return string(s)
}

func Get() VersionInfo {

	return VersionInfo{
		GitVersion:   GitVersion,
		GitCommit:    GitCommit,
		GitTreeState: GitTreeStateStr, // 感觉可以换成项目版本

		BuildDate: BuildDate,
		GoVersion: runtime.Version(),
		Compiler:  runtime.Compiler,
		Platform:  fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	}
}

func NewVersion(gVersion, gCommit, gts, bTime string) *VersionInfo {
	return &VersionInfo{
		GitVersion:   gVersion,
		GitCommit:    gCommit,
		GitTreeState: gts,
		BuildDate:    bTime,
		GoVersion:    runtime.Version(),
		Compiler:     runtime.Compiler,
		Platform:     fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	}
}

// BuildTime 打印此时的编译的时间
func BuildTime() string {
	return time.Now().Format(time.DateTime)
}

// GitTreeState 打印 git 树状态
// gitTreeStateDirty 源代码在编译时处于 Git 树的脏状态
// gitTreeStateClean 编译时的 Git 树状态为干净状态 (无提交的更改)
// gitTreeStateUnknown 无法确定源代码Git状态
func GitTreeState(gtsType int) string {

	switch gtsType {
	case gitTreeStateDirty:
		return "dirty"

	case gitTreeStateClean:
		return "clean"

	case gitTreeStateUnknown:
		return "unknown"

	default:
		return "unknown"
	}

}

/*
$ ./iam-apiserver --version / --version=true
  gitVersion: v0.3.0
   gitCommit: ccc31e292f66e6bad94efb1406b5ced84e64675c
   gitTreeState: dirty
   buildDate: 2020-12-17T12:24:37Z
   goVersion: go1.15.1
    compiler: gc
    platform: linux/amd64
$ ./iam-apiserver --version=raw
version.Info{GitVersion:"v0.3.0", GitCommit:"ccc31e292f66e6bad94efb1406b5ced84e64675c", GitTreeState:"dirty", BuildDate:"2020-12-17T12:24:37Z", GoVersion:"go1.15.1", Compiler:"gc", Platform:"linux/amd64"}
*/
