package filereader

type FileReader interface {
	OpenFile(filePath string) ([][]string, error)
}
