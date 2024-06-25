package parser

import (
	"errors"
	"sync"
)

type NodeMetaDataHandler interface {
	GetMetaDataMap() map[interface{}]interface{}
	SetMetaDataMap(metaDataMap map[interface{}]interface{})
	NewMetaDataMap() map[interface{}]interface{}
}

type DefaultNodeMetaDataHandler struct {
	metaDataMap map[interface{}]interface{}
	mu          sync.RWMutex
}

func (handler *DefaultNodeMetaDataHandler) GetMetaDataMap() map[interface{}]interface{} {
	handler.mu.RLock()
	defer handler.mu.RUnlock()
	return handler.metaDataMap
}

func (handler *DefaultNodeMetaDataHandler) SetMetaDataMap(metaDataMap map[interface{}]interface{}) {
	handler.mu.Lock()
	defer handler.mu.Unlock()
	handler.metaDataMap = metaDataMap
}

func (handler *DefaultNodeMetaDataHandler) NewMetaDataMap() map[interface{}]interface{} {
	return make(map[interface{}]interface{})
}

func (handler *DefaultNodeMetaDataHandler) GetNodeMetaData(key interface{}) interface{} {
	handler.mu.RLock()
	defer handler.mu.RUnlock()
	if handler.metaDataMap == nil {
		return nil
	}
	return handler.metaDataMap[key]
}

func (handler *DefaultNodeMetaDataHandler) GetNodeMetaDataWithFunc(key interface{}, valFn func() interface{}) interface{} {
	if key == nil {
		panic(errors.New("Tried to get/set meta data with null key"))
	}

	handler.mu.Lock()
	defer handler.mu.Unlock()
	if handler.metaDataMap == nil {
		handler.metaDataMap = handler.NewMetaDataMap()
		handler.SetMetaDataMap(handler.metaDataMap)
	}
	if val, ok := handler.metaDataMap[key]; ok {
		return val
	}
	val := valFn()
	handler.metaDataMap[key] = val
	return val
}

func (handler *DefaultNodeMetaDataHandler) CopyNodeMetaData(other NodeMetaDataHandler) {
	otherMetaDataMap := other.GetMetaDataMap()
	if otherMetaDataMap == nil {
		return
	}
	handler.mu.Lock()
	defer handler.mu.Unlock()
	if handler.metaDataMap == nil {
		handler.metaDataMap = handler.NewMetaDataMap()
		handler.SetMetaDataMap(handler.metaDataMap)
	}
	for k, v := range otherMetaDataMap {
		handler.metaDataMap[k] = v
	}
}

func (handler *DefaultNodeMetaDataHandler) SetNodeMetaData(key, value interface{}) {
	if old := handler.PutNodeMetaData(key, value); old != nil {
		panic(errors.New("Tried to overwrite existing meta data"))
	}
}

func (handler *DefaultNodeMetaDataHandler) PutNodeMetaData(key, value interface{}) interface{} {
	if key == nil {
		panic(errors.New("Tried to set meta data with null key"))
	}

	handler.mu.Lock()
	defer handler.mu.Unlock()
	if handler.metaDataMap == nil {
		if value == nil {
			return nil
		}
		handler.metaDataMap = handler.NewMetaDataMap()
		handler.SetMetaDataMap(handler.metaDataMap)
	} else if value == nil {
		return handler.metaDataMap[key]
	}
	return handler.metaDataMap[key]
}

func (handler *DefaultNodeMetaDataHandler) RemoveNodeMetaData(key interface{}) {
	if key == nil {
		panic(errors.New("Tried to remove meta data with null key"))
	}

	handler.mu.Lock()
	defer handler.mu.Unlock()
	if handler.metaDataMap != nil {
		delete(handler.metaDataMap, key)
	}
}

func (handler *DefaultNodeMetaDataHandler) GetNodeMetaDataMap() map[interface{}]interface{} {
	handler.mu.RLock()
	defer handler.mu.RUnlock()
	if handler.metaDataMap == nil {
		return map[interface{}]interface{}{}
	}
	return handler.metaDataMap
}
