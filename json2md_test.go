package main

import (
	"strings"
	"testing"
)

func Test(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tests := []struct {
			name     string
			input    string
			expected string
		}{
			{
				name: "data types",
				input: `[
					{
						"null": null,
						"bool": true,
						"number": "12345.6789",
						"big number": 12345678901234567891234,
						"string": "hello world"
					}
				]`,
				expected: `
| null   | bool | number     | big number              | string      |
|--------|------|------------|-------------------------|-------------|
| <null> | true | 12345.6789 | 12345678901234567891234 | hello world |
				`,
			},
			{
				name: "header longer than content",
				input: `[
					{
						"very long header": "content"
					}
				]`,
				expected: `
| very long header |
|------------------|
| content          |
				`,
			},
			{
				name: "field order is preserved",
				input: `[
					{
						"b": 2,
						"c": 3,
						"a": 1
					},
					{
						"a": 10,
						"b": 20,
						"c": 30
					},
					{
						"c": 300,
						"a": 100,
						"b": 200
					}
				]`,
				expected: `
| b   | c   | a   |
|-----|-----|-----|
| 2   | 3   | 1   |
| 20  | 30  | 10  |
| 200 | 300 | 100 |
				`,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				output, err := json2md(nil, []byte(tt.input))
				received := strings.TrimSpace(string(output))
				expected := strings.TrimSpace(tt.expected)
				if err != nil {
					t.Fatal(err)
				}
				if received != expected {
					t.Fatalf("----- EXPECTED -----\n%s\n----- RECEIVED -----\n%s", expected, received)
				}
			})
		}
	})

	t.Run("errors", func(t *testing.T) {
		tests := []struct {
			name  string
			input string
		}{
			{
				name:  "invalid JSON",
				input: `@`,
			},
			{
				name:  "not array",
				input: `{"a": 1}`,
			},
			{
				name:  "invalid row start",
				input: `[@`,
			},
			{
				name:  "multiple values",
				input: `[{"a": 1}][{"b": 2}]`,
			},
			{
				name:  "row not an object",
				input: `[1]`,
			},
			{
				name:  "unexpected token in key position",
				input: `[{@:1}]`,
			},
			{
				name:  "key mismatch",
				input: `[{"a": 1}, {"b": 2}]`,
			},
			{
				name:  "missing key",
				input: `[{"a": 1, "b": 1}, {"a": 2}]`,
			},
			{
				name:  "missing value",
				input: `[{"a":}]`,
			},
			{
				name:  "array value",
				input: `[{"a":[]}]`,
			},
			{
				name:  "object value",
				input: `[{"a":{}}]`,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				_, err := json2md(nil, []byte(tt.input))
				if err == nil {
					t.Fatal("expected error")
				}
			})
		}
	})
}
