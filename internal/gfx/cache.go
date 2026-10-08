package gfx

import (
	"fmt"
	"image/color"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
)

// SpriteCache manages and memoizes procedural sprites.
type SpriteCache struct {
	mu           sync.RWMutex
	towerSprites map[string]*ebiten.Image
	enemySprites map[string]*ebiten.Image
	iconSprites  map[string]*ebiten.Image
	cpuCore      *ebiten.Image
	spawnerPort  *ebiten.Image
}

var globalCache *SpriteCache
var once sync.Once

// GetCache returns the global singleton SpriteCache.
func GetCache() *SpriteCache {
	once.Do(func() {
		globalCache = &SpriteCache{
			towerSprites: make(map[string]*ebiten.Image),
			enemySprites: make(map[string]*ebiten.Image),
			iconSprites:  make(map[string]*ebiten.Image),
			cpuCore:      GenerateCPUCoreSprite(),
			spawnerPort:  GenerateSpawnerPortSprite(),
		}
	})
	return globalCache
}

// GetTowerSprite returns a cached procedural tower sprite for the given ID, tier, and color.
func (c *SpriteCache) GetTowerSprite(id string, tier int, baseColor color.RGBA) *ebiten.Image {
	key := fmt.Sprintf("tower_%s_%d_%d_%d_%d", id, tier, baseColor.R, baseColor.G, baseColor.B)
	c.mu.RLock()
	img, ok := c.towerSprites[key]
	c.mu.RUnlock()
	if ok {
		return img
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	if img, ok := c.towerSprites[key]; ok {
		return img
	}

	img = GenerateTowerSprite(id, tier, baseColor)
	c.towerSprites[key] = img
	return img
}

// GetEnemySprite returns a cached procedural enemy sprite.
func (c *SpriteCache) GetEnemySprite(id string, isBoss bool, baseColor color.RGBA, radius float32) *ebiten.Image {
	key := fmt.Sprintf("enemy_%s_%t_%d_%d_%d_%.1f", id, isBoss, baseColor.R, baseColor.G, baseColor.B, radius)
	c.mu.RLock()
	img, ok := c.enemySprites[key]
	c.mu.RUnlock()
	if ok {
		return img
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	if img, ok := c.enemySprites[key]; ok {
		return img
	}

	img = GenerateEnemySprite(id, isBoss, baseColor, radius)
	c.enemySprites[key] = img
	return img
}

// GetIconBadge returns a cached icon badge.
func (c *SpriteCache) GetIconBadge(id string, isTower bool, baseColor color.RGBA) *ebiten.Image {
	key := fmt.Sprintf("icon_%s_%t_%d_%d_%d", id, isTower, baseColor.R, baseColor.G, baseColor.B)
	c.mu.RLock()
	img, ok := c.iconSprites[key]
	c.mu.RUnlock()
	if ok {
		return img
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	if img, ok := c.iconSprites[key]; ok {
		return img
	}

	img = GenerateIconBadge(id, isTower, baseColor)
	c.iconSprites[key] = img
	return img
}

// GetCPUCore returns the cached CPU Core sprite.
func (c *SpriteCache) GetCPUCore() *ebiten.Image {
	return c.cpuCore
}

// GetSpawnerPort returns the cached Spawner Port sprite.
func (c *SpriteCache) GetSpawnerPort() *ebiten.Image {
	return c.spawnerPort
}
