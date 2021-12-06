/*
* @Author: Yajun
* @Date:   2021/12/6 12:57
 */

package main

import (
	"abtest"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"unsafe"
)

var (
	storage abtest.Storage
	err     error

	ex1DomainPtr unsafe.Pointer = unsafe.Pointer(&abtest.Domain{})
	ex2DomainPtr unsafe.Pointer = unsafe.Pointer(&abtest.Domain{})

	registry = &abtest.Registry{
		Name: "123",
		Dict: map[string]*unsafe.Pointer{
			"example1": &ex1DomainPtr,
			"example3": &ex2DomainPtr,
		}}
)

func init() {
	storage = abtest.NewLocalFile("example/domain")
	if err = storage.Register(registry); err != nil {
		log.Panicln("RegisterMap failed", err)
	}
	storage.Init()
	go storage.Watch()
}

func DomainHandler(w http.ResponseWriter, r *http.Request) {
	domain := r.URL.Query().Get("domain")
	d, ok := registry.Lookup(domain)
	if !ok {
		fmt.Fprintln(w, "domain not found")
	}
	b, _ := json.Marshal(d)
	fmt.Fprintln(w, string(b))
}

func AllocateHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("domain")
	dd, ok := registry.Lookup(name)
	if !ok {
		fmt.Fprintln(w, fmt.Sprintf("invalid domain: %s", name))
		return
	}
	res := dd.Execute(r.URL.Query().Get("id"))
	b, _ := json.Marshal(res)
	fmt.Fprintln(w, string(b))
}

func main() {
	http.HandleFunc("/domain", DomainHandler)
	http.HandleFunc("/allocate", AllocateHandler)
	if err := http.ListenAndServe(":8001", nil); err != nil {
		log.Panicln(err)
	}
}
