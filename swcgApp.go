
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
	AbilitiesHTML   template.HTML
	Health          int
	Quote           string
	ObjectiveSets   string
	Set             string
	Number          int
}

func CombatIconsHTML(ci *swcg.CardCombatIcons) string {
	if ci == nil {
		return "N/A"
	}
	return fmt.Sprintf(
		"Combat Damage: %d [%d], Tactics: %d [%d], BlastDamage: %d [%d]",
		ci.CombatDamage[0], ci.CombatDamage[1],
		ci.Tactics[0], ci.Tactics[1],
		ci.BlastDamage[0], ci.BlastDamage[1])
}

func ObjectiveSetsHTML(sets *[]swcg.ObjectiveSet) string {
	out := ""
	size := len(*sets)
	for i, set := range *sets {
		out += fmt.Sprintf("%d (%d/6)", set.SetId, set.CardSetNumber)
		if i < size -1 {
			out += ", "
		}
	}
	return out
}

func AbilitiesHTML(abilities *[]swcg.AbilityInterface) string {
	out := ""
	traits := make([]*swcg.CardTrait,0)
	for _, a := range *abilities {
		switch ability := a.(type) {
		case *swcg.CardAbility:
			out += "<div class=\"ability\">"
			out += fmt.Sprintf("<span class=\"abilityName\">%s</span><span class=\"abilityDesc\">%s</span>",
				swcg.AbilityNames[ability.Type], ability.Description)
			out += "</div>"
		case *swcg.SimpleKeyword:
			out += "<div class=\"keyword\">"
			out += fmt.Sprintf("<span class=\"keyword\">%s</span>", swcg.KeywordNames[ability.K])
			out += "</div>"
		case *swcg.ComplexKeyword:
			out += "<div class=\"keyword\">"
			out += fmt.Sprintf("<span class=\"keyword\">%s [%d]</span>", swcg.KeywordNames[ability.K], ability.V)
			out += "</div>"
		case *swcg.ProtectKeywordType:
			out += "<div class=\"keyword\">"
			out += fmt.Sprintf("<span class=\"keyword\">%s from %s</span>", swcg.KeywordNames[ability.K], ability.ProtectedTrait)
			out += "</div>"
		case *swcg.CardTrait:
			traits = append(traits, ability)
		default:
			out += ""
		}
	}
	traitsStr := "<div class=\"traits\">"
	traitSize := len(traits)
	for i, trait := range traits {
		traitsStr += "<span class=\"trait\">"+swcg.TraitNames[trait.Trait]+"</span>"
		if i < traitSize -1 {
			traitsStr += ", "
		}
	}
	traitsStr += "</div>"
	return traitsStr + out
}

func CreateView(c *swcg.Card) *CardView {
	v := new(CardView)

	v.Name       	  = c.Name
	v.Faction    	  = swcg.FactionNames[c.Faction]
	v.Type 	     	  = swcg.CardTypeNames[c.Type.GetType()]
	v.Cost 	     	  = c.Cost
	v.Ressources 	  = c.Ressources
	v.ForceIcons 	  = c.ForceIcons
	v.CardCombatIcons = CombatIconsHTML(c.CardCombatIcons)
	v.AbilitiesHTML   = template.HTML(AbilitiesHTML(&c.Abilities))
	v.Health    	  = c.Health
	v.Quote     	  = c.Quote
	v.ObjectiveSets   = ObjectiveSetsHTML(&c.ObjectiveSets)
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