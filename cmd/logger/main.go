package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/irfansharif/log"
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

	log.SetGlobalLogMode(log.Mode(logModeFlag))
	for _, flm := range logFilterFlag {
		log.SetFileLogMode(flm.fname, flm.fmode)
	}
	for _, tp := range backtracePointFlag {
		log.EnableTracePoint(tp)
	}

	var writer io.Writer
	// writer = ioutil.Discard
	// writer = log.SynchronizedWriter(writer)
	if logDirFlag != "" {
		writer = log.RotatingDirectoryLogger(logDirFlag)
	}
	if logToStderrFlag {
		writer = log.MultiWriter(writer, os.Stderr)
	}

	logger := log.New(writer)
	logger.Info("log-dir:", logDirFlag)
	logger.Info("log-to-stderr:", logToStderrFlag)
	logger.Info("log-mode:", logModeFlag.String())
	logger.Info("log-filter:", logFilterFlag.String())
	logger.Info("log-backtrace-at:", backtracePointFlag.String())
}
