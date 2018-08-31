package worker

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_NewWorker(t *testing.T) {
	assert := assert.New(t)
	var worker = NewWorker(nil, nil)
	assert.NotNil(worker)
	assert.Nil(worker.block)
	assert.Nil(worker.chain)
}
