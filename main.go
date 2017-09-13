package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
)

var sum int = 0

func main() {
	fmt.Println("Calculator is up!")

	var sum int = 0

	http.HandleFunc("/ac", func(w http.ResponseWriter, r *http.Request) {
		sum = 0
	})

	http.HandleFunc("/total", func(w http.ResponseWriter, r *http.Request) {
		resp := struct{ Total int }{Total: sum}

		js, err := json.Marshal(resp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	})

	http.HandleFunc("/add", func(w http.ResponseWriter, r *http.Request) {
		if rand.Intn(100)%2 == 0 {
			http.Error(w, "Failing to showcase test execution Retry logic", http.StatusBadRequest)
			return
		}

		vals, _ := r.URL.Query()["val"]
		for _, val := range vals {
			v, err := strconv.Atoi(val)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
			}
			sum += v
		}
	})

	http.HandleFunc("/rand", func(w http.ResponseWriter, r *http.Request) {
		sum := rand.Intn(100)

		resp := struct{ Total int }{Total: sum}

		js, err := json.Marshal(resp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	})

	http.ListenAndServe(":3000", nil)
}
