package storage

import (
	"haru/engine/utils"
	"os"
	"testing"
)

func TestOpenOrCreateKv(t *testing.T) {
	db, err := openKVStorage("kv_test")
	utils.Expect(t, "<nil>", err)
	db.Close()

	db, err = openKVStorage("kv_test")
	utils.Expect(t, "<nil>", err)
	err = db.Set([]byte("key1"), []byte("value1"))
	utils.Expect(t, "<nil>", err)

	buffer := make([]byte, 100)
	buffer, err = db.Get([]byte("key1"))
	utils.Expect(t, "<nil>", err)
	utils.Expect(t, "value1", string(buffer))

	walFile := db.WALName()
	db.Close()
	os.Remove(walFile)
	os.Remove("kv_test")
}
