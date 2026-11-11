package minify

import (
	"reflect"
	"testing"
)

func TestStrip(t *testing.T) {
	tests := []struct {
		name string
		src  string
		ext  string
		want StripResult
	}{
		{
			name: "go line and block comments",
			src:  "package main\n// hi\nfunc f() {}\n/* block */\nx := 1\n",
			ext:  "go",
			want: StripResult{
				Lines: []string{"package main", "func f() {}", "x := 1"},
				Map:   []int{1, 3, 5},
			},
		},
		{
			name: "markers inside strings kept",
			src:  `s := "// not a comment"` + "\n" + `t := "/* nope */"` + "\n",
			ext:  "go",
			want: StripResult{
				Lines: []string{`s := "// not a comment"`, `t := "/* nope */"`},
				Map:   []int{1, 2},
			},
		},
		{
			name: "go raw string kept",
			src:  "s := `// still code`\n",
			ext:  "go",
			want: StripResult{
				Lines: []string{"s := `// still code`"},
				Map:   []int{1},
			},
		},
		{
			name: "hash comments for py",
			src:  "x = 1\n# comment\ny = 2\n",
			ext:  "py",
			want: StripResult{
				Lines: []string{"x = 1", "y = 2"},
				Map:   []int{1, 3},
			},
		},
		{
			name: "hash not comment for go",
			src:  "x := \"# tag\"\n",
			ext:  "go",
			want: StripResult{
				Lines: []string{`x := "# tag"`},
				Map:   []int{1},
			},
		},
		{
			name: "blank lines dropped",
			src:  "a\n\n\nb\n",
			ext:  "go",
			want: StripResult{
				Lines: []string{"a", "b"},
				Map:   []int{1, 4},
			},
		},
		{
			name: "yaml hash comments",
			src:  "key: 1\n# skip\nother: 2\n",
			ext:  "yaml",
			want: StripResult{
				Lines: []string{"key: 1", "other: 2"},
				Map:   []int{1, 3},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Strip(tt.src, tt.ext)
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("Strip() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestFormatMap(t *testing.T) {
	if got := FormatMap(nil); got != "" {
		t.Fatalf("empty map: got %q", got)
	}
	if got := FormatMap([]int{3, 5}); got != "1>3,2>5" {
		t.Fatalf("FormatMap = %q, want 1>3,2>5", got)
	}
}
