package data

import (
	"encoding/json"
	"fmt"
	"io"
)

type File struct {
	ID 		int
	Name 	string
	Type 	string
}

func (f *File) FromJSON(r io.Reader) error{
	e := json.NewDecoder(r)
	return e.Decode(f)
}

type Files []*File

func (f *Files) ToJSON(w io.Writer) error{
	e := json.NewEncoder(w)
	return e.Encode(f)
}

func GetFiles() Files{
	return fileList
}

func AddFile(f *File)  {
	f.ID = getNextID()
	fileList = append(fileList, f)
}

func DeleteFile(id int, f *File) error{
	_, pos, err := findFile(id)
	if err != nil {
		return err
	}
	f.ID = id
	fileList[pos] = f

	return nil
}

func getNextID() int{
	fp := fileList[len(fileList)-1]
	return fp.ID + 1
}

func findFile(id int) (*File, int, error)  {
	for i, p := range fileList{
		if p.ID == id{
			return p,i,nil
		}
	}
	return nil, -1, ErrProductNotFound
}
var ErrProductNotFound = fmt.Errorf("Product not found")

var fileList = []*File{
	&File{
		ID: 1,
		Name: "testitiedosto",
		Type: ".txt",
	},

	&File{
		ID: 2,
		Name: "testitiedosto2",
		Type: ".txt",
	},


}