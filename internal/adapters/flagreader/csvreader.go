package flagreader

import (
	"encoding/csv"
	"github.com/realfabecker/kevin/internal/core/ports"
	"io"
	"os"
)

type csvFlagReader struct{}

func NewCsvFlagReader() ports.FlagListReader {
	return &csvFlagReader{}
}

func (c *csvFlagReader) Read(filePath string) ([]map[string]string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	csvReader := csv.NewReader(f)
	header, err := csvReader.Read()

	var args = make([]map[string]string, 0)
	for {
		rec, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		var m = make(map[string]string)
		for i, v := range rec {
			m[header[i]] = v
		}
		args = append(args, m)
	}
	return args, nil
}
