// FIXME there must be a better way to pad constants in templates
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"text/template"

	"github.com/twpayne/go-doarama"
)

var tmpl = template.Must(template.New("gat").Parse("" +
	"package doarama\n" +
	"\n" +
	"const (\n" +
	"{{range $constName, $id := .ConstActivityIds}}\t{{$constName}} = {{$id}}\n" +
	"{{end}})\n" +
	"\n" +
	"var DefaultActivityTypes = ActivityTypes{\n" +
	"{{range $at := .ActivityTypes}}\t{Id: {{$at.Id}}, Name: {{$at.Name | printf \"%#v\"}}},\n" +
	"{{end}}}\n" +
	"\n" +
	"//go:generate go run cmd/generate-activity-types/generate-activity-types.go -o {{.Filename}}\n",
))

func constantize(s string) ([]string, error) {
	var result []string
	ss := strings.SplitN(s, " - ", 2)
	switch len(ss) {
	case 1:
		result = append(result, ss[0])
	case 2:
		prefix := ss[0]
		for _, desc := range strings.Split(ss[1], "/") {
			suffix := desc
			suffix = strings.TrimSuffix(suffix, " etc")
			suffix = strings.Replace(suffix, " ", "", -1)
			suffix = strings.TrimSpace(suffix)
			suffix = strings.Replace(suffix, "+", "And", -1)
			result = append(result, prefix+suffix)
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

type ById doarama.ActivityTypes

func (ats ById) Len() int           { return len(ats) }
func (ats ById) Less(i, j int) bool { return ats[i].Id < ats[j].Id }
func (ats ById) Swap(i, j int)      { ats[i], ats[j] = ats[j], ats[i] }

func generateActivityIds(filename string) error {
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
	sort.Sort(ById(activityTypes))
	maxConstLen := 0
	maxNameLen := 0
	constActivityIds := make(map[string]int)
	for _, at := range activityTypes {
		maxNameLen = max(maxNameLen, len(at.Name))
		ss, err := constantize(at.Name)
		if err != nil {
			return err
		}
		for _, s := range ss {
			constActivityIds[s] = at.Id
			maxConstLen = max(maxConstLen, len(s))
		}
	}
	paddedConstActivityIds := make(map[string]int)
	for name, id := range constActivityIds {
		paddedConstActivityIds[pad(name, maxConstLen)] = id
	}
	if err := tmpl.Execute(f, struct {
		ActivityTypes    doarama.ActivityTypes
		ConstActivityIds map[string]int
		Filename         string
	}{
		ActivityTypes:    activityTypes,
		ConstActivityIds: paddedConstActivityIds,
		Filename:         filename,
	}); err != nil {
		return err
	}
	return nil
}

func main() {
	o := flag.String("o", "", "")
	flag.Parse()
	if err := generateActivityIds(*o); err != nil {
		log.Fatal(err)
	}
}
