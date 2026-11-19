package main

import (
	"bytes"
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCommittedFixturesMatchZeroSeed(t *testing.T) {
	good, bad, err := goodAndBad(zeroTestSeed())
	if err != nil {
		t.Fatal(err)
	}
	gotGood, err := os.ReadFile("../../vectors/eventv1/good_ed25519_envelope.json")
	if err != nil {
		t.Fatal(err)
	}
	gotBad, err := os.ReadFile("../../vectors/eventv1/bad_ed25519_envelope.json")
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(gotGood, good) {
		t.Fatalf("good fixture drifted\n got %s\nwant %s", gotGood, good)
	}
	if !bytes.Equal(gotBad, bad) {
		t.Fatalf("bad fixture drifted\n got %s\nwant %s", gotBad, bad)
	}
	kind, reason := judge(gotGood)
	if kind != "PASS" || reason != "" {
		t.Fatalf("good: %s %s", kind, reason)
	}
	kind, reason = judge(gotBad)
	if kind != "FAIL" || reason != "sig-mismatch" {
		t.Fatalf("bad: %s %s", kind, reason)
	}
}

func TestPlaceholderIsIndeterminate(t *testing.T) {
	raw, err := os.ReadFile("../../vectors/eventv1/sample_envelope.json")
	if err != nil {
		t.Fatal(err)
	}
	kind, reason := judge(raw)
	if kind != "INDETERMINATE" || reason != "alg" {
		t.Fatalf("placeholder: %s %s", kind, reason)
	}
	if bytes.Contains(raw, []byte("timestamp")) {
		t.Fatal("placeholder still carries a wall-clock field")
	}
}

func TestTamperMessageFails(t *testing.T) {
	good, _, err := goodAndBad(zeroTestSeed())
	if err != nil {
		t.Fatal(err)
	}
	var env fileEnvelope
	if err := json.Unmarshal(good, &env); err != nil {
		t.Fatal(err)
	}
	env.Payload.Message = env.Payload.Message + " tampered"
	raw, err := json.Marshal(env)
	if err != nil {
		t.Fatal(err)
	}
	kind, reason := judge(raw)
	if kind != "FAIL" || reason != "sig-mismatch" {
		t.Fatalf("tamper: %s %s", kind, reason)
	}
}

func TestTimestampFieldIsIndeterminate(t *testing.T) {
	good, _, err := goodAndBad(zeroTestSeed())
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]any
	if err := json.Unmarshal(good, &m); err != nil {
		t.Fatal(err)
	}
	m["timestamp"] = "2026-01-01T00:00:00Z"
	raw, err := json.Marshal(m)
	if err != nil {
		t.Fatal(err)
	}
	kind, reason := judge(raw)
	if kind != "INDETERMINATE" || reason != "unknown-field" {
		t.Fatalf("clock field: %s %s", kind, reason)
	}
}

func TestDuplicateKeyIndeterminate(t *testing.T) {
	raw := []byte(`{"version":"eventv1","version":"eventv1","type":"fixture","payload":{"message":"x"},"signature":{"alg":"ed25519-test","public_key":"aa","value":"bb"}}`)
	kind, reason := judge(raw)
	if kind != "INDETERMINATE" || reason != "json" {
		t.Fatalf("dup: %s %s", kind, reason)
	}
}

func TestReplacedKeyWithValidSignaturePasses(t *testing.T) {
	seed := bytes.Repeat([]byte{1}, ed25519.SeedSize)
	msg, err := canonical(fixtureVersion, fixtureType, fixtureMessage)
	if err != nil {
		t.Fatal(err)
	}
	priv := ed25519.NewKeyFromSeed(seed)
	body, err := signedFixture(seed, hex.EncodeToString(ed25519.Sign(priv, msg)))
	if err != nil {
		t.Fatal(err)
	}
	kind, reason := judge(body)
	if kind != "PASS" || reason != "" {
		t.Fatalf("other key: %s %s — PASS is key agreement, not operator identity", kind, reason)
	}
}

func TestAlgCaseIsIndeterminate(t *testing.T) {
	good, _, err := goodAndBad(zeroTestSeed())
	if err != nil {
		t.Fatal(err)
	}
	raw := bytes.Replace(good, []byte(`"ed25519-test"`), []byte(`"Ed25519-test"`), 1)
	kind, reason := judge(raw)
	if kind != "INDETERMINATE" || reason != "alg" {
		t.Fatalf("alg case: %s %s", kind, reason)
	}
}

func TestShortSignatureFails(t *testing.T) {
	good, _, err := goodAndBad(zeroTestSeed())
	if err != nil {
		t.Fatal(err)
	}
	var env fileEnvelope
	if err := json.Unmarshal(good, &env); err != nil {
		t.Fatal(err)
	}
	env.Signature.Value = "abcd"
	raw, err := json.Marshal(env)
	if err != nil {
		t.Fatal(err)
	}
	kind, reason := judge(raw)
	if kind != "FAIL" || reason != "sig-length" {
		t.Fatalf("short: %s %s", kind, reason)
	}
}

func TestCanonicalHasNoClock(t *testing.T) {
	b, err := canonical(fixtureVersion, fixtureType, fixtureMessage)
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Contains(b, []byte("timestamp")) || bytes.Contains(b, []byte("T00:")) {
		t.Fatalf("canonical included a clock: %s", b)
	}
	want := `{"message":"operator-desk ed25519-test fixture. not quantumguard. not post-quantum.","type":"fixture","version":"eventv1"}`
	if string(b) != want {
		t.Fatalf("canonical = %s", b)
	}
}

func TestDirectoryFailOutranksPass(t *testing.T) {
	dir := t.TempDir()
	good, bad, err := goodAndBad(zeroTestSeed())
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "good_envelope.json"), good, 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "bad_envelope.json"), bad, 0o644); err != nil {
		t.Fatal(err)
	}
	code, lines := verifyPaths(dir)
	if code != 1 {
		t.Fatalf("exit %d lines %v", code, lines)
	}
}

func TestHTMLNotEscapedInCanonical(t *testing.T) {
	b, err := canonical("eventv1", "fixture", `a<b&c`)
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Contains(b, []byte(`\u003c`)) || bytes.Contains(b, []byte(`\u0026`)) {
		t.Fatalf("html escaped: %s", b)
	}
}

func TestUppercaseHexPasses(t *testing.T) {
	env := zeroFixture(t)
	env.Signature.PublicKey = strings.ToUpper(env.Signature.PublicKey)
	env.Signature.Value = strings.ToUpper(env.Signature.Value)
	kind, reason := judge(mustJSON(t, env))
	if kind != "PASS" || reason != "" {
		t.Fatalf("upper hex: %s %s", kind, reason)
	}
}

func TestExtraWhitespacePasses(t *testing.T) {
	raw, err := json.MarshalIndent(zeroFixture(t), "", "\t")
	if err != nil {
		t.Fatal(err)
	}
	kind, reason := judge(raw)
	if kind != "PASS" || reason != "" {
		t.Fatalf("whitespace: %s %s", kind, reason)
	}
}

func TestTrailingValueIndeterminate(t *testing.T) {
	good, _, err := goodAndBad(zeroTestSeed())
	if err != nil {
		t.Fatal(err)
	}
	kind, reason := judge(append(good, []byte(`{}`)...))
	if kind != "INDETERMINATE" || reason != "json" {
		t.Fatalf("trailing: %s %s", kind, reason)
	}
}

func TestEmptyObjectIndeterminate(t *testing.T) {
	kind, reason := judge([]byte(`{}`))
	if kind != "INDETERMINATE" || reason != "missing-field" {
		t.Fatalf("empty: %s %s", kind, reason)
	}
}

func TestPlainEd25519AlgIsIndeterminate(t *testing.T) {
	good, _, err := goodAndBad(zeroTestSeed())
	if err != nil {
		t.Fatal(err)
	}
	raw := bytes.Replace(good, []byte(`"alg": "ed25519-test"`), []byte(`"alg": "ed25519"`), 1)
	kind, reason := judge(raw)
	if kind != "INDETERMINATE" || reason != "alg" {
		t.Fatalf("plain alg: %s %s raw %s", kind, reason, raw)
	}
}

func TestSwappedPublicKeyFails(t *testing.T) {
	env := zeroFixture(t)
	other := ed25519.NewKeyFromSeed(bytes.Repeat([]byte{2}, ed25519.SeedSize)).Public().(ed25519.PublicKey)
	env.Signature.PublicKey = hex.EncodeToString(other)
	kind, reason := judge(mustJSON(t, env))
	if kind != "FAIL" || reason != "sig-mismatch" {
		t.Fatalf("swapped key: %s %s", kind, reason)
	}
}

func TestSignatureOverFileBytesFails(t *testing.T) {
	good, _, err := goodAndBad(zeroTestSeed())
	if err != nil {
		t.Fatal(err)
	}
	priv := ed25519.NewKeyFromSeed(zeroTestSeed())
	env := zeroFixture(t)
	env.Signature.Value = hex.EncodeToString(ed25519.Sign(priv, good))
	kind, reason := judge(mustJSON(t, env))
	if kind != "FAIL" || reason != "sig-mismatch" {
		t.Fatalf("file-byte sig: %s %s", kind, reason)
	}
}

func TestAngleBracketRoundTrip(t *testing.T) {
	msg := `a<b&c>`
	canon, err := canonical(fixtureVersion, fixtureType, msg)
	if err != nil {
		t.Fatal(err)
	}
	priv := ed25519.NewKeyFromSeed(zeroTestSeed())
	pub := priv.Public().(ed25519.PublicKey)
	env := fileEnvelope{Version: fixtureVersion, Type: fixtureType}
	env.Payload.Message = msg
	env.Signature.Alg = algEd25519Test
	env.Signature.PublicKey = hex.EncodeToString(pub)
	env.Signature.Value = hex.EncodeToString(ed25519.Sign(priv, canon))
	kind, reason := judge(mustJSON(t, env))
	if kind != "PASS" || reason != "" {
		t.Fatalf("bracket: %s %s canon %s", kind, reason, canon)
	}
}

func TestLongMessagePasses(t *testing.T) {
	msg := strings.Repeat("x", 70000)
	canon, err := canonical(fixtureVersion, fixtureType, msg)
	if err != nil {
		t.Fatal(err)
	}
	priv := ed25519.NewKeyFromSeed(zeroTestSeed())
	pub := priv.Public().(ed25519.PublicKey)
	env := fileEnvelope{Version: fixtureVersion, Type: fixtureType}
	env.Payload.Message = msg
	env.Signature.Alg = algEd25519Test
	env.Signature.PublicKey = hex.EncodeToString(pub)
	env.Signature.Value = hex.EncodeToString(ed25519.Sign(priv, canon))
	kind, reason := judge(mustJSON(t, env))
	if kind != "PASS" || reason != "" {
		t.Fatalf("long: %s %s", kind, reason)
	}
}

func zeroFixture(t *testing.T) fileEnvelope {
	t.Helper()
	good, _, err := goodAndBad(zeroTestSeed())
	if err != nil {
		t.Fatal(err)
	}
	var env fileEnvelope
	if err := json.Unmarshal(good, &env); err != nil {
		t.Fatal(err)
	}
	return env
}

func mustJSON(t *testing.T, env fileEnvelope) []byte {
	t.Helper()
	raw, err := json.Marshal(env)
	if err != nil {
		t.Fatal(err)
	}
	return raw
}
