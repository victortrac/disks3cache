package disks3cache

import (
	"io/ioutil"
	"log"

	"github.com/gregjones/httpcache/diskcache"
	"github.com/peterbourgon/diskv"
	"github.com/victortrac/disks3cache/s3cache"
)


type Cache struct {
	disk     *diskcache.Cache
	s3       *s3cache.Cache
}

func (c *Cache) Get(key string) (resp []byte, ok bool) {
	// Check disk first
	resp, ok = c.disk.Get(key)
	if ok == true {
		log.Printf("Found %v in disk cache", key)
		return resp, ok
	}
	resp, ok = c.s3.Get(key)
	if ok == true {
		log.Printf("Found %v in s3 cache", key)
		go c.disk.Set(key, resp)
		return resp, ok
	}
	log.Printf("%v not found in cache", key)
	return []byte{}, ok
}

func (c *Cache) Set(key string, resp []byte) {
	log.Printf("Setting key %v on disk", key)
	go c.disk.Set(key, resp)
	log.Printf("Setting key %v in s3", key)
	go c.s3.Set(key, resp)
}

func (c *Cache) Delete(key string) {
	log.Printf("Deleting key %v", key)
	go c.disk.Delete(key)
	go c.s3.Delete(key)
}

func New(cacheDir string, cacheSize uint64, bucketURL string) *Cache {
	if cacheDir == "" {
		cacheDir, _ = ioutil.TempDir("", "disks3cache")
	}
	dv := diskv.New(diskv.Options{
		BasePath:       cacheDir,
		CacheSizeMax:   cacheSize * 1024 * 1024,
	})
	return &Cache{
		disk:   diskcache.NewWithDiskv(dv),
		s3:     s3cache.New(bucketURL),
	}
}