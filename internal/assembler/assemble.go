package assembler

import (
	"fmt"
	"github.com/ThreeToes/blogvm/internal/executable"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func AssembleFile(filePath string, includePaths []string) (*executable.LoadableFile, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return Assemble(f, includePaths)
}

func AssembleString(input string, includePaths []string) (*executable.LoadableFile, error) {
	reader := strings.NewReader(input)
	return Assemble(reader, includePaths)
}

func Assemble(input io.Reader, includePaths []string) (*executable.LoadableFile, error) {
	firstPassF, err := firstPass(input, 0x100)
	if err != nil {
		return nil, err
	}

	imports, err := assembleImports(firstPassF, includePaths)
	if err != nil {
		return nil, err
	}
	err = firstPassF.merge(imports)
	if err != nil {
		return nil, err
	}
	var errs []error
	for _, v := range firstPassF.symbolTable {
		switch v.symbolType {
		case MTDF:
			errs = append(errs, fmt.Errorf("duplicate symbol %q line %d: %s", v.label, v.relativeLineNumber, v.sourceLine))
		case INVALID:
			errs = append(errs, fmt.Errorf("invalid line %d: %q", v.relativeLineNumber, v.sourceLine))
		}
	}
	if err != nil {
		errString := &strings.Builder{}
		errString.WriteString("following errors found:")
		for _, err := range errs {
			errString.WriteString(fmt.Sprintf("\n\t* %v", err))
		}
		return nil, fmt.Errorf(errString.String())
	}
	return secondPass(firstPassF)
}

func assembleImports(records *firstPassFile, includePaths []string) (*firstPassFile, error) {
	ret := newFirstPassFile()
	for _, rec := range records.records {
		if rec.symbolType != IMPORT {
			continue
		}
		f, err := findFile(rec.sourceLine, includePaths)
		if err != nil {
			return nil, err
		}
		pass, err := firstPass(f, 0)
		f.Close()
		if err != nil {
			return nil, err
		}
		err = ret.merge(pass)
		if err != nil {
			return nil, err
		}
		recs, err := assembleImports(pass, includePaths)
		if err != nil {
			return nil, err
		}
		err = ret.merge(recs)
	}
	return ret, nil
}

func findFile(fileName string, searchPath []string) (*os.File, error) {
	fn := strings.TrimPrefix(fileName, "IMPORT ")
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
