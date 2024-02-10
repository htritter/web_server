package main

// import "fmt"
import "log"
import "net/http"
import "os"
import "html/template"

type Page struct {
	Title string
	Body []byte
}

func loadPage(title string) (*Page, error){
	filename := title + ".txt"
	body, err := os.ReadFile(filename)
	if err != nil{
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

// func handler(w http.ResponseWriter, r *http.Request){
// 	fmt.Fprint(w, "Hi there, %s!", r.URL.Path[1:])
// }

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page){
	t, err := template.ParseFiles(tmpl + ".html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	// t, _ := template.ParseFiles(tmpl + ".html")
	// t.Execute(w, p)
}

func viewHandler(w http.ResponseWriter, r *http.Request){
	title := r.URL.Path[len("/view/"):]
	p, _ := loadPage(title)
	if err != err {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
	}
	// t, _ := template.ParseFiles("view.html")
	// fmt.Fprintf(w, "<h1>%s</h><div>%s</div>", p.Title, p.Body)
	// t.Execute(w, p)
	renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request){
	title := r.URL.Path[len("/edit"):]
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	// t, _ := template.ParseFiles("edit.html")
	// t.Execute(w, p)
	renderTemplate(w, "edit", p)
	// fmt.Fprintf(w, "<h1>Editing %s</h1>"+
    //     "<form action=\"/save/%s\" method=\"POST\">"+
    //     "<textarea name=\"body\">%s</textarea><br>"+
    //     "<input type=\"submit\" value=\"Save\">"+
    //     "</form>",
    //     p.Title, p.Title, p.Body)
}

func saveHandler(w http.ResponseWriter, r *http.Request){
	title := r.URL.Path[len("/save/"):]
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func main() {
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/edit/", editHandler)
	// http.HandleFunc("/save/", saveHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}