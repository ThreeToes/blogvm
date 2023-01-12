package assembler

import (
	"fmt"
	"github.com/ThreeToes/blogvm/internal/executable"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_secondPass(t *testing.T) {
	type args struct {
		firstPass *firstPassFile
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
				firstPass: &firstPassFile{
					symbolTable: symbols{},
					records: []*symbol{
						{
							symbolType:         REL,
							label:              "",
							relativeLineNumber: 0,
							sourceLine:         "ADD R1 R2",
							assemblyLink:       opcodeTable["ADD"],
						},
					},
				},
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
			got, err := secondPass(tt.args.firstPass)
			if !tt.wantErr(t, err, fmt.Sprintf("secondPass(%v)", tt.args.firstPass)) {
				return
			}
			assert.Equalf(t, tt.want, got, "secondPass(%v)", tt.args.firstPass)
		})
	}
}
