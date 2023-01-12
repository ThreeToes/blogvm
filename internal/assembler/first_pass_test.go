package assembler

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"strings"
	"testing"
)

func Test_firstPassLine(t *testing.T) {
	t.Run("test op code record no label", func(t *testing.T) {
		const line = "ADD 0x10 R1"
		rec, err := firstPassLine(10, line)
		if !assert.NoError(t, err) {
			return
		}
		assert.Equal(t, "", rec.label)
		assert.Equal(t, opcodeTable["ADD"], rec.assemblyLink)
		assert.Equal(t, uint32(10), rec.relativeLineNumber)
		assert.Equal(t, line, rec.sourceLine)
	})
	t.Run("test op code record with label", func(t *testing.T) {
		const line = "YEEHAW ADD 0x10 R1"
		rec, err := firstPassLine(10, line)
		if !assert.NoError(t, err) {
			return
		}
		assert.Equal(t, "YEEHAW", rec.label)
		assert.Equal(t, opcodeTable["ADD"], rec.assemblyLink)
		assert.Equal(t, uint32(10), rec.relativeLineNumber)
		assert.Equal(t, line, rec.sourceLine)
	})
	t.Run("test directive no label", func(t *testing.T) {
		const line = "WORD 0x7FFF"
		rec, err := firstPassLine(10, line)
		if !assert.NoError(t, err) {
			return
		}
		assert.Equal(t, "", rec.label)
		assert.Equal(t, directiveTable["WORD"], rec.assemblyLink)
		assert.Equal(t, uint32(10), rec.relativeLineNumber)
		assert.Equal(t, line, rec.sourceLine)
	})
	t.Run("test directive with label", func(t *testing.T) {
		const line = "IMPORTANTNUMBER WORD 0x7FFF"
		rec, err := firstPassLine(10, line)
		if !assert.NoError(t, err) {
			return
		}
		assert.Equal(t, "IMPORTANTNUMBER", rec.label)
		assert.Equal(t, directiveTable["WORD"], rec.assemblyLink)
		assert.Equal(t, uint32(10), rec.relativeLineNumber)
		assert.Equal(t, line, rec.sourceLine)
	})
	t.Run("test comment without label", func(t *testing.T) {
		const line = ";IMPORTANTNUMBER WORD 0x7FFF"
		rec, err := firstPassLine(10, line)
		if !assert.NoError(t, err) {
			return
		}
		assert.Equal(t, "", rec.label)
		assert.Equal(t, COMMENT, rec.symbolType)
		assert.Equal(t, uint32(10), rec.relativeLineNumber)
		assert.Equal(t, line, rec.sourceLine)
	})
	t.Run("test comment with label", func(t *testing.T) {
		const line = "SOMECOMMENT ;this is a comment"
		rec, err := firstPassLine(10, line)
		if !assert.NoError(t, err) {
			return
		}
		assert.Equal(t, "SOMECOMMENT", rec.label)
		assert.Equal(t, COMMENT, rec.symbolType)
		assert.Nil(t, rec.assemblyLink)
		assert.Equal(t, uint32(10), rec.relativeLineNumber)
		assert.Equal(t, line, rec.sourceLine)
	})
}

func Test_firstPass(t *testing.T) {
	type args struct {
		sourceFile io.Reader
	}
	tests := []struct {
		name     string
		args     args
		passFile *firstPassFile
		wantErr  assert.ErrorAssertionFunc
	}{
		{
			name: "single line",
			args: args{
				sourceFile: strings.NewReader("ADD R0 R1"),
			},
			passFile: &firstPassFile{
				symbolTable: symbols{},
				records: []*symbol{
					{
						symbolType:         REL,
						label:              "",
						relativeLineNumber: 0x100,
						sourceLine:         "ADD R0 R1",
						assemblyLink:       opcodeTable["ADD"],
					},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "multiple lines",
			args: args{
				sourceFile: strings.NewReader(`DEADBEEF WORD 0xDEADBEEF
READ DEADBEEF R0`),
			},
			passFile: &firstPassFile{
				symbolTable: symbols{
					"DEADBEEF": {
						symbolType:         REL,
						label:              "DEADBEEF",
						relativeLineNumber: 0x100,
						sourceLine:         "DEADBEEF WORD 0xDEADBEEF",
						assemblyLink:       directiveTable["WORD"],
					},
				},
				records: []*symbol{
					{
						symbolType:         REL,
						label:              "DEADBEEF",
						relativeLineNumber: 0x100,
						sourceLine:         "DEADBEEF WORD 0xDEADBEEF",
						assemblyLink:       directiveTable["WORD"],
					},
					{
						symbolType:         REL,
						label:              "",
						relativeLineNumber: 0x101,
						sourceLine:         "READ DEADBEEF R0",
						assemblyLink:       opcodeTable["READ"],
					},
				},
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got1, err := firstPass(tt.args.sourceFile, 0x100)
			if !tt.wantErr(t, err, fmt.Sprintf("firstPass(%v)", tt.args.sourceFile)) {
				return
			}
			assert.Equalf(t, tt.passFile, got1, "firstPass(%v)", tt.args.sourceFile)
		})
	}
}
