package chikador

type eventQueue struct {
	head *eventNode
	tail *eventNode
}

type eventNode struct {
	value *Message
	next  *eventNode
}

func (queue *eventQueue) append(event *Message) {
	node := &eventNode{value: event}
	if queue.tail != nil {
		queue.tail.next = node
	}
	if queue.head == nil {
		queue.head = node
	}
	queue.tail = node
}

func (queue *eventQueue) poll() *Message {
	if queue.head == nil {
		return nil
	}
	head := queue.head
	queue.head = queue.head.next
	return head.value
}
