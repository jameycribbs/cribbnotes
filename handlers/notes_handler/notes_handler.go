package notes_handler

import (
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
	TotalRecs    int
	StatusMsg    string
}

type TemplateData struct {
	SearchString string
	Rec          *models.Note
	CsrfToken    string
	TotalRecs    int
	StatusMsg    string
}

func Create(w http.ResponseWriter, r *http.Request, throwaway string, dataDir string) {
	title := r.FormValue("title")
	text := r.FormValue("text")

	if title == "" || text == "" {
		msg := "Missing note title or text!"
		rec := models.Note{Title: title, Text: text}

		templateData := TemplateData{Rec: &rec, CsrfToken: nosurf.Token(r), TotalRecs: models.NotesCount(dataDir), StatusMsg: msg}

		renderTemplate(w, "new", &templateData)
	} else {
		rec := models.Note{Title: title, Text: text, CreatedAt: time.Now(), UpdatedAt: time.Now()}

		_, err := models.CreateNote(dataDir, &rec)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func Delete(w http.ResponseWriter, r *http.Request, fileId string, dataDir string) {
	var rec *models.Note

	rec, err := models.FindNote(dataDir, fileId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	templateData := TemplateData{Rec: rec, CsrfToken: nosurf.Token(r), TotalRecs: models.NotesCount(dataDir)}

	renderTemplate(w, "delete", &templateData)
}

func Destroy(w http.ResponseWriter, r *http.Request, throwaway string, dataDir string) {
	fileId := r.FormValue("fileId")

	err := models.DeleteNote(dataDir, fileId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

func Edit(w http.ResponseWriter, r *http.Request, fileId string, dataDir string) {
	var rec *models.Note

	rec, err := models.FindNote(dataDir, fileId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	templateData := TemplateData{Rec: rec, CsrfToken: nosurf.Token(r), TotalRecs: models.NotesCount(dataDir)}

	renderTemplate(w, "edit", &templateData)
}

func Index(w http.ResponseWriter, r *http.Request, throwAway string, dataDir string) {
	var err error

	templateData := IndexTemplateData{}

	templateData.SearchString = r.FormValue("searchString")

	templateData.Notes, err = models.FindNotes(dataDir, templateData.SearchString)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	templateData.CsrfToken = nosurf.Token(r)

	templateData.TotalRecs = models.NotesCount(dataDir)

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

func New(w http.ResponseWriter, r *http.Request, throwaway string, dataDir string) {
	rec := models.Note{}

	templateData := TemplateData{Rec: &rec, CsrfToken: nosurf.Token(r), TotalRecs: models.NotesCount(dataDir)}

	renderTemplate(w, "new", &templateData)
}

func Update(w http.ResponseWriter, r *http.Request, throwaway string, dataDir string) {
	fileId := r.FormValue("fileId")
	title := r.FormValue("title")
	text := r.FormValue("text")

	if title == "" || text == "" {
		msg := "Missing note title or text!"
		rec := models.Note{FileId: fileId, Title: title, Text: text}

		templateData := TemplateData{Rec: &rec, CsrfToken: nosurf.Token(r), TotalRecs: models.NotesCount(dataDir), StatusMsg: msg}

		renderTemplate(w, "edit", &templateData)
	} else {
		rec, err := models.FindNote(dataDir, fileId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		rec.FileId = fileId
		rec.Title = title
		rec.Text = text
		rec.UpdatedAt = time.Now()

		err = models.UpdateNote(dataDir, rec, fileId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusFound)
	}
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
