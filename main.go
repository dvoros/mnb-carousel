package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

func checkError(message string, err error) {
	if err != nil {
		log.Fatal(message, err)
	}
}

func errorJsonResp(w http.ResponseWriter, errorMsg string) {
	w.WriteHeader(500)
	_, err := w.Write([]byte(fmt.Sprintf(`{"Error":%s}`, errorMsg)))
	if err != nil {
		fmt.Printf("Ooooh! Error writing error response: %v\n", err)
	}
}

func errorHtmlResp(w http.ResponseWriter, errorMsg string) {
	w.WriteHeader(500)
	_, err := w.Write([]byte(fmt.Sprintf(`error: %s`, errorMsg)))
	if err != nil {
		fmt.Printf("Ooooh! Error writing error response: %v\n", err)
	}
}

func main() {
	if len(os.Args) != 3 {
		fmt.Printf("please provide CSV path and HTML path as arguments")
	}

	file, err := os.Open(os.Args[1])
	checkError("Cannot create file", err)
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()

	rand.Seed(time.Now().UnixNano())

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Add("Content-Type", "text/html; charset=UTF-8")
		htmlFile, err := os.Open(os.Args[2])
		if err != nil {
			errorHtmlResp(w, fmt.Sprintf("error opening HTML file: %v", err))
			return
		}
		bytes, err := ioutil.ReadAll(htmlFile)
		if err != nil {
			errorHtmlResp(w, fmt.Sprintf("error reading HTML file: %v", err))
			return
		}
		_, err = w.Write(bytes)
		if err != nil {
			fmt.Printf("Error writing response: %v", err)
			return
		}
	})


	http.HandleFunc("/random", func (w http.ResponseWriter, req *http.Request) {
		w.Header().Add("Content-Type", "application/json; charset=UTF-8")

		i := rand.Intn(len(records))
		randomLine := records[i]
		if len(randomLine) != 4 {
			errorJsonResp(w, "found a malformed record, sorry")
			return
		}

		response := struct{
			Response interface{} `json:"response"`
		}{
			struct{
				Id   string `json:"id"`
				Path string `json:"path"`
				Mime string `json:"mime"`
				Url  string `json:"url"`
			}{
				Id:   randomLine[0],
				Path: randomLine[1],
				Mime: randomLine[2],
				Url:  randomLine[3],
			},
		}

		bytes, err := json.Marshal(response)
		if err != nil {
			errorJsonResp(w, "couldn't marshal json response")
			return
		}

		_, err = w.Write(bytes)
		if err != nil {
			fmt.Printf("Error writing response: %v", err)
			return
		}
	})

	http.ListenAndServe(":8090", nil)
}