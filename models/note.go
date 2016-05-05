package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Note struct {
	FileId    string    `json:"-"y`
	Title     string    `json:"title"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Notes []Note

func (slice Notes) Len() int {
	return len(slice)
}

func (slice Notes) Less(i, j int) bool {
	return strings.ToLower(slice[i].Title) < strings.ToLower(slice[j].Title)
}

func (slice Notes) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

func (note *Note) FormattedCreatedAt() string {
	return note.CreatedAt.Format(time.RFC822)
}

func (note *Note) FormattedUpdatedAt() string {
	return note.UpdatedAt.Format(time.RFC822)
}

func (note *Note) FormattedText() template.HTML {
	return template.HTML(strings.Replace(note.Text, "\n", "</br>", -1))
}

func CreateNote(rec *Note) (string, error) {
	fileId, err := nextAvailableFileId()
	if err != nil {
		return "", err
	}

	err = writeRec(rec, fileId)
	if err != nil {
		return "", err
	}

	return fileId, nil
}

func DeleteNote(fileId string) error {
	filename := filePath(fileId)

	err := os.Remove(filename)
	if err != nil {
		return err
	}

	return nil
}

func FindNote(fileId string) (*Note, error) {
	var rec *Note

	filename := filePath(fileId)

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return rec, err
	}

	err = json.Unmarshal(data, &rec)
	if err != nil {
		return rec, err
	}

	rec.FileId = fileId

	return rec, nil
}

func FindNotes(searchString string) ([]Note, error) {
	var results Notes
	var rec Note
	var valuesFound int

	searchValues := strings.Split(strings.ToLower(searchString), " ")
	searchValuesCount := len(searchValues)

	for _, fileId := range fileIdsInDataDir() {
		filename := filePath(fileId)

		data, err := ioutil.ReadFile(filename)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(data, &rec)
		if err != nil {
			return nil, err
		}

		rec.FileId = fileId

		valuesFound = 0

		for _, searchValue := range searchValues {
			if searchValue == "" || strings.Contains(strings.ToLower(rec.Title), searchValue) || strings.Contains(strings.ToLower(rec.Text),
				searchValue) {

				valuesFound++
			} else {
				break
			}
		}

		if valuesFound == searchValuesCount {
			results = append(results, rec)
		}
	}

	sort.Sort(results)
	return results, nil
}

func NotesCount() int {
	return len(fileIdsInDataDir())
}

func UpdateNote(rec *Note, fileId string) error {
	if stringInSlice(fileId, fileIdsInDataDir()) {
		err := writeRec(rec, fileId)
		if err != nil {
			return err
		}
	} else {
		return errors.New("File ID not found")
	}
	return nil
}

//*****************************************************************************
// Private Methods
//*****************************************************************************

// fileIdsInDataDir returns all file ids in the data directory.
func fileIdsInDataDir() []string {
	var ids []string

	files, _ := ioutil.ReadDir("data")
	for _, file := range files {
		if !file.IsDir() {
			if path.Ext(file.Name()) == ".json" {
				ids = append(ids, file.Name()[:len(file.Name())-5])
			}
		}
	}

	return ids
}

// filePath returns a file name for a file id.
func filePath(fileId string) string {
	return fmt.Sprintf("data/%v.json", fileId)
}

// loadRec reads a json file into the supplied Note struct.
func loadRec(rec Note, fileId string) error {
	filename := filePath(fileId)

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, rec)

	return err
}

// nextAvailableFileId returns the next ascending available file id in a
// directory.
func nextAvailableFileId() (string, error) {
	var fileIds []int
	var nextFileId string

	for _, f := range fileIdsInDataDir() {
		fileId, err := strconv.Atoi(f)
		if err != nil {
			return "", err
		}

		fileIds = append(fileIds, fileId)
	}

	if len(fileIds) == 0 {
		nextFileId = "1"
	} else {
		sort.Ints(fileIds)
		lastFileId := fileIds[len(fileIds)-1]

		nextFileId = strconv.Itoa(lastFileId + 1)
	}

	return nextFileId, nil
}

func stringInSlice(s string, list []string) bool {
	for _, x := range list {
		if x == s {
			return true
		}
	}
	return false
}

func writeRec(rec *Note, fileId string) error {
	marshalledRec, err := json.Marshal(rec)

	if err != nil {
		return err
	}

	filename := filePath(fileId)

	err = ioutil.WriteFile(filename, marshalledRec, 0600)
	if err != nil {
		return err
	}

	return nil
}
