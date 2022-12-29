package assembler

import (
	"fmt"
	"github.com/ThreeToes/blogvm/internal/executable"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_secondPass(t *testing.T) {
	type args struct {
		firstPass   firstPassFile
		symbolTable symbolTableType
	}
	tests := []struct {
		name    string
		args    args
		want    *executable.LoadableFile
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "no symbols",
			args: args{
				firstPass: firstPassFile{
					{
						label:        "",
						recordType:   instructionRecord,
						line:         0,
						directivePtr: nil,
						opCodePtr:    opcodeTable["ADD"],
						source:       "ADD R1 R2",
					},
				},
				symbolTable: nil,
			},
			want: &executable.LoadableFile{
				BlockCount: 0x01,
				Flags:      0x00,
				Blocks: []*executable.MemoryBlock{
					{
						Address:   0x100,
						BlockSize: 0x01,
						Words: []uint32{
							0x04120000,
						},
					},
				},
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := secondPass(tt.args.firstPass, tt.args.symbolTable)
			if !tt.wantErr(t, err, fmt.Sprintf("secondPass(%v, %v)", tt.args.firstPass, tt.args.symbolTable)) {
				return
			}
			assert.Equalf(t, tt.want, got, "secondPass(%v, %v)", tt.args.firstPass, tt.args.symbolTable)
		})
	}
}
