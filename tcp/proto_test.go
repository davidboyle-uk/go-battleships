package tcp

import (
	"testing"
)

func TestProtoToString(t *testing.T) {
	for name, tt := range map[string]struct {
		p        Proto
		expected string
	}{
		"control": {
			p: Proto{
				Action: "test",
				Player: 1,
				Body:   `{"key":false}`,
			},
			expected: `test|1|{"key":false}` + "\n",
		},
		"hello": {
			p: Proto{
				Action: "hello",
			},
			expected: `hello|0|` + "\n",
		},
	} {
		t.Run(name, func(t *testing.T) {
			actual := tt.p.String()

			if tt.expected != actual {
				t.Fatalf("Expected `%s` Got `%s`", tt.expected, actual)
			}
		})
	}
}

func TestParseMessage(t *testing.T) {
	for name, tt := range map[string]struct {
		p        string
		expected Proto
		err      bool
	}{
		"control": {
			p: `test|1|{"key":false}` + "\n",
			expected: Proto{
				Action: "test",
				Player: 1,
				Body:   `{"key":false}`,
			},
		},
		"hello": {
			p: `hello|0|` + "\n",
			expected: Proto{
				Action: "hello",
			},
		},
		"empty": {
			err: true,
		},
	} {
		t.Run(name, func(t *testing.T) {
			actual, err := ParseMessage(tt.p)
			if err != nil && !tt.err {
				t.Fatal(err)
			}

			if tt.expected != actual {
				t.Fatalf("Expected `%s` Got `%s`", tt.expected, actual)
			}
		})
	}
}
