package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/irfansharif/log"
	"github.com/irfansharif/log/cmd/logger/pkg"
	"github.com/irfansharif/log/cmd/logger/pkg/subpkg"
)

// $ <binary-name> -help
// Usage of <binary-name>:
//   -log-dir string
//         Write log files in this directory.
//   -log-to-stderr
//         Log to standard error.
//   -log-mode value
//         Log mode for logs emitted globally (can be overrode using -log-filter).
//   -log-filter value
//         Comma-separated list of filename:level settings for file-filtered logging modes.
//   -log-backtrace-at value
//         Comma-separated list of filename:N settings, when any logging statement at
//         the specified locations are executed, a stack trace will be emitted.

func main() {
	var logDirFlag string
	var logToStderrFlag bool
	var logModeFlag logMode
	var logFilterFlag logFilter
	var backtracePointFlag backtracePoints

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.StringVar(&logDirFlag, "log-dir", "",
		"Write log files in this directory.")
	flag.BoolVar(&logToStderrFlag, "log-to-stderr", false,
		"Log to standard error.")
	flag.Var(&logModeFlag, "log-mode",
		"Log mode for logs emitted globally (can be overrode using -log-filter)")
	flag.Var(&logFilterFlag, "log-filter",
		"Comma-separated list of pattern:level settings for file-filtered logging modes.")
	flag.Var(&backtracePointFlag, "log-backtrace-at",
		"Comma-separated list of filename:N settings, when any logging statement at "+
			"the specified locations are executed, a stack trace will be emitted.")

	flag.Parse()

	if logModeFlag.set {
		log.SetGlobalLogMode(log.Mode(logModeFlag.m))
	}
	for _, flm := range logFilterFlag {
		log.SetFileLogMode(flm.fname, flm.fmode)
	}
	for _, tp := range backtracePointFlag {
		log.SetTracePoint(tp)
	}

	writer := ioutil.Discard
	if logDirFlag != "" {
		writer = log.LogRotationWriter(logDirFlag, 50<<20 /* 50 MiB */)
	}
	if logToStderrFlag {
		writer = log.MultiWriter(writer, os.Stderr)
	}
	writer = log.SynchronizedWriter(writer)

	logf := log.Ldate | log.Ltime | log.Lmicroseconds | log.Llongfile | log.LUTC | log.Lmode
	logger := log.New(log.Writer(writer), log.Flags(logf), log.SkipBasePath())

	logger.Debug("log-dir:", logDirFlag)
	logger.Debug("log-to-stderr:", logToStderrFlag)
	logger.Debug("log-mode:", logModeFlag.String())
	logger.Debug("log-filter:", logFilterFlag.String())
	logger.Debug("log-backtrace-at:", backtracePointFlag.String())

	logger.Info("from main!")
	pkg.Log(logger)
	subpkg.Log(logger)
}
