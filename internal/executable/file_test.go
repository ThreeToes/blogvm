package executable

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"io"
	"reflect"
	"testing"
)

func uintToBytes(i uint32) []byte {
	b0 := byte(i >> 24 & 0xFF)
	b1 := byte(i >> 16 & 0xFF)
	b2 := byte(i >> 8 & 0xFF)
	b3 := byte(i & 0xFF)
	return []byte{b0, b1, b2, b3}
}

func uintsToBytes(is ...uint32) []byte {
	var byteSlice []byte

	for _, i := range is {
		byteSlice = append(byteSlice, uintToBytes(i)...)
	}

	return byteSlice
}

func Test_nextWord(t *testing.T) {
	type args struct {
		bs io.ByteReader
	}
	tests := []struct {
		name    string
		args    args
		want    uint32
		wantErr bool
	}{
		{
			name: "happy path",
			args: args{
				bs: bytes.NewReader([]byte{0x01, 0x02, 0x03, 0x04}),
			},
			want:    0x01020304,
			wantErr: false,
		},
		{
			name: "not enough bytes in reader",
			args: args{
				bs: bytes.NewReader([]byte{0x02, 0x03, 0x04}),
			},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := nextWord(tt.args.bs)
			if (err != nil) != tt.wantErr {
				t.Errorf("nextWord() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("nextWord() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_loadBlock(t *testing.T) {
	type args struct {
		bs io.ByteReader
	}
	tests := []struct {
		name    string
		args    args
		want    *MemoryBlock
		wantErr bool
	}{
		{
			name: "single word block",
			args: args{
				bs: bytes.NewReader(uintsToBytes(
					0x100,
					0x01,
					0x12345678,
				)),
			},
			want: &MemoryBlock{
				Address:   0x100,
				BlockSize: 0x01,
				Words:     []uint32{0x12345678},
			},
			wantErr: false,
		},
		{
			name: "multi word block",
			args: args{
				bs: bytes.NewReader(uintsToBytes(
					0x100,
					0x05,
					0x1234,
					0x5678,
					0x9876,
					0x5432,
					0x1234,
					0x4567,
				)),
			},
			want: &MemoryBlock{
				Address:   0x100,
				BlockSize: 0x05,
				Words: []uint32{
					0x1234,
					0x5678,
					0x9876,
					0x5432,
					0x1234,
				},
			},
			wantErr: false,
		},
		{
			name: "no address",
			args: args{
				bs: bytes.NewReader(uintsToBytes()),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "no size",
			args: args{
				bs: bytes.NewReader(uintsToBytes(0x100)),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "not enough blocks",
			args: args{
				bs: bytes.NewReader(uintsToBytes(
					0x100,
					0x02,
					0x1234,
				)),
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := loadBlock(tt.args.bs)
			if (err != nil) != tt.wantErr {
				t.Errorf("loadBlock() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("loadBlock() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_loadBlocks(t *testing.T) {
	type args struct {
		blockCount uint32
		bs         io.ByteReader
	}
	tests := []struct {
		name    string
		args    args
		want    []*MemoryBlock
		wantErr bool
	}{
		{
			name: "single block",
			args: args{
				blockCount: 1,
				bs: bytes.NewReader(uintsToBytes(
					0x100,
					0x01,
					0x12345678,
				)),
			},
			want: []*MemoryBlock{
				{
					Address:   0x100,
					BlockSize: 0x01,
					Words:     []uint32{0x12345678},
				},
			},
			wantErr: false,
		},
		{
			name: "multi block",
			args: args{
				blockCount: 3,
				bs: bytes.NewReader(uintsToBytes(
					0x100,
					0x01,
					0x12345678,
					0x101,
					0x02,
					0x12345678,
					0x12345678,
					0x103,
					0x01,
					0x12345678,
				)),
			},
			want: []*MemoryBlock{
				{
					Address:   0x100,
					BlockSize: 0x01,
					Words:     []uint32{0x12345678},
				},
				{
					Address:   0x101,
					BlockSize: 0x02,
					Words:     []uint32{0x12345678, 0x12345678},
				},
				{
					Address:   0x103,
					BlockSize: 0x01,
					Words:     []uint32{0x12345678},
				},
			},
			wantErr: false,
		},
		{
			name: "bad block",
			args: args{
				blockCount: 1,
				bs: bytes.NewReader(uintsToBytes(
					0x100,
					0x01,
				)),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "no blocks",
			args: args{
				blockCount: 0,
				bs: bytes.NewReader(uintsToBytes(
					0x100,
					0x01,
				)),
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := loadBlocks(tt.args.blockCount, tt.args.bs)
			if (err != nil) != tt.wantErr {
				t.Errorf("loadBlocks() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("loadBlocks() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLoad(t *testing.T) {
	type args struct {
		bs io.ByteReader
	}
	tests := []struct {
		name    string
		args    args
		want    *LoadableFile
		wantErr bool
	}{
		{
			name: "single block",
			args: args{
				bs: bytes.NewReader(uintsToBytes(0x01, 0x00, 0x100, 0x01, 0x1234)),
			},
			want: &LoadableFile{
				BlockCount: 1,
				Flags:      0,
				Blocks: []*MemoryBlock{
					{
						Address:   0x100,
						BlockSize: 1,
						Words:     []uint32{0x1234},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "multi block",
			args: args{
				bs: bytes.NewReader(uintsToBytes(0x02, 0x00, 0x100, 0x01, 0x1234, 0x200, 0x02, 0x1234, 0x5678, 0x09)),
			},
			want: &LoadableFile{
				BlockCount: 2,
				Flags:      0,
				Blocks: []*MemoryBlock{
					{
						Address:   0x100,
						BlockSize: 1,
						Words:     []uint32{0x1234},
					},
					{
						Address:   0x200,
						BlockSize: 2,
						Words: []uint32{
							0x1234,
							0x5678,
						},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Load(tt.args.bs)
			if (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Load() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_wordToBytes(t *testing.T) {
	assert.Equal(t, []byte{0x01, 0x02, 0x03, 0x04}, wordToBytes(0x01020304))
}

func TestLoadableFile_Save(t *testing.T) {
	t.Run("successful save 1 block", func(t *testing.T) {
		buf := bytes.NewBuffer([]byte{})
		file := &LoadableFile{
			BlockCount: 1,
			Flags:      0,
			Blocks: []*MemoryBlock{
				{
					Address:   0x100,
					BlockSize: 2,
					Words: []uint32{
						0x1234,
						0x5678,
					},
				},
			},
		}
		err := file.Save(buf)
		if !assert.NoError(t, err) {
			return
		}
		assert.Equal(t, uintsToBytes(0x01, 0x0, 0x100, 0x02, 0x1234, 0x5678), buf.Bytes())
	})
	t.Run("successful save multi blocks", func(t *testing.T) {

		buf := bytes.NewBuffer([]byte{})
		file := &LoadableFile{
			BlockCount: 3,
			Flags:      0,
			Blocks: []*MemoryBlock{
				{
					Address:   0x100,
					BlockSize: 2,
					Words: []uint32{
						0x1234,
						0x5678,
					},
				},
				{
					Address:   0x200,
					BlockSize: 1,
					Words: []uint32{
						0x1234,
					},
				},
				{
					Address:   0x300,
					BlockSize: 3,
					Words: []uint32{
						0x1234,
						0x5678,
						0x90,
					},
				},
			},
		}
		err := file.Save(buf)
		if !assert.NoError(t, err) {
			return
		}
		assert.Equal(t, uintsToBytes(0x03, 0x0, 0x100, 0x02, 0x1234, 0x5678, 0x200, 0x01, 0x1234, 0x300, 0x03, 0x1234, 0x5678, 0x90), buf.Bytes())
	})
}
