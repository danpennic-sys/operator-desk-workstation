// Package main – EventV1 CLI scaffold for Operator Desk
// Primary language: Go 1.21–1.22
// Commands: generate | envelope | verify
//
// verify is fail-closed. This is not an EventV1 v1.1.2 implementation.
package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("eventv1 – Operator Desk CLI")
		fmt.Println("Usage: eventv1 <generate|envelope|verify> [flags]")
		os.Exit(1)
	}

	cmd := os.Args[1]
	switch cmd {
	case "generate":
		generateCmd(os.Args[2:])
	case "envelope":
		envelopeCmd(os.Args[2:])
	case "verify":
		os.Exit(runVerify(os.Args[2:], os.Stderr))
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", cmd)
		os.Exit(1)
	}
}

func generateCmd(args []string) {
	fs := flag.NewFlagSet("generate", flag.ExitOnError)
	out := fs.String("out", "vectors/eventv1/", "output directory")
	_ = fs.Parse(args)
	fmt.Printf("[generate] writing sample vectors to %s\n", *out)
	// Real implementation replaces this scaffold
}

func envelopeCmd(args []string) {
	fs := flag.NewFlagSet("envelope", flag.ExitOnError)
	_ = fs.Parse(args)
	fmt.Println("[envelope] building EventV1 envelope (scaffold)")
}

// runVerify rejects every input. Verification is not implemented.
// The process status is never 0. The reason code is E_NOT_IMPLEMENTED.
func runVerify(args []string, stderr io.Writer) int {
	fs := flag.NewFlagSet("verify", flag.ContinueOnError)
	fs.SetOutput(stderr)
	vector := fs.String("vector", "vectors/eventv1/", "vector path")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	fmt.Fprintf(stderr, "E_NOT_IMPLEMENTED\n")
	fmt.Fprintf(stderr, "[verify] %s: verification is not implemented\n", *vector)
	return 1
}
