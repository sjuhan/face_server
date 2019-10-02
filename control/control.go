package control

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"unsafe"
)

var (
	user32               = syscall.MustLoadDLL("user32.dll")
	procEnumWindows      = user32.MustFindProc("EnumWindows")
	procGetWindowTextW   = user32.MustFindProc("GetWindowTextW")
	procGetClassNameW    = user32.MustFindProc("GetClassNameW")
	procEnumChildWindows = user32.MustFindProc("EnumChildWindows")
	procSendMessage      = user32.MustFindProc("SendMessageW")
)

func sendMessage(hwnd syscall.Handle, uMsg uint32, wParam uintptr, lParam uintptr) (lResult uintptr) {
	r0, _, _ := syscall.Syscall6(procSendMessage.Addr(), 4, uintptr(hwnd), uintptr(uMsg), uintptr(wParam), uintptr(lParam), 0, 0)
	lResult = uintptr(r0)
	return
}

func wmTextlengh(hwnd syscall.Handle) uintptr {
	//b := make([]uint32, 200)
	return sendMessage(hwnd, 0xE, 0, 0)
}

func wmGetText(hwnd syscall.Handle) string {

	b := make([]uint16, wmTextlengh(hwnd)+1)
	syscall.Syscall6(procSendMessage.Addr(), 4, uintptr(hwnd), uintptr(0xD), uintptr(len(b)+1), uintptr(unsafe.Pointer(&b[0])), 0, 0)

	return syscall.UTF16ToString(b)
}

//GetClassName 1
func GetClassName(hwnd syscall.Handle) (name string, err error) {
	n := make([]uint16, 256)
	p := &n[0]
	r0, _, e1 := syscall.Syscall(procGetClassNameW.Addr(), 3, uintptr(hwnd), uintptr(unsafe.Pointer(p)), uintptr(len(n)))
	if r0 == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
		return
	}
	name = syscall.UTF16ToString(n)
	return
}

//EnumWindows 1
func EnumWindows(enumFunc uintptr, lparam uintptr) (err error) {
	r1, _, e1 := syscall.Syscall(procEnumWindows.Addr(), 2, uintptr(enumFunc), uintptr(lparam), 0)
	if r1 == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

//EnumChildWindows 1
func EnumChildWindows(hwnd syscall.Handle, enumFunc uintptr, lparam uintptr) (err error) {
	r1, _, e1 := syscall.Syscall(procEnumChildWindows.Addr(), 3, uintptr(hwnd), uintptr(enumFunc), uintptr(lparam))
	if r1 == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

//GetWindowText 1
func GetWindowText(hwnd syscall.Handle, str *uint16, maxCount int32) (len int32, err error) {
	r0, _, e1 := syscall.Syscall(procGetWindowTextW.Addr(), 3, uintptr(hwnd), uintptr(unsafe.Pointer(str)), uintptr(maxCount))
	len = int32(r0)
	if len == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

//FindWindow 1
func FindWindow(title string) (syscall.Handle, string, error) {
	var hwnd syscall.Handle
	cb := syscall.NewCallback(func(h syscall.Handle, p uintptr) uintptr {
		b := make([]uint16, 200)
		_, err := GetWindowText(h, &b[0], int32(len(b)))
		if err != nil {
			// ignore the error
			return 1 // continue enumeration
		}
		if strings.Contains(syscall.UTF16ToString(b), title) {
			// note the window
			hwnd = h
			title = syscall.UTF16ToString(b)

			return 0 // stop enumeration
		}
		return 1 // continue enumeration
	})
	EnumWindows(cb, 0)
	if hwnd == 0 {
		return 0, "0", fmt.Errorf("No window with title '%s' found", title)
	}
	return hwnd, title, nil
}

//FindChildWindow 1
func FindChildWindow(handle syscall.Handle, class string, n int) (syscall.Handle, string, error) {
	var hwnd syscall.Handle
	var num int = 0
	//var tt []string
	cb := syscall.NewCallback(func(h syscall.Handle, p uintptr) uintptr {

		t, err := GetClassName(h)
		//tt = append(tt, t)
		if err != nil {
			// ignore the error
			return 1 // continue enumeration
		}
		if strings.Contains(t, class) {
			// note the window
			num++
			if num == n {

				hwnd = h
				class = t
				return 0
			}

			return 1 // stop enumeration
		}
		return 1 // continue enumeration
	})
	EnumChildWindows(handle, cb, 0)
	//fmt.Println(tt)
	if hwnd == 0 {
		return 0, "0", fmt.Errorf("No window with title '%s' found", class)
	}
	return hwnd, class, nil
}

var (
	kernel32DLL          = syscall.NewLazyDLL("kernel32.dll")
	procCreateJobObjectA = kernel32DLL.NewProc("CreateJobObjectA")
)

// CreateJobObject uses the CreateJobObjectA Windows API Call to create and return a Handle to a new JobObject
func CreateJobObject(attr *syscall.SecurityAttributes, name string) (syscall.Handle, error) {
	r1, _, err := procCreateJobObjectA.Call(
		uintptr(unsafe.Pointer(attr)),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(name))),
	)
	if err != syscall.Errno(0) {
		return 0, err
	}
	return syscall.Handle(r1), nil
}

//GetText 오토핫키의 ControlGetText 구현
func GetText(win string, control string) string {
	숫자 := regexp.MustCompile("[0-9]+$")
	문자 := regexp.MustCompile("[a-zA-Z]+[0-9]?[a-zA-Z]+")
	n := 숫자.FindAllString(control, -1)
	c := 문자.FindAllString(control, -1)
	h, _, _ := FindWindow(win)
	rn, _ := strconv.Atoi(n[0])
	hh, _, _ := FindChildWindow(h, c[0], rn)

	return wmGetText(hh)
}

func main() {
	fmt.Println(GetText("종합검진", "ThunderRT6TextBox34"))
	fmt.Println(GetText("종합검진", "ThunderRT6TextBox33"))
}
