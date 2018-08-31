package worker

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_NewWorker(t *testing.T) {
	var worker = NewWorker(nil, nil)
	assert.NotNil(t, worker)
	assert.Nil(t, worker.block)
	assert.Nil(t, worker.chain)
}
