package logg

import (
	"bufio"
	"encoding/json"
	"fgw_web_admin_panel/pkg/msg"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

const (
	levelInfo  LogLevel = "INFO"
	levelWarn  LogLevel = "WARN"
	levelError LogLevel = "ERROR"

	defaultPathToLog                          = "logs/"
	filenameForLog                            = "fgw"
	formatFileForLog                          = "json"
	defaultFilePermissions        os.FileMode = 0644
	maxQuantityFilesForLog                    = 7
	timePeriodicCleaningForLog                = 1 * time.Hour
	timePeriodicResetBufferForLog             = 3 * time.Second

	separatorIpAddress    = " | "
	defaultMaxStackFrames = 15

	colorGreen = "\033[32m"
)

type LogLevel string

type InfoPCEntry struct {
	Domain string `json:"domain"`
	IPAddr string `json:"ipAddr"`
}

type MessageEntry struct {
	Code    string  `json:"code"`
	Message string  `json:"msg"`
	Error   *string `json:"error,omitempty"`
}

type ResponseEntry struct {
	StatusCode int    `json:"statusCode"`
	MethodHTTP string `json:"methodHTTP"`
	URL        string `json:"url"`
}

type DetailEntry struct {
	FuncName   string `json:"funcName"`
	FileName   string `json:"fileName"`
	LineNumber int    `json:"lineNumber"`
	PathToFile string `json:"pathToFile"`
}

type LogEntry struct {
	DateTime string         `json:"dateTime"`
	InfoPC   InfoPCEntry    `json:"infoPC"`
	LevelLog LogLevel       `json:"levelLog"`
	Message  MessageEntry   `json:"message"`
	Response *ResponseEntry `json:"response"`
	Detail   *DetailEntry   `json:"detail"`
}

type Logger struct {
	file        *os.File
	currentDate string
	writer      *bufio.Writer
	ticker      *time.Ticker
	done        chan struct{}
	wg          sync.WaitGroup
	mu          sync.RWMutex
}

// NewLogger создает новый экземпляр.
func NewLogger() (*Logger, error) {
	if err := ensureLogDir(defaultPathToLog); err != nil {
		return nil, fmt.Errorf("%s: %s. %w", msg.EL5009, defaultPathToLog, err)
	}

	today := time.Now().Format("2006-01-02")
	logger := &Logger{
		currentDate: today,
		done:        make(chan struct{}),
	}

	filename := createFileNameForLog(defaultPathToLog, filenameForLog, today)

	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, defaultFilePermissions)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", msg.EL5000, err)
	}

	logger.file = file
	logger.writer = bufio.NewWriter(file)
	logger.ticker = time.NewTicker(timePeriodicResetBufferForLog)

	logger.wg.Add(1)
	go logger.flushPeriodically()

	go periodicCleaningLogFiles()

	return logger, nil
}

// periodicCleaning периодически очищает старые лог файлы.
func periodicCleaningLogFiles() {
	ticker := time.NewTicker(timePeriodicCleaningForLog)
	defer ticker.Stop()

	for range ticker.C {
		if err := cleanOldLogs(defaultPathToLog, maxQuantityFilesForLog); err != nil {
			log.Printf("%s: %s. %v", msg.EL5008, defaultPathToLog, err)
		}
	}
}

// flushPeriodically периодически сбрасывает буфер на диск.
func (l *Logger) flushPeriodically() {
	defer l.wg.Done()
	for {
		select {
		case <-l.ticker.C:
			l.flush()
		case <-l.done:
			l.flush()
			return
		}
	}
}

// flush сбрасывает буфер и синхронизирует файл (с защитой).
func (l *Logger) flush() {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.writer != nil {
		if err := l.writer.Flush(); err != nil {
			log.Printf("%s: %s. %v", msg.EL5012, defaultPathToLog, err)
		}
	}

	if l.file != nil {
		if err := l.file.Sync(); err != nil {
			log.Printf("%s: %s. %v", msg.EL5005, defaultPathToLog, err)
		}
	}
}

// rotateIfNeeded проверяет необходимость смены файла лога и выполняет ротацию.
func (l *Logger) rotateIfNeeded() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	today := time.Now().Format("2006-01-02")

	if l.currentDate == today && l.file != nil {
		return nil
	}

	if l.writer != nil {
		if err := l.writer.Flush(); err != nil {
			log.Printf("%s: %v", msg.EL5012, err)
		}
	}

	if l.file != nil {
		if err := l.file.Close(); err != nil {
			return fmt.Errorf("%s: %v", msg.EL5005, err)
		}
		l.file = nil
		l.writer = nil
	}

	filename := createFileNameForLog(defaultPathToLog, filenameForLog, today)

	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, defaultFilePermissions)
	if err != nil {
		return fmt.Errorf("%s: %w", msg.EL5000, err)
	}

	l.file = file
	l.writer = bufio.NewWriter(file)
	l.currentDate = today

	return nil
}

// loggCustom формирует и записывает структурированную запись лога в формате JSON.
func (l *Logger) loggCustom(levelLog LogLevel, msgEntry string, errMsg error, response *ResponseEntry, skipNumOfStackFrames int) {
	if err := l.rotateIfNeeded(); err != nil {
		log.Printf("%s: %s. %v", msg.EL5000, defaultPathToLog, err)
	}

	entry := &LogEntry{
		DateTime: time.Now().Format("2006-01-02 15:04:05"),
		InfoPC: InfoPCEntry{
			Domain: l.hostName(),
			IPAddr: l.ipAddr(),
		},
		LevelLog: levelLog,
		Message:  l.createdMsg(msgEntry, errMsg),
		Response: response,
		Detail:   getStackFrameInfo(skipNumOfStackFrames),
	}

	if err := l.writeEntry(entry); err != nil {
		log.Println(err)
	}
}

// LogI логирует информационное сообщение с указанием уровнем детализации стека.
//
// Параметры:
//   - message: текст информационного сообщения;
//   - skipNumOfStack: кол-во пропускаемых кадров стека.
func (l *Logger) LogI(message string, skipNumOfStack int) {
	fmt.Println(time.Now().Format(time.DateTime), levelInfo, message)
	l.loggCustom(levelInfo, message, nil, nil, skipNumOfStack)
}

// LogIf логирует информационное сообщение с указанием уровнем детализации стека с форматированием.
//
// Параметры:
//   - skipNumOfStack: количество пропускаемых кадров стека;
//   - format: строка формата;
//   - args: аргументы для форматирования.
func (l *Logger) LogIf(skipNumOfStack int, format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	fmt.Println(time.Now().Format(time.DateTime), levelInfo, message)
	l.loggCustom(levelInfo, message, nil, nil, skipNumOfStack)
}

// LogW логирует предупреждающие сообщения с указанием уровнем детализации стека.
//
// Параметры:
//   - message: текст предупреждающего сообщения;
//   - skipNumOfStack: кол-во пропускаемых кадров стека.
func (l *Logger) LogW(message string, skipNumOfStack int) {
	fmt.Println(time.Now().Format(time.DateTime), levelWarn, message)
	l.loggCustom(levelWarn, message, nil, nil, skipNumOfStack)
}

// LogWf логирует предупреждающие сообщения с указанием уровнем детализации стека с форматированием.
//
// Параметры:
//   - skipNumOfStack: количество пропускаемых кадров стека;
//   - format: строка формата;
//   - args: аргументы для форматирования.
func (l *Logger) LogWf(skipNumOfStack int, format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	fmt.Println(time.Now().Format(time.DateTime), levelWarn, message)
	l.loggCustom(levelWarn, message, nil, nil, skipNumOfStack)
}

// LogE логирует сообщение об ошибки с указанным уровнем детализации стека.
//
// Параметры:
//   - message: текстовое описание ошибки;
//   - errMsg: объект ошибки;
//   - skipNumOfStack: кол-во пропускаемых кадров стека.
func (l *Logger) LogE(message string, errMsg error, skipNumOfStack int) {
	fmt.Println(time.Now().Format(time.DateTime), levelError, message, errMsg)
	l.loggCustom(levelError, message, errMsg, nil, skipNumOfStack)
}

// LogEf логирует сообщение об ошибки с указанным уровнем детализации стека с форматированием.
//
// Параметры:
//   - skipNumOfStack: кол-во пропускаемых кадров стека;
//   - errMsg: объект ошибки;
//   - format: строка формата;
//   - args: аргументы для форматирования.
func (l *Logger) LogEf(skipNumOfStack int, errMsg error, format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	fmt.Println(time.Now().Format(time.DateTime), levelError, message, errMsg)
	l.loggCustom(levelError, message, errMsg, nil, skipNumOfStack)
}

// LogHttpI логирует успешный HTTP-запрос.
//
// Параметры:
//   - statusCode: HTTP статус ответа (< 400);
//   - methodHTTP: HTTP методы (PUT, POST, DELETE, GET);
//   - url: запрашиваемый url;
//   - message: текстовое описание;
//   - skipNumOfStack: кол-во пропускаемых кадров стека.
func (l *Logger) LogHttpI(statusCode int, methodHTTP, url, message string, skipNumOfStack int) {
	responseEntry := &ResponseEntry{
		StatusCode: statusCode,
		MethodHTTP: methodHTTP,
		URL:        url,
	}
	fmt.Println(time.Now().Format(time.DateTime), levelInfo, message, responseEntry)
	l.loggCustom(levelInfo, message, nil, responseEntry, skipNumOfStack)
}

// LogHttpIf логирует успешный HTTP-запрос с форматированием.
//
// Параметры:
//   - statusCode: HTTP статус ответа (< 400);
//   - methodHTTP: HTTP методы (PUT, POST, DELETE, GET);
//   - url: запрашиваемый url;
//   - skipNumOfStack: кол-во пропускаемых кадров стека.
//   - format: строка формата;
//   - args: аргументы для форматирования.
func (l *Logger) LogHttpIf(statusCode int, methodHTTP, url string, skipNumOfStack int, format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	responseEntry := &ResponseEntry{
		StatusCode: statusCode,
		MethodHTTP: methodHTTP,
		URL:        url,
	}
	fmt.Println(time.Now().Format(time.DateTime), levelInfo, message, responseEntry)
	l.loggCustom(levelInfo, message, nil, responseEntry, skipNumOfStack)
}

// LogHttpE логирует ошибочный HTTP-запрос.
//
// Параметры:
//   - statusCode: HTTP статус ответа (>= 400);
//   - methodHTTP: HTTP методы (PUT, POST, DELETE, GET);
//   - url: запрашиваемый url;
//   - message: текстовое описание ошибки;
//   - errMsg: объект ошибки (может быть nil);
//   - skipNumOfStack: кол-во пропускаемых кадров стека.
func (l *Logger) LogHttpE(statusCode int, methodHTTP, url, message string, errMsg error, skipNumOfStack int) {
	responseEntry := &ResponseEntry{
		StatusCode: statusCode,
		MethodHTTP: methodHTTP,
		URL:        url,
	}
	fmt.Println(time.Now().Format(time.DateTime), levelError, message, errMsg, responseEntry)
	l.loggCustom(levelError, message, errMsg, responseEntry, skipNumOfStack)
}

// LogHttpEf логирует ошибочный HTTP-запрос с форматированием.
//
// Параметры:
//   - statusCode: HTTP статус ответа (>= 400);
//   - methodHTTP: HTTP методы (PUT, POST, DELETE, GET);
//   - url: запрашиваемый url;
//   - errMsg: объект ошибки (может быть nil);
//   - skipNumOfStack: кол-во пропускаемых кадров стека;
//   - format: строка формата;
//   - args: аргументы для форматирования.
func (l *Logger) LogHttpEf(statusCode int, methodHTTP, url string, errMsg error, skipNumOfStack int, format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	responseEntry := &ResponseEntry{
		StatusCode: statusCode,
		MethodHTTP: methodHTTP,
		URL:        url,
	}
	fmt.Println(time.Now().Format(time.DateTime), levelError, message, errMsg, responseEntry)
	l.loggCustom(levelError, message, errMsg, responseEntry, skipNumOfStack)
}

// writeEntry записывает структурированный лог в файл в формате JSON.
//
// Параметры:
//   - entry: указатель на структуру лога, содержащую все данные для записи.
func (l *Logger) writeEntry(entry *LogEntry) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	data, err := json.MarshalIndent(entry, "", " ")
	if err != nil {
		return fmt.Errorf("%s: %w", msg.EL5003, err)
	}

	data = append(data, ',', '\n')

	if _, err := l.writer.Write(data); err != nil {
		return fmt.Errorf("%s: %w", msg.EL5004, err)
	}

	if _, err := l.file.Write(data); err != nil {
		return fmt.Errorf("%s: %w", msg.EL5004, err)
	}

	if err := l.writer.Flush(); err != nil {
		log.Printf("%s: %v", msg.EL5012, err)
	}

	return nil
}

// hostName возвращает имя текущего хоста (компьютера/сервера).
func (l *Logger) hostName() string {
	hostname, err := os.Hostname()
	if err != nil {
		return fmt.Sprintf("%s: %v", msg.EL5001, err)
	}

	return hostname
}

// ipAddr возвращает строковое представление всех IP-адресов текущего хоста.
func (l *Logger) ipAddr() string {
	ips, err := net.LookupIP(l.hostName())
	if err != nil {
		return fmt.Sprintf("%s: %v", msg.EL5002, err)
	}

	ipsStr := make([]string, 0)
	for _, ip := range ips {
		ipsStr = append(ipsStr, ip.String())
	}

	return strings.Join(ipsStr, separatorIpAddress)
}

// createdMsg создает структуру MessageEntry на основе входной строки и ошибки.
func (l *Logger) createdMsg(msgEntry string, err error) MessageEntry {
	code, message := l.splitCodeMessage(msgEntry)

	var errStr *string
	if err != nil {
		errMsg := err.Error()
		errStr = &errMsg
	}

	return MessageEntry{
		Code:    code,
		Message: message,
		Error:   errStr,
	}
}

// splitCodeMessage разделяет входную строку на код и сообщения.
func (l *Logger) splitCodeMessage(message string) (string, string) {
	if message == "" {
		return "", fmt.Sprintf(msg.WL4000)
	}

	spaceIndex := strings.Index(message, " ")
	if spaceIndex == -1 {
		return message, fmt.Sprintf(msg.WL4001)
	}

	code := message[:spaceIndex]
	msgWithoutCode := strings.TrimSpace(message[spaceIndex+1:])

	return code, msgWithoutCode
}

// getStackFrameInfo возвращает информацию о месте вызова в стеке.
//
// Параметры:
//   - skipNumOfStack: количество кадров для пропуска (0 - сама функция, 1 - вызывающая и т.д.)
func getStackFrameInfo(skipNumOfStack int) *DetailEntry {
	pc := make([]uintptr, defaultMaxStackFrames)
	frameCount := runtime.Callers(skipNumOfStack, pc)
	if frameCount == 0 {
		return &DetailEntry{
			FuncName:   fmt.Sprint("неизвестно"),
			FileName:   fmt.Sprint("неизвестно"),
			LineNumber: 0,
			PathToFile: "",
		}
	}

	frames := runtime.CallersFrames(pc[:frameCount])
	frame, ok := frames.Next()
	if !ok {
		return &DetailEntry{
			FuncName:   fmt.Sprint("неизвестно"),
			FileName:   fmt.Sprint("неизвестно"),
			LineNumber: 0,
			PathToFile: "",
		}
	}

	idxFile := strings.LastIndexByte(frame.File, '/')
	fileName := frame.File[idxFile+1:]

	return &DetailEntry{
		FuncName:   frame.Function,
		FileName:   fileName,
		LineNumber: frame.Line,
		PathToFile: frame.File,
	}
}

// ensureLogDir обеспечивает наличие директории для лог-файла.
func ensureLogDir(filePath string) error {
	dir := filepath.Dir(filePath)

	var currentDir = "." // текущая директория

	if dir != currentDir && dir != "" {
		if err := os.MkdirAll(dir, defaultFilePermissions); err != nil {
			return fmt.Errorf("%s: %w", msg.EL5006, err)
		}
	}

	return nil
}

// createFileNameForLog создает имя файла для журнала.
func createFileNameForLog(dir, filename string, currentDate string) string {
	return filepath.Join(dir, fmt.Sprintf("%s_%s.%s", filename, currentDate, formatFileForLog))
}

// cleanOldLogs очищает старого файла с логами.
func cleanOldLogs(dir string, maxFiles int) error {
	files, err := filepath.Glob(filepath.Join(dir, "*.json"))
	if err != nil {
		return fmt.Errorf("%s: %w", msg.EL5007, err)
	}

	if len(files) >= maxFiles {
		err = os.Remove(files[0])
		if err != nil {
			return fmt.Errorf("%s: %w", msg.EL5007, err)
		}
		log.Printf("%s: %s", msg.IL2001, files[0])
	}

	return nil
}

func (l *Logger) Close() {
	if l.done != nil {
		close(l.done)
	}

	l.wg.Wait()

	var errs []string

	l.mu.Lock()
	defer l.mu.Unlock()

	if l.writer != nil {
		if err := l.writer.Flush(); err != nil {
			errs = append(errs, fmt.Sprintf("%s: %v", msg.EL5012, err))
		}
		l.writer = nil
	}

	if l.file != nil {
		if err := l.file.Sync(); err != nil {
			errs = append(errs, fmt.Sprintf("%s: %v", msg.EL5011, err))
		}
		if err := l.file.Close(); err != nil {
			errs = append(errs, fmt.Sprintf("%s: %v", msg.EL5005, err))
		}
		l.file = nil
	}

	if len(errs) > 0 {
		log.Printf("%s: %s", msg.EL5013, strings.Join(errs, ", "))
	} else {
		log.Printf("%s%s", colorGreen, msg.IL2000)
	}
}
