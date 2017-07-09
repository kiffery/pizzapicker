package main

import (
	"bufio"
	"encoding/csv"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("web/*"))
}

func index(res http.ResponseWriter, req *http.Request) {

	fn := req.FormValue("first")
	ln := req.FormValue("last")
	num, _ := strconv.ParseInt(req.FormValue("number"), 10, 64)
	meatNum, _ := strconv.ParseInt(req.FormValue("meatNum"), 10, 64)
	cheeseNum, _ := strconv.ParseInt(req.FormValue("cheeseNum"), 10, 64)
	vegNum, _ := strconv.ParseInt(req.FormValue("vegNum"), 10, 64)

	number := int64(len(fn)+len(ln)) + num

	err := tpl.ExecuteTemplate(res, "index.html", createPizza(number, meatNum, cheeseNum, vegNum))
	if err != nil {
		http.Error(res, err.Error(), 500)
		log.Fatal(err)
	}
}

func main() {
	http.HandleFunc("/", index)
	http.ListenAndServe(":8080", nil)
}

func createPizza(seed int64, meats int64, cheeses int64, vegs int64) Pizza {
	var pizza Pizza

	pizza.Crust = chooseItem("options/crust.txt", seed)
	pizza.Sauce = chooseItem("options/sauce.txt", seed)
	pizza.Cheeses = chooseItems("options/extra_cheese.txt", cheeses, seed)
	pizza.Meats = chooseItems("options/meat_toppings.txt", meats, seed)
	pizza.Veggies = chooseItems("options/veggie_toppings.txt", vegs, seed)

	return pizza
}

func chooseItem(itemList string, seed int64) string {
	file, err := os.Open(itemList)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	r := csv.NewReader(bufio.NewReader(file))
	r.TrimLeadingSpace = true
	options, err := r.ReadAll()

	rand.Seed(seed)
	option := options[rand.Intn(len(options))][0]
	return option
}

func chooseItems(itemList string, amount int64, seed int64) []string {
	file, err := os.Open(itemList)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	r := csv.NewReader(bufio.NewReader(file))
	r.TrimLeadingSpace = true
	options, err := r.ReadAll()
	var items []string
	var i int64
	for i < amount {
		i++
		seed += 5
		rand.Seed(seed)
		items = append(items, options[rand.Intn(len(options))][0])
	}
	return items
}
