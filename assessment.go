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
	tmpl := `
	<!DOCTYPE html>
	<html>
	<head>
		<title>Elderly Fall Risk Assessment</title>
	</head>
	<body>
		<h1>Elderly Fall Risk Assessment</h1>
		<form action="/submit" method="POST">
			<label>How would you rate your leg strength (1-5)?</label><br>
			<input type="number" name="leg_strength" min="1" max="5" required><br><br>

			<label>How would you rate your eyesight (1-5)?</label><br>
			<input type="number" name="vision" min="1" max="5" required><br><br>

			<label>How would you rate your balance (1-5)?</label><br>
			<input type="number" name="balance" min="1" max="5" required><br><br>

			<label>Are you currently taking medication that affects your balance? (1 for Yes, 0 for No)</label><br>
			<input type="number" name="medication" min="0" max="1" required><br><br>

			<label>Have you had a fall in the last year? (1 for Yes, 0 for No)</label><br>
			<input type="number" name="history_of_falls" min="0" max="1" required><br><br>

			<label>Do you have a knee injury? (1 for Yes, 0 for No)</label><br>
			<input type="number" name="knee_injury" min="0" max="1" required><br><br>

			<input type="submit" value="Submit">
		</form>
		{{if .}}
			<h2>Your Fall Risk Assessment Result: {{.}}</h2>
		{{end}}
	</body>
	</html>
	`
	t, _ := template.New("form").Parse(tmpl)
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
