package repository

import (
	"crypto/md5"
	"crypto/sha1"
	"sync"
)

const bloomFilterSize = 100
const numHashFunctions = 7

type BloomFilter struct {
	bitArray     []bool
	newBitArray  []bool
	migrateIndex int
	migrating    bool
	mu           sync.Mutex
}

// 返回布隆
func NewBloomFilter() *BloomFilter {
	return &BloomFilter{
		bitArray:     make([]bool, bloomFilterSize),
		newBitArray:  nil,
		migrateIndex: 0,
		migrating:    false,
	}
}

// Add 方法将给定的URL添加到布隆过滤器中。
// 它通过计算多个哈希值并将相应的位设置为true来实现。
func (bf *BloomFilter) Add(url string) {
	bf.mu.Lock()
	if bf.migrating {
		bf.incrementalMigrate()
	}
	hashes := bf.getHashIndices(url)
	for _, index := range hashes {
		bf.bitArray[index] = true
	}
	if bf.isFull() {
		bf.startMigration()
	}
	bf.mu.Unlock()
}

func (bf *BloomFilter) startMigration() {
	if bf.migrating {
		return
	}
	newSize := len(bf.bitArray) * 2
	bf.newBitArray = make([]bool, newSize)
	bf.migrateIndex = 0
	bf.migrating = true
}

func (bf *BloomFilter) incrementalMigrate() {
	if !bf.migrating || bf.newBitArray == nil {
		return
	}
	oldSize := len(bf.bitArray)
	newSize := len(bf.newBitArray)
	halfSize := oldSize / 2
	endIndex := bf.migrateIndex + halfSize
	if endIndex > oldSize {
		endIndex = oldSize
	}
	for i := bf.migrateIndex; i < endIndex; i++ {
		if bf.bitArray[i] {
			newIndex := bf.rehashIndex(i, newSize)
			bf.newBitArray[newIndex] = true
		}
	}
	bf.migrateIndex = endIndex
	if bf.migrateIndex >= oldSize {
		bf.bitArray = bf.newBitArray
		bf.newBitArray = nil
		bf.migrating = false
		bf.migrateIndex = 0
	}
}

func (bf *BloomFilter) rehashIndex(oldIndex int, newSize int) int {
	hash := md5.Sum([]byte(string(rune(oldIndex))))
	return int(hash[0])<<8 | int(hash[1])%newSize
}

// MightContain 方法检查布隆过滤器是否可能包含给定的URL。
// 如果所有计算出的哈希位都为true，则返回true，表示可能包含。
// 请注意，这个方法可能会产生假阳性，但不会产生假阴性。
func (bf *BloomFilter) MightContain(url string) bool {
	bf.mu.Lock()
	if bf.migrating {
		bf.incrementalMigrate()
	}
	hashes := bf.getHashIndices(url)
	for _, index := range hashes {
		if !bf.bitArray[index] {
			bf.mu.Unlock()
			return false
		}
	}
	bf.mu.Unlock()
	return true
}

// getHashIndices 方法根据给定的URL计算出多个哈希索引。
// 它使用MD5和SHA1哈希函数组合来生成这些索引，并确保它们在位数组的范围内。
func (bf *BloomFilter) getHashIndices(url string) []int {
	hashes := make([]int, numHashFunctions)
	hash1 := md5.Sum([]byte(url))
	hash2 := sha1.Sum([]byte(url))
	size := len(bf.bitArray)
	for i := 0; i < numHashFunctions; i++ {
		combinedHash := hash1[i%len(hash1)] + hash2[(i*2)%len(hash2)]
		index := int(combinedHash) % size
		if index < 0 {
			index = -index
		}
		hashes[i] = index
	}
	return hashes
}

// isFull 检查位数组是否已满（所有位都为 true）。
func (bf *BloomFilter) isFull() bool {
	for _, b := range bf.bitArray {
		if !b {
			return false
		}
	}
	return true
}

// Bloom 是一个全局的布隆过滤器实例
var Bloom = NewBloomFilter()
