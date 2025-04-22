package custodian

import (
	"bytes"
	"github.com/hootuu/eggcone/database/pgx"
	"github.com/hootuu/eggcone/eggdbx"
	"github.com/hootuu/gelato/crtpto/aesx"
	"github.com/hootuu/gelato/crtpto/ed25519x"
	"github.com/hootuu/gelato/errors"
	"github.com/hootuu/gelato/idx"
	"github.com/hootuu/gelato/logger"
	"go.uber.org/zap"
	"sync"
	"time"
)

const (
	bufferSize     = 100
	prefixLen      = 20
	reloadInterval = int64(10 * 24 * 60 * 60)
	syncInterval   = int64(1 * 60 * 60)
)

type key struct {
	idx   []byte
	key   []byte
	usage uint64
}

var (
	gBuf            [bufferSize]*key
	gFill           = 0
	gLstSelectedIdx = 0
	gLstLoadTime    = int64(0)
	gLstSyncTime    = int64(0)
	gLock           sync.RWMutex
)

func EncryptWithPwd(src []byte, pwdBytes []byte) ([]byte, *errors.Error) {
	pwd := Password(pwdBytes)
	coverSrc, err := pwd.Cover(src)
	if err != nil {
		return nil, err
	}
	return Encrypt(coverSrc)
}

func DecryptWithPwd(src []byte, pwdBytes []byte) ([]byte, *errors.Error) {
	pwd := Password(pwdBytes)
	uncoverSrc, err := pwd.Uncover(src)
	if err != nil {
		return nil, err
	}
	return Decrypt(uncoverSrc)
}

func Encrypt(src []byte) ([]byte, *errors.Error) {
	selKey, err := doSelect()
	if err != nil {
		return nil, err
	}
	encBytes, err := aesx.Encrypt(src, selKey.key)
	if err != nil {
		return nil, err
	}
	fullBytes := append(selKey.idx, encBytes...)
	return fullBytes, nil
}

func Decrypt(src []byte) ([]byte, *errors.Error) {
	if len(src) < prefixLen {
		return nil, errors.E("src len to short")
	}
	idxBytes := src[:prefixLen]
	encBytes := src[prefixLen:]

	var willUsePriKey []byte
	for _, wrap := range gBuf {
		if bytes.Equal(wrap.idx, idxBytes) {
			willUsePriKey = wrap.key
		}
	}
	if willUsePriKey == nil {
		kcIdx := string(idxBytes)
		custodianM, err := pgx.MustGet[Model](eggdbx.EggPgDB(), "idx = ?", kcIdx)
		if err != nil {
			logger.Error.Error("custodian.Decrypt : db error", zap.Error(err))
			return nil, err
		}
		willUsePriKey = custodianM.PrivateKey[:32]
	}
	decBytes, err := aesx.Decrypt(encBytes, willUsePriKey)
	if err != nil {
		return nil, err
	}
	return decBytes, nil
}

func genNew() (*Model, *errors.Error) {
	_, privateKey, err := ed25519x.NewRandom()
	if err != nil {
		logger.Error.Error("custodian.genNew: ed25519x.NewRandom() err", zap.Error(err))
		return nil, err
	}
	newKc := &Model{
		Idx:        idx.New(),
		PrivateKey: privateKey,
		Usage:      0,
		Available:  true,
	}
	return newKc, nil
}

func doSelect() (*key, *errors.Error) {
	gLock.Lock()
	defer gLock.Unlock()
	if gFill == 0 || time.Now().Unix()-gLstLoadTime > reloadInterval {
		err := reload()
		if err != nil {
			return nil, err
		}
	}
	if gLstSelectedIdx == len(gBuf)-1 {
		gLstSelectedIdx = 0
	} else {
		gLstSelectedIdx++
	}
	curKw := gBuf[gLstSelectedIdx]
	curKw.usage++
	trySync()
	return curKw, nil
}

func trySync() {
	if gLstSyncTime == 0 {
		gLstSyncTime = time.Now().Unix()
	}
	if time.Now().Unix()-gLstSyncTime < syncInterval {
		return
	}
	defer logger.Elapse("custodian.trySync", logger.Logger, func() []zap.Field {
		return []zap.Field{zap.Int64("gLstSyncTime", gLstSyncTime)}
	})()
	for _, item := range gBuf {
		err := updateUsage(string(item.idx), item.usage)
		if err != nil {
			logger.Error.Error("[ignore] custodian.trySync err", zap.Error(err))
			continue
		}
	}
	gLstSyncTime = time.Now().Unix()
}

func doLoad() *errors.Error {
	defer logger.Elapse("custodian.doLoad", logger.Logger, func() []zap.Field {
		return nil
	})()
	gFill = 0
	arr, err := loadAvailable(bufferSize)
	if err != nil {
		return err
	}
	for i, item := range arr {
		gBuf[i] = &key{
			[]byte(item.Idx),
			item.PrivateKey[:32],
			item.Usage,
		}
		gFill++
	}
	gLstLoadTime = time.Now().Unix()
	return nil
}

func reload() *errors.Error {
	err := doLoad()
	if err != nil {
		return err
	}
	if gFill < bufferSize {
		err = batchGen(bufferSize - gFill)
		if err != nil {
			return err
		}
		err = doLoad()
		if err != nil {
			return err
		}
	}
	return nil
}

func batchGen(size int) *errors.Error {
	defer logger.Elapse("custodian.batchGen", logger.Logger, func() []zap.Field {
		return []zap.Field{zap.Int("size", size)}
	})()
	var arr []*Model
	for i := 0; i < size; i++ {
		newKc, err := genNew()
		if err != nil {
			return err
		}
		arr = append(arr, newKc)
	}
	err := multiCreate(arr)
	if err != nil {
		logger.Error.Error("custodian.batchGen: multiCreate err", zap.Error(err))
		return err
	}
	return nil
}
