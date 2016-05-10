package notes_handler

import (
	"github.com/jameycribbs/cribbnotes/models/note_model"
	"html/template"
	"net/http"
	"path"
	"time"
)

type IndexTemplateData struct {
	SearchString      string
	FocusSearchString bool
	Recs              []note_model.Record
	TotalRecs         int
	StatusMsg         string
}

type TemplateData struct {
	SearchString      string
	FocusSearchString bool
	Rec               *note_model.Record
	TotalRecs         int
	StatusMsg         string
}

func Create(w http.ResponseWriter, r *http.Request, throwaway string, dataDir string) {
	title := r.FormValue("title")
	text := r.FormValue("text")

	if title == "" || text == "" {
		msg := "Missing note title or text!"
		rec := note_model.Record{Title: title, Text: text}

		templateData := TemplateData{Rec: &rec, TotalRecs: note_model.Count(dataDir), StatusMsg: msg}

		renderTemplate(w, "new", &templateData)
	} else {
		rec := note_model.Record{Title: title, Text: text, CreatedAt: time.Now(), UpdatedAt: time.Now()}

		_, err := note_model.Create(dataDir, &rec)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func Delete(w http.ResponseWriter, r *http.Request, fileId string, dataDir string) {
	var rec *note_model.Record

	rec, err := note_model.Find(dataDir, fileId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	templateData := TemplateData{Rec: rec, TotalRecs: note_model.Count(dataDir)}

	renderTemplate(w, "delete", &templateData)
}

func Destroy(w http.ResponseWriter, r *http.Request, throwaway string, dataDir string) {
	fileId := r.FormValue("fileId")

	err := note_model.Delete(dataDir, fileId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

func Edit(w http.ResponseWriter, r *http.Request, fileId string, dataDir string) {
	var rec *note_model.Record

	rec, err := note_model.Find(dataDir, fileId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	templateData := TemplateData{Rec: rec, TotalRecs: note_model.Count(dataDir)}

	renderTemplate(w, "edit", &templateData)
}

func Index(w http.ResponseWriter, r *http.Request, throwAway string, dataDir string) {
	var err error

	templateData := IndexTemplateData{}

	templateData.SearchString = r.FormValue("searchString")
	templateData.FocusSearchString = true

	templateData.Recs, err = note_model.Search(dataDir, templateData.SearchString)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	templateData.TotalRecs = note_model.Count(dataDir)

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
	rec := note_model.Record{}

	templateData := TemplateData{Rec: &rec, TotalRecs: note_model.Count(dataDir)}

	renderTemplate(w, "new", &templateData)
}

func Update(w http.ResponseWriter, r *http.Request, throwaway string, dataDir string) {
	fileId := r.FormValue("fileId")
	title := r.FormValue("title")
	text := r.FormValue("text")

	if title == "" || text == "" {
		msg := "Missing note title or text!"
		rec := note_model.Record{FileId: fileId, Title: title, Text: text}

		templateData := TemplateData{Rec: &rec, TotalRecs: note_model.Count(dataDir), StatusMsg: msg}

		renderTemplate(w, "edit", &templateData)
	} else {
		rec, err := note_model.Find(dataDir, fileId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		rec.FileId = fileId
		rec.Title = title
		rec.Text = text
		rec.UpdatedAt = time.Now()

		err = note_model.Update(dataDir, rec, fileId)
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
