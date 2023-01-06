package assembler

type symbolType uint8

const (
	REL symbolType = iota
	ABS
	MTDF
	INVALID
	IMPORT
	COMMENT
)

type symbol struct {
	symbolType         symbolType
	label              string
	relativeLineNumber uint32
	sourceLine         string
	assemblyLink       assemblable
}

func (s *symbol) assemble(symbolTable symbols) ([]uint32, error) {
	if s.assemblyLink != nil {
		return s.assemblyLink.assemble(s.sourceLine, symbolTable)
	}
	return nil, nil
}

type symbols map[string]*symbol

type relocatableFile struct {
	symbolTable symbols
	records     []*symbol
}

func (r *relocatableFile) merge(other *relocatableFile) error {
	newSymbolTable := symbols{}
	originalLength := len(r.records)
	newLength := originalLength + len(other.records)
	if newLength == 0 {
		return nil
	} else if newLength == originalLength {
		return nil
	}
	newRecordList := make([]*symbol, newLength)
	lineOffset := uint32(0)
	for idx, s := range r.records {
		newRecordList[idx] = s
		if s.label != "" {
			newSymbolTable[s.label] = s
		}
		if idx == originalLength-1 {
			lineOffset = s.relativeLineNumber + s.assemblyLink.calculateSize(s.sourceLine)
		}
	}
	offset := uint32(originalLength)
	for idx, rec := range other.records {
		newLineNum := rec.relativeLineNumber + lineOffset
		recCopy := &symbol{
			symbolType:         rec.symbolType,
			label:              rec.label,
			relativeLineNumber: newLineNum,
			sourceLine:         rec.sourceLine,
			assemblyLink:       rec.assemblyLink,
		}
		newRecordList[offset+uint32(idx)] = recCopy
		if rec.label != "" {
			if retrieved, ok := newSymbolTable[rec.label]; ok {
				retrieved.symbolType = MTDF
			} else {
				newSymbolTable[rec.label] = recCopy
			}
		}
	}

	r.records = newRecordList
	r.symbolTable = newSymbolTable
	return nil
}

func newRelocatableFile() *relocatableFile {
	return &relocatableFile{
		symbolTable: symbols{},
		records:     nil,
	}
}
