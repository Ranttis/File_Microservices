package data

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
)

type File struct {
	ID 		int 	`json:"id"`
	Name 	string 	`json:"name"`
	Type 	string  `json:"type"`
}

var n int
var createFilePath = "/"
var deleteFilePath = "./"

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
	if len(fileList) != 0{
		for _,  f := range fileList{
			n = n + 1
			contentString := "ID: " + strconv.Itoa(f.ID) + "\n" + "name: "+ f.Name+ "\n"+ "type: " + f.Type
			content := []byte(contentString)
			testFile, err := os.Create(f.Name+f.Type)
			_, _ = testFile.Write(content)
			if err !=nil{
				println(err)
			}
		}
	}
	return fileList
}

func AddFile(f *File)  {
	f.ID = getNextID()
	fileList = append(fileList, f)

	contentString := "ID: " + strconv.Itoa(f.ID) + "\n" + "name: "+ f.Name+ "\n"+ "type: " + f.Type
	content := []byte(contentString)
	testFile, err := os.Create(f.Name+f.Type)
	if err != nil{
		println(err)
	}
	_, _ = testFile.Write(content)
}

func DeleteFile(id int, f *File) error{
	ff, pos, err := findFile(id)
	if err != nil {
		return err
	}
	var pathToDelete = deleteFilePath+ff.Name+ff.Type
	println(pathToDelete)
	//os.Remove("C:/Users/roope/go/src/File_Microservices/Main/data/"+ff.Name+ff.Type)
	//println("C:/Users/roope/go/src/File_Microservices/Main/data/"+ff.Name+ff.Type)
	var er = os.Remove(pathToDelete)
	if er != nil{
		return err
	}

	fileList[pos] = f
	fileList = RemoveIndex(fileList,id)

	return nil
}

func RemoveIndex(f []*File, index int) []*File{
	f = append(f[index:], f[index+1:]...)
	return f
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