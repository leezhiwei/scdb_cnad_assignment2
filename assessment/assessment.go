package main

import (
	"fmt"
	"html/template"
	"net/http"
)

type Assessment struct {
	LegStrength    int
	Vision         int
	Balance        int
	Medication     bool
	HistoryOfFalls bool
	KneeInjury     bool
}

var result string

func main() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/submit", submitAssessment)
	http.ListenAndServe(":8080", nil)
}

func homePage(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("assessment.html")
	if err != nil {
		fmt.Println("Error loading template:", err)
		http.Error(w, "Error Loading Assessment Form. Please Try Again.", http.StatusInternalServerError)
		return
	}
	t.Execute(w, result)
}

func submitAssessment(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		legStrength := r.FormValue("leg_strength")
		vision := r.FormValue("vision")
		balance := r.FormValue("balance")
		medication := r.FormValue("medication")
		historyOfFalls := r.FormValue("history_of_falls")
		kneeInjury := r.FormValue("knee_injury")

		assessment := Assessment{
			LegStrength:    parseToInt(legStrength),
			Vision:         parseToInt(vision),
			Balance:        parseToInt(balance),
			Medication:     medication == "1",
			HistoryOfFalls: historyOfFalls == "1",
			KneeInjury:     kneeInjury == "1",
		}

		result = calculateRisk(assessment)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func parseToInt(value string) int {
	var result int
	fmt.Sscanf(value, "%d", &result)
	return result
}

func calculateRisk(assessment Assessment) string {
	score := 0

	score += (6 - assessment.LegStrength)
	score += (6 - assessment.Vision)
	score += (6 - assessment.Balance)
	if assessment.Medication {
		score += 1
	}
	if assessment.HistoryOfFalls {
		score += 2
	}
	if assessment.KneeInjury {
		score += 2
	}

	if score <= 2 {
		return "Low Risk"
	} else if score <= 4 {
		return "Moderate Risk"
	} else {
		return "High Risk - Please consult a healthcare professional."
	}
}
