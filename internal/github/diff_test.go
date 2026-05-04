package github

import (
	"testing"
)

func TestParsePatchLines(t *testing.T) {
	tests := []struct {
		name      string
		patch     string
		wantLines map[int]bool
	}{
		{
			name:      "empty patch",
			patch:     "",
			wantLines: map[int]bool{},
		},
		{
			name: "new file",
			patch: `@@ -0,0 +1,3 @@
+line 1
+line 2
+line 3`,
			wantLines: map[int]bool{1: true, 2: true, 3: true},
		},
		{
			name: "context and additions",
			patch: `@@ -10,4 +10,6 @@ some function
 context line
 context line
+added line 1
+added line 2
 context line
 context line`,
			wantLines: map[int]bool{
				10: true, 11: true,
				12: true, 13: true,
				14: true, 15: true,
			},
		},
		{
			name: "deletions do not increment line counter",
			patch: `@@ -5,4 +5,3 @@
 context
-removed
+replaced
 context`,
			wantLines: map[int]bool{5: true, 6: true, 7: true},
		},
		{
			name: "multiple hunks",
			patch: `@@ -1,3 +1,3 @@
 line 1
-old line 2
+new line 2
 line 3
@@ -10,3 +10,4 @@
 line 10
+inserted
 line 11
 line 12`,
			wantLines: map[int]bool{
				1: true, 2: true, 3: true,
				10: true, 11: true, 12: true, 13: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParsePatchLines(tt.patch)

			if len(got) != len(tt.wantLines) {
				t.Fatalf("ParsePatchLines() returned %d lines, want %d\ngot:  %v\nwant: %v",
					len(got), len(tt.wantLines), got, tt.wantLines)
			}

			for line := range tt.wantLines {
				if !got[line] {
					t.Errorf("missing expected line %d", line)
				}
			}
			for line := range got {
				if !tt.wantLines[line] {
					t.Errorf("unexpected line %d", line)
				}
			}
		})
	}
}
