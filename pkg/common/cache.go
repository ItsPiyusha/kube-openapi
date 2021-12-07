/*
Copyright 2021 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package common

import (
	"crypto/sha512"
	"fmt"
	"sync"
)

func computeETag(data []byte) string {
	if data == nil {
		return ""
	}
	return fmt.Sprintf("\"%X\"", sha512.Sum512(data))
}

type HandlerCache struct {
	BuildCache func() ([]byte, error)
	once       sync.Once
	bytes      []byte
	etag       string
	err        error
}

func (c *HandlerCache) Get() ([]byte, string, error) {
	c.once.Do(func() {
		bytes, err := c.BuildCache()
		// if there is an error updating the cache, there can be situations where
		// c.bytes contains a valid value (carried over from the previous update)
		// but c.err is also not nil; the cache user is expected to check for this
		c.err = err
		if c.err == nil {
			// don't override previous spec if we had an error
			c.bytes = bytes
			c.etag = computeETag(c.bytes)
		}
	})
	return c.bytes, c.etag, c.err
}

func (c *HandlerCache) New(cacheBuilder func() ([]byte, error)) HandlerCache {
	return HandlerCache{
		bytes:      c.bytes,
		etag:       c.etag,
		BuildCache: cacheBuilder,
	}
}
