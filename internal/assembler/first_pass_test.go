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
		assert.Equal(t, opcodeTable["ADD"], rec.opCodePtr)
		assert.Nil(t, rec.directivePtr)
		assert.Equal(t, uint32(10), rec.line)
		assert.Equal(t, line, rec.source)
	})
	t.Run("test op code record with label", func(t *testing.T) {
		const line = "YEEHAW ADD 0x10 R1"
		rec, err := firstPassLine(10, line)
		if !assert.NoError(t, err) {
			return
		}
		assert.Equal(t, "YEEHAW", rec.label)
		assert.Equal(t, opcodeTable["ADD"], rec.opCodePtr)
		assert.Nil(t, rec.directivePtr)
		assert.Equal(t, uint32(10), rec.line)
		assert.Equal(t, line, rec.source)
	})
	t.Run("test directive no label", func(t *testing.T) {
		const line = "WORD 0x7FFF"
		rec, err := firstPassLine(10, line)
		if !assert.NoError(t, err) {
			return
		}
		assert.Equal(t, "", rec.label)
		assert.Equal(t, directiveTable["WORD"], rec.directivePtr)
		assert.Nil(t, rec.opCodePtr)
		assert.Equal(t, uint32(10), rec.line)
		assert.Equal(t, line, rec.source)
	})
	t.Run("test directive with label", func(t *testing.T) {
		const line = "IMPORTANTNUMBER WORD 0x7FFF"
		rec, err := firstPassLine(10, line)
		if !assert.NoError(t, err) {
			return
		}
		assert.Equal(t, "IMPORTANTNUMBER", rec.label)
		assert.Equal(t, directiveTable["WORD"], rec.directivePtr)
		assert.Nil(t, rec.opCodePtr)
		assert.Equal(t, uint32(10), rec.line)
		assert.Equal(t, line, rec.source)
	})
	t.Run("test comment without label", func(t *testing.T) {
		const line = ";IMPORTANTNUMBER WORD 0x7FFF"
		rec, err := firstPassLine(10, line)
		if !assert.NoError(t, err) {
			return
		}
		assert.Equal(t, "", rec.label)
		assert.Equal(t, commentRecord, rec.recordType)
		assert.Nil(t, rec.directivePtr)
		assert.Nil(t, rec.opCodePtr)
		assert.Equal(t, uint32(10), rec.line)
		assert.Equal(t, line, rec.source)
	})
	t.Run("test comment with label", func(t *testing.T) {
		const line = "SOMECOMMENT ;this is a comment"
		rec, err := firstPassLine(10, line)
		if !assert.NoError(t, err) {
			return
		}
		assert.Equal(t, "SOMECOMMENT", rec.label)
		assert.Equal(t, commentRecord, rec.recordType)
		assert.Nil(t, rec.directivePtr)
		assert.Nil(t, rec.opCodePtr)
		assert.Equal(t, uint32(10), rec.line)
		assert.Equal(t, line, rec.source)
	})
}

func Test_firstPass(t *testing.T) {
	type args struct {
		sourceFile io.Reader
	}
	tests := []struct {
		name        string
		args        args
		symbolTable symbolTableType
		passFile    firstPassFile
		wantErr     assert.ErrorAssertionFunc
	}{
		{
			name: "single line",
			args: args{
				sourceFile: strings.NewReader("ADD R0 R1"),
			},
			symbolTable: symbolTableType{},
			passFile: firstPassFile{
				{
					label:        "",
					recordType:   instructionRecord,
					line:         0,
					directivePtr: nil,
					opCodePtr:    opcodeTable["ADD"],
					source:       "ADD R0 R1",
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
			symbolTable: symbolTableType{
				"DEADBEEF": {
					label:        "DEADBEEF",
					recordType:   directiveRecord,
					line:         0,
					directivePtr: directiveTable["WORD"],
					opCodePtr:    nil,
					source:       "DEADBEEF WORD 0xDEADBEEF",
				},
			},
			passFile: firstPassFile{
				{
					label:        "DEADBEEF",
					recordType:   directiveRecord,
					line:         0,
					directivePtr: directiveTable["WORD"],
					opCodePtr:    nil,
					source:       "DEADBEEF WORD 0xDEADBEEF",
				},
				{
					label:        "",
					recordType:   instructionRecord,
					line:         1,
					directivePtr: nil,
					opCodePtr:    opcodeTable["READ"],
					source:       "READ DEADBEEF R0",
				},
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := firstPass(tt.args.sourceFile)
			if !tt.wantErr(t, err, fmt.Sprintf("firstPass(%v)", tt.args.sourceFile)) {
				return
			}
			assert.Equalf(t, tt.symbolTable, got, "firstPass(%v)", tt.args.sourceFile)
			assert.Equalf(t, tt.passFile, got1, "firstPass(%v)", tt.args.sourceFile)
		})
	}
}
