package main

import (
	"bytes"
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
)

const (
	algEd25519Test = "ed25519-test"
	fixtureVersion = "eventv1"
	fixtureType    = "fixture"
	fixtureMessage = "operator-desk ed25519-test fixture. not quantumguard. not post-quantum."
)

// canonical is the only byte string that is signed.
// Field order is fixed. HTML is not escaped. No timestamp. No algorithm name.
func canonical(version, typ, message string) ([]byte, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	err := enc.Encode(struct {
		Message string `json:"message"`
		Type    string `json:"type"`
		Version string `json:"version"`
	}{
		Message: message,
		Type:    typ,
		Version: version,
	})
	if err != nil {
		return nil, err
	}
	return bytes.TrimSuffix(buf.Bytes(), []byte("\n")), nil
}

func zeroTestSeed() []byte {
	return bytes.Repeat([]byte{0}, ed25519.SeedSize)
}

func seedFromEnv() ([]byte, error) {
	raw := os.Getenv("PQC_TEST_SEED")
	if raw == "" {
		fmt.Fprintln(os.Stderr, "eventv1: PQC_TEST_SEED unset; generate uses the documented 32-byte zero test seed")
		return zeroTestSeed(), nil
	}
	b, err := hex.DecodeString(raw)
	if err != nil || len(b) != ed25519.SeedSize {
		return nil, fmt.Errorf("PQC_TEST_SEED must be 32-byte hex")
	}
	return b, nil
}

type fileEnvelope struct {
	Version string `json:"version"`
	Type    string `json:"type"`
	Payload struct {
		Message string `json:"message"`
	} `json:"payload"`
	Signature struct {
		Alg       string `json:"alg"`
		PublicKey string `json:"public_key"`
		Value     string `json:"value"`
	} `json:"signature"`
}

func signedFixture(seed []byte, sigHex string) ([]byte, error) {
	priv := ed25519.NewKeyFromSeed(seed)
	pub := priv.Public().(ed25519.PublicKey)
	env := fileEnvelope{
		Version: fixtureVersion,
		Type:    fixtureType,
	}
	env.Payload.Message = fixtureMessage
	env.Signature.Alg = algEd25519Test
	env.Signature.PublicKey = hex.EncodeToString(pub)
	env.Signature.Value = sigHex
	body, err := json.MarshalIndent(env, "", "  ")
	if err != nil {
		return nil, err
	}
	return append(body, '\n'), nil
}

func goodAndBad(seed []byte) (good, bad []byte, err error) {
	priv := ed25519.NewKeyFromSeed(seed)
	msg, err := canonical(fixtureVersion, fixtureType, fixtureMessage)
	if err != nil {
		return nil, nil, err
	}
	sig := hex.EncodeToString(ed25519.Sign(priv, msg))
	good, err = signedFixture(seed, sig)
	if err != nil {
		return nil, nil, err
	}
	bad, err = signedFixture(seed, flipLastNibble(sig))
	if err != nil {
		return nil, nil, err
	}
	return good, bad, nil
}

func flipLastNibble(sigHex string) string {
	b := []byte(sigHex)
	if len(b) == 0 {
		return sigHex
	}
	last := b[len(b)-1]
	var n byte
	switch {
	case last >= '0' && last <= '9':
		n = last - '0'
	case last >= 'a' && last <= 'f':
		n = last - 'a' + 10
	default:
		n = 0
	}
	n ^= 1
	if n < 10 {
		b[len(b)-1] = '0' + n
	} else {
		b[len(b)-1] = 'a' + (n - 10)
	}
	return string(b)
}

// judge returns PASS, FAIL, or INDETERMINATE.
// PASS means the embedded ed25519-test public key signed the canonical message.
// It does not identify an operator, and it is not a QuantumGuard verdict.
func judge(raw []byte) (kind, reason string) {
	if err := rejectDupAndTrailing(raw); err != nil {
		return "INDETERMINATE", "json"
	}
	var top map[string]json.RawMessage
	if err := json.Unmarshal(raw, &top); err != nil {
		return "INDETERMINATE", "json"
	}
	for _, key := range []string{"version", "type", "payload", "signature"} {
		if _, ok := top[key]; !ok {
			return "INDETERMINATE", "missing-field"
		}
	}
	if len(top) != 4 {
		return "INDETERMINATE", "unknown-field"
	}
	payloadKeys, err := objectKeys(top["payload"])
	if err != nil {
		return "INDETERMINATE", "json"
	}
	if _, ok := payloadKeys["message"]; !ok {
		return "INDETERMINATE", "missing-field"
	}
	if len(payloadKeys) != 1 {
		return "INDETERMINATE", "unknown-field"
	}
	sigKeys, err := objectKeys(top["signature"])
	if err != nil {
		return "INDETERMINATE", "json"
	}
	for _, key := range []string{"alg", "public_key", "value"} {
		if _, ok := sigKeys[key]; !ok {
			return "INDETERMINATE", "missing-field"
		}
	}
	if len(sigKeys) != 3 {
		return "INDETERMINATE", "unknown-field"
	}

	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.DisallowUnknownFields()
	var env fileEnvelope
	if err := dec.Decode(&env); err != nil {
		return "INDETERMINATE", "json"
	}
	if env.Version != fixtureVersion {
		return "INDETERMINATE", "version"
	}
	if env.Type == "" || env.Payload.Message == "" {
		return "INDETERMINATE", "empty"
	}
	if env.Signature.Alg != algEd25519Test {
		return "INDETERMINATE", "alg"
	}
	pub, errPub := hex.DecodeString(env.Signature.PublicKey)
	sig, errSig := hex.DecodeString(env.Signature.Value)
	if errPub != nil || errSig != nil || len(pub) != ed25519.PublicKeySize || len(sig) != ed25519.SignatureSize {
		return "FAIL", "sig-length"
	}
	msg, err := canonical(env.Version, env.Type, env.Payload.Message)
	if err != nil {
		return "INDETERMINATE", "json"
	}
	if !ed25519.Verify(ed25519.PublicKey(pub), msg, sig) {
		return "FAIL", "sig-mismatch"
	}
	return "PASS", ""
}

func objectKeys(raw json.RawMessage) (map[string]json.RawMessage, error) {
	var m map[string]json.RawMessage
	dec := json.NewDecoder(bytes.NewReader(raw))
	if err := dec.Decode(&m); err != nil {
		return nil, err
	}
	if dec.More() {
		return nil, fmt.Errorf("trailing")
	}
	return m, nil
}

func rejectDupAndTrailing(raw []byte) error {
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.UseNumber()
	if err := walkValue(dec); err != nil {
		return err
	}
	var extra any
	err := dec.Decode(&extra)
	if err == nil {
		return fmt.Errorf("trailing")
	}
	if errors.Is(err, io.EOF) {
		return nil
	}
	return err
}

func walkValue(dec *json.Decoder) error {
	tok, err := dec.Token()
	if err != nil {
		return err
	}
	delim, ok := tok.(json.Delim)
	if !ok {
		return nil
	}
	switch delim {
	case '{':
		seen := map[string]bool{}
		for dec.More() {
			kt, err := dec.Token()
			if err != nil {
				return err
			}
			key, ok := kt.(string)
			if !ok {
				return fmt.Errorf("bad-key")
			}
			if seen[key] {
				return fmt.Errorf("duplicate-key")
			}
			seen[key] = true
			if err := walkValue(dec); err != nil {
				return err
			}
		}
		end, err := dec.Token()
		if err != nil {
			return err
		}
		if end != json.Delim('}') {
			return fmt.Errorf("bad-object")
		}
		return nil
	case '[':
		for dec.More() {
			if err := walkValue(dec); err != nil {
				return err
			}
		}
		end, err := dec.Token()
		if err != nil {
			return err
		}
		if end != json.Delim(']') {
			return fmt.Errorf("bad-array")
		}
		return nil
	default:
		return fmt.Errorf("bad-delim")
	}
}

func summarize(kinds []string) int {
	if len(kinds) == 0 {
		return 2
	}
	fail := false
	indet := false
	for _, k := range kinds {
		switch k {
		case "FAIL":
			fail = true
		case "INDETERMINATE":
			indet = true
		case "PASS":
		default:
			indet = true
		}
	}
	if fail {
		return 1
	}
	if indet {
		return 2
	}
	return 0
}

func judgeFile(path string) (kind, reason string) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return "INDETERMINATE", "json"
	}
	return judge(raw)
}

func collect(path string) ([]string, [][3]string, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, [][3]string{{"INDETERMINATE", path, "json"}}, nil
	}
	if !info.IsDir() {
		kind, reason := judgeFile(path)
		return []string{kind}, [][3]string{{kind, path, reason}}, nil
	}
	matches, err := filepath.Glob(filepath.Join(path, "*_envelope.json"))
	if err != nil {
		return nil, nil, err
	}
	sort.Strings(matches)
	if len(matches) == 0 {
		return []string{"INDETERMINATE"}, [][3]string{{"INDETERMINATE", path, "empty"}}, nil
	}
	kinds := make([]string, 0, len(matches))
	rows := make([][3]string, 0, len(matches))
	for _, name := range matches {
		kind, reason := judgeFile(name)
		kinds = append(kinds, kind)
		rows = append(rows, [3]string{kind, name, reason})
	}
	return kinds, rows, nil
}

func formatRow(kind, path, reason string) string {
	if reason == "" {
		return kind + "\t" + path
	}
	return kind + "\t" + path + "\t" + reason
}

func verifyPaths(path string) (int, []string) {
	kinds, rows, err := collect(path)
	if err != nil {
		return 2, []string{formatRow("INDETERMINATE", path, "json")}
	}
	lines := make([]string, 0, len(rows))
	for _, row := range rows {
		lines = append(lines, formatRow(row[0], row[1], row[2]))
	}
	return summarize(kinds), lines
}
