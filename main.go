package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type response map[string]interface{}

var (
	name     = flag.String("name", "Unknown", "base class name")
	topArray = flag.Bool("t", false, "top level array")
)

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

	dec := json.NewDecoder(r.Body)
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

	var j data
	for len(cache) > 0 {
		j, cache = cache[0], cache[1:]
		parse(j)
	}
}

func parse(j data) {
	fmt.Println(fmt.Sprintf("@AutoValue public abstract class %s {", j.name))

	for k, v := range j.r {
		name := normalize(k)
		s := fmt.Sprintf(`    @SerializedName("%s") abstract %s `, k, getType(v, name))
		s += name + "();"
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
