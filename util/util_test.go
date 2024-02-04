package util

import (
	"testing"
)

func Test_RuneToInt(t *testing.T) {
	for name, tt := range map[string]struct {
		val    rune
		expect int
	}{
		"a": {
			rune("a"[0]),
			1,
		},
		"b": {
			rune("b"[0]),
			2,
		},
		"c": {
			rune("c"[0]),
			3,
		},
	} {
		t.Run(name, func(t *testing.T) {
			result := RuneToInt(tt.val)
			if result != tt.expect {
				t.Fatalf("expected %v, got %v", tt.expect, result)
			}
		})
	}
}
