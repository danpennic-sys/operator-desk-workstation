package main

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestPreCommitUsesFailClosedGofmtCheck(t *testing.T) {
	hook, err := os.ReadFile(filepath.Join(repoRoot(t), ".githooks", "pre-commit"))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(hook, []byte("scripts/gofmt-check.sh")) {
		t.Fatal("pre-commit does not call scripts/gofmt-check.sh")
	}
	if bytes.Contains(hook, []byte("gofmt -l")) {
		t.Fatal("pre-commit still pipes gofmt -l instead of calling scripts/gofmt-check.sh")
	}
}

func TestGofmtCheckRejectsDirty(t *testing.T) {
	dir := t.TempDir()
	// gofmt requires a space before the brace. This file is dirty.
	body := "package p\n\nfunc F(){println(1)}\n"
	if err := os.WriteFile(filepath.Join(dir, "dirty.go"), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command(filepath.Join(repoRoot(t), "scripts", "gofmt-check.sh"), dir)
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("dirty tree reported clean: %s", out)
	}
	ee, ok := err.(*exec.ExitError)
	if !ok || ee.ExitCode() == 0 {
		t.Fatalf("expected exit 1, err=%v output=%s", err, out)
	}
	if !bytes.Contains(out, []byte("dirty.go")) {
		t.Fatalf("dirty file not reported: %s", out)
	}
}

func TestGofmtCheckAcceptsClean(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "clean.go")
	body := "package p\n\nfunc F() { println(1) }\n"
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	fmtCmd := exec.Command("gofmt", "-w", path)
	if out, err := fmtCmd.CombinedOutput(); err != nil {
		t.Fatalf("gofmt -w: %v %s", err, out)
	}
	cmd := exec.Command(filepath.Join(repoRoot(t), "scripts", "gofmt-check.sh"), dir)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("clean tree rejected: %v %s", err, out)
	}
	if !bytes.Contains(out, []byte("[gofmt] clean")) {
		t.Fatalf("missing clean marker: %s", out)
	}
}

// The old hook idiom is not fail-closed once set -o pipefail is on.
// grep -q exits on the first match, a large writer dies with SIGPIPE,
// and the inverted status takes the clean branch. A tiny listing can
// still be caught, which is why this repository hid the defect.
func TestPipefailGrepQFalseCleansLargeOutput(t *testing.T) {
	script := `set -o pipefail
if ! python3 -c 'print("f.go\n"*80000)' | grep -q .; then
  echo FALSE_CLEAN
  exit 0
fi
echo CAUGHT
exit 1
`
	cmd := exec.Command("bash", "-c", script)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("old idiom unexpectedly failed closed: %v %s", err, out)
	}
	if !bytes.Contains(out, []byte("FALSE_CLEAN")) {
		t.Fatalf("old idiom output: %s", out)
	}
}

func TestVerifyEntrypointsUseFailClosedSubcommand(t *testing.T) {
	root := repoRoot(t)
	checks := []struct {
		path string
		want string
		ban  string
	}{
		{
			path: filepath.Join(root, ".github", "actions", "eventv1-verify", "action.yml"),
			want: "eventv1-verify verify -vector",
			ban:  "eventv1-verify --vector",
		},
		{
			path: filepath.Join(root, "scripts", "run-eventv1-verify.sh"),
			want: "eventv1-verify verify -vector",
			ban:  "eventv1-verify --vector",
		},
	}
	for _, c := range checks {
		body, err := os.ReadFile(c.path)
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Contains(body, []byte(c.want)) {
			t.Fatalf("%s does not call %q", c.path, c.want)
		}
		if bytes.Contains(body, []byte(c.ban)) {
			t.Fatalf("%s still uses %q", c.path, c.ban)
		}
	}
}
