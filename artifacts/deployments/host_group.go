package deployments

const (
	JumpboxHostGroupType = "jumpbox"
)

type HostGroup struct {
	Name      string
	GroupType string
	Hosts     []Host
}
