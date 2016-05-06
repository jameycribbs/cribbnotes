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

func CreateNote(dataDir string, rec *Note) (string, error) {
	fileId, err := nextAvailableFileId(dataDir)
	if err != nil {
		return "", err
	}

	err = writeRec(dataDir, rec, fileId)
	if err != nil {
		return "", err
	}

	return fileId, nil
}

func DeleteNote(dataDir string, fileId string) error {
	filename := filePath(dataDir, fileId)

	err := os.Remove(filename)
	if err != nil {
		return err
	}

	return nil
}

func FindNote(dataDir string, fileId string) (*Note, error) {
	var rec *Note

	filename := filePath(dataDir, fileId)

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

func FindNotes(dataDir string, searchString string) ([]Note, error) {
	var results Notes
	var rec Note
	var valuesFound int

	searchValues := strings.Split(strings.ToLower(searchString), " ")
	searchValuesCount := len(searchValues)

	for _, fileId := range fileIdsInDataDir(dataDir) {
		filename := filePath(dataDir, fileId)

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

func NotesCount(dataDir string) int {
	return len(fileIdsInDataDir(dataDir))
}

func UpdateNote(dataDir string, rec *Note, fileId string) error {
	if stringInSlice(fileId, fileIdsInDataDir(dataDir)) {
		err := writeRec(dataDir, rec, fileId)
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
func fileIdsInDataDir(dataDir string) []string {
	var ids []string

	files, _ := ioutil.ReadDir(dataDir)
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
func filePath(dataDir string, fileId string) string {
	return fmt.Sprintf("%v/%v.json", dataDir, fileId)
}

// loadRec reads a json file into the supplied Note struct.
func loadRec(dataDir string, rec Note, fileId string) error {
	filename := filePath(dataDir, fileId)

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, rec)

	return err
}

// nextAvailableFileId returns the next ascending available file id in a
// directory.
func nextAvailableFileId(dataDir string) (string, error) {
	var fileIds []int
	var nextFileId string

	for _, f := range fileIdsInDataDir(dataDir) {
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

func writeRec(dataDir string, rec *Note, fileId string) error {
	marshalledRec, err := json.Marshal(rec)

	if err != nil {
		return err
	}

	filename := filePath(dataDir, fileId)

	err = ioutil.WriteFile(filename, marshalledRec, 0600)
	if err != nil {
		return err
	}

	return nil
}
