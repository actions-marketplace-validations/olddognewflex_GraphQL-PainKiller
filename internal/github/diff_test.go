package github

import (
	"testing"
)

func TestParsePatchLines(t *testing.T) {
	tests := []struct {
		name      string
		patch     string
		wantLines map[int]int
	}{
		{
			name:      "empty patch",
			patch:     "",
			wantLines: map[int]int{},
		},
		{
			name: "new file",
			patch: `@@ -0,0 +1,3 @@
+line 1
+line 2
+line 3`,
			wantLines: map[int]int{1: 1, 2: 2, 3: 3},
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
			wantLines: map[int]int{
				10: 1, 11: 2,
				12: 3, 13: 4,
				14: 5, 15: 6,
			},
		},
		{
			name: "deletions do not increment line counter",
			patch: `@@ -5,4 +5,3 @@
 context
-removed
+replaced
 context`,
			wantLines: map[int]int{5: 1, 6: 2, 7: 3},
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
			wantLines: map[int]int{
				1: 1, 2: 2, 3: 3,
				10: 4, 11: 5, 12: 6, 13: 7,
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

			for line, wantPos := range tt.wantLines {
				gotPos, ok := got[line]
				if !ok {
					t.Errorf("missing expected line %d", line)
					continue
				}
				if gotPos != wantPos {
					t.Errorf("line %d position = %d, want %d", line, gotPos, wantPos)
				}
			}
			for line := range got {
				if _, ok := tt.wantLines[line]; !ok {
					t.Errorf("unexpected line %d", line)
				}
			}
		})
	}
}
