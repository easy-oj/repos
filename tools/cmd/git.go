package cmd

import "github.com/easy-oj/common/settings"

func Git(args ...string) *cmd {
	return New(settings.Path.Git, args...)
}
