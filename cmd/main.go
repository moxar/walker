package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/moxar/walker"

	"github.com/go-yaml/yaml"
)

// nolint
type Alias struct {
	Canon, Alias walker.Vertex
	Links        []walker.Vertex
}

func debug(v interface{}) {
	raw, _ := yaml.Marshal(v)
	fmt.Println(string(raw))
}

func main() {

	raw, err := ioutil.ReadFile("schema.yml")
	if err != nil {
		log.Fatalln(err)
	}

	schema := walker.NewSchema()
	if err := yaml.Unmarshal(raw, &schema); err != nil {
		log.Fatalln(err)
	}

	var aliases struct {
		Aliases map[string]Alias
	}
	if err := yaml.Unmarshal(raw, &aliases); err != nil {
		log.Fatalln(err)
	}
	for _, al := range aliases.Aliases {
		schema.Alias(al.Canon, al.Alias, al.Links...)
	}

	debug(schema)

	g, err := walker.NewGraph(schema)
	if err != nil {
		log.Fatalln(err)
	}

	for _, c := range [][2]string{
		{"users", "customers"},
		{"customers", "users"},
		{"users", "projects"},
		{"projects", "customers"},
		{"projects", "owners"},
		{"projects", "users"},
		{"customers", "owners"},
		{"customer_user", "customer_user2"},
		{"customer_user_group_user", "users"},
		{"user_family", "users"},
	} {
		p, err := g.From(c[0], c[1])
		if err != nil {
			fmt.Println(c, "->", err)
			continue
		}
		fmt.Println(c)
		for _, l := range p {
			fmt.Println(l)
		}
		fmt.Println()
	}

	p, err := g.From("users", "customers", "projects")
	if err != nil {
		fmt.Println(err)
	}
	for _, l := range p {
		fmt.Println(l)
	}
	fmt.Println()

}
