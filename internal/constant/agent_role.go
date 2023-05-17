package constant

type Role int32

const (
	SUPERVISOR Role = iota
	AGENT
	ADMIN
	NONE
)

func AgentRole(s int32) string {
	switch s {
	case 0:
		return "SUPERVISOR"
	case 1:
		return "AGENT"
	case 2:
		return "ADMIN"
	case 3:
		return "NONE"
	}
	return ""
}
