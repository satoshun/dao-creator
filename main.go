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

var (
	index    = 0
	name     = flag.String("name", "", "base class name")
	topArray = flag.Bool("t", false, "top level array")
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: dao-creator https://hoge.com")
		return
	}

	flag.Parse()

	url := os.Args[len(os.Args)-1]
	r, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	defer r.Body.Close()

	var j map[string]interface{}
	dec := json.NewDecoder(r.Body)
	if *topArray {
		var jj []map[string]interface{}
		if err := dec.Decode(&jj); err != nil {
			panic(err)
			return
		}
		j = jj[0]
	} else {
		var jj map[string]interface{}
		if err := dec.Decode(&jj); err != nil {
			panic(err)
			return
		}
		j = jj
	}
	parse(j)
}

func parse(j map[string]interface{}) {
	if index == 0 {
		fmt.Println(fmt.Sprintf("@AutoValue public abstract class %s {", *name))
	} else {
		fmt.Println(fmt.Sprintf("@AutoValue public abstract class %s%d {", *name, index))
	}

	for k, v := range j {
		s := fmt.Sprintf(`    @SerializedName("%s") abstract %s `, k, getType(v))
		s += normalize(k) + "();"
		fmt.Println(s)
	}
	fmt.Println("}")
	index++
}

func getType(t interface{}) string {
	switch v := t.(type) {
	case float64:
		s := strconv.FormatFloat(v, 'f', -1, 64)
		if strings.Contains(s, ".") {
			return "double"
		} else {
			return "int"
		}
	case bool:
		return "boolean"
	case string:
		return "String"
	case []byte:
		return "String"
	default:
		// TODO; check array
		return "unknown"
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
