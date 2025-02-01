package pythonx

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

func Scan(scanner *bufio.Scanner, excludes []string) (map[string]int, error) {
	maps := make(map[string]int)

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

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return maps, nil
}

func Lockfile(maps map[string]int) *bytes.Buffer {
	buff := new(bytes.Buffer)

	for pkg := range maps {
		buff.WriteString(pkg + "\n")
	}

	return buff
}
