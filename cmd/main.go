package main

import (
	"encoding/json"
	"github.com/PuerkitoBio/goquery"
	"github.com/claudioontheweb/escherwyss-api/models"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"regexp"
	"strings"
)

func main() {

	r := mux.NewRouter()
	r.HandleFunc("/", GetMenuHandler)
	log.Fatal(http.ListenAndServe("localhost:8080", r))

}

func GetMenuHandler(w http.ResponseWriter, r *http.Request) {
	w. WriteHeader(http.StatusOK)

	menus := scrapeSite()

	err := json.NewEncoder(w).Encode(menus)
	if err != nil {
		panic(err)
	}

}

func scrapeSite() []models.Menu {
	response, err := http.Get("https://www.escherwyss.com/de/restaurant/men/tagesmen")
	if err != nil {
		panic(err)
	}

	defer response.Body.Close()

	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		panic(err)
	}

	var menus []models.Menu
	document.Find(".foodmenu").Each(func(i int, s *goquery.Selection) {
		menus = append(menus, createMenu(s.Text()))
	})

	return menus
}

func createMenu(menuString string) models.Menu {
	var menu models.Menu

	titleRegexp := regexp.MustCompile("Menü.*")
	menu.Title = titleRegexp.FindString(menuString)

	// Get menu Price
	priceRegexp := regexp.MustCompile("CHF.*")
	menu.Price = priceRegexp.FindString(menuString)

	// Get menu Content
	r := strings.NewReplacer(menu.Title, "", menu.Price, "", "\n", "")
	content := r.Replace(menuString)
	menu.Content = strings.TrimSpace(content)

	return menu
}