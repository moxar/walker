package walker_test

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/moxar/walker"
	"gopkg.in/yaml.v2"
)

// nolint
type Alias struct {
	Canon, Alias walker.Vertex
	Links        []walker.Vertex
}

type Case struct {
	Schema  *walker.Schema
	Aliases []Alias
	In      []string
	Want    string
	Err     string
}

func TestGraph_From(t *testing.T) {
	infos, err := ioutil.ReadDir("resources/tests")
	if err != nil {
		t.Error(err)
		return
	}
	for _, info := range infos {
		file := info.Name()
		if !strings.HasSuffix(file, ".yml") {
			continue
		}

		fixture := strings.TrimRight(strings.ReplaceAll(file, "_", " "), ".yml")
		t.Run(fixture, func(t *testing.T) {
			raw, err := ioutil.ReadFile("resources/tests/" + file)
			if err != nil {
				t.Error(err)
				return
			}
			var c Case
			c.Schema = walker.NewSchema()
			if err := yaml.Unmarshal(raw, &c); err != nil {
				t.Error(err)
				return
			}

			for _, al := range c.Aliases {
				c.Schema.Alias(al.Canon, al.Alias, al.Links...)
			}

			g, err := walker.NewGraph(c.Schema)
			if err != nil {
				t.Error(err)
				return
			}

			if len(c.In) == 0 {
				t.Error(`empty input field "In"`)
				return
			}
			got, err := g.From(c.In[0], c.In[1:]...)
			if (c.Err == "") != (err == nil) || err != nil && err.Error() != c.Err {
				t.Log("want err:", c.Err)
				t.Log("got err: ", err)
				t.Fail()
				return
			}

			if got != c.Want {
				t.Log("got:  ", got)
				t.Log("want: ", c.Want)
				t.Fail()
			}
		})
	}
}
