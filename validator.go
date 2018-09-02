package validator

import (
	"bytes"
	"fmt"
	"github.com/DSiSc/blockchain"
	"github.com/DSiSc/craft/types"
	"github.com/DSiSc/txpool/log"
	"github.com/DSiSc/validator/common"
	"github.com/DSiSc/validator/tools/account"
	"github.com/DSiSc/validator/tools/signature"
	"github.com/DSiSc/validator/worker"
)

type Validator struct {
	Account *account.Account
}

func NewValidator(account *account.Account) *Validator {
	return &Validator{
		Account: account,
	}
}

func (self *Validator) ValidateBlock(block *types.Block) (*types.Header, error) {
	chain, err := blockchain.NewLatestStateBlockChain()
	if err != nil {
		return nil, fmt.Errorf("New Latest State BlockChain Error:%s.", err)
	}
	// new worker to verify the block
	work := worker.NewWorker(chain, block)
	err = work.VerifyBlock()
	if nil != err {
		return nil, err
	}
	// sign the block
	hash := common.BlockHash(block)
	sign, ok := signature.Sign(self.Account, hash[:])
	if ok != nil {
		return nil, fmt.Errorf("[Signature],Sign error:%s.", ok)
	}
	notSigned := true
	for _, value := range block.Header.SigData {
		if bytes.Equal(value, sign) {
			notSigned = false
			log.Warn("Duplicate sign")
		}
	}
	if notSigned {
		block.Header.SigData = append(block.Header.SigData, sign)
	}

	return block.Header, nil
}
