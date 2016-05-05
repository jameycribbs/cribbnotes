package notes_handler

import (
	"github.com/jameycribbs/cribbnotes/global_vars"
	"github.com/jameycribbs/cribbnotes/models"
	"github.com/justinas/nosurf"
	"html/template"
	"net/http"
	"path"
	"time"
)

type IndexTemplateData struct {
	SearchString string
	Notes        []models.Note
	CsrfToken    string
}

type TemplateData struct {
	SearchString string
	Rec          *models.Note
	CsrfToken    string
}

func Index(w http.ResponseWriter, r *http.Request, throwAway string, gv *global_vars.GlobalVars) {
	var err error

	templateData := IndexTemplateData{}

	templateData.SearchString = r.FormValue("searchString")

	templateData.Notes, err = models.FindNotes(templateData.SearchString)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	templateData.CsrfToken = nosurf.Token(r)

	lp := path.Join("templates", "layouts", "layout.html")
	fp := path.Join("templates", "notes", "index.html")

	tmpl := template.New("idx")

	tmpl, err = tmpl.ParseFiles(lp, fp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, "layout", templateData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func New(w http.ResponseWriter, r *http.Request, throwaway string, gv *global_vars.GlobalVars) {
	templateData := TemplateData{CsrfToken: nosurf.Token(r)}

	renderTemplate(w, "new", &templateData)
}

func Create(w http.ResponseWriter, r *http.Request, throwaway string, gv *global_vars.GlobalVars) {
	title := r.FormValue("title")
	text := r.FormValue("text")

	rec := models.Note{Title: title, Text: text, CreatedAt: time.Now(), UpdatedAt: time.Now()}

	_, err := models.CreateNote(&rec)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

func Edit(w http.ResponseWriter, r *http.Request, fileId string, gv *global_vars.GlobalVars) {
	var rec *models.Note

	rec, err := models.FindNote(fileId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	templateData := TemplateData{Rec: rec, CsrfToken: nosurf.Token(r)}

	lp := path.Join("templates", "layouts", "layout.html")
	fp := path.Join("templates", "notes", "edit.html")

	tmpl := template.New("edt")

	tmpl, err = tmpl.ParseFiles(lp, fp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, "layout", templateData)
}

func Update(w http.ResponseWriter, r *http.Request, throwaway string, gv *global_vars.GlobalVars) {
	var rec *models.Note

	fileId := r.FormValue("fileId")
	title := r.FormValue("title")
	text := r.FormValue("text")

	rec, err := models.FindNote(fileId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rec.FileId = fileId
	rec.Title = title
	rec.Text = text
	rec.UpdatedAt = time.Now()

	err = models.UpdateNote(rec, fileId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

func Delete(w http.ResponseWriter, r *http.Request, fileId string, gv *global_vars.GlobalVars) {
	var rec *models.Note

	rec, err := models.FindNote(fileId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	templateData := TemplateData{Rec: rec, CsrfToken: nosurf.Token(r)}
	renderTemplate(w, "delete", &templateData)
}

func Destroy(w http.ResponseWriter, r *http.Request, throwaway string, gv *global_vars.GlobalVars) {
	fileId := r.FormValue("fileId")

	err := models.DeleteNote(fileId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

//=============================================================================
// Helper Functions
//=============================================================================
func renderTemplate(w http.ResponseWriter, templateName string, templateData *TemplateData) {
	lp := path.Join("templates", "layouts", "layout.html")
	fp := path.Join("templates", "notes", templateName+".html")

	tmpl, _ := template.ParseFiles(lp, fp)
	err := tmpl.ExecuteTemplate(w, "layout", templateData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
