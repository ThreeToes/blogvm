package assembler

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

type recordType uint8

const (
	instructionRecord recordType = iota
	directiveRecord
	commentRecord
	invalidRecord
	importRecord
	duplicateSymbolRecord
)

type symbolTableType map[string]*record

type firstPassFile []*record

func firstPass(sourceFile io.Reader, lineNum uint32, symbolTable symbolTableType) (firstPassFile, uint32, error) {
	ln := lineNum
	src := bufio.NewReader(sourceFile)
	var firstPass firstPassFile
	for {
		line, _, err := src.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, ln, err
		}
		rec, err := firstPassLine(ln, string(line))
		if rec.label != "" {
			_, ok := symbolTable[rec.label]
			if ok {
				symbolTable[rec.label] = &record{
					label:        rec.label,
					recordType:   duplicateSymbolRecord,
					line:         ln,
					directivePtr: nil,
					opCodePtr:    nil,
					source:       string(line),
				}
			} else {
				symbolTable[rec.label] = rec
			}
		}
		firstPass = append(firstPass, rec)
		if rec.opCodePtr != nil {
			ln += rec.opCodePtr.calculateSize(string(line))
		} else if rec.directivePtr != nil {
			ln += rec.directivePtr.calculateSize(string(line))
		}
	}
	return firstPass, ln, nil
}

func firstPassLine(lineNo uint32, line string) (*record, error) {
	cols := strings.Split(line, " ")
	if len(cols) == 0 {
		return nil, nil
	}
	op, ok := opcodeTable[cols[0]]
	if ok {
		return &record{
			label:        "",
			recordType:   instructionRecord,
			line:         lineNo,
			directivePtr: nil,
			opCodePtr:    op,
			source:       line,
		}, nil
	}
	dir, ok := directiveTable[cols[0]]
	if ok {
		return &record{
			label:        "",
			recordType:   directiveRecord,
			line:         lineNo,
			directivePtr: dir,
			opCodePtr:    nil,
			source:       line,
		}, nil
	}
	if cols[0] == "IMPORT" {
		if len(cols) < 2 {
			return nil, fmt.Errorf("import statement is not valid")
		}
		return &record{
			label:        "",
			recordType:   importRecord,
			line:         lineNo,
			directivePtr: nil,
			opCodePtr:    nil,
			source:       "",
			importFile:   cols[1],
		}, nil
	}
	if line[0] == ';' {
		return &record{
			label:        "",
			recordType:   commentRecord,
			line:         lineNo,
			directivePtr: nil,
			opCodePtr:    nil,
			source:       line,
		}, nil
	}

	if len(cols) == 1 {
		return &record{
			label:        "",
			recordType:   invalidRecord,
			line:         lineNo,
			directivePtr: nil,
			opCodePtr:    nil,
			source:       line,
		}, nil
	}

	label := cols[0]
	op, ok = opcodeTable[cols[1]]
	if ok {
		return &record{
			label:        label,
			recordType:   instructionRecord,
			line:         lineNo,
			directivePtr: nil,
			opCodePtr:    op,
			source:       line,
		}, nil
	}
	dir, ok = directiveTable[cols[1]]
	if ok {
		return &record{
			label:        label,
			recordType:   directiveRecord,
			line:         lineNo,
			directivePtr: dir,
			opCodePtr:    nil,
			source:       line,
		}, nil
	}
	if cols[1][0] == ';' {
		return &record{
			label:        label,
			recordType:   commentRecord,
			line:         lineNo,
			directivePtr: nil,
			opCodePtr:    nil,
			source:       line,
		}, nil
	}

	return &record{
		label:        "",
		recordType:   invalidRecord,
		line:         lineNo,
		directivePtr: nil,
		opCodePtr:    nil,
		source:       line,
	}, nil
}
