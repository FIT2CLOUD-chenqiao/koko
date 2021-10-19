package srvconn

import (
	gossh "golang.org/x/crypto/ssh"
	"io"
	"log"
	"regexp"
	"time"
)

func NewSSHMultipleConnection(sess *gossh.Session, opts ...SSHOption) (*SSHConnection, error) {
	con, err := NewSSHConnection(sess, opts...)
	if err != nil {
		return nil, err
	}
	err = con.session.Start("su;exit")
	if err != nil {
		return nil, err
	}
	/*
		i：		su && exit\r
		o：		密码|password|
		i：

	*/
	steps := []stepItem{
		//{
		//	"eric@jumpserver-qa",
		//	"su;exit\r",
		//},
		{
			"密码|password|Password",
			"Calong@2020\r",
		},
		{

			"root@jumpserver-qa",
			"\r\n",
		},
	}
	i := 0
	for i < len(steps) {
		if ok, err := expectFunc(con, steps[i].expect); err == nil && ok {
			_, _ = con.Write([]byte(steps[i].Input))
			log.Printf("step %d :pass\n", i)
			i++
		} else {
			time.Sleep(time.Second)
		}
	}

	return con, nil
}

type stepItem struct {
	expect string
	Input  string
}

func expectFunc(read io.ReadCloser, expect string) (bool, error) {

	buf := make([]byte, 8192)
	log.Println("read start ")
	nr, err := read.Read(buf)

	if err != nil {
		return false, err
	}
	log.Println("read end ")
	return regexp.Match(expect, buf[:nr])
}


