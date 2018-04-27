package etcd

import (
	"log"

	v3 "github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/clientv3/concurrency"
)

//DMutex DMutex
type DMutex struct {
	mutex   *concurrency.Mutex
	session *concurrency.Session
	client  *v3.Client
	pfx     string
}

//GetMutex get a mutex
func GetMutex(cli *v3.Client, key string, ttl int) (*DMutex, error) {
	result := new(DMutex)
	result.pfx = key
	ops := concurrency.WithTTL(ttl)
	s1, err := concurrency.NewSession(cli, ops)
	if err != nil {
		return nil, err
	}
	result.client = cli
	result.session = s1
	m1 := concurrency.NewMutex(s1, result.pfx)
	result.mutex = m1
	return result, nil
}

//Lock Lock
func (m *DMutex) Lock() bool {
	if m.mutex == nil {
		return false
	}
	if err := m.mutex.Lock(m.client.Ctx()); err != nil {
		log.Println(" DMutex Lock err:", err)
		return false
	}
	return true
}

//UnLock UnLock
func (m *DMutex) UnLock() {
	if m.client == nil {
		return
	}
	if m.session == nil {
		return
	}
	if m.mutex == nil {
		return
	}
	//defer m.client.Close()
	defer m.session.Close()
	defer m.mutex.Unlock(m.client.Ctx())
}

//IsLocked IsLocked
func (m *DMutex) IsLocked() bool {
	if m.mutex == nil {
		return false
	}
	resp, err := m.client.Get(m.client.Ctx(), m.pfx, v3.WithLastRev()[0])
	if err != nil {
		return false
	}
	if len(resp.Kvs) > 0 {
		return resp.Kvs[0].Lease == int64(m.session.Lease())
	}
	return false
}

//Lease Lease
func (m *DMutex) Lease() v3.LeaseID {
	return m.session.Lease()
}
