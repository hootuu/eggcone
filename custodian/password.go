package custodian

import (
	"crypto/aes"
	"crypto/sha256"
	"github.com/hootuu/gelato/crtpto/aesx"
	"github.com/hootuu/gelato/errors"
	"github.com/hootuu/gelato/logger"
	"go.uber.org/zap"
)

type Password []byte

func (pwd Password) Cover(src []byte) ([]byte, *errors.Error) {
	hash := sha256.Sum256(pwd)
	bytes, err := aesx.Encrypt(src, hash[:])
	if err != nil {
		logger.Error.Error("custodian.pwd.Cover err", zap.Error(err))
		return nil, err
	}
	return bytes, nil
}
func (pwd Password) Uncover(src []byte) ([]byte, *errors.Error) {
	if len(src)%aes.BlockSize != 0 {
		return nil, errors.Verify("Uncover: invalid ciphertext length")
	}
	hash := sha256.Sum256(pwd)
	bytes, err := aesx.Decrypt(src, hash[:])
	if err != nil {
		logger.Error.Error("custodian.pwd.Uncover err", zap.Error(err))
		return nil, err
	}
	return bytes, nil
}
