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

func (m *buildLockManager) Unlock(appName string) {
	lockInterface, ok := m.locks.Load(appName)
	if !ok {
		return
	}

	lock := lockInterface.(*sync.Mutex)
	// 某些场景下，为了防止多次调用 Unlock 导致 panic，
	// 可以在这里配合状态位，但通常建议由调用方保证 Lock/Unlock 成对出现
	lock.Unlock()

	// 如果确定该应用构建彻底结束且以后不再需要，可以考虑删除 Key
	m.locks.Delete(appName)
}
