package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
)

var classes = make(map[string]map[string]float64)
var students = make(map[string]map[string]float64)

func main() {
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", serveMainPage)
	http.HandleFunc("/index", servePages)
	http.HandleFunc("/createScore", servePages)
	http.HandleFunc("/studentAvg", servePages)
	http.HandleFunc("/classAvg", servePages)
	http.HandleFunc("/generalAvg", handleGenAvg)
	http.HandleFunc("/createStudentGrade", handleGradeCreation)
	http.HandleFunc("/studentAvgInfo", handleStdAvgInf)
	http.HandleFunc("/classAvgInfo", handleClassAvgInf)

	log.Println("Listening on :3000...")
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func handleGenAvg(res http.ResponseWriter, req *http.Request) {
	var message string
	switch req.Method {
	case "POST":
		fmt.Fprintf(res, "Cannot POST /generalAvg") //Only allow get request for this
	case "GET":
		avgScore := float64(0)
		classesForStudent := float64(0)
		if len(students) > 0 {
			for _, sClasses := range students { //Iterating through all of the students
				for _, grade := range sClasses { //Iterating through all of the classes for a specific student
					avgScore += grade
				}
				classesForStudent += float64(len(sClasses)) //Adding the # of classes for one student to get our global average
			}
			averageScore := avgScore / classesForStudent
			message = "The score across all classes for all students is: " + fmt.Sprintf("%f", averageScore)
		} else {
			message = "Nothing has been registered yet"
		}
		main := filepath.Join("templates", "generalAvg.html")
		tmpl, err := template.ParseFiles(main)
		if err != nil {
			panic(err)
		}
		err = tmpl.Execute(res, message)
	}
}

func handleStdAvgInf(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "POST":
		if err := req.ParseMultipartForm(32 << 20); err != nil {
			fmt.Fprintf(res, "ParseForm() error %v", err)
			return
		}

		//<------------------------Backend handler---------------------------->
		studentVal := strings.Title(req.FormValue("student-name"))
		sClasses, studentExists := students[studentVal]
		if studentExists { //Checking if the student exists first
			avgScore := float64(0)
			for _, grade := range sClasses { //Iterating through all of the clases for that student
				avgScore += grade
			}
			avgStudScore := avgScore / float64(len(sClasses)) //Dividing the amount of classes for said student with the grade sum of score across all classes.
			res.Write([]byte("Student score is: " + fmt.Sprintf("%f", avgStudScore)))
		} else {
			res.Write([]byte("Student: " + studentVal + " doesn't exist")) //Student was not found
		}
	case "GET":
		fmt.Fprintf(res, "Cannot GET /studentAvgInfo")
	}
}

func handleClassAvgInf(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "POST":
		if err := req.ParseMultipartForm(32 << 20); err != nil {
			fmt.Fprintf(res, "ParseForm() error %v", err)
			return
		}

		// <-------------------------- Backend handler ----------------------------------->
		classVal := strings.Title(req.FormValue("class-name"))
		cStudents, classExists := classes[classVal]
		if classExists { //Checking if the class exists first
			avgScore := float64(0)
			for _, grade := range cStudents {
				avgScore += grade
			}
			averageScore := avgScore / float64(len(cStudents)) //Dividing the amount of students for said class.
			res.Write([]byte("Student score is: " + fmt.Sprintf("%f", averageScore)))
		} else {
			res.Write([]byte("Class: " + classVal + " doesn't exist")) //Class was not found
		}
	case "GET":
		fmt.Fprintf(res, "Cannot GET /classAvgInfo")
	}
}

func handleGradeCreation(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "POST":

		if err := req.ParseMultipartForm(32 << 20); err != nil {
			fmt.Fprintf(res, "ParseForm() error %v", err)
			return
		}
		//<------------------ Backend handling info------------------------>
		class := strings.Title(req.FormValue("class-name"))
		fmt.Println("Received request for class: " + class)
		student := strings.Title(req.FormValue("student-name"))
		grade, err := strconv.ParseFloat(req.FormValue("grade"), 64) //Invoking an attribute multiple times is more costly than assigning the value to a variable
		if err != nil {
			panic(err)
		}

		_, classExists := classes[class] //Checking if the class already exists
		if classExists {                 //Class found checking if student exists next
			_, studentExists := classes[class][student]
			if studentExists { //A student was found returning error message
				res.Write([]byte("Student " + student + " already has a grade for: " + class))
				return
			}
			//We don't need to create a new hash for this entry since we did found something in the map already meaning there's at least one entry in here.
			classes[class][student] = grade //else not required here since if it goes in the previous statement the control would end there.
		} else {
			//Creating hash for class since this is the first entry
			classes[class] = make(map[string]float64)
			classes[class][student] = grade
		}
		_, studentExists := students[student] //Checking if the student exists
		if !studentExists {                   //If it doesn't we need to create the hash for the classes since this is the first entry
			students[student] = make(map[string]float64)
		}
		students[student][class] = grade
		res.Write([]byte("Grade added correctly"))

	case "GET": //Don't allow get request specifically for this.
		fmt.Fprintf(res, "Cannot GET /createStudentGrade")
	}
}

func serveMainPage(res http.ResponseWriter, req *http.Request) {
	layout := filepath.Join("templates", "layout.html")
	main := filepath.Join("templates", "dashboard.html")
	tmpl, err := template.ParseFiles(layout, main)
	if err != nil {
		panic(err)
	}
	err = tmpl.ExecuteTemplate(res, "template", nil)
	if err != nil {
		panic(err)
	}
}

func servePages(res http.ResponseWriter, req *http.Request) {
	layout := filepath.Join("templates", "layout.html")
	main := filepath.Join("templates", filepath.Clean(req.URL.Path)+".html")

	tmpl, err := template.ParseFiles(layout, main)
	if err != nil {
		panic(err)
	}
	err = tmpl.ExecuteTemplate(res, "template", nil)
	if err != nil {
		panic(err)
	}
}
