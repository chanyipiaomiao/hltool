package hltool

import (
	"fmt"
	"os"
	"path"

	"github.com/boltdb/bolt"
)

var (
	db *bolt.DB
)

// BoltDB DB操作
type BoltDB struct {
	DBPath    string // 数据库路径名称
	TableName string // 表名
}

// NewBoltDB 初始化数据库对象
func NewBoltDB(dbPath, tableName string) (*BoltDB, error) {
	if dbPath == "" {
		return nil, fmt.Errorf("database path required")
	}
	dirname := path.Dir(dbPath)
	if !IsExist(dirname) {
		err := os.MkdirAll(dirname, os.ModePerm)
		if err != nil {
			return nil, fmt.Errorf("create dir(%s) error: %s", dirname, err)
		}
	}
	return &BoltDB{DBPath: dbPath, TableName: tableName}, nil
}

// table 获取表 表不存在则创建
func (btb *BoltDB) table() error {
	var err error
	db, err = bolt.Open(btb.DBPath, 0600, nil)
	if err != nil {
		return fmt.Errorf("open db error: %s", err)
	}

	return db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(btb.TableName))
		if err != nil {
			return fmt.Errorf("create table error: %s", err)
		}
		return nil
	})
}

// Set 设置值
// kv 键值对
func (btb *BoltDB) Set(kv map[string][]byte) error {
	err := btb.table()
	if err != nil {
		return err
	}
	defer db.Close()

	return db.Batch(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(btb.TableName))
		var err error
		for k, v := range kv {
			err = b.Put([]byte(k), v)
			if err != nil {
				return err
			}
		}
		return err
	})
}

// Get 根据键名数组获取各自的值
// keys 键名数组
func (btb *BoltDB) Get(keys []string) (map[string][]byte, error) {
	err := btb.table()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	values := make(map[string][]byte)
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(btb.TableName))
		for _, k := range keys {
			result := b.Get([]byte(k))
			tmp := make([]byte, len(result))
			copy(tmp, result)
			values[k] = tmp
		}
		return nil
	})
	return values, err
}

// Delete 删除键值
func (btb *BoltDB) Delete(keys []string) error {
	err := btb.table()
	if err != nil {
		return err
	}
	defer db.Close()

	return db.Batch(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(btb.TableName))
		var err error
		for _, key := range keys {
			err = b.Delete([]byte(key))
			if err != nil {
				return err
			}
		}
		return err
	})
}

// Backup 备份数据库文件
func (btb *BoltDB) Backup(filepath string) error {
	db, err := bolt.Open(btb.DBPath, 0600, nil)
	if err != nil {
		return fmt.Errorf("open db error: %s", err)
	}

	err = db.View(func(tx *bolt.Tx) error {
		err := tx.CopyFile(filepath, 0644)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return err
	}

	return db.Close()
}
