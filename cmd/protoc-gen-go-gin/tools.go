package main

import "strings"

func snakeCase(s string) string {
	var b []byte
	for i := 0; i < len(s); i++ { // proto identifiers are always ASCII
		c := s[i]
		if isASCIIUpper(c) {
			if i > 0 && b != nil && b[len(b)-1] != '_' {
				b = append(b, '_')
			}
			c += 'a' - 'A' // convert to lowercase
		}
		b = append(b, c)
	}
	return string(b)
}

func isASCIIUpper(c byte) bool {
	return 'A' <= c && c <= 'Z'
}

func getPackageName(importName string) string {
	return importName[strings.LastIndex(importName, "/")+1:]
}

func stringsLastIndex(s, sep string) string {
	if i := strings.LastIndex(s, sep); i != -1 {
		return s[i+len(sep):]
	}
	return s
}
