package gquiz

import (
	yaml "gopkg.in/yaml.v2"
)

type QuizBuilder struct {
}

func (qb QuizBuilder) BuildQGraph(content []byte) (QGraph, error) {
	var qNodes []QNode
	err := yaml.Unmarshal(content, &qNodes)
	if err != nil {
		return QGraph{}, err
	}
	return QGraph{
		QNodes: qNodes,
	}, nil
}
