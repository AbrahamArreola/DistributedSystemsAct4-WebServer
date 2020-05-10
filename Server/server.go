package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

var students = make(map[string]map[string]float32)
var subjects = make(map[string]map[string]float32)

func loadHTML(path string) string {
	html, _ := ioutil.ReadFile(path)

	return string(html)
}

func executeResponse(response http.ResponseWriter, request *http.Request, html string, responseString string) {
	response.Header().Set(
		"Content-Type",
		"text/html",
	)
	fmt.Fprintf(
		response,
		loadHTML(html),
		responseString,
	)
}

func addScore(response http.ResponseWriter, request *http.Request) {
	responseString := ""

	switch request.Method {
	case "POST":
		if err := request.ParseForm(); err != nil {
			fmt.Fprintf(response, "ParseForm() error %v", err)
			return
		}

		student := request.FormValue("student")
		subject := request.FormValue("subject")
		score, err := strconv.ParseFloat(request.FormValue("score"), 32)
		if err != nil {
			fmt.Fprintf(response, "Error to parse string to float", err)
			return
		}

		if _, exists := subjects[subject]; !exists {
			subjects[subject] = make(map[string]float32)
		}

		if _, exists := students[student]; !exists {
			students[student] = make(map[string]float32)
		}

		if _, exists := subjects[subject][student]; !exists {
			students[student][subject] = float32(score)
			subjects[subject][student] = float32(score)
			responseString = "Calificación agregada!"
		} else {
			responseString = "Error: calificación existente"
		}
	}
	executeResponse(response, request, "../Pages/add_score.html", responseString)
}

func studentAverage(response http.ResponseWriter, request *http.Request) {
	responseString := ""

	switch request.Method {
	case "POST":
		if err := request.ParseForm(); err != nil {
			fmt.Fprintf(response, "ParseForm() error %v", err)
			return
		}

		student := request.FormValue("student")
		if _, exists := students[student]; exists {
			var average float32 = 0

			for _, v := range students[student] {
				average += v
			}

			average /= float32(len(students[student]))
			responseString = "Promedio de alumno: " + fmt.Sprintf("%.2f", average)
		} else {
			responseString = "Error: alumno inexistente"
		}
	}
	executeResponse(response, request, "../Pages/student_average.html", responseString)
}

func subjectAverage(response http.ResponseWriter, request *http.Request) {
	responseString := ""

	switch request.Method {
	case "POST":
		if err := request.ParseForm(); err != nil {
			fmt.Fprintf(response, "ParseForm() error %v", err)
			return
		}

		subject := request.FormValue("subject")
		if _, exists := subjects[subject]; exists {
			var average float32 = 0

			for _, v := range subjects[subject] {
				average += v
			}

			average /= float32(len(subjects[subject]))
			responseString = "Promedio de materia: " + fmt.Sprintf("%.2f", average)
		} else {
			responseString = "Error: materia inexistente"
		}
	}
	executeResponse(response, request, "../Pages/subject_average.html", responseString)
}

func totalAverage(response http.ResponseWriter, request *http.Request) {
	responseString := ""

	switch request.Method {
	case "GET":
		if len(students) > 0 {
			var studentAverage float32 = 0
			var totalAverage float32 = 0

			for _, value := range students {
				for _, v := range value {
					studentAverage += v
				}

				studentAverage /= (float32)(len(value))
				totalAverage += studentAverage
				studentAverage = 0
			}

			totalAverage /= (float32)(len(students))
			responseString = "Promedio general: " + fmt.Sprintf("%.2f", totalAverage)
		} else {
			responseString = "Error: Ninguna calificación capturada"
		}
		executeResponse(response, request, "../Pages//total_average.html", responseString)
	}
}

func storageOptions(response http.ResponseWriter, request *http.Request) {
	response.Header().Set(
		"Content-Type",
		"text/html",
	)
	fmt.Fprint(
		response,
		loadHTML("../Pages/storage.html"),
	)
}

func saveInFile(response http.ResponseWriter, request *http.Request) {
	outFile, err := os.Create("../Storage/alumnos.txt")
	if err != nil {
		fmt.Println("Error al convertir a JSON", err.Error())
		return
	}
	codified := json.NewEncoder(outFile)
	codified.SetIndent("", "    ")
	if err := codified.Encode(students); err != nil {
		fmt.Println("Error al convertir a JSON", err.Error())
		return
	}
	outFile.Close()

	outFile, err = os.Create("../Storage/materias.txt")
	if err != nil {
		fmt.Println("Error al convertir a JSON", err.Error())
		return
	}
	codified = json.NewEncoder(outFile)
	codified.SetIndent("", "    ")
	if err := codified.Encode(subjects); err != nil {
		fmt.Println("Error al convertir a JSON", err.Error())
		return
	}
	outFile.Close()

	storageOptions(response, request)
}

func loadFromFile(response http.ResponseWriter, request *http.Request) {

}

func main() {
	host := "127.0.0.1:9000"
	http.HandleFunc("/add_score", addScore)
	http.HandleFunc("/student_average", studentAverage)
	http.HandleFunc("/total_average", totalAverage)
	http.HandleFunc("/subject_average", subjectAverage)
	http.HandleFunc("/storage_options", storageOptions)
	http.HandleFunc("/save_in_file", saveInFile)
	http.HandleFunc("/load_from_file", loadFromFile)
	fmt.Println("Servidor corriendo en:", host)
	http.ListenAndServe(host, nil)
}
