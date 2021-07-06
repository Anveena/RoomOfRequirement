package ezNetworking

import (
	"fmt"
	"syscall"
)

func setIOBufferSize(c syscall.RawConn, bufferSize int, opt int) *EZNetError {
	rsErr := &EZNetError{}
	err := c.Control(func(fd uintptr) {
		if err := syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, opt, bufferSize); err != nil {
			rsErr.IsSettingError = true
			rsErr.errorDesc = err.Error()
		}
	})
	if err != nil {
		rsErr.IsSettingError = true
		rsErr.errorDesc = err.Error()
		return rsErr
	}
	if rsErr.errorDesc != "" {
		return rsErr
	}
	err = c.Control(func(fd uintptr) {
		if actuallyValue, err := syscall.GetsockoptInt(int(fd), syscall.SOL_SOCKET, opt); err != nil {
			rsErr.IsLoadingError = true
			rsErr.errorDesc = err.Error()
			return
		} else if actuallyValue < bufferSize {
			rsErr.IsCheckingError = true
			rsErr.errorDesc = fmt.Sprintf("set buffer failed,wanted result:%v,actually:%v", bufferSize, actuallyValue)
			return
		}
	})
	if err != nil {
		rsErr.IsLoadingError = true
		rsErr.errorDesc = err.Error()
		return rsErr
	}
	if rsErr.errorDesc != "" {
		return rsErr
	}
	return nil
}
