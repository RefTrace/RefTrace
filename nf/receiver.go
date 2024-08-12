package nf

import (
	"fmt"
	"reft-go/parser"
)

type Receiver[T any] struct {
	type_  parser.IClassNode
	object bool
	data   T
}

func MakeReceiver[T any](type_ parser.IClassNode) *Receiver[T] {
	if type_ == nil {
		type_ = parser.OBJECT_TYPE.GetPlainNodeReference()
	}
	return &Receiver[T]{
		type_:  type_,
		object: true,
	}
}

func NewReceiver[T any](type_ parser.IClassNode) *Receiver[T] {
	return NewReceiverWithDataAndObject[T](type_, true, *new(T))
}

func NewReceiverWithData[T any](type_ parser.IClassNode, data T) *Receiver[T] {
	return NewReceiverWithDataAndObject[T](type_, true, data)
}

func NewReceiverWithDataAndObject[T any](type_ parser.IClassNode, object bool, data T) *Receiver[T] {
	if type_ == nil {
		panic("type cannot be nil")
	}
	return &Receiver[T]{
		type_:  type_,
		object: object,
		data:   data,
	}
}

func (r *Receiver[T]) GetData() T {
	return r.data
}

func (r *Receiver[T]) IsObject() bool {
	return r.object
}

func (r *Receiver[T]) GetType() parser.IClassNode {
	return r.type_
}

func (r *Receiver[T]) String() string {
	objectPrefix := ""
	if !r.object {
		objectPrefix = "*"
	}
	return fmt.Sprintf("Receiver{data=%v, type=%s%s}", r.data, objectPrefix, r.type_.ToStringWithRedirect(false))
}
