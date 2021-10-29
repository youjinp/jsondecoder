package jsondecoder

import (
	"bytes"
	"compress/gzip"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

type SmallReader struct {
	r io.Reader
	n int
}

func (sm *SmallReader) next() int {
	sm.n = (sm.n + 3) % 5
	if sm.n < 1 {
		sm.n++
	}
	return sm.n
}

func (sm *SmallReader) Read(buf []byte) (int, error) {
	return sm.r.Read(buf[:min(sm.next(), len(buf))])
}

func fixture(tb testing.TB, path string) *bytes.Reader {
	f, err := os.Open(filepath.Join("testdata", path+".json.gz"))
	check(tb, err)
	defer f.Close()
	gz, err := gzip.NewReader(f)
	check(tb, err)
	buf, err := ioutil.ReadAll(gz)
	check(tb, err)
	return bytes.NewReader(buf)
}

func check(tb testing.TB, err error) {
	if err != nil {
		tb.Helper()
		tb.Fatal(err)
	}
}

var inputs = []struct {
	path       string
	tokens     int // decoded tokens
	alltokens  int // raw tokens, includes : and ,
	whitespace int // number of whitespace chars
}{
	// from https://github.com/miloyip/nativejson-benchmark
	{"canada", 223236, 334373, 33},
	{"citm_catalog", 85035, 135990, 1227563},
	{"twitter", 29573, 55263, 167931},
	{"code", 217707, 396293, 3},

	// from https://raw.githubusercontent.com/mailru/easyjson/master/benchmark/example.json
	{"example", 710, 1297, 4246},

	// from https://github.com/ultrajson/ultrajson/blob/master/tests/sample.json
	{"sample", 5276, 8677, 518549},
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
