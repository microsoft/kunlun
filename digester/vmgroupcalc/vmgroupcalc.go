package vmgroupcalc

import (
	"github.com/Microsoft/kunlun/digester/common"
)

type Requirement struct {
	ConcurrentUserNumber int
}

func Calc(r Requirement) common.Infra {
	res := common.Infra{
		Size: common.SizeSmall,
	}
	x := r.ConcurrentUserNumber
	if x >= 1000 {
		res.Size = common.SizeMedium
	}
	if x >= 2000 {
		res.Size = common.SizeLarge
	}
	if x >= 4000 {
		res.Size = common.SizeMaximum
	}

	return res
}
