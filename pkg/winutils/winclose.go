package winutils

import (
	"errors"
	"fmt"

	"github.com/gonutz/w32/v2"
	"github.com/samber/lo"
)

/*
Source: https://gist.github.com/SCP002/ab863ef9ffbacedc2c0b1b4d30e80805
*/

func CloseWindow(pid int, wait bool) error {
	wnd, isUWP, err := GetWindow(pid)
	if err != nil {
		return err
	}
	var ok bool
	message := lo.Ternary(isUWP, w32.WM_QUIT, w32.WM_CLOSE)
	if wait {
		ok = w32.SendMessage(wnd, uint32(message), 0, 0) == 0
	} else {
		ok = w32.PostMessage(wnd, uint32(message), 0, 0)
	}
	if !ok {
		return errors.New("Failed to close the window with PID " + fmt.Sprint(pid))
	}
	return nil
}

func GetWindow(pid int) (w32.HWND, bool, error) {
	var wnd w32.HWND
	var isUWP bool
	w32.EnumWindows(func(hwnd w32.HWND) bool {
		_, currentPid := w32.GetWindowThreadProcessId(hwnd)

		if int(currentPid) == pid {
			if IsUWPApp(hwnd) {
				isUWP = true
				wnd = hwnd
				// Stop enumerating.
				return false
			}
			if IsMainWindow(hwnd) {
				wnd = hwnd
				// Stop enumerating.
				return false
			}
		}
		// Continue enumerating.
		return true
	})
	if wnd != 0 {
		return wnd, isUWP, nil
	} else {
		return wnd, isUWP, errors.New("No window found for PID " + fmt.Sprint(pid))
	}
}

func IsUWPApp(hwnd w32.HWND) bool {
	info, _ := w32.GetWindowInfo(hwnd)
	return info.AtomWindowType == 49223
}

func IsMainWindow(hwnd w32.HWND) bool {
	return w32.GetWindow(hwnd, w32.GW_OWNER) == 0 && w32.IsWindowVisible(hwnd)
}
