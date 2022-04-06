package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

const corpus = "./testcorpus"

// TestConvertCorpus goes through all the files in ./testcorpus and runs the test for them.
func TestConvertCorpus(t *testing.T) {
	entries, err := os.ReadDir(corpus)
	if err != nil {
		t.Fatalf("Error reading %q: %v", corpus, err)
	}

	for _, e := range entries {
		if filepath.Ext(e.Name()) != ".md" {
			continue
		}

		name := e.Name()

		t.Run(name, func(t *testing.T) {
			f, err := os.Open(filepath.Join(corpus, name))
			if err != nil {
				t.Errorf("Error reading file %+q: %v", name, err)
				return
			}
			defer f.Close()
			output := &bytes.Buffer{}

			if err := convert(f, output); err != nil {
				t.Errorf("Error converting: %v", err)
				return
			}

			expectedFileName := filepath.Join(corpus, strings.Replace(name, ".md", ".html", -1))
			expected, err := os.ReadFile(expectedFileName)
			if err != nil {
				t.Errorf("Unable to read expected output %+q: %v", expectedFileName, err)
			}

			if diff := cmp.Diff(string(expected), output.String()); diff != "" {
				t.Errorf("convert(%v) mismatch (-want +got):\n%s", name, diff)
			}
		})
	}
}
