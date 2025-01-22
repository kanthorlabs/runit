package python

import (
	"bufio"
	"bytes"
	_ "embed"
	"regexp"
	"slices"
	"strings"
)

//go:embed packages.txt
var packages string
var PackageSystem []string
var PackageRegex = regexp.MustCompile(`^\s*(import|from)\s+([\w.]+)`)

func init() {
	PackageSystem = strings.Split(packages, "\n")
}

func ParseImports(codes []byte, excludes []string) []string {
	maps := make(map[string]int)

	scanner := bufio.NewScanner(bytes.NewReader(codes))
	for scanner.Scan() {
		line := scanner.Text()
		if matches := PackageRegex.FindStringSubmatch(line); matches != nil {
			parts := strings.Split(matches[2], ".")

			if len(excludes) > 0 && slices.Contains(excludes, parts[0]) {
				continue
			}

			if count, exist := maps[parts[0]]; exist {
				maps[parts[0]] = count + 1
			} else {
				maps[parts[0]] = 1
			}
		}
	}

	var imports []string
	for pkg := range maps {
		imports = append(imports, pkg)
	}

	return imports
}
