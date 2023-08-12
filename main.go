package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

type Cards struct {
	Data []Card `json:"cards"`
}

type Card struct {
	ID          int    `json:"id"`
	TypeCard    string `json:"type_card"`
	Description string `json:"description"`
}

func generateCard(w http.ResponseWriter, r *http.Request) {

	jsonFile, err := os.Open("card_data.json")

	if err != nil {
		fmt.Println(err)
	}

	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)

	var res Cards

	json.Unmarshal(byteValue, &res)

	number := randomizeNumber(res, 0, len(res.Data))

	generateHtml(w, r, res, number)
}

func randomizeNumber(res Cards, min, max int) int {
	rand.Seed(time.Now().UnixNano())
	randomNum := min + rand.Intn((max-1)-min+1)

	return randomNum
}

func generateHtml(w http.ResponseWriter, r *http.Request, res Cards, number int) {
	var tmpl, err = template.ParseFiles("views/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var data = map[string]interface{}{
		"title":       "YUK NGOBROL",
		"type_card":   res.Data[number].TypeCard,
		"description": res.Data[number].Description,
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handleRequests() {
	http.HandleFunc("/generate", generateCard)
	http.Handle("/static/",
		http.StripPrefix("/static/",
			http.FileServer(http.Dir("assets"))))
	log.Println("Service is listening at port 3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}

func main() {
	handleRequests()
}
