package web

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestDockerfileCopiesEveryGoPackage guards the image build stage. The builder
// stage compiles from the copied tree only, so a package that is left out of
// the COPY lines makes `go build` fail inside the image even though the local
// build passes.
func TestDockerfileCopiesEveryGoPackage(t *testing.T) {
	root := filepath.Join("..", "..")
	raw, err := os.ReadFile(filepath.Join(root, "Dockerfile"))
	if err != nil {
		t.Fatalf("read Dockerfile: %v", err)
	}

	var dests []string
	for _, line := range strings.Split(string(raw), "\n") {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "COPY ") || strings.Contains(line, "--from=") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 3 {
			continue
		}
		dests = append(dests, filepath.ToSlash(strings.TrimPrefix(fields[len(fields)-1], "./")))
	}

	covered := func(rel string) bool {
		rel = filepath.ToSlash(rel)
		for _, dest := range dests {
			if dest == "" || dest == "." {
				// COPY into the workdir covers whatever is listed on the left.
				continue
			}
			if rel == dest || strings.HasPrefix(rel, dest+"/") {
				return true
			}
		}
		return false
	}

	for _, dir := range []string{"cmd", "core", "internal"} {
		base := filepath.Join(root, dir)
		err := filepath.Walk(base, func(path string, info os.FileInfo, walkErr error) error {
			if walkErr != nil || !info.IsDir() {
				return walkErr
			}
			entries, err := os.ReadDir(path)
			if err != nil {
				return err
			}
			hasGo := false
			for _, entry := range entries {
				name := entry.Name()
				if entry.IsDir() || !strings.HasSuffix(name, ".go") || strings.HasSuffix(name, "_test.go") {
					continue
				}
				hasGo = true
				break
			}
			if !hasGo {
				return nil
			}
			rel, err := filepath.Rel(root, path)
			if err != nil {
				return err
			}
			if !covered(rel) {
				t.Errorf("Dockerfile does not copy %s, so the builder stage cannot compile it", filepath.ToSlash(rel))
			}
			return nil
		})
		if err != nil {
			t.Fatalf("walk %s: %v", dir, err)
		}
	}
}
