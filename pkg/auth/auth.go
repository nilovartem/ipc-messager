package auth

import (
	"os"
	"os/user"
	"strconv"
)

type Data struct {
	PID
	UID
	GID
}

// Current get's auth data for client and server
func Current() Data {
	u, err := user.Lookup(os.ExpandEnv("$USER"))
	if err != nil {
		data := Data{}
		data.PID = (PID)(0)
		data.UID = (UID)(0)
		data.GID = (GID)(0)
		return data
	}
	data := Data{}
	data.PID = (PID)(os.Getpid())
	data.UID = (UID)(u.Uid)
	data.GID = (GID)(u.Gid)
	return data
}

type UID string

func (uid UID) String() string {
	return string(uid)
}

type PID int

func (pid PID) String() string {
	return strconv.Itoa(int(pid))
}

type GID string

func (gid GID) String() string {
	return string(gid)
}
