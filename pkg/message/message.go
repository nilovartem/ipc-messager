package message

type IMessage interface {
	Unmarshall([]byte) error
	Marshall() ([]byte, error)
	CreateMessage([]byte)
}
