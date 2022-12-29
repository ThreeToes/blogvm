package assembler

import (
	"fmt"
	"github.com/ThreeToes/blogvm/internal/executable"
	"github.com/ThreeToes/blogvm/internal/machine"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"runtime"
	"testing"
)

func TestAssembleFile(t *testing.T) {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)
	type args struct {
		filePath string
	}
	tests := []struct {
		name    string
		args    args
		want    *executable.LoadableFile
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "load simple add",
			args: args{
				filePath: filepath.Join(basepath, "test_files", "simple_add.bs"),
			},
			want: &executable.LoadableFile{
				BlockCount: 0x01,
				Flags:      0x00,
				Blocks: []*executable.MemoryBlock{
					{
						Address:   0x100,
						BlockSize: 0x06,
						Words: []uint32{
							0x03F00005,
							0x03F10005,
							0x04010000,
							0x021F0105,
							0x00000000,
							0x00000000,
						},
					},
				},
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := AssembleFile(tt.args.filePath)
			if !tt.wantErr(t, err, fmt.Sprintf("AssembleFile(%v)", tt.args.filePath)) {
				return
			}
			assert.Equalf(t, tt.want, got, "AssembleFile(%v)", tt.args.filePath)
		})
	}
}

func TestRunScenarios(t *testing.T) {
	_, b, _, _ := runtime.Caller(0)
	testingFilePath := filepath.Join(filepath.Dir(b), "test_files")
	t.Run("simple add", func(t *testing.T) {
		addFile := filepath.Join(testingFilePath, "simple_add.bs")
		assembledFile, err := AssembleFile(addFile)
		if !assert.NoError(t, err) {
			return
		}
		mem := machine.NewMemory()
		bus := machine.NewBus(mem)
		registers := machine.NewRegisterBank()
		cpu := machine.NewCPU(registers, bus)
		err = mem.Load(assembledFile)
		if !assert.NoError(t, err) {
			return
		}
		sr, err := registers.GetRegister(machine.SR)
		if !assert.NoError(t, err) {
			return
		}
		for sr.Value&machine.STATUS_HALT == 0 {
			err = cpu.Tick()
			if !assert.NoError(t, err) {
				return
			}
		}
		writtenMem, err := mem.Read(0x0105)
		if !assert.NoError(t, err) {
			return
		}

		assert.Equal(t, uint32(0x0A), writtenMem)
	})
}
