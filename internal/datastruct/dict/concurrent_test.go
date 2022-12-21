package dict

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/lixvyang/Go-redis/internal/utils"
)

func TestConcurrentPut(t *testing.T) {
	d := MakeConcurrent(0)
	count := 100
	var wg sync.WaitGroup
	wg.Add(count)
	for i := 0; i < count; i++ {
		go func(i int) {
			// insert
			key := "k" + strconv.Itoa(i)
			ret := d.Put(key, i)
			if ret != 1 { // insert 1
				t.Error("put test failed: expected result 1, actual: " + strconv.Itoa(ret) + ", key: " + key)
			}
			val, ok := d.Get(key)
			if ok {
				intVal, _ := val.(int)
				if intVal != i {
					t.Error("put test failed: expected " + strconv.Itoa(i) + ", actual: " + strconv.Itoa(intVal) + ", key: " + key)
				}
			} else {
				_, ok := d.Get(key)
				t.Error("put test failed: expected true, actual: false, key: " + key + ", retry: " + strconv.FormatBool(ok))
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
}

func TestConcurrentPutIfAbsent(t *testing.T) {
	d := MakeConcurrent(0)
	count := 100
	var wg sync.WaitGroup
	wg.Add(count)
	for i := 0; i < count; i++ {
		go func(i int) {
			// insert
			key := "k" + strconv.Itoa(i)
			ret := d.PutIfAbsent(key, i)
			if ret != 1 { // insert 1
				t.Error("put test failed: expected result 1, actual: " + strconv.Itoa(ret) + ", key: " + key)
			}
			val, ok := d.Get(key)
			if ok {
				intVal, _ := val.(int)
				if intVal != i {
					t.Error("put test failed: expected " + strconv.Itoa(i) + ", actual: " + strconv.Itoa(intVal) +
						", key: " + key)
				}
			} else {
				_, ok := d.Get(key)
				t.Error("put test failed: expected true, actual: false, key: " + key + ", retry: " + strconv.FormatBool(ok))
			}

			// update
			ret = d.PutIfAbsent(key, i*10)
			if ret != 0 { // no update
				t.Error("put test failed: expected result 0, actual: " + strconv.Itoa(ret))
			}
			val, ok = d.Get(key)
			if ok {
				intVal, _ := val.(int)
				if intVal != i {
					t.Error("put test failed: expected " + strconv.Itoa(i) + ", actual: " + strconv.Itoa(intVal) + ", key: " + key)
				}
			} else {
				t.Error("put test failed: expected true, actual: false, key: " + key)
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
}

func TestConcurrentPutIfExists(t *testing.T) {
	d := MakeConcurrent(0)
	count := 100
	var wg sync.WaitGroup
	wg.Add(count)
	for i := 0; i < count; i++ {
		go func(i int) {
			// insert
			key := "k" + strconv.Itoa(i)
			// insert
			ret := d.PutIfExists(key, i)
			if ret != 0 { // insert
				t.Error("put test failed: expected result 0, actual: " + strconv.Itoa(ret))
			}

			d.Put(key, i)
			d.PutIfExists(key, 10*i)
			val, ok := d.Get(key)
			if ok {
				intVal, _ := val.(int)
				if intVal != 10*i {
					t.Error("put test failed: expected " + strconv.Itoa(10*i) + ", actual: " + strconv.Itoa(intVal))
				}
			} else {
				_, ok := d.Get(key)
				t.Error("put test failed: expected true, actual: false, key: " + key + ", retry: " + strconv.FormatBool(ok))
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
}

func TestConcurrentRemove(t *testing.T) {
	d := MakeConcurrent(0)
	totalCount := 100
	// remove head node
	for i := 0; i < totalCount; i++ {
		// insert
		key := "k" + strconv.Itoa(i)
		d.Put(key, i)
	}
	if d.Len() != totalCount {
		t.Error("put test failed: expected len is 100, actual: " + strconv.Itoa(d.Len()))
	}
	for i := 0; i < totalCount; i++ {
		key := "k" + strconv.Itoa(i)

		val, ok := d.Get(key)
		if ok {
			intVal, _ := val.(int)
			if intVal != i {
				t.Error("put test failed: expected " + strconv.Itoa(i) + ", actual: " + strconv.Itoa(intVal))
			}
		} else {
			t.Error("put test failed: expected true, actual: false")
		}

		ret := d.Remove(key)
		if ret != 1 {
			t.Error("remove test failed: expected result 1, actual: " + strconv.Itoa(ret) + ", key:" + key)
		}
		if d.Len() != totalCount-i-1 {
			t.Error("put test failed: expected len is 99, actual: " + strconv.Itoa(d.Len()))
		}
		_, ok = d.Get(key)
		if ok {
			t.Error("remove test failed: expected true, actual false")
		}
		ret = d.Remove(key)
		if ret != 0 {
			t.Error("remove test failed: expected result 0 actual: " + strconv.Itoa(ret))
		}
		if d.Len() != totalCount-i-1 {
			t.Error("put test failed: expected len is 99, actual: " + strconv.Itoa(d.Len()))
		}
	}

	// remove tail node
	d = MakeConcurrent(0)
	for i := 0; i < 100; i++ {
		// insert
		key := "k" + strconv.Itoa(i)
		d.Put(key, i)
	}
	for i := 9; i >= 0; i-- {
		key := "k" + strconv.Itoa(i)

		val, ok := d.Get(key)
		if ok {
			intVal, _ := val.(int)
			if intVal != i {
				t.Error("put test failed: expected " + strconv.Itoa(i) + ", actual: " + strconv.Itoa(intVal))
			}
		} else {
			t.Error("put test failed: expected true, actual: false")
		}

		ret := d.Remove(key)
		if ret != 1 {
			t.Error("remove test failed: expected result 1, actual: " + strconv.Itoa(ret))
		}
		_, ok = d.Get(key)
		if ok {
			t.Error("remove test failed: expected true, actual false")
		}
		ret = d.Remove(key)
		if ret != 0 {
			t.Error("remove test failed: expected result 0 actual: " + strconv.Itoa(ret))
		}
	}

	// remove middle node
	d = MakeConcurrent(0)
	d.Put("head", 0)
	for i := 0; i < 10; i++ {
		// insert
		key := "k" + strconv.Itoa(i)
		d.Put(key, i)
	}
	d.Put("tail", 0)
	for i := 9; i >= 0; i-- {
		key := "k" + strconv.Itoa(i)

		val, ok := d.Get(key)
		if ok {
			intVal, _ := val.(int)
			if intVal != i {
				t.Error("put test failed: expected " + strconv.Itoa(i) + ", actual: " + strconv.Itoa(intVal))
			}
		} else {
			t.Error("put test failed: expected true, actual: false")
		}

		ret := d.Remove(key)
		if ret != 1 {
			t.Error("remove test failed: expected result 1, actual: " + strconv.Itoa(ret))
		}
		_, ok = d.Get(key)
		if ok {
			t.Error("remove test failed: expected true, actual false")
		}
		ret = d.Remove(key)
		if ret != 0 {
			t.Error("remove test failed: expected result 0 actual: " + strconv.Itoa(ret))
		}
	}
}

func TestConcurrentRandomKey(t *testing.T) {
	d := MakeConcurrent(0)
	count := 100
	for i := 0; i < count; i++ {
		key := "k" + strconv.Itoa(i)
		d.Put(key, i)
	}
	fetchSize := 10
	result := d.RandomKeys(fetchSize)
	if len(result) != fetchSize {
		t.Errorf("expect %d random keys acturally %d", fetchSize, len(result))
	}
	result = d.RandomDistinctKeys(fetchSize)
	distinct := make(map[string]struct{})
	for _, key := range result {
		distinct[key] = struct{}{}
	}
	if len(result) != fetchSize {
		t.Errorf("expect %d random keys acturally %d", fetchSize, len(result))
	}
	if len(result) > len(distinct) {
		t.Errorf("get duplicated keys in result")
	}
}

func TestConcurrentDict_Keys(t *testing.T) {
	d := MakeConcurrent(0)
	size := 1000
	for i := 0; i < size; i++ {
		d.Put(utils.RandString(5), utils.RandString(5))
	}
	if len(d.Keys()) != size {
		t.Errorf("expect %d keys, actual: %d", size, len(d.Keys()))
	}
}

var times int = 10000000

// BenchmarkTestConcurrentMap 测试ConcurrentDictMap
func BenchmarkTestConcurrentDictMap(b *testing.B) {
	for k := 0; k < b.N; k++ {
		b.StopTimer()
		// 产生10000个不重复的键值对(string -> int)
		testKV := map[string]int{}
		for i := 0; i < 10000; i++ {
			testKV[strconv.Itoa(i)] = i
		}

		pMap := MakeConcurrent(0)

		// set到map中
		for k, v := range testKV {
			pMap.Put(k, v)
		}

		// 开始计时
		b.StartTimer()

		wg := sync.WaitGroup{}
		wg.Add(2)

		// 赋值
		go func() {
			// 对随机key,赋值times次
			for i := 0; i < times; i++ {
				index := rand.Intn(times)
				pMap.Put(strconv.Itoa(index), index+1)
			}
			wg.Done()
		}()

		// 读取
		go func() {
			// 对随机key,读取times次
			for i := 0; i < times; i++ {
				index := rand.Intn(times)
				pMap.Get(strconv.Itoa(index))
			}
			wg.Done()
		}()

		// 等待两个协程处理完毕
		wg.Wait()
	}
}
