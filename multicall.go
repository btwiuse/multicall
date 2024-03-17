package multicall

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	match "github.com/alexpantyukhin/go-pattern-match"
)

type RunnerFunc func([]string) error

func (r RunnerFunc) Run(x []string) error {
	return r(x)
}

type Runner interface {
	Run([]string) error
}

type RunnerFuncMap map[string]RunnerFunc

func keys(m RunnerFuncMap) (keys []string) {
	for k := range m {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	return
}

func expand(s string) (result []any) {
	parts := strings.Split(s, "/")
	for _, p := range parts {
		result = append(result, p)
	}
	return
}

func (cmdRun RunnerFuncMap) Run(args []string) error {
	exe := strings.TrimSuffix(filepath.Base(os.Args[0]), ".exe")

	osargs := append([]string{exe}, args...)

	// arg parse using rust-style match
	// https://github.com/ylxdzsw/v2socks/blob/master/src/main.rs
	// https://github.com/alexpantyukhin/go-pattern-match
	matcher := match.Match(osargs)

	for cmd, runf := range cmdRun {
		matcher = matcher.
			When(
				append(append([]any{}, expand(cmd)...), match.ANY),
				func() error {
					return runf(osargs[len(expand(cmd)):])
				},
			).
			When(
				append(append(append([]any{}, match.ANY), expand(cmd)...), match.ANY),
				func() error {
					return runf(osargs[len(expand(cmd))+1:])
				},
			)
	}

	ok, err := matcher.Result()
	if !ok {
		usage(cmdRun)
	}

	if err != nil {
		if len(args) == 0 {
			return fmt.Errorf("multicall: %w", err)
		}
		return fmt.Errorf("multicall: %s: %w", args[0], err)
	}

	return nil
}

func usage(cmdRun RunnerFuncMap) {
	fmt.Println("please specify one of the subcommands:")
	for _, c := range keys(cmdRun) {
		fmt.Println("-", c)
	}
	os.Exit(1)
}
