package assembler

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

func firstPass(sourceFile io.Reader, lineNum uint32) (*firstPassFile, error) {
	ln := lineNum
	src := bufio.NewReader(sourceFile)
	reloc := newFirstPassFile()
	for {
		line, _, err := src.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		rec, err := firstPassLine(ln, string(line))
		if rec.label != "" {
			_, ok := reloc.symbolTable[rec.label]
			if ok {
				reloc.symbolTable[rec.label] = &symbol{
					label:              rec.label,
					symbolType:         MTDF,
					relativeLineNumber: ln,
					sourceLine:         string(line),
				}
			} else {
				reloc.symbolTable[rec.label] = rec
			}
		}
		reloc.records = append(reloc.records, rec)
		if rec.assemblyLink != nil {
			ln += rec.assemblyLink.calculateSize(string(line))
		}
	}
	return reloc, nil
}

func firstPassLine(lineNo uint32, line string) (*symbol, error) {
	cols := strings.Split(line, " ")
	if len(cols) == 0 {
		return nil, nil
	}
	op, ok := opcodeTable[cols[0]]
	if ok {
		return &symbol{
			symbolType:         REL,
			label:              "",
			relativeLineNumber: lineNo,
			sourceLine:         line,
			assemblyLink:       op,
		}, nil
	}
	dir, ok := directiveTable[cols[0]]
	if ok {
		return &symbol{
			symbolType:         REL,
			label:              "",
			relativeLineNumber: lineNo,
			sourceLine:         line,
			assemblyLink:       dir,
		}, nil
	}
	if cols[0] == "IMPORT" {
		if len(cols) < 2 {
			return nil, fmt.Errorf("import statement is not valid")
		}
		return &symbol{
			symbolType:         IMPORT,
			label:              "",
			relativeLineNumber: lineNo,
			sourceLine:         line,
			assemblyLink:       nil,
		}, nil
	}
	if line[0] == ';' {
		return &symbol{
			symbolType:         COMMENT,
			label:              "",
			relativeLineNumber: lineNo,
			sourceLine:         line,
			assemblyLink:       nil,
		}, nil
	}

	if len(cols) == 1 {
		return &symbol{
			symbolType:         INVALID,
			label:              "",
			relativeLineNumber: lineNo,
			sourceLine:         line,
			assemblyLink:       nil,
		}, nil
	}

	label := cols[0]
	op, ok = opcodeTable[cols[1]]
	if ok {
		return &symbol{
			symbolType:         REL,
			label:              label,
			relativeLineNumber: lineNo,
			sourceLine:         line,
			assemblyLink:       op,
		}, nil
	}
	dir, ok = directiveTable[cols[1]]
	if ok {
		return &symbol{
			symbolType:         REL,
			label:              label,
			relativeLineNumber: lineNo,
			sourceLine:         line,
			assemblyLink:       dir,
		}, nil
	}
	if cols[1][0] == ';' {
		return &symbol{
			symbolType:         COMMENT,
			label:              label,
			relativeLineNumber: lineNo,
			sourceLine:         line,
			assemblyLink:       nil,
		}, nil
	}

	return &symbol{
		symbolType:         INVALID,
		label:              "",
		relativeLineNumber: lineNo,
		sourceLine:         line,
		assemblyLink:       nil,
	}, nil
}
