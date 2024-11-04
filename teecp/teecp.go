package teecp

type teeCPNode struct {
	next   *teeCPNode
	accept MsgReceiver
}

type TeeCPList struct {
	first *teeCPNode
	last  *teeCPNode
}

type MsgReceiver func(msg string) bool

func removeNext(curr *teeCPNode, teecp *TeeCPList) func() {
	return func() {
		if curr.next == teecp.last {
			teecp.last = curr
		} else {
			curr.next = curr.next.next
		}

	}
}

func (teecp *TeeCPList) Broadcast(msg string) {
	it := teecp.first
	remove := func() {
		if teecp.first == teecp.last {
			teecp.first = nil
			teecp.last = nil
		}
		teecp.first = teecp.first.next
	}
	for it != nil {
		couldWrite := it.accept(msg)
		if !couldWrite {
			remove()
		}
		remove = removeNext(it, teecp)
		it = it.next
	}
}

func Attach(lista *TeeCPList, msgRecv MsgReceiver) *TeeCPList {
	var nodo teeCPNode = teeCPNode{accept: msgRecv}
	if lista.last == nil {
		lista.first = &nodo
		lista.last = &nodo
	} else {
		lista.last.next = &nodo
		lista.last = &nodo
	}
	return lista
}
