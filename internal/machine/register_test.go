package machine

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRegisterBank_GetRegister(t *testing.T) {
	t.Run("Get known register", func(t *testing.T) {
		rb := NewRegisterBank()
		r0, err := rb.GetRegister(R0)
		if !assert.NoError(t, err) {
			t.FailNow()
		}
		assert.True(t, r0 == rb.registerMap[R0])
	})
	t.Run("Get invalid register", func(t *testing.T) {
		rb := NewRegisterBank()
		r, err := rb.GetRegister(0xFF)
		assert.Nil(t, r)
		assert.Error(t, err)
	})
}
