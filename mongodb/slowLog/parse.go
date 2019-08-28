package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

func main() {
	if len(os.Args) == 1 {
		panic("必须有参数,表示日志路径")
	}
	if file, err := os.Open(os.Args[1]); err != nil {
		panic(err.Error())
	} else {
		elements := make(map[string][]*Element, 100)
		defer file.Close()
		reader := bufio.NewReader(file)
		for {
			if line, _, err := reader.ReadLine(); err != nil {
				if err == io.EOF {
					break
				} else {
					panic(err.Error())
				}
			} else {
				if bytes.Contains(line,[]byte(`I WRITE`)) || bytes.Contains(line,[]byte(`I READ`)){
					if element, err := parseLine(line); err != nil {
						panic(err.Error())
					} else {
						elements[element.Id] = append(elements[element.Id], element)
					}
				}

			}
		}
		for id, conElements := range elements {
			sort.Slice(conElements,func(i,j int)bool{
				return conElements[i].Start.Before(conElements[j].Start)

			} )
			fmt.Printf("线程%s:\n", id)
			for _, element := range conElements {
				fmt.Printf("\t%s\n", element.String())
			}
			fmt.Println("----------------")
		}
	}
}

//297573ms 2019-04-04T10:55:50.236+0800 I WRITE    [conn280462] update trade.user_assets command: { q: { uid: 6, freeze_asset.usdt.m_zec_usdt_1554201224549521160: { $lt: 1E-22 } }, u: { $unset: { freeze_asset.usdt.m_zec_usdt_1554201224549521160: 1 } }, multi: false, upsert: false } planSummary: IXSCAN { uid: 1, asset.xmr: 1 } keysExamined:1 docsExamined:1 fromMultiPlanner:1 nMatched:1 nModified:1 writeConflicts:1219 numYields:1230 locks:{ Global: { acquireCount: { r: 1231, w: 1231 } }, Database: { acquireCount: { w: 1231 }, acquireWaitCount: { w: 2 }, timeAcquiringMicros: { w: 46578 } }, Collection: { acquireCount: { w: 1231 } } } 297573ms
func parseLine(data []byte) (elem *Element, err error) {
	parses := []Parse{
		parseTakes,
		parseStart,
		parseId,
		parseOperation,
		parseCollection,
		parseCommand,
		parsePlanSummary,
		parseExamined,
		parseLock,
	}
	elem = &Element{}

	for _, parse := range parses {
		if data, err = parse(data, elem); err != nil {
			return
		}
	}
	return
}

type Parse func([]byte, *Element) ([]byte, error)

func parseTakes(data []byte, elem *Element) (result []byte, err error) {
	if end := bytes.Index(data, []byte(` `)); end != -1 {
		elem.Takes, err = time.ParseDuration(string(data[:end]))
		err = errors.Wrap(err, string(data))
		result = data[end+1:]
	}
	return
}

func parseStart(data []byte, elem *Element) (result []byte, err error) {
	if end := bytes.Index(data, []byte(` `)); end != -1 {
		elem.Start, err = time.Parse("2006-01-02T15:04:05.000+0800", string(data[:end]))
		result = data[end+1:]
	}
	return
}

func parseId(data []byte, elem *Element) (result []byte, err error) {
	if start := bytes.Index(data, []byte(`[`)); start != -1 {
		if end := bytes.Index(data, []byte(`]`)); end != -1 {
			elem.Id = string(data[start+5 : end])
			result = data[end+1:]
		}
	}
	return
}

func parseOperation(data []byte, elem *Element) (result []byte, err error) {
	data = bytes.TrimSpace(data)
	if end := bytes.Index(data, []byte(` `)); end != -1 {
		elem.Operation = string(data[:end])
		result = data[end+1:]
	}
	return
}

func parseCollection(data []byte, elem *Element) (result []byte, err error) {
	data = bytes.TrimSpace(data)
	if end := bytes.Index(data, []byte(` `)); end != -1 {
		elem.Collection = string(data[:end])
		result = data[end+1:]
	}
	return
}

//command: { q: { uid: 6, freeze_asset.usdt.m_zec_usdt_1554201224549521160: { $lt: 1E-22 } }, u: { $unset: { freeze_asset.usdt.m_zec_usdt_1554201224549521160: 1 } }, multi: false, upsert: false } planSummary: IXSCAN { uid: 1, asset.xmr: 1 } keysExamined:1 docsExamined:1 fromMultiPlanner:1 nMatched:1 nModified:1 writeConflicts:1219 numYields:1230 locks:{ Global: { acquireCount: { r: 1231, w: 1231 } }, Database: { acquireCount: { w: 1231 }, acquireWaitCount: { w: 2 }, timeAcquiringMicros: { w: 46578 } }, Collection: { acquireCount: { w: 1231 } } } 297573ms
func parseCommand(data []byte, elem *Element) (result []byte, err error) {
	data = bytes.TrimSpace(data)
	count := 0
	start := 0
	end := 0
	for i, char := range data {
		switch char {
		case byte('{'):
			count++
			if count == 1 {
				start = i
			}
		case byte('}'):
			count--
			if count == 0 {
				end = i
				elem.Command = string(data[start:end])
				result = data[end+1:]
				return

			}
		}
	}
	return

}

func parsePlanSummary(data []byte, elem *Element) (result []byte, err error) {
	key := []byte(`planSummary:`)

	if start := bytes.Index(data, key); start != -1 {
		data = data[start+len(key):]
		data = bytes.TrimSpace(data)
		if end := bytes.Index(data, []byte(`}`)); end != -1 {
			elem.Plan = string(data[:end+1])
			result = data[end+1:]
		}
	}
	return
}
func parseExamined(data []byte, elem *Element) (result []byte, err error) {
	data, err = parseExaminedEach(data, []byte(`keysExamined:`), &elem.KeysExamined)
	if err != nil {
		return
	}
	data, err = parseExaminedEach(data, []byte(`docsExamined:`), &elem.DocsExamined)
	data, err = parseExaminedEach(data, []byte(`nMatched:`), &elem.Matched)
	data, err = parseExaminedEach(data, []byte(`nModified:`), &elem.Modified)
	data, err = parseExaminedEach(data, []byte(`writeConflicts:`), &elem.WriteConflicts)
	result, err = parseExaminedEach(data, []byte(`numYields:`), &elem.Yields)
	return
}

func parseExaminedEach(data []byte, key []byte, value *int) (result []byte, err error) {
	if start := bytes.Index(data, key); start != -1 {
		data = data[start+len(key):]
		data = bytes.TrimSpace(data)
		if end := bytes.Index(data, []byte(` `)); end != -1 {
			*value, _ = strconv.Atoi(string(data[:end]))
			result = data[end+1:]
		}
	}
	return
}

//locks:{ Global: { acquireCount: { r: 1231, w: 1231 } }, Database: { acquireCount: { w: 1231 }, acquireWaitCount: { w: 2 }, timeAcquiringMicros: { w: 46578 } }, Collection: { acquireCount: { w: 1231 } } } 297573ms
func parseLock(data []byte, elem *Element) (result []byte, err error) {
	key := []byte(`locks:`)
	count := 0
	singStart := 0
	if start := bytes.Index(data, key); start != -1 {
		for i, char := range data[start:] {
			switch char {
			case byte('{'):
				count++
				if count == 1 {
					singStart = i + 1
				}
			case byte('}'):
				count--
				switch count {
				case 0:
					result = data[i+1:]
					return
				case 1:
					if _, err = parseSingleLock(data[singStart:i+1], elem); err != nil {
						return
					}
					singStart = i + 1

				}
			}
		}
	}
	return
}

// Global: { acquireCount: { r: 1231, w: 1231 } }
func parseSingleLock(data []byte, elem *Element) (result []byte, err error) {
	data = bytes.TrimLeft(data, ", ")
	data = bytes.TrimSpace(data)
	//println("single", string(data))
	if semi := bytes.Index(data, []byte(":")); semi != -1 {
		switch string(data[:semi]) {
		case "Global":
			elem.Global = LockInfo{}
			parseSingleLockInfo(data[semi+1:], &elem.Global)
		case "Database":
			parseSingleLockInfo(data[semi+1:], &elem.DataBase)

		case "Collection":
			parseSingleLockInfo(data[semi+1:], &elem.CollectionL)

		}
	}
	return nil, nil
}

//{ acquireCount: { r: 1231, w: 1231 } }
func parseSingleLockInfo(data []byte, info *LockInfo) {
	data = bytes.TrimSpace(data)
	for {
		//println("lockInfo", string(data))

		if semi := bytes.Index(data, []byte(":")); semi != -1 {
			empty := bytes.Index(data, []byte(" "))
			end := bytes.Index(data, []byte(`}`))
			switch strings.TrimSpace(string(data[empty+1 : semi])) {
			case "acquireCount":
				info.AcquiredCount = LockDetail{}
				parseSingleLockDetail(data[semi+1:end+1], &info.AcquiredCount)
			case "acquireWaitCount":
				info.AcquiredWaitCount = LockDetail{}
				parseSingleLockDetail(data[semi+1:end+1], &info.AcquiredWaitCount)
			case "timeAcquiringMicros":
				info.AcquireTakes = LockDetail{}
				parseSingleLockDetail(data[semi+1:end+1], &info.AcquireTakes)
			case "deadlockCount":
				info.DeadLockTakes = LockDetail{}
				parseSingleLockDetail(data[semi+1:end+1], &info.DeadLockTakes)
			default:
				//println("lockInfo-default", string(data[:semi]))
			}
			data = data[end+1:]
		} else {

			return
		}
	}

}

//{ r: 1231, w: 1231 } }
func parseSingleLockDetail(data []byte, detail *LockDetail) {
	//println("detail", string(data))
	data = bytes.TrimSpace(data)
	var target *int64
	for {
		if semi := bytes.Index(data, []byte(":")); semi != -1 {
			start := bytes.Index(data, []byte(` `))
			switch string(data[start+1 : semi]) {
			case "w":
				target = &detail.TryWrite

			case "r":
				target = &detail.TryRead

			case "W":
				target = &detail.Write

			case "R":
				target = &detail.Read
			default:
				println("detail-default", string(data[start+1:semi]))
			}
			end := bytes.IndexAny(data, ",}")
			*target, _ = strconv.ParseInt(string(bytes.TrimSpace(data[semi+1:end])), 10, 64)
			data = data[end+1:]
		} else {
			break
		}
	}

}

type Element struct {
	Id             string
	Takes          time.Duration
	Start          time.Time
	Operation      string
	Command        string
	Collection     string
	Plan           string
	KeysExamined   int
	DocsExamined   int
	Matched        int
	Modified       int
	WriteConflicts int
	Yields         int
	Global         LockInfo
	DataBase       LockInfo
	CollectionL    LockInfo
}

func (e Element) String() string {
	builder := &strings.Builder{}
	fmt.Fprintf(builder, "\t%s 执行 %s 操作,耗时%s", e.Start.String(), e.Operation, e.Takes.String())
	fmt.Fprintf(builder, "\t数据库[%s],命令[%s]\\n", e.Collection, e.Command)
	fmt.Fprintf(builder, "\t plan[%s],检查Key 数量[%d],检查文档数量:[%d]\\n", e.Plan, e.KeysExamined, e.DocsExamined)
	fmt.Fprintf(builder, "\t 匹配数量[%d],修改数量[%d]\\n", e.Matched, e.Modified)
	fmt.Fprintf(builder, "\t 写冲突[%d],放弃[%d]\\n", e.WriteConflicts, e.Yields)
	fmt.Fprintf(builder, "\t全局[%s]\\n", e.Global.String())
	fmt.Fprintf(builder, "\t数据库[%s]\\n", e.DataBase.String())
	fmt.Fprintf(builder, "\t集合[%s]\\n", e.CollectionL.String())
	return builder.String()
}

type PlanSummary struct {
	IndexScan string
}
type LockLevel int

const (
	Global LockLevel = iota
	DataBase
	Collection
)

type LockKind int

const (
	Read LockKind = iota
	Write
	TryRead
	TryWrite
)

type LockInfo struct {
	Level             LockLevel
	AcquiredCount     LockDetail
	AcquiredWaitCount LockDetail
	AcquireTakes      LockDetail
	DeadLockTakes     LockDetail
}

func (l LockInfo) String() string {
	builder := &strings.Builder{}
	if str := l.AcquiredCount.String(false); str != "" {
		fmt.Fprintf(builder, "获取锁次数[%s]", str)
	}

	if str := l.AcquiredWaitCount.String(false); str != "" {
		fmt.Fprintf(builder, "获取锁等待次数[%s]", str)
	}
	if str := l.AcquireTakes.String(true); str != "" {
		fmt.Fprintf(builder, "获取锁耗时[%s]", str)
	}
	if str := l.DeadLockTakes.String(true); str != "" {
		fmt.Fprintf(builder, "死锁耗时[%s]", str)
	}
	return builder.String()
}

type LockDetail struct {
	Read     int64
	Write    int64
	TryRead  int64
	TryWrite int64
}

func (e LockDetail) String(takes bool) string {
	if e.Read+e.Write+e.TryWrite+e.TryRead == 0 {
		return ""
	}
	if takes {
		return fmt.Sprintf("等待读取锁耗时[%s],等待写锁耗时[%s],尝试获取读取锁耗时[%s],尝试获取写锁耗时[%s]", microString(e.Read), microString(e.Write), microString(e.TryRead), microString(e.TryWrite))
	} else {
		return fmt.Sprintf("获取读取锁[%d],获取写锁[%d],尝试获取读取锁[%d],尝试获取写锁[%d]", e.Read, e.Write, e.TryRead, e.TryWrite)
	}
}

func microString(value int64) string {
	return time.Duration(int64(time.Microsecond) * value).String()
}
