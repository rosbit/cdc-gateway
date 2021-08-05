package cache

import (
	"sync"
	"fmt"
	"log"
	"time"
)

type cacheVal struct {
	content []byte
	expire  int64
}

var (
	cache = map[string]*cacheVal{}
	expire2key = map[int64][]string{}
	mutex = &sync.RWMutex{}
)

func Get(app, date, cdMark string, seq uint64) ([]byte) {
	key := makeKey(app, date, cdMark, seq)
	log.Printf("key to Get: %s\n", key)

	mutex.RLock()
	defer mutex.RUnlock()

	c, ok := cache[key]
	if !ok || c == nil {
		return nil
	}
	now := time.Now().Unix()
	log.Printf(" + now: %d, c.expire: %d\n", now, c.expire)
	if now > c.expire {
		return nil
	}
	return c.content
}

func Set(app, date, cdMark string, seq uint64, content []byte, liveTimeInSecs int64) {
	key := makeKey(app, date, cdMark, seq)

	mutex.Lock()
	defer mutex.Unlock()

	expire := time.Now().Unix() + liveTimeInSecs
	log.Printf("key to Set: %s, with expire time %d\n", key, expire)
	cache[key] = &cacheVal{
		content: content,
		expire: expire,
	}

	if ks, ok := expire2key[expire]; ok {
		for _, k := range ks {
			if k == key {
				return
			}
		}
		expire2key[expire] = append(ks, key)
	} else {
		expire2key[expire] = []string{key}
	}
}

func makeKey(app, date, cdMark string, seq uint64) string {
	return fmt.Sprintf("%s_%s_%s_%d", app, date, cdMark, seq)
}

func StartCleaningThread() {
	go func() {
		t := time.NewTicker(1*time.Minute)
		for now := range t.C {
			log.Printf("cleaning thread wakeup\n")
			clearExpiredContents(now.Unix())
		}
	}()
}

func clearExpiredContents(now int64) {
		mutex.Lock()
		defer mutex.Unlock()

		expireTime := []int64{}
		for expire, _ := range expire2key {
			if now < expire {
				break
			}
			expireTime = append(expireTime, expire)
		}

		if len(expireTime) == 0 {
			log.Printf(" + no cache removed\n")
			return
		}
		for _, expire := range expireTime {
			ks := expire2key[expire]
			delete(expire2key, expire)
			log.Printf(" + cache for expiredTime %d removed\n", expire)

			for _, key := range ks {
				c, ok := cache[key]
				if ok && c.expire <= expire {
					delete(cache, key)
					log.Printf("   - key %s removed\n", key)
				}
			}
		}
}
