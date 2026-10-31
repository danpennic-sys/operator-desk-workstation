package main

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestVerifySampleEnvelopeFailsClosed(t *testing.T) {
	vector := filepath.Join(repoRoot(t), "vectors", "eventv1", "sample_envelope.json")
	if _, err := os.Stat(vector); err != nil {
		t.Fatalf("sample envelope missing: %v", err)
	}
	var stderr bytes.Buffer
	code := runVerify([]string{"-vector", vector}, &stderr)
	if code == 0 {
		t.Fatal("verify returned 0")
	}
	out := stderr.String()
	if !strings.Contains(out, "E_NOT_IMPLEMENTED") {
		t.Fatalf("stderr missing E_NOT_IMPLEMENTED: %q", out)
	}
	if strings.Contains(out, "always pass") || strings.Contains(strings.ToUpper(out), "PASS") {
		t.Fatalf("stderr claims success: %q", out)
	}
}

func TestVerifyProcessExitNonZero(t *testing.T) {
	if os.Getenv("EVENTV1_CHILD") == "1" {
		os.Exit(runVerify([]string{"-vector", os.Getenv("EVENTV1_VECTOR")}, os.Stderr))
	}
	vector := filepath.Join(repoRoot(t), "vectors", "eventv1", "sample_envelope.json")
	cmd := exec.Command(os.Args[0], "-test.run=^TestVerifyProcessExitNonZero$")
	cmd.Env = append(os.Environ(), "EVENTV1_CHILD=1", "EVENTV1_VECTOR="+vector)
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("process exited 0; output: %s", out)
	}
	ee, ok := err.(*exec.ExitError)
	if !ok || ee.ExitCode() == 0 {
		t.Fatalf("expected nonzero exit, err=%v output=%s", err, out)
	}
	if !bytes.Contains(out, []byte("E_NOT_IMPLEMENTED")) {
		t.Fatalf("output missing E_NOT_IMPLEMENTED: %s", out)
	}
}

func repoRoot(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	return filepath.Clean(filepath.Join(filepath.Dir(file), "..", ".."))
}
