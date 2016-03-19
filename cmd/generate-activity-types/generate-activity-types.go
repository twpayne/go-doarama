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
	"// Activity types\n" +
	"const (\n" +
	"{{range $constName, $id := .ConstActivityIDs}}\t{{$constName}} = {{$id}}\n" +
	"{{end}})\n" +
	"\n" +
	"// DefaultActivityTypes contains the default activity types.\n" +
	"var DefaultActivityTypes = ActivityTypes{\n" +
	"{{range $at := .ActivityTypes}}\t{ID: {{$at.ID}}, Name: {{$at.Name | printf \"%#v\"}}},\n" +
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
	}
	return y
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

type byID doarama.ActivityTypes

func (ats byID) Len() int           { return len(ats) }
func (ats byID) Less(i, j int) bool { return ats[i].ID < ats[j].ID }
func (ats byID) Swap(i, j int)      { ats[i], ats[j] = ats[j], ats[i] }

func generateActivityIDs(filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	client := doarama.NewClient(doarama.APIURL, "", "")
	activityTypes, err := client.ActivityTypes()
	if err != nil {
		return err
	}
	sort.Sort(byID(activityTypes))
	maxConstLen := 0
	maxNameLen := 0
	constActivityIDs := make(map[string]int)
	for _, at := range activityTypes {
		maxNameLen = max(maxNameLen, len(at.Name))
		ss, err := constantize(at.Name)
		if err != nil {
			return err
		}
		for _, s := range ss {
			constActivityIDs[s] = at.ID
			maxConstLen = max(maxConstLen, len(s))
		}
	}
	paddedConstActivityIDs := make(map[string]int)
	for name, id := range constActivityIDs {
		paddedConstActivityIDs[pad(name, maxConstLen)] = id
	}
	if err := tmpl.Execute(f, struct {
		ActivityTypes    doarama.ActivityTypes
		ConstActivityIDs map[string]int
		Filename         string
	}{
		ActivityTypes:    activityTypes,
		ConstActivityIDs: paddedConstActivityIDs,
		Filename:         filename,
	}); err != nil {
		return err
	}
	return nil
}

func main() {
	o := flag.String("o", "", "")
	flag.Parse()
	if err := generateActivityIDs(*o); err != nil {
		log.Fatal(err)
	}
}
