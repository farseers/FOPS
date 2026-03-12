package apps

import (
	"sync"
)

// buildLockManager 应用构建锁管理器
type buildLockManager struct {
	locks sync.Map // key: appName, value: *sync.Mutex
}

var lockManager = &buildLockManager{}

// TryLock 尝试获取应用的构建锁
// 返回 true 表示成功获取锁，false 表示该应用正在构建中
func (m *buildLockManager) TryLock(appName string) bool {
	// 获取或创建该应用的锁
	lockInterface, _ := m.locks.LoadOrStore(appName, &sync.Mutex{})
	lock := lockInterface.(*sync.Mutex)

	// 尝试获取锁（非阻塞）
	return lock.TryLock()
}

// Unlock 释放应用的构建锁
func (m *buildLockManager) Unlock(appName string) {
	if lockInterface, ok := m.locks.Load(appName); ok {
		lock := lockInterface.(*sync.Mutex)
		lock.Unlock()
	}
}
