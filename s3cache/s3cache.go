// Package s3cache provides an implementation of httpcache.Cache that stores and
// retrieves data using Amazon S3.
// Adapted from "sourcegraph.com/sourcegraph/s3cache" to use aws SDK
package s3cache

import (
	"crypto/md5"
	"bytes"
	"encoding/hex"
	"io"
	"io/ioutil"
	"regexp"

	"github.com/going/toolkit/log"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// Cache objects store and retrieve data using Amazon S3.
type Cache struct {
	// s3.S3 client
	S3 *s3.S3

	// name of the bucket
	Bucket string
}

func (c *Cache) Get(key string) (resp []byte, ok bool) {
	params := &s3.GetObjectInput{
		Bucket:		aws.String(c.Bucket),
		Key:		aws.String(CacheKeyToObjectKey(key)),
	}

	object, err := c.S3.GetObject(params)

	if err != nil {
		log.Debugf("Error getting s3 object: %v", err)
		return []byte{}, false
	}

	data, err := ioutil.ReadAll(object.Body)
	if err != nil {
		log.Errorf("Error reading s3 object: %v", err.Error())
		return []byte{}, false
	}
	return data, err == nil
}

func (c *Cache) Set(key string, resp []byte) {
	params := &s3.PutObjectInput{
		Bucket:		aws.String(c.Bucket),
		Key:		aws.String(CacheKeyToObjectKey(key)),
		Body:		bytes.NewReader(resp),
	}

	_, err := c.S3.PutObject(params)

	if err != nil {
		log.Errorf("Error Putting s3 object: %v", err.Error())
		return
	}
}

func (c *Cache) Delete(key string) {
	params := &s3.DeleteObjectInput{
		Bucket:		aws.String(c.Bucket),
		Key:		aws.String(CacheKeyToObjectKey(key)),
	}

	_, err := c.S3.DeleteObject(params)

	if err != nil {
		log.Errorf("Error deleting s3 object: %v", err.Error())
		return
	}
}

func CacheKeyToObjectKey(key string) string {
	h := md5.New()
	io.WriteString(h, key)
	return hex.EncodeToString(h.Sum(nil))
}

func New(bucketURL string) *Cache {
	// Parse bucket string into region and bucketname
	re := regexp.MustCompile(`//(s3-)?([\w\-)]+)..*/([\w\-]+$)`).FindStringSubmatch(bucketURL)
	region := re[2]
	bucket := re[3]
	maxRetries := 10
	session := session.New(&aws.Config{Region: aws.String(region), MaxRetries: &maxRetries})
	log.Printf("s3cache: S3 Connection - Region: %v, Bucket: %v", region, bucket)
	return &Cache{
		S3: s3.New(session),
		Bucket: bucket,
	}
}
