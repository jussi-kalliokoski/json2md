package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
)

func main() {
	input, err := io.ReadAll(os.Stdin)
	fatal(err)
	output, err := json2md(nil, input)
	fatal(err)
	_, err = os.Stdout.Write(output)
	fatal(err)
}

func fatal(err error) {
	if err == nil {
		return
	}
	code := 1
	fmt.Fprintf(os.Stderr, "error: %v\n", err)
	os.Exit(code)
}

func json2md(dst []byte, data []byte) ([]byte, error) {
	var header []string
	var rows [][]string
	colsByKey := make(map[string]int)
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()
	if t, err := dec.Token(); err != nil {
		return nil, fmt.Errorf("error reading JSON data: %w", err)
	} else if delim, ok := t.(json.Delim); !ok || delim != '[' {
		return nil, fmt.Errorf("expected JSON array, got %T", t)
	}
	for {
		t, err := dec.Token()
		if err != nil {
			return nil, fmt.Errorf("error reading JSON row: %w", err)
		}
		if delim, ok := t.(json.Delim); ok && delim == ']' {
			if t, err := dec.Token(); err == nil {
				return nil, fmt.Errorf("expected end of JSON data, got %T", t)
			}
			return appendTable(dst, header, rows), nil
		}
		if delim, ok := t.(json.Delim); !ok || delim != '{' {
			return nil, fmt.Errorf("expected JSON object, got %T", t)
		}
		row := make([]string, len(header))
		hasHeader := rows != nil
		valuesFound := 0
		for {
			t, err := dec.Token()
			if err != nil {
				return nil, fmt.Errorf("error reading JSON object key: %w", err)
			}
			if delim, ok := t.(json.Delim); ok && delim == '}' {
				if valuesFound < len(header) {
					return nil, fmt.Errorf("row does not contain same values as the first row")
				}
				break
			}
			// the token can only be a string as enforced by the JSON decoder
			key, _ := t.(string)
			if !hasHeader {
				colsByKey[key] = len(header)
				header = append(header, key)
				row = append(row, "")
			}
			idx, keyFound := colsByKey[key]
			if !keyFound {
				return nil, fmt.Errorf("unexpected JSON object key (not included in the first row): %q", key)
			}
			t, err = dec.Token()
			if err != nil {
				return nil, fmt.Errorf("error reading JSON object value: %w", err)
			}
			if delim, ok := t.(json.Delim); ok {
				return nil, fmt.Errorf("expected a JSON primitive value, got delimiter: %v", delim)
			}
			val, err := tokenToString(t)
			if err != nil {
				return nil, err
			}
			row[idx] = val
			valuesFound++
		}
		rows = append(rows, row)
	}
}

func tokenToString(token json.Token) (string, error) {
	switch v := token.(type) {
	case string:
		return v, nil
	case bool:
		return strconv.FormatBool(v), nil
	case json.Number:
		return v.String(), nil
	case nil:
		return "<null>", nil
	default:
		return "", fmt.Errorf("unexpected JSON token: %T", v)
	}
}

func appendTable(dst []byte, header []string, rows [][]string) []byte {
	lens := make([]int, len(header))
	for i, col := range header {
		lens[i] = len(col)
	}
	for _, row := range rows {
		for i, col := range row {
			if len(col) > lens[i] {
				lens[i] = len(col)
			}
		}
	}
	dst = appendRow(dst, lens, header)
	dst = appendDivider(dst, lens)
	for _, row := range rows {
		dst = appendRow(dst, lens, row)
	}
	return dst
}

func appendRow(dst []byte, lens []int, row []string) []byte {
	dst = append(dst, '|')
	for i, col := range row {
		dst = append(dst, ' ')
		dst = append(dst, col...)
		dst = appendRepeated(dst, ' ', lens[i]-len(col))
		dst = append(dst, " |"...)
	}
	dst = append(dst, '\n')
	return dst
}

func appendDivider(dst []byte, lens []int) []byte {
	dst = append(dst, '|')
	for _, l := range lens {
		dst = appendRepeated(dst, '-', l+2)
		dst = append(dst, '|')
	}
	dst = append(dst, '\n')
	return dst
}

func appendRepeated(dst []byte, b byte, times int) []byte {
	for i := 0; i < times; i++ {
		dst = append(dst, b)
	}
	return dst
}
