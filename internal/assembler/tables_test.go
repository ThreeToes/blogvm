package assembler

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_opCode_instructionMask(t *testing.T) {
	type fields struct {
		mnemonic string
		opcode   uint8
		hasI1    bool
		hasI2    bool
	}
	tests := []struct {
		name   string
		fields fields
		want   uint32
	}{
		{
			name: "test arbitrary",
			fields: fields{
				mnemonic: "ADD",
				opcode:   0x04,
				hasI1:    true,
				hasI2:    true,
			},
			want: 0x04000000,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &opCode{
				mnemonic: tt.fields.mnemonic,
				opcode:   tt.fields.opcode,
				hasI1:    tt.fields.hasI1,
				hasI2:    tt.fields.hasI2,
			}
			if got := o.instructionMask(); got != tt.want {
				t.Errorf("instructionMask() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_directive_size_calcs(t *testing.T) {
	// string terminates with a null char (0x00)
	t.Run("test WORD directive", func(t *testing.T) {
		assert.Equal(t, uint32(1), directiveTable["WORD"].sizeCalc("WORD 0x1234"))
	})
	t.Run("test STRING directive no label", func(t *testing.T) {
		assert.Equal(t, uint32(len("hello, world!")+1), directiveTable["STRING"].sizeCalc("STRING hello, world!"))
	})
	t.Run("test STRING directive with label", func(t *testing.T) {
		assert.Equal(t, uint32(len("hello, world!")+1), directiveTable["STRING"].sizeCalc("ABC STRING hello, world!"))
	})
}

func Test_opCode_assemble(t *testing.T) {
	type args struct {
		sourceLine  string
		symbolTable symbolTableType
	}
	tests := []struct {
		name    string
		opCode  *opCode
		args    args
		want    []uint32
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:   "HALT",
			opCode: opcodeTable["HALT"],
			args: args{
				sourceLine:  "HALT",
				symbolTable: symbolTableType{},
			},
			want:    []uint32{0x00000000},
			wantErr: assert.NoError,
		},
		// Add has two inputs, other similar commands should be fine
		{
			name:   "ADD no label",
			opCode: opcodeTable["ADD"],
			args: args{
				sourceLine:  "ADD R1 0x50",
				symbolTable: symbolTableType{},
			},
			want:    []uint32{0x041F0050},
			wantErr: assert.NoError,
		},
		{
			name:   "ADD with label",
			opCode: opcodeTable["ADD"],
			args: args{
				sourceLine:  "MAGIC ADD 0x50 R3",
				symbolTable: symbolTableType{},
			},
			want:    []uint32{0x04F30050},
			wantErr: assert.NoError,
		},
		{
			name:   "ADD with label and decimal value",
			opCode: opcodeTable["ADD"],
			args: args{
				sourceLine:  "MAGIC ADD 15 R3",
				symbolTable: symbolTableType{},
			},
			want:    []uint32{0x04F3000F},
			wantErr: assert.NoError,
		},
		{
			name:   "ADD with label and octal value",
			opCode: opcodeTable["ADD"],
			args: args{
				sourceLine:  "MAGIC ADD 012 R3",
				symbolTable: symbolTableType{},
			},
			want:    []uint32{0x04F3000A},
			wantErr: assert.NoError,
		},
		{
			name:   "ADD with label and binary value",
			opCode: opcodeTable["ADD"],
			args: args{
				sourceLine:  "MAGIC ADD 0b101 R3",
				symbolTable: symbolTableType{},
			},
			want:    []uint32{0x04F30005},
			wantErr: assert.NoError,
		},
		{
			name:   "ADD with unsupported symbol",
			opCode: opcodeTable["ADD"],
			args: args{
				sourceLine: "MAGIC ADD BIGDOG R3",
				symbolTable: symbolTableType{
					"BIGDOG": {
						label:        "BIGDOG",
						recordType:   instructionRecord,
						line:         0x123,
						directivePtr: nil,
						opCodePtr:    opcodeTable["READ"],
						source:       "READ 0x1234",
					},
				},
			},
			want:    nil,
			wantErr: assert.Error,
		},
		// Try JMP for symbol resolution in I1
		{
			name:   "JMP with no symbols",
			opCode: opcodeTable["JMP"],
			args: args{
				sourceLine:  "JMP 0x1234",
				symbolTable: symbolTableType{},
			},
			want:    []uint32{0x0CF01234},
			wantErr: assert.NoError,
		},
		{
			name:   "JMP with symbols",
			opCode: opcodeTable["JMP"],
			args: args{
				sourceLine: "JMP BIGPIG",
				symbolTable: symbolTableType{
					"BIGPIG": {
						label:        "BIGPIG",
						recordType:   instructionRecord,
						line:         0x1234,
						directivePtr: nil,
						opCodePtr:    opcodeTable["READ"],
						source:       "READ BIGDOG",
					},
				},
			},
			want:    []uint32{0x0CF01234},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := tt.opCode
			got, err := o.assemble(tt.args.sourceLine, tt.args.symbolTable)
			if !tt.wantErr(t, err, fmt.Sprintf("assemble(%v)", tt.args.sourceLine)) {
				return
			}
			assert.Equalf(t, tt.want, got, "assemble(%v)", tt.args.sourceLine)
		})
	}
}

func Test_directive_assemble(t *testing.T) {
	type fields struct {
		mnemonic     string
		sizeCalc     func(sourceLine string) uint32
		assembleFunc func(sourceLine string, symbolTable symbolTableType) ([]uint32, error)
	}
	type args struct {
		sourceLine  string
		symbolTable symbolTableType
	}
	tests := []struct {
		name      string
		directive *directive
		args      args
		want      []uint32
		wantErr   assert.ErrorAssertionFunc
	}{
		{
			name:      "word no label",
			directive: directiveTable["WORD"],
			args: args{
				sourceLine:  "WORD 0x1234",
				symbolTable: nil,
			},
			want:    []uint32{0x1234},
			wantErr: assert.NoError,
		},
		{
			name:      "word with label",
			directive: directiveTable["WORD"],
			args: args{
				sourceLine:  "BIGDOG WORD 0x1234",
				symbolTable: nil,
			},
			want:    []uint32{0x1234},
			wantErr: assert.NoError,
		},
		{
			name:      "string no label",
			directive: directiveTable["STRING"],
			args: args{
				sourceLine:  "STRING 0x1234",
				symbolTable: nil,
			},
			want:    []uint32{0x30, 0x78, 0x31, 0x32, 0x33, 0x34, 0x0},
			wantErr: assert.NoError,
		},
		{
			name:      "string with label",
			directive: directiveTable["STRING"],
			args: args{
				sourceLine:  "BIGDOG STRING hello world",
				symbolTable: nil,
			},
			want:    []uint32{0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x20, 0x77, 0x6f, 0x72, 0x6c, 0x64, 0x0},
			wantErr: assert.NoError,
		},
		{
			name:      "address valid symbol",
			directive: directiveTable["ADDRESS"],
			args: args{
				sourceLine: "ADDRESS HELLO R1",
				symbolTable: symbolTableType{
					"HELLO": {
						label:        "HELLO",
						recordType:   directiveRecord,
						line:         0x123,
						directivePtr: directiveTable["WORD"],
						opCodePtr:    nil,
						source:       "WORD 0x33",
					},
				},
			},
			want:    []uint32{0x03F10123},
			wantErr: assert.NoError,
		},
		{
			name:      "address invalid symbol",
			directive: directiveTable["ADDRESS"],
			args: args{
				sourceLine: "ADDRESS GOODBYE R1",
				symbolTable: symbolTableType{
					"HELLO": {
						label:        "HELLO",
						recordType:   directiveRecord,
						line:         0x123,
						directivePtr: directiveTable["WORD"],
						opCodePtr:    nil,
						source:       "WORD 0x33",
					},
				},
			},
			want:    nil,
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := tt.directive
			got, err := d.assemble(tt.args.sourceLine, tt.args.symbolTable)
			if !tt.wantErr(t, err, fmt.Sprintf("assemble(%v, %v)", tt.args.sourceLine, tt.args.symbolTable)) {
				return
			}
			assert.Equalf(t, tt.want, got, "assemble(%v, %v)", tt.args.sourceLine, tt.args.symbolTable)
		})
	}
}
