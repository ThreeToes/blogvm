package assembler

type record struct {
	label        string
	recordType   recordType
	line         uint32
	directivePtr *directive
	opCodePtr    *opCode
	source       string
	importFile   string
}

func (r *record) assemble(symbolTable symbols) ([]uint32, error) {
	if r.opCodePtr != nil {
		return r.opCodePtr.assemble(r.source, symbolTable)
	} else if r.directivePtr != nil {
		return r.directivePtr.assemble(r.source, symbolTable)
	}
	return nil, nil
}
