package core

import (
	"bytes"
	"github.com/boltdb/bolt"
	"strconv"
	"sync"
	"errors"
)

var ins *DbHandler
var once sync.Once

const (
	blocksBucket     = "blocks"
	BLOCKCHAIN_INDEX = "l"
	BLOCK_INDEX      = "b"
)

type DbHandler struct {
	Path         string
	BlocksBucket string
	Db           *bolt.DB
}

func GetDbInstance(path string) *DbHandler {
	once.Do(func() {
		db, err := bolt.Open(path, 0600, nil)
		if err != nil {
			logger.Error(err)
		}
		ins = &DbHandler{Path: path, BlocksBucket: blocksBucket, Db: db}
	})
	return ins
}

func (dh *DbHandler) BlockChainIsExist() bool {
	isExist := false
	err := dh.Db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		if b == nil {
			isExist = false
		} else {
			isExist = true
		}
		return nil
	})
	if err != nil {
		logger.Error(err)
		return isExist
	}
	return isExist
}

func (dh *DbHandler) PutBlock(block *Block) error {
	err := dh.Db.Update(func(tx *bolt.Tx) error {
		var err error
		b := tx.Bucket([]byte(blocksBucket))
		// if b is not exist ,create one
		if b == nil {
			b, err = tx.CreateBucket([]byte(blocksBucket))
			if err != nil {
				logger.Error(err)
				return err
			}
		}
		res := b.Get(bytes.Join([][]byte{[]byte(BLOCK_INDEX), []byte(strconv.FormatInt(block.Header.Height, 10))}, nil))
		if res != nil {
			return errors.New("Block has been exist.")
		}

		err = b.Put(bytes.Join([][]byte{[]byte(BLOCK_INDEX), []byte(strconv.FormatInt(block.Header.Height, 10))}, nil), block.Serialize())
		if err != nil {
			logger.Error(err)
			return err
		}
		err = b.Put([]byte(BLOCKCHAIN_INDEX), []byte(strconv.FormatInt(block.Header.Height, 10)))
		if err != nil {
			logger.Error(err)
			return err
		}
		return nil
	})
	return err
}

func (dh *DbHandler) ResumeBlock() (int64, error) {
	if dh.BlockChainIsExist() == false {
		return -1, nil
	}
	var index int64
	err := dh.Db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		res := b.Get([]byte(BLOCKCHAIN_INDEX))
		index, _ = strconv.ParseInt(string(res), 10, 64)
		return nil
	})
	if err != nil {
		return -1, err
	}
	return index, nil
}

func (dh *DbHandler) GetBlock(height int64) *Block {
	if dh.BlockChainIsExist() == false {
		return nil
	}
	var blockBytes []byte
	err := dh.Db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		blockBytes = b.Get(bytes.Join([][]byte{[]byte(BLOCK_INDEX), []byte(strconv.FormatInt(height, 10))}, nil))
		return nil
	})
	if err != nil {
		return nil
	}
	if blockBytes != nil {
		block := DeserializeBlock(blockBytes)
		return block
	}
	return nil
}
