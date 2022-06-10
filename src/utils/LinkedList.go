package utils

type Node[T any] struct {
	content T
	next    *Node[T]
}

type LinkedList[T any] struct {
	firstElement *Node[T]
	lastElement  *Node[T]
	length       int
}

func CreateLinkedList[T any](initialElements ...T) *LinkedList[T] {
	list := LinkedList[T]{}
	for _, element := range initialElements {
		list.Add(element)
	}
	return &list
}

func (list *LinkedList[T]) Add(element T) {
	newElement := Node[T]{
		content: element,
	}

	if list.length == 0 {
		list.firstElement = &newElement
	} else {
		list.lastElement.next = &newElement
	}
	list.lastElement = &newElement
	list.length++
}

func (list *LinkedList[T]) Length() int {
	return list.length
}

func (list *LinkedList[T]) Get(i int) T {
	if i == -1 || i == list.length-1 {
		return list.lastElement.content
	} else if i < -1 || i >= list.length {
		panic("Index out of bounds")
	}
	index := 0
	element := list.firstElement
	for i != index {
		element = element.next
		index++
	}
	return element.content
}

func (list *LinkedList[T]) Iterator() Iterator[T] {
	return Iterator[T]{
		index:          0,
		currentElement: list.firstElement,
	}
}

type Iterator[T any] struct {
	index          int
	currentElement *Node[T]
}

func (iter *Iterator[T]) HasNext() bool {
	return iter.currentElement != nil
}

func (iter *Iterator[T]) Next() T {
	if !iter.HasNext() {
		panic("Iterator has no next element")
	}
	element := iter.currentElement
	iter.currentElement = iter.currentElement.next
	return element.content
}
