package keyring

/*
#cgo CXXFLAGS:-O2 -std=c++11
#cgo CFLAGS: -W -Wall -fstack-protector-all -fPIC
#cgo LDFLAGS:-ldl
#include "common.h"
#include <stdlib.h>
*/
import "C"
import (
	"bufio"
	"fmt"
	"os"
	"path"
	"sync"
	"unsafe"
	"encoding/hex"
	"io/ioutil"
	dutils "github.com/linuxdeepin/go-lib/utils"
)

var fileLocker sync.Mutex

func ucharToArrayByte(value *C.uchar) string {
	data := C.GoString((*C.char)(unsafe.Pointer(value)))
	return hex.EncodeToString([]byte(data))
}

func ucharToString(value *C.uchar) string {
	return C.GoString((*C.char)(unsafe.Pointer(value)))
}

func writeFile(filename, data string) error {
	fileLocker.Lock()
	defer fileLocker.Unlock()

	fp, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("failed to open %q: %v", filename, err)
		return err
	}
	defer fp.Close()

	_, err = fp.WriteString(data + "\n")
	if err != nil {
		fmt.Println("failed to WriteString err :", err)
		return err
	}

	err = fp.Sync()
	if err != nil {
		fmt.Println("fp Sync err :", err)
		return err
	}

	return nil
}

func loadFile(filename string) ([]string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var lines []string
	scanner := bufio.NewScanner(bufio.NewReader(f))
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}
	if scanner.Err() != nil {
		return nil, scanner.Err()
	}

	return lines, nil
}

func createWhiteBoxUFile(dir, filePath string) error {
	fmt.Println("filePath: ", filePath)
	if !dutils.IsFileExist(dir) {
		err := os.MkdirAll(dir, 0600)
		if err != nil {
			fmt.Println("Mkdir err : ", err)
			return err
		}
	}
	if !dutils.IsFileExist(filePath) {
		err := dutils.CreateFile(filePath)
		if err != nil {
			fmt.Println("CreateFile err : ", err)
			return err
		}
	}
	return nil
}

func CreateWhiteBoxUFile(name string) error {
	C.set_debug_flag(0)

	UKEK := C.GoString(C.generate_random_len(C.MASTER_KEY_LEN))
	UKEKIV := C.GoString(C.generate_random_len(C.MASTER_KEY_LEN))
	fmt.Println("[CreateWhiteBoxUFile] UKEK generate_random_str   : ", UKEK)
	fmt.Println("[CreateWhiteBoxUFile] UKEKIV generate_random_str : ", UKEKIV)

	UKEK_ := unsafe.Pointer(C.CString(UKEK))
	defer C.free(UKEK_)

	UKEKIV_ := unsafe.Pointer(C.CString(UKEKIV))
	defer C.free(UKEKIV_)

	key := ""
	//key：16个0 白盒加密 UKEK --> WB_UKEK
	WB_UKEK := C.deepin_wb_encrypt((*C.uchar)(UKEK_), (*C.uchar)(unsafe.Pointer(C.CString(key))), false)

	////key：16个0 sm4解密 UKEK --> UKEK
	//out1 := C.sm4_crypt((*C.uchar)(WB_UKEK), (*C.uchar)(unsafe.Pointer(C.CString(key))), 0)
	//defer C.free(unsafe.Pointer(out1))
	//fmt.Println(">>>>>>>>>dec WB_UKEK -> UKEK : ", ucharToString(out1))

	//key：UKEK 白盒加密 UKEKIV --> CIPHER_UKEKIV
	CIPHER_UKEKIV := C.deepin_wb_encrypt((*C.uchar)(UKEKIV_), (*C.uchar)(UKEK_), true)
	defer C.free(unsafe.Pointer(CIPHER_UKEKIV))
	//fmt.Println(">>>>>>>>>UKEKIV enc UKEK -> CIPHER_UKEKIV : ", ucharToString(CIPHER_UKEKIV))

	//key：UKEKIV sm4解密 CIPHER_UKEKIV --> UKEK
	//out2 := C.sm4_crypt((*C.uchar)(CIPHER_UKEKIV), (*C.uchar)(UKEK_), 0)
	//defer C.free(unsafe.Pointer(out2))
	//fmt.Println(">>>>>>>>>UKEKIV dec CIPHER_UKEKIV -> UKEKIV : ", ucharToString(out2))

	//创建新增账户WB_UFile文件
	dir := fmt.Sprintf("/var/lib/keyring/%s", name)
	filePath := path.Join(dir, "WB_UFile")
	createWhiteBoxUFile(dir, filePath)

	writeFile(filePath, hex.EncodeToString([]byte(ucharToString(WB_UKEK))))
	writeFile(filePath, hex.EncodeToString([]byte(ucharToString(CIPHER_UKEKIV))))

	// 读取WB_UFILE文件
	lines, err := loadFile(filePath)
	if err != nil {
		fmt.Println("loadFile err : ", err)
		return err
	}
	fmt.Println("loadFile lines : ", lines)

	return nil
}

func DeleteWhiteBoxUFile(name string) error {
	filePath := fmt.Sprintf("/var/lib/keyring/%s", name)
	if dutils.IsFileExist(filePath) {
		dirs, err := ioutil.ReadDir(filePath)
		if err != nil {
			fmt.Println("ReadDir err : ", err)
			return err
		}
		for _, dir := range dirs {
			err = os.RemoveAll(path.Join(filePath, dir.Name()))
			if err != nil {
				fmt.Println("RemoveAll dir err : ", err)
				return err
			}
		}
		err = os.Remove(filePath)
		if err != nil {
			fmt.Println("Remove filePath err : ", err)
			return err
		}
	}
	return nil
}