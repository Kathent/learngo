package file_read_analysis

import (
	"io/ioutil"
	"log"
	"time"
	"os"
	"path/filepath"
	"fmt"
)

const (
	TIME_FORMATE = "2006-01-02 15:04:05"
)

func LoadFiles(dir string, startDate, endDate string) {
	startTime, err := time.Parse(TIME_FORMATE, startDate)
	if err != nil {
		panic(err)
	}
	endTime, err := time.Parse(TIME_FORMATE, endDate)
	if err != nil {
		panic(err)
	}

	resultFile, err := os.Create("result.txt")
	if err != nil {
		panic(err)
	}

	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if info == nil {
			return nil
		}

		if info.IsDir() {
			return nil
		}
		bytes, err := ioutil.ReadFile(path)
		if err != nil {
			log.Println(err)
		}

		start, end := 0, 0
		bl := len(bytes)
		for k := range bytes{
			if k + 19 >= bl {
				break
			}

			str := string(bytes[k:k+19])
			predicateTime, err := time.Parse(TIME_FORMATE, string(str))
			if err != nil {
				continue
			}

			//if end > 0 {
			//	resultFile.Write(bytes[start: k + 19])
			//}

			if predicateTime.After(startTime) && predicateTime.Before(endTime) {
				if start == 0 {
					start = k
				}
			}else if start != 0 && end == 0{
				end = k + 19
				break
			}
		}

		if start != 0 && end == 0{
			end = len(bytes) - 1
		}


		log.Println(fmt.Sprintf("start:%d, end:%d", start, end))
		resultFile.Write(bytes[start:end+1])
		return nil
	})
}

const(
	TIME_STR_LENGTH = 19
)

func LoadFilesWithReader(dir string, startDate, endDate string) error{
	startTime, err := time.Parse(TIME_FORMATE, startDate)
	if err != nil {
		return err
	}
	endTime, err := time.Parse(TIME_FORMATE, endDate)
	if err != nil {
		return err
	}

	reader, err := NewFilesReaderWithDir(dir)
	if err != nil {
		return err
	}

	resultFile, err := os.Create("result.txt")
	if err != nil {
		panic(err)
	}

	oneCommandArray := make([]byte, 0)
	forTimeArray := make([]byte, 0, TIME_STR_LENGTH)
	var needAppendOneCommand bool
	for {
		b, err := reader.ReadByte()
		if err != nil {
			break
		}
		if len(forTimeArray) == 0 || len(forTimeArray) < TIME_STR_LENGTH{
			forTimeArray = append(forTimeArray, b)
		}else {
			forTimeArray = append(forTimeArray[1:], b)
		}

		if needAppendOneCommand {
			oneCommandArray = append(oneCommandArray, b)
		}

		if len(forTimeArray) == 19 {
			parse, err := time.Parse(TIME_FORMATE, string(forTimeArray))
			if err == nil {
				if timeProper(startTime, endTime, parse) {//当前时间合适
					if !needAppendOneCommand {
						needAppendOneCommand = true
						oneCommandArray = forTimeArray
					}else {
						err := flushOneCommandToFile(resultFile, oneCommandArray)
						if err != nil{
							return err
						}
						oneCommandArray = oneCommandArray[:0]
					}
				}else {
					if needAppendOneCommand {
						needAppendOneCommand = false
						flushOneCommandToFile(resultFile, oneCommandArray)
					}
				}
			}
		}
	}

	if needAppendOneCommand {
		flushOneCommandToFile(resultFile, oneCommandArray)
	}

	return reader.Close()
}

func flushOneCommandToFile(f *os.File, bytes []byte) error{
	_, err := f.Write(bytes)
	return err
}

func timeProper(start time.Time, end time.Time, parse time.Time) bool{
	return parse.After(start) && parse.Before(end)
}
