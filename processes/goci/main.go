package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"time"

  "os/signal"
  "syscall"
)

type executer interface {
	execute() (string, error)
}

func run(proj string, out io.Writer) error {

	if proj == "" {
		return fmt.Errorf("Project directory is required: %w", ErrValidation)
	}
	pipeline := make([]executer, 4)
	pipeline[0] = newStep(
		"go build",
		"go",
		"Go Build: SUCCESS",
		proj,
		[]string{"build", ".", "errors"},
	)
	pipeline[1] = newStep(
		"go test",
		"go",
		"Go Test: SUCCESS",
		proj,
		[]string{"test", "-v"},
	)
	pipeline[2] = newExceptionStep(
		"go fmt",
		"gofmt",
		"Gofmt: SUCCESS",
		proj,
		[]string{"-l", "."},
	)
	pipeline[3] = newTimeoutStep(
		"git push",
		"git",
		"Git Push: SUCCESS",
		proj,
		[]string{"push", "origin", "master"},
		10*time.Second,
	)
  sig := make(chan os.Signal, 1)
  errCh := make(chan error)
  done := make(chan struct{})

  signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

  go func() {
	  for _, s := range pipeline {
  		msg, err := s.execute()
  		if err != nil {
  			return err
  		}
  		_, err = fmt.Fprintln(out, msg)
  		if err != nil {
  			return err
  		}
  	}
    close(done)
  }()

	return nil
}

func main() {
	proj := flag.String("p", "", "Project directory")
	flag.Parse()

	if err := run(*proj, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
