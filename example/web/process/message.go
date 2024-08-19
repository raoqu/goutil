package process

const MSG_LOG = "log"
const MSG_KEEP_ALIVE = "<<<alive>>>"
const MSG_GROUP = "group"
const MSG_MESSAGE = "msg"
const MSG_STATUS_UPDATE = "status"

func ComposeMessage(msgType string, msg string) string {
	return msgType + ":" + msg
}
