package main

import (
	"fmt"
	"syscall"
	"unsafe"
)

var (
	moduser32   = syscall.NewLazyDLL("user32.dll")
	modwin32u   = syscall.NewLazyDLL("win32u.dll")
	ntoskrnlDLL = syscall.NewLazyDLL("ntoskrnl.exe")
)

var (
	procNtMapVisualRelativePoints = modwin32u.NewProc("NtMapVisualRelativePoints")
	procPsLookupProcessByProcessId = ntoskrnlDLL.NewProc("PsLookupProcessByProcessId")
	procMmCopyVirtualMemory         = ntoskrnlDLL.NewProc("MmCopyVirtualMemory")
	procObfDereferenceObject        = ntoskrnlDLL.NewProc("ObfDereferenceObject")
	procMmGetPhysicalAddress        = ntoskrnlDLL.NewProc("MmGetPhysicalAddress")
)

type Data struct {
	FromPid   uint64
	ToPid     uint64
	FromAddr  unsafe.Pointer
	ToAddr    unsafe.Pointer
	Size      uintptr
}

func mapped(data *Data) syscall.Errno {
	var fromProcess, toProcess unsafe.Pointer
	ret, _, _ := procPsLookupProcessByProcessId.Call(
		data.FromPid,
		uintptr(unsafe.Pointer(&fromProcess)),
	)
	if ret != 0 {
		return syscall.Errno(ret)
	}

	ret, _, _ = procPsLookupProcessByProcessId.Call(
		data.ToPid,
		uintptr(unsafe.Pointer(&toProcess)),
	)
	if ret != 0 {
		return syscall.Errno(ret)
	}

	var size uintptr
	ret, _, _ = procMmCopyVirtualMemory.Call(
		uintptr(fromProcess),
		uintptr(data.FromAddr),
		uintptr(toProcess),
		uintptr(data.ToAddr),
		data.Size,
		1,
		uintptr(unsafe.Pointer(&size)),
	)
	if ret != 0 {
		return syscall.Errno(ret)
	}

	procObfDereferenceObject.Call(uintptr(fromProcess))
	procObfDereferenceObject.Call(uintptr(toProcess))

	return nil
}

func main() {
	hUser32 := moduser32.Handle()
	defer syscall.FreeLibrary(hUser32)

	hWin32u := modwin32u.Handle()
	defer syscall.FreeLibrary(hWin32u)

	fnNtMapVisualRelativePoints, err := procNtMapVisualRelativePoints.Find()
	if err != nil {
		fmt.Println("[X] Error, ", err)
		return
	}

	var testVariableOne = int32(0xDEAD)
	var testVariableTwo = int32(0xBEEF)

	data := Data{
		FromPid:  uint64(syscall.GetCurrentProcessId()),
		ToPid:    uint64(syscall.GetCurrentProcessId()),
		FromAddr: unsafe.Pointer(&testVariableOne),
		ToAddr:   unsafe.Pointer(&testVariableTwo),
		Size:     uintptr(unsafe.Sizeof(testVariableOne)),
	}

	fmt.Printf("[+] before: %X\n", testVariableTwo)
	result, _, _ := syscall.Syscall(fnNtMapVisualRelativePoints, 1, uintptr(unsafe.Pointer(&data)), 0, 0)
	if result != 0 {
		fmt.Println("[!] Mapped func broke:", result)
	} else {
		fmt.Printf("[+] after: %X\n", testVariableTwo)
	}

	fmt.Scanln()
}
