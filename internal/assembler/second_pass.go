package assembler

import "github.com/ThreeToes/blogvm/internal/executable"

func secondPass(firstPass firstPassFile, symbolTable symbolTableType) (*executable.LoadableFile, error) {
	ret := &executable.LoadableFile{
		BlockCount: 0x1,
		Flags:      0,
		Blocks:     nil,
	}
	b := &executable.MemoryBlock{
		Address:   0x100,
		BlockSize: 0,
		Words:     nil,
	}
	for _, rec := range firstPass {
		words, err := rec.assemble(symbolTable)
		if err != nil {
			return nil, err
		}
		b.Words = append(b.Words, words...)
	}
	ret.Blocks = append(ret.Blocks, b)
	b.BlockSize = uint32(len(b.Words))
	return ret, nil
}
