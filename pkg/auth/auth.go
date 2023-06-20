package auth

import (
	"fmt"
	"os"
	"os/user"
	"strconv"
)

type Data struct {
	PID
	UID
	GID
}

func Current() Data { //вернуть еще и ошибку
	//fmt.Println(os.ExpandEnv("$USER"))
	u, err := user.Lookup(os.ExpandEnv("$USER"))
	if err != nil {
		fmt.Printf("AUTH.CURRENT: %s", err)
	}
	data := Data{}
	data.PID = (PID)(os.Getpid())
	//	fmt.Println("uid", u.Uid)
	//	fmt.Println("uname", u.Name)
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
