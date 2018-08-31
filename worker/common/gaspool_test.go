package common

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_NewGasPool(t *testing.T) {
	assert := assert.New(t)
	gasInit := GasPool(100)
	gasInitPoint := &gasInit
	assert.NotNil(gasInitPoint)

	assert.Equal(uint64(100), gasInitPoint.Gas())

	addGas := gasInitPoint.AddGas(uint64(20))
	assert.Equal(GasPool(120), *addGas)
	assert.Equal(GasPool(120), *gasInitPoint)

	err := gasInitPoint.SubGas(uint64(20))
	assert.Nil(err)
	assert.Equal(GasPool(100), *gasInitPoint)

	assert.Equal("100", gasInitPoint.String())
}
