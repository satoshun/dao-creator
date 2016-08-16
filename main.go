package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type response map[string]interface{}

var (
	name     = flag.String("name", "Unknown", "base class name")
	topArray = flag.Bool("t", false, "top level array")
	verbose  = flag.Bool("v", false, "verbose mode")
	lib      = flag.String("l", "", "kind of library(autovalue or other)")
)

type printer interface {
	header(name string) string
	item(t, v, name string) string
}

type standardPrinter struct {
}

func (a standardPrinter) header(name string) string {
	return fmt.Sprintf("public class %s {", name)
}

func (a standardPrinter) item(t, v, name string) string {
	return fmt.Sprintf(`    @SerializedName("%s") %s %s;`, t, v, name)
}

type autovalue struct {
}

func (a autovalue) header(name string) string {
	return fmt.Sprintf("@AutoValue public abstract class %s {", name)
}

func (a autovalue) item(t, v, name string) string {
	return fmt.Sprintf(`    @SerializedName("%s") abstract %s %s();`, t, v, name)
}

type data struct {
	r    response
	name string
}

var cache []data

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: dao-creator https://hoge.com")
		return
	}

	flag.Parse()

	cache = make([]data, 0)

	url := os.Args[len(os.Args)-1]
	r, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	br := bytes.NewReader(body)
	dec := json.NewDecoder(br)
	if *topArray {
		var jj []map[string]interface{}
		if err := dec.Decode(&jj); err != nil {
			panic(err)
		}
		d := data{r: jj[0], name: *name}
		cache = append(cache, d)
	} else {
		var jj map[string]interface{}
		if err := dec.Decode(&jj); err != nil {
			panic(err)
		}
		d := data{r: jj, name: *name}
		cache = append(cache, d)
	}

	var p printer = standardPrinter{}
	if *lib == "autovalue" {
		p = autovalue{}
	}

	var j data
	for len(cache) > 0 {
		j, cache = cache[0], cache[1:]
		parse(p, j, string(body))
	}
}

type ri struct {
	index int
	key   string
}

func parse(p printer, j data, body string) {
	fmt.Println(p.header(j.name))

	values := make([]ri, 0, len(j.r))
	for k := range j.r {
		index := strings.Index(body, "\""+k+"\"")
		r := ri{index: index, key: k}

		var f bool
		for i, v := range values {
			if v.index > r.index {
				values = append(values[0:i], append([]ri{r}, values[i:]...)...)
				f = true
				break
			}
		}
		if !f {
			values = append(values, r)
		}
	}

	for _, r := range values {
		v, k := j.r[r.key], r.key
		name := normalize(k)
		s := p.item(k, getType(v, name), name)
		if *verbose {
			s += " // " + fmt.Sprintf("%v", v)
		}
		fmt.Println(s)
	}
	fmt.Println("}")
}

func getType(t interface{}, name string) string {
	switch v := t.(type) {
	case float64:
		s := strconv.FormatFloat(v, 'f', -1, 64)
		if strings.Contains(s, ".") {
			return "double"
		}
		return "int"
	case bool:
		return "boolean"
	case string:
		// TODO: more smart logic
		if strings.Contains(v, "-") &&
			strings.Contains(v, ":") &&
			strings.Contains(v, "T") {
			return "Date"
		}
		return "String"
	case map[string]interface{}:
		n := strings.Title(name)
		cache = append(cache, data{r: v, name: n})
		return n
	case []interface{}:
		if len(v) == 0 {
			// unknown type
			return "List<Object>"
		}
		return "List<" + getType(v[0], name) + ">"
	default:
		return "Object"
	}
}

func normalize(s string) string {
	ss := ""
	skip := false
	for i, b := range s {
		if skip {
			skip = false
			continue
		}
		if i != 0 && b == '_' {
			skip = true
			ss += strings.ToUpper(string(s[i+1]))
			continue
		}
		ss += string(b)
	}
	return ss
}
