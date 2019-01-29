package validator

import (
	"bytes"
	"fmt"
	"github.com/DSiSc/blockchain"
	"github.com/DSiSc/craft/log"
	"github.com/DSiSc/craft/types"
	"github.com/DSiSc/validator/tools/account"
	"github.com/DSiSc/validator/tools/signature"
	"github.com/DSiSc/validator/worker"
)

type Validator struct {
	Account  *account.Account
	Receipts types.Receipts
}

func NewValidator(account *account.Account) *Validator {
	return &Validator{
		Account: account,
	}
}

func (self *Validator) ValidateBlock(block *types.Block, signVerify bool) (*types.Header, error) {
	chain, err := blockchain.NewLatestStateBlockChain()
	if err != nil {
		return nil, fmt.Errorf("get NewLatestStateBlockChain error:%s ", err)
	}
	// new worker to verify the block
	work := worker.NewWorker(chain, block, signVerify)
	err = work.VerifyBlock()
	if nil != err {
		return nil, err
	}
	self.Receipts = work.GetReceipts()
	// sign the block
	sign, ok := signature.Sign(self.Account, block.Header.MixDigest[:])
	if ok != nil {
		return nil, fmt.Errorf("sign block failed with error: %s", ok)
	}

	notSigned := true
	for _, value := range block.Header.SigData {
		if bytes.Equal(value, sign) {
			notSigned = false
			log.Warn("Duplicate sign.")
			break
		}
	}
	if notSigned {
		block.Header.SigData = append(block.Header.SigData, sign)
		log.Debug("Validator add sign %x to block %d.", sign, block.Header.Height)
	}

	return block.Header, nil
}
