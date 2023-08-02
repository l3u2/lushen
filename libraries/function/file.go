package function

import (
	"archive/zip"
	"atsjkhelper/libraries/crc"
	"bufio"
	"errors"
	"hash/crc64"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

func Md5FileContent(file string) string {
	content, _ := ioutil.ReadFile(file)
	md5Content := Md5(string(content))

	return md5Content
}

func IsDir(filename string) bool {
	fd, err := os.Stat(filename)
	if err != nil {
		return false
	}
	return fd.IsDir()
}

func IsFile(filename string) bool {
	return !IsDir(filename)
}

func WriteFile(filename, str string) (err error) {
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return
	}
	defer f.Close()
	outputWriter := bufio.NewWriter(f)
	_, err = outputWriter.WriteString(str)
	if err != nil {
		return
	}
	err = outputWriter.Flush()
	if err != nil {
		return
	}
	return
}

func GetDirAllFilePathFollowSymlink(dirname string) ([]string, error) {
	dirname = strings.TrimSuffix(dirname, string(os.PathSeparator))

	infos, err := ioutil.ReadDir(dirname)
	if err != nil {
		return nil, err
	}

	paths := make([]string, 0, len(infos))
	for _, info := range infos {
		path := dirname + string(os.PathSeparator) + info.Name()
		realInfo, err := os.Stat(path)
		if err != nil {
			return nil, err
		}
		if realInfo.IsDir() {
			tmp, err := GetDirAllFilePathFollowSymlink(path)
			if err != nil {
				return nil, err
			}
			paths = append(paths, tmp...)
			continue
		}
		paths = append(paths, path)
	}

	return paths, nil
}

func MkdirAll(dirname string, mode os.FileMode) error {
	dirExist := IsDir(dirname)

	if dirExist {
		return nil
	}

	return os.MkdirAll(dirname, mode)
}

func DownloadFile(ossUrl string, downloadPath string) error {
	// Get the data
	resp, err := http.Get(ossUrl)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 创建一个文件用于保存
	out, err := os.Create(downloadPath)
	if err != nil {
		return err
	}
	defer out.Close()

	// 然后将响应流和文件流对接起来
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func Unzip(src string, dest string) error {
	reader, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer reader.Close()
	for _, file := range reader.File {
		filePath := path.Join(dest, file.Name)
		if file.FileInfo().IsDir() {
			os.MkdirAll(filePath, os.ModePerm)
		} else {
			if err = os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
				return err
			}
			inFile, err := file.Open()
			if err != nil {
				return err
			}
			defer inFile.Close()
			outFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
			if err != nil {
				return err
			}
			defer outFile.Close()
			_, err = io.Copy(outFile, inFile)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func Basename(path string) string {
	return filepath.Base(path)
}

func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	if err == nil {
		return true
	} else {
		return false
	}
}

func ReadFile(filename string) (data []byte, err error) {
	file, err := os.Open(filename)
	if err != nil {
		return
	}
	defer file.Close()
	data, err = ioutil.ReadAll(file)
	if err != nil {
		return
	}
	return
}

func Crc64Sum(fileName string, chunkNum int) (string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return "", err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return "", err
	}

	if int64(chunkNum) > stat.Size() {
		return "", errors.New("chunkNum bigger than fileName")
	}
	var chunkN = (int64)(chunkNum)
	var dataSize int64
	tab := crc64.MakeTable(crc64.ECMA)
	var crcval uint64
	for i := int64(0); i < chunkN; i++ {
		if i == chunkN-1 {
			dataSize = stat.Size()/chunkN + stat.Size()%chunkN
		} else {
			dataSize = stat.Size() / chunkN
		}
		calc := crc.NewCRC(tab, 0)

		byteSlice := make([]byte, dataSize)
		_, err := io.ReadFull(file, byteSlice)
		if err != nil {
			return "", err
		}
		calc.Write(byteSlice)
		crcval = crc.CRC64Combine(crcval, calc.Sum64(), (uint64)(dataSize))
	}
	return strconv.FormatUint(crcval, 10), nil
}

func FileSize(filename string) (int64, error) {
	info, err := os.Stat(filename)
	if err != nil && os.IsNotExist(err) {
		return 0, err
	}
	return info.Size(), nil
}
