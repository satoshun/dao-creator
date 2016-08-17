package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/satoshun/dao-creator/dao"
)

var (
	name     = flag.String("name", "Unknown", "base class name")
	topArray = flag.Bool("t", false, "top level array")
	verbose  = flag.Bool("v", false, "verbose mode")
	lib      = flag.String("l", "", "kind of library(autovalue or other)")
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
		d := dao.Data{Res: jj[0], Name: *name}
		dao.AddData(&d)
	} else {
		var jj map[string]interface{}
		if err := dec.Decode(&jj); err != nil {
			panic(err)
		}
		d := dao.Data{Res: jj, Name: *name}
		dao.AddData(&d)
	}

	var p dao.Printer
	if *lib == "" {
		p = dao.GetDefaultPrinter()
	} else {
		p = dao.GetPrinter(*lib)
	}

	for j := dao.PopData(); j != nil; j = dao.PopData() {
		dao.Parse(p, j, string(body), *verbose)
	}
}
