package constant

type MsgType int32

const (
	TEXT MsgType = iota
	AUDIO_VIDEO
)

//var MessageType = map[string]int32{
//	"TEXT":        0,
//	"AUDIO_VIDEO": 1,
//}

func MessageType(s string) MsgType {
	switch s {
	case "TEXT":
		return 0
	case "AUDIO_VIDEO":
		return 1
	}
	return 2
}
