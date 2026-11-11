package collect

import "testing"

func TestMatcher(t *testing.T) {
	m := newMatcher([]string{
		"vendor/",
		"*.min.js",
		"bin",
		"# comment rule",
		"",
		"dist/",
	})

	tests := []struct {
		name  string
		path  string
		isDir bool
		want  bool
	}{
		{"vendor dir", "vendor", true, true},
		{"under vendor root", "vendor/pkg.go", false, true},
		// Dir rules with trailing "/" do not segment-match nested file paths;
		// WalkDir skips via Match(dir, true) on the ancestor instead.
		{"nested vendor file not matched by vendor/", "pkg/vendor/x.go", false, false},
		{"nested vendor dir skipped", "pkg/vendor", true, true},
		{"min.js basename", "app.min.js", false, true},
		{"min.js nested", "src/app.min.js", false, true},
		{"bin segment dir", "bin", true, true},
		{"bin segment anywhere", "tools/bin/run.go", false, true},
		{"dist dir", "dist", true, true},
		{"under dist", "dist/out.js", false, true},
		{"keep main", "main.go", false, false},
		{"keep nested main", "src/main.go", false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := m.match(tt.path, tt.isDir); got != tt.want {
				t.Fatalf("match(%q, %v) = %v, want %v", tt.path, tt.isDir, got, tt.want)
			}
		})
	}
}
