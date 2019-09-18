package margui

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime/debug"
	"strings"

	"golang.org/x/xerrors"
)

// FilenameWithoutExtension убирает расширение из пути fn
func FilenameWithoutExtension(fn string) string {
	return strings.TrimSuffix(fn, path.Ext(fn))
}

// PanicRecover ловит панику и сообщает о ней в лог
func PanicRecover() interface{} {
	if r := recover(); r != nil {
		logGeneric(4, logFStr, "PanicRecover", r)
		log.Println(StackTrace(1))
		return r
	}
	return nil
}

// PanicRecoverCallback ловит панику и сообщает о ней в лог, а также вызывает переданную функцию-обработчик
func PanicRecoverCallback(callback func()) interface{} {
	if r := recover(); r != nil {
		logGeneric(4, logFStr, "PanicRecoverCallback", r)
		log.Println(StackTrace(2))
		if callback != nil {
			callback()
		}
		return r
	}
	return nil
}

// PanicRecoverErr ловит панику, сообщает о ней в лог и отдаем ошибку
func PanicRecoverErr(err *error) {
	if r := recover(); r != nil {
		LogF("PanicRecoverErr")
		if err != nil {

			// find out exactly what the error was and set err
			switch x := r.(type) {
			case string:
				*err = xerrors.New(x)
			case error:
				*err = x
			default:
				// Fallback err (per specs, error strings should be lowercase w/o punctuation
				*err = xerrors.New("unknown panic")
			}
		}
	}
	return
}

func Close(closer io.Closer) {
	if closer == nil {
		return
	}
	err := closer.Close()
	if err == nil {
		return
	}
	logGeneric(1, logWStr, "Close:", err)
}

func CloseNoLog(closer io.Closer) {
	if closer == nil {
		return
	}
	_ = closer.Close()
}

// SafeWriteFile is a drop-in replacement for ioutil.WriteFile;
// but SafeWriteFile writes data to a temporary file first and
// only upon success renames that file to filename.
func SafeWriteFile(filename string, data []byte, perm os.FileMode) error {

	dir := filepath.Dir(filename)

	err := os.MkdirAll(dir, perm)
	if err != nil {
		return err
	}

	// open temp file
	f, err := ioutil.TempFile(dir, "tmp")
	if err != nil {
		return err
	}
	_ = f.Chmod(perm)

	tmpname := f.Name()

	// write data to temp file
	n, err := f.Write(data)
	if err == nil && n < len(data) {
		err = io.ErrShortWrite
	}
	if err1 := f.Close(); err == nil {
		err = err1
	}
	if err != nil {
		return err
	}

	return os.Rename(tmpname, filename)
}

//IsFileExisting существует ли файл или директория в файловой системе
func IsFileExisting(fullPath string) bool {
	_, err := os.Stat(fullPath)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

//FirstExistingFile Возвращает индекс первого файла из указанного списка, который существует, либо -1.
func FirstExistingFile(filenames []string) int {
	for i := range filenames {
		if IsFileExisting(filenames[i]) {
			return i
		}
	}
	return -1
}

func CountFilesInDirectory(path string, suffix string) (n int) {
	dir, err := os.Open(path)
	if err != nil {
		return
	}
	defer Close(dir)
	list, err := dir.Readdirnames(-1)
	if err != nil {
		return
	}
	for _, s := range list {
		if strings.HasSuffix(s, suffix) {
			n++
		}
	}
	return
}

func StackTrace(skipFrames int) string {
	skipFrames += 2               // пропускаем содержимое debug.Stack() и текущей функции
	skipLines := skipFrames*2 + 1 // на каждый фрейм приходится по 2 строки, плюс первая строка - имя горутины
	stack := strings.SplitN(string(debug.Stack()), "\n", skipLines+1)
	if len(stack) < skipLines+1 {
		return ""
	}
	goroutineInfo := stack[0]
	requestedStackPart := stack[skipLines]
	return goroutineInfo + "\n" + requestedStackPart
}

func ToMultilineJsonString(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}
