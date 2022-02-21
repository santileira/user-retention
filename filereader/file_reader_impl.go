package filereader

import (
	"encoding/csv"
	"fmt"
	"os"
)

type FileReaderImpl struct{}

func NewFileReaderImpl() *FileReaderImpl {
	return &FileReaderImpl{}
}

func (f *FileReaderImpl) OpenFile(filePath string) ([][]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error opening the file %s, err: %s", filePath, err.Error())
	}
	defer file.Close()

	csvReader := csv.NewReader(file)
	rows, err := csvReader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("error reading the file %s, err: %s", filePath, err.Error())
	}

	return rows, nil
}
