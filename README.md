disks3cache provides an implementation of httpcache.Cache that
caches locally on disk using diskv and also in Amazon S3.

If the object isn't on local disk, it checks S3. If it's found in S3,
then the object is then cached locally on disk going forward.

S3 provides a permanent and shared storage mechanism for multiple
caching instances.


Use something like this:
```
package main

import (
    "ioutil"

    "github.com/victortrac/disks3cache"
)

func main() {
    cacheDir := ioutil.TempDir("", "myTempDir")
    cacheSize := 512 // in megabytes
    s3CacheURL := "s3://s3-us-west-2.amazonaws.com/my-bucket"

    c = disks3cache.New(cacheDir, cacheSize, s3CacheURL)

    // do something with your new 2 layer cache

}
```

See https://github.com/gregjones/httpcache for other httpcache backends.
