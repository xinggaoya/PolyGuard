package db

import (
	"PolyGuard/consts"
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/boltdb/bolt"
	"os"
)

/**
  @author: XingGao
  @date: 2023/8/22
**/

type Database struct {
	db *bolt.DB
}

func NewDatabase(dbFilePath string) (*Database, error) {
	var dirPath = consts.DBPath
	// 如果文件夹不存在则创建
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		err = os.Mkdir(dirPath, os.ModePerm)
		if err != nil {
			return nil, err
		}
	}
	db, err := bolt.Open(dirPath+"/"+dbFilePath, 0600, nil)
	if err != nil {
		return nil, err
	}

	return &Database{db: db}, nil
}

func (d *Database) Close() error {
	return d.db.Close()
}

func Set(key string, value interface{}) error {

	d, err := NewDatabase(consts.DBName)
	if err != nil {
		return err
	}
	defer d.Close()
	return d.db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(consts.DBBucketName))
		if err != nil {
			return err
		}

		var encodedValue bytes.Buffer
		enc := gob.NewEncoder(&encodedValue)
		if err := enc.Encode(value); err != nil {
			return err
		}

		return bucket.Put([]byte(key), encodedValue.Bytes())
	})
}

func Get(key string, targetValue interface{}) error {
	d, err := NewDatabase(consts.DBName)
	if err != nil {
		return err
	}
	defer d.Close()
	return d.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(consts.DBBucketName))
		if bucket == nil {
			return fmt.Errorf("bucket not found")
		}

		encodedValue := bucket.Get([]byte(key))
		buffer := encodedValue
		dec := gob.NewDecoder(bytes.NewReader(buffer))
		if err := dec.Decode(targetValue); err != nil {
			return err
		}

		return nil
	})
}
