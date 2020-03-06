package stat

import (
	"fmt"
	"sync"

	"github.com/alibaba/sentinel-golang/core/base"
	sbase "github.com/alibaba/sentinel-golang/core/stat/base"
)

type ResourceNode struct {
	BaseStatNode

	resourceName string
	resourceType base.ResourceType
	// key is "sampleCount/intervalInMs"
	readOnlyStats map[string]*sbase.SlidingWindowMetric
	updateLock    sync.RWMutex
}

// NewResourceNode creates a new resource node with given name and classification.
func NewResourceNode(resourceName string, resourceType base.ResourceType) *ResourceNode {
	return &ResourceNode{
		// TODO: make this configurable
		BaseStatNode:  *NewBaseStatNode(base.DefaultSampleCount, base.DefaultIntervalMs),
		resourceName:  resourceName,
		resourceType:  resourceType,
		readOnlyStats: make(map[string]*sbase.SlidingWindowMetric),
	}
}

func (n *ResourceNode) ResourceType() base.ResourceType {
	return n.resourceType
}

func (n *ResourceNode) ResourceName() string {
	return n.resourceName
}

func (n *ResourceNode) GetSlidingWindowMetric(key string) *sbase.SlidingWindowMetric {
	n.updateLock.RLock()
	defer n.updateLock.RUnlock()
	return n.readOnlyStats[key]
}

func (n *ResourceNode) GetOrCreateSlidingWindowMetric(sampleCount, intervalInMs uint32) *sbase.SlidingWindowMetric {
	key := fmt.Sprintf("%d/%d", sampleCount, intervalInMs)
	fastVal := n.GetSlidingWindowMetric(key)
	if fastVal != nil {
		return fastVal
	}

	newSlidingWindow := sbase.NewSlidingWindowMetric(sampleCount, intervalInMs, n.arr)

	n.updateLock.Lock()
	defer n.updateLock.Unlock()

	n.readOnlyStats[key] = newSlidingWindow
	// TODO clean unused entity in readOnlyStats.
	return newSlidingWindow
}
