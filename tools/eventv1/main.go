// Command eventv1 is the Operator Desk fixture tool.
// verify prints one of PASS, FAIL, INDETERMINATE.
// It is not qg verify. The algorithm is ed25519-test, not post-quantum.
package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func main() {
	args := os.Args[1:]
	cmd := "verify"
	if len(args) > 0 && args[0] != "" && args[0][0] != '-' {
		cmd = args[0]
		args = args[1:]
	}
	switch cmd {
	case "verify":
		os.Exit(verifyCmd(args))
	case "generate":
		os.Exit(generateCmd(args))
	case "envelope":
		fmt.Fprintln(os.Stderr, "eventv1: envelope is not a verifier and is not implemented")
		os.Exit(2)
	default:
		fmt.Fprintf(os.Stderr, "eventv1: unknown command %q\n", cmd)
		os.Exit(2)
	}
}

func verifyCmd(args []string) int {
	fs := flag.NewFlagSet("verify", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	vector := fs.String("vector", "", "envelope file or directory")
	if err := fs.Parse(args); err != nil || *vector == "" {
		fmt.Fprintln(os.Stderr, "usage: eventv1 verify --vector <file|dir>")
		return 2
	}
	fmt.Fprintln(os.Stderr, "eventv1: ed25519-test desk verdict. not qg. not post-quantum.")
	code, lines := verifyPaths(*vector)
	for _, line := range lines {
		fmt.Println(line)
	}
	return code
}

func generateCmd(args []string) int {
	fs := flag.NewFlagSet("generate", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	out := fs.String("out", "vectors/eventv1", "output directory")
	if err := fs.Parse(args); err != nil {
		fmt.Fprintln(os.Stderr, "usage: eventv1 generate --out <dir>")
		return 2
	}
	seed, err := seedFromEnv()
	if err != nil {
		fmt.Fprintln(os.Stderr, "eventv1:", err)
		return 2
	}
	good, bad, err := goodAndBad(seed)
	if err != nil {
		fmt.Fprintln(os.Stderr, "eventv1:", err)
		return 2
	}
	if err := os.MkdirAll(*out, 0o755); err != nil {
		fmt.Fprintln(os.Stderr, "eventv1:", err)
		return 2
	}
	goodPath := filepath.Join(*out, "good_ed25519_envelope.json")
	badPath := filepath.Join(*out, "bad_ed25519_envelope.json")
	if err := os.WriteFile(goodPath, good, 0o644); err != nil {
		fmt.Fprintln(os.Stderr, "eventv1:", err)
		return 2
	}
	if err := os.WriteFile(badPath, bad, 0o644); err != nil {
		fmt.Fprintln(os.Stderr, "eventv1:", err)
		return 2
	}
	fmt.Printf("wrote %s\n", goodPath)
	fmt.Printf("wrote %s\n", badPath)
	fmt.Fprintln(os.Stderr, "eventv1: ed25519-test fixtures. not post-quantum. not qg verify.")
	return 0
}
