package dict

import (
	"math/rand"
	"sync"
	"time"
)

type Smap struct {
	s     sync.Map
	count int
}

func MakeSync() *Smap {
	return &Smap{
		s: sync.Map{},
	}
}

func (s *Smap) Get(key string) (val interface{}, exists bool) {
	if s == nil {
		panic("dict is nil")
	}
	val, exists = s.s.Load(key)
	return
}

// Len returns the number of dict
func (s *Smap) Len() (l int) {
	if s == nil {
		panic("dict is nil")
	}
	s.s.Range(func(k, v interface{}) bool {
		l++
		return true
	})
	return
}

// Put puts key value into dict and returns the number of new inserted key-value
func (s *Smap) Put(key string, val interface{}) (result int) {
	if s == nil {
		panic("dict is nil")
	}
	s.addCount()

	if _, ok := s.s.Load(key); ok {
		s.s.Store(key, val)
		return 0
	}
	s.s.Store(key, val)
	return 1
}

// PutIfAbsent puts value if the key is not exists and returns the number of updated key-value
func (s *Smap) PutIfAbsent(key string, val interface{}) (result int) {
	if s == nil {
		panic("dict is nil")
	}

	if _, ok := s.s.Load(key); ok {
		return 0
	}
	s.s.Load(key)
	s.addCount()
	return 1
}

// PutIfExists puts value if the key is exist and returns the number of inserted key-value
func (s *Smap) PutIfExists(key string, val interface{}) (result int) {
	if s == nil {
		panic("dict is nil")
	}
	if _, ok := s.s.Load(key); ok {
		s.s.Store(key, val)
		return 1
	}
	return 0
}

// Remove removes the key and return the number of deleted key-value
func (s *Smap) Remove(key string) (result int) {
	if s == nil {
		panic("dict is nil")
	}

	if _, ok := s.s.Load(key); ok {
		s.s.Delete(key)
		s.decreaseCount()
		return 1
	}
	return 0
}

// ForEach traversal the dict
// it may not visits new entry inserted during traversal
func (s *Smap) ForEach(consumer Consumer) {
	if s == nil {
		panic("dict is nil")
	}

	s.s.Range(func(key, value interface{}) bool {
		continues := consumer(key.(string), value)
		if !continues {
			return false
		}
		return true
	})
}

// Keys returns all keys in dict
func (s *Smap) Keys() []string {
	keys := make([]string, s.Len())
	i := 0
	s.ForEach(func(key string, val interface{}) bool {
		if i < len(keys) {
			keys[i] = key
			i++
		} else {
			keys = append(keys, key)
		}
		return true
	})

	return keys
}

// RandomKeys randomly returns keys of the given number, may contain duplicated key
func (s *Smap) RandomKeys(limit int) []string {
	size := s.Len()
	if limit >= size {
		return s.Keys()
	}

	result := make([]string, limit)
	keys := s.Keys()
	nR := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < limit; {
		random := nR.Intn(len(keys))
		if keys[random] != "" {
			result[i] = keys[random]
			i++
		}
	}
	return result
}

// RandomDistinctKeys randomly returns keys of the given number, won't contain duplicated key
func (s *Smap) RandomDistinctKeys(limit int) []string {
	size := s.Len()
	if limit >= size {
		return s.Keys()
	}

	result := make([]string, limit)
	keys := s.Keys()
	nR := rand.New(rand.NewSource(time.Now().UnixNano()))
	randoms := make([]int, 0, limit)
	for i := 0; i < limit; {
		random := nR.Intn(len(keys))
		for i := 0; i < len(randoms); i++ {
			if keys[random] == keys[i] {
				random = nR.Intn(len(keys))
			}
		}
		randoms = append(randoms, random)
		if keys[random] != "" {
			result[i] = keys[random]
			i++
		}
	}
	return result
}

func (s *Smap) addCount() int32 {
	s.count++
	return int32(s.count)
}

func (s *Smap) decreaseCount() int32 {
	s.count--
	return int32(s.count)
}

// Clear removes all keys in dict
func (s *Smap) Clear() {
	*s = *MakeSync()
}
