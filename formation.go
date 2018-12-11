package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"sort"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		os.Exit(1)
	}
}

type Table struct {
	Name      string
	DependsOn string
	HashKey   string
	RangeKey  string
	Ttl       string
	Indexes   []Index `yaml:"-"`

	RawKey     string            `yaml:"key"`
	RawIndexes map[string]string `yaml:"indexes"`
}

type Index struct {
	Name     string
	HashKey  string
	RangeKey string
}

type Tables map[string]Table

func run() error {
	data, err := ioutil.ReadFile("tables.yml")
	if err != nil {
		return err
	}

	var tm map[string]Table

	if err := yaml.Unmarshal(data, &tm); err != nil {
		return err
	}

	ts := []Table{}

	for k, v := range tm {
		v.Name = k

		p := strings.Split(v.RawKey, ",")

		v.HashKey = strings.TrimSpace(p[0])

		if len(p) > 1 {
			v.RangeKey = strings.TrimSpace(p[1])
		}

		v.Indexes = []Index{}

		for k, vi := range v.RawIndexes {
			i := Index{Name: k}

			p := strings.Split(vi, ",")

			i.HashKey = strings.TrimSpace(p[0])

			if len(p) > 1 {
				i.RangeKey = strings.TrimSpace(p[1])
			}

			v.Indexes = append(v.Indexes, i)
		}

		sort.Slice(v.Indexes, func(i, j int) bool { return v.Indexes[i].Name < v.Indexes[j].Name })

		ts = append(ts, v)
	}

	sort.Slice(ts, func(i, j int) bool { return ts[i].Name < ts[j].Name })

	for i := range ts {
		if i > 0 {
			ts[i].DependsOn = ts[i-1].Name
		}
	}

	t, err := template.New("formation").Funcs(helpers()).ParseFiles("formation.json.tmpl")
	if err != nil {
		return err
	}

	params := map[string]interface{}{
		"Tables": ts,
	}

	// fmt.Printf("multiplier = %+v\n", multiplier)
	// fmt.Printf("t = %+v\n", t)
	// fmt.Printf("params = %+v\n", params)

	if err := t.Execute(os.Stdout, params); err != nil {
		return err
	}

	return nil
}

func (t *Table) Attributes() []string {
	ah := map[string]bool{}

	ah[t.HashKey] = true
	ah[t.RangeKey] = true

	for _, i := range t.Indexes {
		ah[i.HashKey] = true
		ah[i.RangeKey] = true
	}

	as := []string{}

	for k := range ah {
		if k != "" {
			as = append(as, k)
		}
	}

	sort.Strings(as)

	return as
}

func helpers() map[string]interface{} {
	return map[string]interface{}{
		"type": func(s string) string {
			switch s {
			case "github-id":
				return "N"
			default:
				return "S"
			}
		},
		"times": func(i, j int) int {
			return i * j
		},
		"upper": func(name string) string {
			if name == "" {
				return ""
			}
			us := strings.ToUpper(name[0:1]) + name[1:]
			for {
				i := strings.Index(us, "-")
				if i == -1 {
					break
				}
				s := us[0:i]
				if len(us) > i+1 {
					s += strings.ToUpper(us[i+1 : i+2])
				}
				if len(us) > i+2 {
					s += us[i+2:]
				}
				us = s
			}
			return us
		},
	}
}
