package day9

type DLinkNode struct {
	key, value int
	pre, next  *DLinkNode
}

// NewDLinkedNode 初始化新节点
func NewDLinkedNode(key, value int) *DLinkNode {
	return &DLinkNode{
		key: key, value: value, pre: nil, next: nil,
	}
}

type LRUCache struct {
	size     int                // used 个数
	capacity int                // max
	cache    map[int]*DLinkNode // data 双向链表
	Head     *DLinkNode
	Tail     *DLinkNode
}

func NewLRUCache(capacity int) LRUCache {
	l := LRUCache{
		size:     0,
		capacity: capacity,
		cache:    make(map[int]*DLinkNode, capacity/3),
		Head:     NewDLinkedNode(0, 0),
		Tail:     NewDLinkedNode(0, 0),
	}
	l.Head.next = l.Tail
	l.Tail.pre = l.Head
	return l
}

func (l *LRUCache) Get(key int) int {
	// 缓存未命中， 返回-1
	if _, ok := l.cache[key]; !ok {
		return -1
	}
	// 命中缓存，返回值。将该值放在链头，即最近访问
	node := l.cache[key]
	l.UpdateToHead(node)
	return node.value
}

// Put 如果key存在，则变更数值；若不存在则插入。若插入导致size 将大于 capacity, 则逐出最久未使用key
func (l *LRUCache) Put(key int, value int) {
	if _, ok := l.cache[key]; !ok {
		// 没找到，尝试插入
		node := NewDLinkedNode(key, value)
		for l.size >= l.capacity {
			l.DeleteLast()
		}
		l.cache[key] = node
		l.InsertNewHead(node)
	} else {
		node := l.cache[key]
		node.value = value
		l.UpdateToHead(node)
	}
}

// UpdateToHead 更新到链头，用于key命中，不改变缓存size
func (l *LRUCache) UpdateToHead(node *DLinkNode) {
	node.pre.next = node.next
	node.next.pre = node.pre
	tmp := l.Head.next
	l.Head.next = node
	node.pre = l.Head
	node.next = tmp
	tmp.pre = node
}

// DeleteLast 删除链尾元素
func (l *LRUCache) DeleteLast() {
	node := l.Tail.pre
	l.Tail.pre = node.pre
	node.pre.next = node.next
	node.pre = nil
	node.next = nil
	l.size--
	delete(l.cache, node.key)
}

// InsertNewHead 添加新的key
func (l *LRUCache) InsertNewHead(node *DLinkNode) {
	tmp := l.Head.next
	l.Head.next = node
	node.pre = l.Head
	node.next = tmp
	tmp.pre = node
	l.size++
}
