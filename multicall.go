package multicall

import (
	"fmt"
	"os"
	"path/filepath"
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

func (cmdRun RunnerFuncMap) Run(args []string) error {
	exe := strings.TrimSuffix(filepath.Base(os.Args[0]), ".exe")

	osargs := append([]string{exe}, args...)

	// arg parse using rust-style match
	// https://github.com/ylxdzsw/v2socks/blob/master/src/main.rs
	// https://github.com/alexpantyukhin/go-pattern-match
	matcher := match.Match(osargs)

	for cmd := range cmdRun {
		subcmd := cmd
		runf, _ := cmdRun[cmd]
		// log.Println(subcmd)
		matcher = matcher.
			When(
				[]interface{}{
					subcmd,
					match.ANY,
				},
				func() error {
					// log.Println(subcmd)
					return runf(osargs[1:])
				},
			).
			When(
				[]interface{}{
					match.ANY,
					subcmd,
					match.ANY,
				},
				func() error {
					// log.Println(subcmd)
					return runf(osargs[2:])
				},
			)
	}

	ok, err := matcher.Result()
	if !ok {
		usage(cmdRun)
	}

	return err.(error)
}

func usage(cmdRun RunnerFuncMap) {
	fmt.Println("please specify one of the subcommands:")
	for c := range cmdRun {
		fmt.Println("-", c)
	}
	os.Exit(1)
}
