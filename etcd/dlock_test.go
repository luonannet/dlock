package etcd

import (
	"testing"
	"time"

	"github.com/coreos/etcd/clientv3"
)

var userid int64

//plusUserID plus user id
func plusUserID(step int64) int64 {
	userid += step
	return userid
}

func TestDlock(t *testing.T) {
	cli, err := clientv3.New(clientv3.Config{Endpoints: []string{"192.168.2.112:2379"}})
	if err != nil {
		t.Fatalf("TestDlock clientv3.New err:%s", err)
	}
	//开启10个协程 ，每个协程都去抢锁 并执行自己的运算
	for i := 1; i <= 10; i++ {
		myMutex, mutexErr := GetMutex(cli, "testkey", 10)
		if mutexErr != nil {
			t.Fatalf("GetMutex %s:%s", myMutex.mutex.Key(), mutexErr)
		}
		go func() {
			if myMutex.Lock() == false {
				t.Fatalf("myMutex%s.Lock false", myMutex.mutex.Key())
			}
			t.Logf(" %s 得到锁，计算结果:%d", myMutex.mutex.Key(), plusUserID(1))
			t.Logf(" %s 释放锁", myMutex.mutex.Key())
			myMutex.UnLock()
		}()
	}
	time.Sleep(time.Second * 1)
}
func BenchmarkDlock(b *testing.B) {
	cli, err := clientv3.New(clientv3.Config{Endpoints: []string{"192.168.2.112:2379"}})
	if err != nil {
		b.Fatalf("BenchmarkDlock clientv3.New err:%s", err)

	}
	// b.SetParallelism(10)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			myMutex, mutexErr := GetMutex(cli, "testkey", 10)
			if mutexErr != nil {
				b.Fatalf("GetMutex %s:%s", myMutex.mutex.Key(), mutexErr)
			}
			if myMutex.Lock() == false {
				b.Fatalf("myMutex%s.Lock false", myMutex.mutex.Key())
			}
			plusUserID(1)
			myMutex.UnLock()
			if myMutex.IsLocked() {
				//
			}
		}
	})
}
