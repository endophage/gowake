package main

import (
	"log"
	"os"
	"os/user"
	"path/filepath"

	"github.com/boltdb/bolt"
)

var (
	db         *bolt.DB
	bucketName = []byte("locations")
)

func init() {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err.Error())
	}
	dbDir := filepath.Join(usr.HomeDir, ".wake")
	err = os.MkdirAll(dbDir, 0755)
	if err != nil {
		log.Fatal(err.Error())
	}
	dbPath := filepath.Join(dbDir, "wake.db")
	db, err = bolt.Open(dbPath, 0755, nil)
	if err != nil {
		log.Fatal(err.Error())
	}
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(bucketName)
		return err
	})
	if err != nil {
		log.Fatal(err.Error())
	}
}

func lookupName(name string) ([]byte, error) {
	var macHex string
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)
		macHex = string(b.Get([]byte(name)))
		return nil
	})
	if err != nil {
		return nil, err
	}
	return parseMAC(macHex)
}

func saveName(name, mac string) error {
	var err error
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)
		err = b.Put([]byte(name), []byte(mac))
		return err
	})
	return err
}

func listNames() (map[string]string, error) {
	res := make(map[string]string)
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)
		return b.ForEach(func(k, v []byte) error {
			res[string(k)] = string(v)
			return nil
		})
	})
	return res, err
}
