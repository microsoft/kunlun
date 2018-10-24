package common

import (
	"fmt"
	"strings"
)

type InfraSize string

const (
	SizeSmall        InfraSize = "small"
	SizeMedium       InfraSize = "medium"
	SizeLarge        InfraSize = "large"
	SizeMaximum      InfraSize = "maximum"
	UnknownInfraSize InfraSize = "unknown"
)

func ParseInfraSize(is string) (InfraSize, error) {
	is = strings.ToLower(is)
	switch is {
	case "small":
		return SizeSmall, nil
	case "medium":
		return SizeMedium, nil
	case "large":
		return SizeLarge, nil
	case "maximum":
		return SizeMaximum, nil
	default:
		return UnknownInfraSize, fmt.Errorf("Not support the size")
	}
}
