package storage

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"
	"strings"
	"neftac/storage/internal/db"
)

var Root = "/data/buckets"

func BucketPath(b string) string { return filepath.Join(Root, b) }
func ObjectPath(b, k string) string { return filepath.Join(BucketPath(b), k) }

func ensureBucket(b string) error {
	var bucket db.Bucket
	if err := db.DB.Where("name = ?", b).First(&bucket).Error; err != nil {
		return db.DB.Create(&db.Bucket{Name: b}).Error
	}
	return nil
}

func WriteObject(b, k string, r io.Reader) (string, int64, error) {
	if err := ensureBucket(b); err != nil { return "", 0, err }
	path := ObjectPath(b, k)
	os.MkdirAll(filepath.Dir(path), 0755)
	f, _ := os.Create(path)
	defer f.Close()

	h := md5.New()
	tee := io.TeeReader(r, h)
	n, _ := io.Copy(f, tee)
	etag := hex.EncodeToString(h.Sum(nil))

	var bucket db.Bucket
	db.DB.Where("name = ?", b).First(&bucket)
	db.DB.Create(&db.Object{BucketID: bucket.ID, Key: k, Size: n, ETag: etag})

	return etag, n, nil
}

func ReadObject(b, k string) (io.ReadCloser, error) { return os.Open(ObjectPath(b, k)) }
func DeleteObject(b, k string) error {
	os.Remove(ObjectPath(b, k))
	var bucket db.Bucket
	db.DB.Where("name = ?", b).First(&bucket)
	db.DB.Where("bucket_id = ? AND key = ?", bucket.ID, k).Delete(&db.Object{})
	return nil
}

func CopyObject(srcB, srcK, dstB, dstK string) error {
	f, _ := ReadObject(srcB, srcK)
	defer f.Close()
	_, _, err := WriteObject(dstB, dstK, f)
	return err
}

func ListObjects(b, p string) ([]string, error) {
	var keys []string
	root := BucketPath(b)
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() { return err }
		rel, _ := filepath.Rel(root, path)
		if p == "" || strings.HasPrefix(rel, p) { keys = append(keys, rel) }
		return nil
	})
	return keys, nil
}
