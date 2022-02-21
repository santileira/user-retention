package filereader

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
)

func TestOpenFileShouldReturnAnErrorWhenFilePathIsEmpty(t *testing.T) {
	// Initialization
	filePath := ""

	fileReader := NewFileReaderImpl()
	// Operation
	output, err := fileReader.OpenFile(filePath)

	// Validation
	assert.Nil(t, output)
	assert.EqualValues(t, fmt.Errorf("error opening the file , err: open : no such file or directory"), err)
}

func TestOpenFileShouldReturnTheContentOfTheFile(t *testing.T) {
	// Initialization
	file, _ := ioutil.TempFile("", "test.csv")
	_, _ = file.Write([]byte("column1,column2,column3\ncolumn11,column22,column33"))
	fileReader := NewFileReaderImpl()

	// Operation
	output, err := fileReader.OpenFile(file.Name())

	// Validation
	assert.EqualValues(t, [][]string{{"column1", "column2", "column3"}, {"column11", "column22", "column33"}}, output)
	assert.Nil(t, err)

	_ = os.Remove(file.Name())
}
