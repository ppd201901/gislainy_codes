package main

import (
	"github.com/gorilla/mux"
	"github.com/gorilla/rpc"
	"github.com/gorilla/rpc/json"
	"log"
	"net/http"
)
type Args struct {
	Height float64
	Gender string
}

type IdealWeight float64

type Result float64

func (t *IdealWeight) Calculate(r *http.Request, args *Args, result *Result) error {
	log.Printf("Calculate %F with %s\n", args.Height, args.Gender)
	 
	var weight Result
	if args.Gender == "male" {
		weight = Result((72.7 * (args.Height)-58))
	} else {
			weight = Result((62.1 * (args.Height)-44.7))
	}
	*result = weight
	return nil
}

func main() {
	s := rpc.NewServer()
	s.RegisterCodec(json.NewCodec(), "application/json")
	s.RegisterCodec(json.NewCodec(), "application/json;charset=UTF-8")
	idealweight := new(IdealWeight)
	s.RegisterService(idealweight, "")
	r := mux.NewRouter()
	r.Handle("/rpc", s)
	http.ListenAndServe(":1234", r)
}
