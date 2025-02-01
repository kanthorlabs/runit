package pythonx

import (
	"bufio"
	"slices"
	"strings"
	"testing"

	"github.com/samber/lo"
)

func TestScan(t *testing.T) {
	packages := map[string][]string{
		"os":      {"", "os"},
		"os.path": {"", ""},
		"xml.dom": {"minidom", "xml"},
		"pandas":  {"", "pandas"},
	}

	var expected []string
	var codes string
	for pkg, input := range packages {
		if input[0] == "" {
			codes += "import " + pkg + "\n"
		} else {
			codes += "from " + pkg + " import " + input[0] + "\n"
		}

		if input[1] != "" {
			expected = append(expected, input[1])
		}
	}
	codes += "# some comment\n"

	scanner := bufio.NewScanner(strings.NewReader(codes))
	maps, err := Scan(scanner, nil)
	if err != nil {
		t.Fatal(err)
	}
	actual := lo.Keys(maps)

	slices.Sort(expected)
	slices.Sort(actual)
	if slices.Compare(expected, actual) != 0 {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}

func TestScan_ExternalOnly(t *testing.T) {
	packages := map[string][]string{
		"os":      {"", "os"},
		"os.path": {"", ""},
		"xml.dom": {"minidom", "xml"},
		"pandas":  {"", "pandas"},
	}

	expected := []string{"pandas"}
	var codes string
	for pkg, input := range packages {
		if input[0] == "" {
			codes += "import " + pkg + "\n"
		} else {
			codes += "from " + pkg + " import " + input[0] + "\n"
		}
	}
	codes += "# some comment\n"

	scanner := bufio.NewScanner(strings.NewReader(codes))
	maps, err := Scan(scanner, PackageSystem)
	if err != nil {
		t.Fatal(err)
	}
	actual := lo.Keys(maps)

	slices.Sort(expected)
	slices.Sort(actual)
	if slices.Compare(expected, actual) != 0 {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}
