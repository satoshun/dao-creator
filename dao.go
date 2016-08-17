package dao

import (
	"fmt"
	"strconv"
	"strings"
)

type response map[string]interface{}

// Data is basic data
type Data struct {
	Res  response
	Name string
}

var cache []*Data

func AddData(d *Data) {
	cache = append(cache, d)
}

func PopData() *Data {
	if len(cache) == 0 {
		return nil
	}
	var j *Data
	j, cache = cache[0], cache[1:]
	return j
}

func GetPrinter(t string) Printer {
	if t == "autovalue" {
		return autovalue{}
	}
	return standardPrinter{}
}

// GetDefaultPrinter returns default printer
func GetDefaultPrinter() Printer {
	return standardPrinter{}
}

func init() {
	cache = make([]*Data, 0, 1)
}

// Printer is printer
type Printer interface {
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

type ri struct {
	index int
	key   string
}

func Parse(p Printer, j *Data, body string, verbose bool) {
	fmt.Println(p.header(j.Name))

	values := make([]ri, 0, len(j.Res))
	for k := range j.Res {
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
		v, k := j.Res[r.key], r.key
		name := normalize(k)
		s := p.item(k, getType(v, name), name)
		if verbose {
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
		cache = append(cache, &Data{Res: v, Name: n})
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
