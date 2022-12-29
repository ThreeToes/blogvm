package assembler

import (
	"github.com/ThreeToes/blogvm/internal/executable"
	"io"
	"os"
	"strings"
)

func AssembleFile(filePath string) (*executable.LoadableFile, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return Assemble(f)
}

func AssembleString(input string) (*executable.LoadableFile, error) {
	reader := strings.NewReader(input)
	return Assemble(reader)
}

func Assemble(input io.Reader) (*executable.LoadableFile, error) {
	symbolTable, firstPassF, err := firstPass(input)
	if err != nil {
		return nil, err
	}
	return secondPass(firstPassF, symbolTable)
}
