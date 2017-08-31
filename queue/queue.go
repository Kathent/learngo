package queue



type E interface{}

type Node struct {
	data E
	next *Node
}

type Queue interface {
	Add(e E)
	Element() E
	//Offer(e E)
	//Peak() E
	Poll() E
	Remove() E
}


type QueueIml struct {
	head *Node
}

func New() *QueueIml{
	return &QueueIml{}
}

func (qi *QueueIml) Add(e E){
	tmp := &Node{data:e}
	if qi.head == nil {
		qi.head = tmp
		return
	}

	var n *Node
	for n = qi.head; n != nil;{
		if n.next != nil {
			n = n.next
		}
	}

	n.next = tmp
}

func (qi QueueIml) Element() E {
	return qi.head.data
}

//func (qi QueueIml) Offer(e E){
//	qi.Add(e)
//}
//
//func (qi QueueIml) Peak() E{
//	return qi.Element()
//}

func (qi QueueIml) Remove() E {
	result := qi.head.data
	qi.head = qi.head.next
	return result
}

