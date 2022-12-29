package assembler

import (
	"fmt"
	"github.com/ThreeToes/blogvm/internal/executable"
	"io"
	"os"
	"path/filepath"
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
	symbolTable := symbolTableType{}
	firstPassF, lineNum, err := firstPass(input, 0x100, symbolTable)
	if err != nil {
		return nil, err
	}

	imports, _, err := assembleImports(firstPassF, lineNum, symbolTable)
	if err != nil {
		return nil, err
	}
	firstPassF = append(firstPassF, imports...)
	return secondPass(firstPassF, symbolTable)
}

func assembleImports(records []*record, lineNum uint32, symbolTable symbolTableType) ([]*record, uint32, error) {
	var ret []*record
	ln := lineNum
	wd, err := os.Getwd()
	if err != nil {
		return nil, ln, err
	}
	libPath := filepath.Join(wd, "lib")
	for _, rec := range records {
		if rec.recordType != importRecord {
			continue
		}
		f, err := findFile(rec.importFile, []string{libPath})
		if err != nil {
			return nil, ln, err
		}
		pass, lineNo, err := firstPass(f, ln, symbolTable)
		f.Close()
		if err != nil {
			return nil, ln, err
		}
		ln = lineNo
		ret = append(ret, pass...)
		recs, lineNo, err := assembleImports(pass, lineNo, symbolTable)
		if err != nil {
			return nil, 0, err
		}
		ln = lineNo
		ret = append(ret, recs...)
	}
	return ret, ln, nil
}

func findFile(fileName string, searchPath []string) (*os.File, error) {
	fn := fileName
	if !strings.HasSuffix(fn, ".bs") {
		fn = fmt.Sprintf("%s.bs", fn)
	}
	for _, p := range searchPath {
		search := filepath.Join(p, fn)
		if _, err := os.Stat(search); err == nil {
			return os.Open(search)
		}
	}
	return nil, fmt.Errorf("could not find file %s on search path", fn)
}
