package parser

import "errors"

var _ NodeMetaDataHandler = (*DefaultNodeMetaDataHandler)(nil)

type NodeMetaDataHandler interface {
	GetMetaDataMap() map[interface{}]interface{}
	SetMetaDataMap(metaDataMap map[interface{}]interface{})
	NewMetaDataMap() map[interface{}]interface{}
	GetNodeMetaData(key interface{}) interface{}
	SetNodeMetaData(key, value interface{})
}

type DefaultNodeMetaDataHandler struct {
	metaDataMap map[interface{}]interface{}
}

func (handler *DefaultNodeMetaDataHandler) GetMetaDataMap() map[interface{}]interface{} {
	return handler.metaDataMap
}

func (handler *DefaultNodeMetaDataHandler) SetMetaDataMap(metaDataMap map[interface{}]interface{}) {
	handler.metaDataMap = metaDataMap
}

func (handler *DefaultNodeMetaDataHandler) NewMetaDataMap() map[interface{}]interface{} {
	return make(map[interface{}]interface{})
}

func (handler *DefaultNodeMetaDataHandler) GetNodeMetaData(key interface{}) interface{} {
	if handler.metaDataMap == nil {
		return nil
	}
	return handler.metaDataMap[key]
}

func (handler *DefaultNodeMetaDataHandler) GetNodeMetaDataWithFunc(key interface{}, valFn func() interface{}) interface{} {
	if key == nil {
		panic(errors.New("Tried to get/set meta data with null key"))
	}

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

	if handler.metaDataMap == nil {
		if value == nil {
			return nil
		}
		handler.metaDataMap = handler.NewMetaDataMap()
		handler.SetMetaDataMap(handler.metaDataMap)
	} else if value == nil {
		return handler.metaDataMap[key]
	}
	oldValue := handler.metaDataMap[key]
	handler.metaDataMap[key] = value
	return oldValue
}

func (handler *DefaultNodeMetaDataHandler) RemoveNodeMetaData(key interface{}) {
	if key == nil {
		panic(errors.New("Tried to remove meta data with null key"))
	}

	if handler.metaDataMap != nil {
		delete(handler.metaDataMap, key)
	}
}

func (handler *DefaultNodeMetaDataHandler) GetNodeMetaDataMap() map[interface{}]interface{} {
	if handler.metaDataMap == nil {
		return map[interface{}]interface{}{}
	}
	return handler.metaDataMap
}
