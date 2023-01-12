package assembler

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_relocatableFile_merge(t *testing.T) {
	t.Run("empty everything", func(t *testing.T) {
		r1 := newFirstPassFile()
		r2 := newFirstPassFile()

		expected := newFirstPassFile()
		err := r1.merge(r2)
		assert.NoError(t, err)
		assert.Equal(t, expected, r1)
	})
	t.Run("empty mergee", func(t *testing.T) {
		testSymbol := &symbol{
			symbolType:         REL,
			label:              "TEST",
			relativeLineNumber: 0,
			sourceLine:         "TEST COPY 0x01 R0",
			assemblyLink:       opcodeTable["COPY"],
		}
		r1 := &firstPassFile{
			symbolTable: symbols{
				"TEST": testSymbol,
			},
			records: []*symbol{
				testSymbol,
			},
		}
		r2 := newFirstPassFile()
		err := r1.merge(r2)
		assert.NoError(t, err)
		assert.Equal(t, symbols{"TEST": testSymbol}, r1.symbolTable)
		assert.Equal(t, []*symbol{testSymbol}, r1.records)
		assert.Equal(t, newFirstPassFile(), r2)
	})
	t.Run("empty merger", func(t *testing.T) {
		testSymbol := &symbol{
			symbolType:         REL,
			label:              "TEST",
			relativeLineNumber: 0,
			sourceLine:         "TEST COPY 0x01 R0",
			assemblyLink:       opcodeTable["COPY"],
		}
		r2 := &firstPassFile{
			symbolTable: symbols{
				"TEST": testSymbol,
			},
			records: []*symbol{
				testSymbol,
			},
		}
		r1 := newFirstPassFile()
		err := r1.merge(r2)
		assert.NoError(t, err)
		assert.Equal(t, symbols{"TEST": testSymbol}, r1.symbolTable)
		assert.Equal(t, []*symbol{testSymbol}, r1.records)
		assert.Equal(t, r1, r2)
	})
	t.Run("offset data", func(t *testing.T) {
		s1 := &symbol{
			symbolType:         REL,
			label:              "",
			relativeLineNumber: 0,
			sourceLine:         "HALT",
			assemblyLink:       opcodeTable["HALT"],
		}
		s2 := &symbol{
			symbolType:         REL,
			label:              "GREETING",
			relativeLineNumber: 1,
			sourceLine:         "GREETING STRING hello",
			assemblyLink:       directiveTable["STRING"],
		}
		s3 := &symbol{
			symbolType:         REL,
			label:              "MAGICNUMBER",
			relativeLineNumber: 0,
			sourceLine:         "MAGICNUMBER WORD 0xDEADBEEF",
			assemblyLink:       directiveTable["WORD"],
		}

		expectedS3 := &symbol{
			symbolType:         REL,
			label:              "MAGICNUMBER",
			relativeLineNumber: 7,
			sourceLine:         "MAGICNUMBER WORD 0xDEADBEEF",
			assemblyLink:       directiveTable["WORD"],
		}

		r1 := &firstPassFile{
			symbolTable: symbols{
				"GREETING": s2,
			},
			records: []*symbol{s1, s2},
		}
		r2 := &firstPassFile{
			symbolTable: symbols{
				"MAGICNUMBER": s3,
			},
			records: []*symbol{s3},
		}

		expectedR := &firstPassFile{
			symbolTable: symbols{
				"GREETING":    s2,
				"MAGICNUMBER": expectedS3,
			},
			records: []*symbol{s1, s2, expectedS3},
		}

		err := r1.merge(r2)
		assert.NoError(t, err)
		assert.Equal(t, expectedR.symbolTable, r1.symbolTable)
		assert.Equal(t, expectedR.records, r1.records)
		assert.Equal(t, expectedR, r1)
	})
}
