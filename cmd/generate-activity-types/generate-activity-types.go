// FIXME there must be a better way to pad constants in templates
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"

	"github.com/twpayne/go-doarama"
)

var tmpl = template.Must(template.New("gat").Parse("" +
	"//go:generate go run cmd/generate-activity-types/generate-activity-types.go -o {{.Filename}}\n" +
	"package doarama\n" +
	"\n" +
	"const (\n" +
	"{{range $constName, $id := .ConstActivityTypes}}\t{{$constName}} = {{$id}}\n" +
	"{{end}})\n" +
	"\n" +
	"var (\n" +
	"\tActivityTypes = map[string]int{\n" +
	"{{range $name, $id := .FormattedActivityTypes}}\t\t{{$name}} {{$id}},\n" +
	"{{end}}\t}\n" +
	")\n",
))

func constantize(s string) ([]string, error) {
	var result []string
	ss := strings.SplitN(s, " - ", 2)
	switch len(ss) {
	case 1:
		result = append(result, strings.ToUpper(ss[0]))
	case 2:
		prefix := strings.ToUpper(ss[0])
		for _, desc := range strings.Split(ss[1], "/") {
			suffix := desc
			suffix = strings.TrimSuffix(suffix, " etc")
			suffix = strings.TrimSpace(suffix)
			suffix = strings.Replace(suffix, " ", "_", -1)
			suffix = strings.Replace(suffix, "+", "AND", -1)
			suffix = strings.ToUpper(suffix)
			result = append(result, prefix+"_"+suffix)
		}
	default:
		return nil, fmt.Errorf("unable to parse %v", s)
	}
	return result, nil
}

func max(x, y int) int {
	if x > y {
		return x
	} else {
		return y
	}
}

func pad(s string, n int) string {
	m := len(s)
	switch {
	case m < n:
		return s + strings.Repeat(" ", n-m)
	case m == n:
		return s
	default:
		return s[:m]
	}
}

func generateActivityTypes(filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	client := doarama.NewClient(doarama.API_URL, "", "")
	activityTypes, err := client.ActivityTypes()
	if err != nil {
		return err
	}
	/*
		activityTypes := map[string]int{
			"Boat - Kayak/Canoe/Row etc": 10,
			"Boat - Motor":               9,
			"Drive - Car/Truck/Bus etc":  24,
			"Snowboard":                  1,
		}
	*/
	maxConstLen := 0
	maxNameLen := 0
	constActivityTypes := make(map[string]int)
	for name, id := range activityTypes {
		maxNameLen = max(maxNameLen, len(name))
		ss, err := constantize(name)
		if err != nil {
			return err
		}
		for _, s := range ss {
			constActivityTypes[s] = id
			maxConstLen = max(maxConstLen, len(s))
		}
	}
	paddedFormattedActivityTypes := make(map[string]int)
	for name, id := range activityTypes {
		paddedFormattedActivityTypes[pad("\""+name+"\":", maxNameLen+3)] = id
	}
	paddedConstActivityTypes := make(map[string]int)
	for name, id := range constActivityTypes {
		paddedConstActivityTypes[pad(name, maxConstLen)] = id
	}
	if err := tmpl.Execute(f, struct {
		FormattedActivityTypes map[string]int
		ConstActivityTypes     map[string]int
		Filename               string
	}{
		FormattedActivityTypes: paddedFormattedActivityTypes,
		ConstActivityTypes:     paddedConstActivityTypes,
		Filename:               filename,
	}); err != nil {
		return err
	}
	return nil
}

func main() {
	o := flag.String("o", "", "")
	flag.Parse()
	if err := generateActivityTypes(*o); err != nil {
		log.Fatal(err)
	}
}
