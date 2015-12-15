package main

import (
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"regexp"
)

type Page struct {
	Title string
	Body  []byte
}

/*var editView, viewView = getTemplateFilePath()*/

// template file parse ONCE when it init.
var editPath, _ = filepath.Abs("./tmpl/edit.html")
var viewPath, _ = filepath.Abs("./tmpl/view.html")

var templates = template.Must(template.ParseFiles(editPath, viewPath))

// request uri validation
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

func (p *Page) save() error {
	filename := p.Title + ".txt"
	thePath, _ := filepath.Abs("./data/" + filename)
	return ioutil.WriteFile(thePath, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	thePath, _ := filepath.Abs("./data/" + filename)

	body, err := ioutil.ReadFile(thePath)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func getTitle(w http.ResponseWriter, r *http.Request) (string, error) {
	m := validPath.FindStringSubmatch(r.URL.Path)
	fmt.Printf(" url match result: %v", m)

	if m == nil {
		http.NotFound(w, r)
		return "", errors.New("Invalid Page Title")
	}

	return m[2], nil // The title is the second subexpression
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	fmt.Printf("title: %s \n", title)

	p, err := loadPage(title)
	if err != nil {
		fmt.Printf("loadPage func returns Error \n")
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}

	renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	fmt.Printf("title: %s \n", title)

	p, err := loadPage(title)
	if err != nil {
		fmt.Printf("loadPage func returns Error \n")
		p = &Page{Title: title}
	}

	renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body")
	fmt.Printf("title: %s", title)
	fmt.Printf("body: %v", body)

	p := &Page{Title: title, Body: []byte(body)}

	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		fmt.Printf("requst URI : %s", m)

		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}

func main() {
	/*
		p1 := &Page{Title: "view", Body: []byte("Great Power Takes Greate Responsibility.")}
		p1.save()

		fmt.Printf(" editPath: %v \n", editPath)
		fmt.Printf(" viewPath: %v \n", viewPath)
	*/

	http.Handle("/", http.FileServer(http.Dir(".")))

	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))

	http.ListenAndServe(":9000", nil)
}
