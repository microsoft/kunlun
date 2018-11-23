package gquiz

const (
	Root = "root"
)

type QGraph struct {
	QNodes []QNode
}

func (qg QGraph) FindRootNode() *QNode {
	return qg.FindNode(Root)
}

func (qg QGraph) FindNode(name string) *QNode {
	for i := range qg.QNodes {
		if qg.QNodes[i].Name == name {
			return &qg.QNodes[i]
		}
	}
	return nil
}
