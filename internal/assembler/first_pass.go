package assembler

import (
	"bufio"
	"io"
	"strings"
)

type recordType uint8

const (
	instructionRecord recordType = iota
	directiveRecord
	commentRecord
	invalidRecord
	duplicateSymbolRecord
)

type symbolTableType map[string]*record

type firstPassFile []*record

func firstPass(sourceFile io.Reader) (symbolTableType, firstPassFile, error) {
	ret := symbolTableType{}
	src := bufio.NewReader(sourceFile)
	// TODO: Find a better way to do this to allow better linking
	lineNum := uint32(0x100)
	var firstPass firstPassFile
	for {
		line, _, err := src.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, nil, err
		}
		rec, err := firstPassLine(lineNum, string(line))
		if rec.label != "" {
			_, ok := ret[rec.label]
			if ok {
				ret[rec.label] = &record{
					label:        rec.label,
					recordType:   duplicateSymbolRecord,
					line:         lineNum,
					directivePtr: nil,
					opCodePtr:    nil,
					source:       string(line),
				}
			} else {
				ret[rec.label] = rec
			}
		}
		firstPass = append(firstPass, rec)
		if rec.opCodePtr != nil {
			lineNum += rec.opCodePtr.calculateSize(string(line))
		} else if rec.directivePtr != nil {
			lineNum += rec.directivePtr.calculateSize(string(line))
		}
	}
	return ret, firstPass, nil
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
