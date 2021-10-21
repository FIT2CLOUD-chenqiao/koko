package srvconn

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/jumpserver/koko/pkg/logger"
)

func LoginToSu(sc *SSHConnection) error {
	steps := make([]stepItem, 0, 2)
	steps = append(steps,
		stepItem{
			Input:         sc.options.sudoCommand,
			ExpectPattern: passwordMatchPattern,
			IsCommand:     true,
		},
		stepItem{
			Input:         sc.options.sudoPassword,
			ExpectPattern: fmt.Sprintf("%s@", sc.options.sudoUsername),
		},
	)
	for i := 0; i < len(steps); i++ {
		if err := executeStep(&steps[i], sc); err != nil {
			return err
		}
	}
	_, _ = sc.Write([]byte("\r\n"))
	return nil
}

func executeStep(step *stepItem, sc *SSHConnection) error {
	return step.Execute(sc)
}

const (
	LinuxSuCommand = "su - %s; exit"

	passwordMatchPattern = "(?i)password|密码"
)

var ErrorTimeout = errors.New("time out")

type stepItem struct {
	Input         string
	ExpectPattern string
	IsCommand     bool
}

func (s *stepItem) Execute(sc *SSHConnection) error {
	success := make(chan struct{}, 1)
	errorChan := make(chan error, 1)
	matchReg, err := regexp.Compile(s.ExpectPattern)
	if err != nil {
		logger.Error(err)
	}
	if s.IsCommand {
		_ = sc.session.Start(s.Input)
	} else {
		_, _ = sc.Write([]byte(s.Input + "\r\n"))
	}
	go func() {
		buf := make([]byte, 8192)
		var revStr strings.Builder
		for {
			nr, err := sc.Read(buf)
			if err != nil {
				errorChan <- err
				return
			}
			revStr.Write(buf[:nr])
			result := revStr.String()
			if matchReg != nil && matchReg.MatchString(result) {
				success <- struct{}{}
				return
			}
		}
	}()
	ticker := time.NewTicker(time.Second * 30)
	defer ticker.Stop()
	select {
	case <-success:
	case err := <-errorChan:
		return err

	case <-ticker.C:
		return ErrorTimeout
	}
	return nil
}
