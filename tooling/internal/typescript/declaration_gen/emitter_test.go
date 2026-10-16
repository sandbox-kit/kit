package declarationgen

import "testing"

func TestCommonIndentIgnoresBlankLines(t *testing.T) {
	lines := []string{"    first", "", "      second"}
	if got := commonIndent(lines); got != 4 {
		t.Fatalf("commonIndent() = %d, want 4", got)
	}
}
