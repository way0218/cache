package cache

import "container/list"

// Cache is a LRU cache. It is not safe for concurrent access.
type Cache struct {
	maxBytes int64 // 允许使用的最大内存
	nbytes   int64 // 当前已使用的内存
	ll       *list.List
	cache    map[string]*list.Element
	// optional and executed when an entry is purged.
	OnEvicted func(key string, value Value) // 某条记录被移除时的回调函数，可以为 nil
}

type entry struct {
	key   string
	value Value
}

// Value 用于返回值所占用的内存大小
type Value interface {
	Len() int
}

// NewCache 初始化 Cache
func NewCache(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

// Get 查找
func (c *Cache) Get(key string) (value Value, ok bool) {
	if el, ok := c.cache[key]; ok {
		c.ll.MoveToFront(el)
		kv := el.Value.(*entry)
		return kv.value, ok
	}
	return
}

// RemoveOldest 移除最近最少访问的节点（队首） --- 缓存淘汰
func (c *Cache) RemoveOldest() {
	el := c.ll.Back()
	if el != nil {
		c.ll.Remove(el)
		kv := el.Value.(*entry)
		delete(c.cache, kv.key)
		c.nbytes -= int64(len(kv.key)) + int64(kv.value.Len())

		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

// Add 新增/修改
// 如果键存在，则更新对应节点的值，并将该节点移到队尾。
// 不存在则是新增场景，首先队尾添加新节点 &entry{key, value}, 并字典中添加 key 和节点的映射关系。
// 更新 c.nbytes，如果超过了设定的最大值 c.maxBytes，则移除最少访问的节点。
func (c *Cache) Add(key string, value Value) {
	if el, ok := c.cache[key]; ok {
		c.ll.MoveToFront(el)
		kv := el.Value.(*entry)
		c.nbytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else {
		el := c.ll.PushFront(&entry{key, value})
		c.cache[key] = el
		c.nbytes += int64(len(key)) + int64(value.Len())
	}

	for c.maxBytes != 0 && c.maxBytes < c.nbytes {
		c.RemoveOldest()
	}
}

// Len 获取添加了多少条数据
func (c *Cache) Len() int {
	return c.ll.Len()
}
