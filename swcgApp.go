
package main

import (
        "swcg"
	"html/template"
	"net/http"
	"strconv"
	"fmt"
)

// Templates  -----------------------------------------------------------------

var templates = template.Must(template.ParseFiles("card.html"))

func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	err := templates.ExecuteTemplate(w, tmpl+".html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Card View ------------------------------------------------------------------

const CardViewPath = "/card/"

type CardView struct {
	Name            string
	Faction         string
	Type            string
	Cost            int
	Ressources      int
	ForceIcons      int
	CardCombatIcons string
	Abilities       string
	Health          int
	Quote           string
	ObjectiveSets   string
	Set             string
	Number          int
}

func CreateView(c *swcg.Card) *CardView {
	v := new(CardView)

	v.Name       	  = c.Name
	v.Faction    	  = swcg.FactionNames[c.Faction]
	v.Type 	     	  = swcg.CardTypeNames[c.Type.GetType()]
	v.Cost 	     	  = c.Cost
	v.Ressources 	  = c.Ressources
	v.ForceIcons 	  = c.ForceIcons
	v.CardCombatIcons = "TODO"
	v.Abilities 	  = "TODO"
	v.Health    	  = c.Health
	v.Quote     	  = c.Quote
	v.ObjectiveSets   = "TODO"
	v.Set    	  = swcg.SetNames[c.Set]
	v.Number 	  = c.Number

	return v
}


func cardViewHandler(w http.ResponseWriter, r *http.Request) {
	index, err := strconv.Atoi(r.URL.Path[len(CardViewPath):])
	if err != nil {
		http.Error(w, "Invalid Card ID Number", http.StatusInternalServerError)
		return
	}
	if (*CardDataCache.CardMap)[index] == nil {
		http.Error(w, fmt.Sprintf("Card Number %d not present in DataBase", index), http.StatusInternalServerError)
		return
	}
	data := CreateView((*CardDataCache.CardMap)[index])
	renderTemplate(w, "card", data)
}

var AllCards      []swcg.Card
var CardDataCache *swcg.DataCache

func main() {
	AllCards, CardDataCache = swcg.AnalyzeDB(swcg.CreateDB())
	http.HandleFunc(CardViewPath, cardViewHandler)
	http.ListenAndServe(":8080", nil)
}