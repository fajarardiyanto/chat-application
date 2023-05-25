package constant

type Role int32

const (
	SUPERVISOR Role = iota
	AGENT
	ADMIN
)

func AgentRole(s int32) string {
	switch s {
	case 0:
		return "SUPERVISOR"
	case 1:
		return "AGENT"
	case 2:
		return "ADMIN"
	}
	return ""
}
