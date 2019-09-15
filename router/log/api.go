package logfile

import (
	"time"
)

// Write log data with log level was Debug.
func (l *LogWrite) WriteDebug(args ...interface{}) error {
	if l.w.level == LevelDebug {
		return l.put(LevelDebug, args)
	}
	return nil
}

// Write log data with log level was Info.
func (l *LogWrite) WriteInfo(args ...interface{}) error {
	if l.w.level <= LevelInfo {
		return l.put(LevelInfo, args)
	}
	return nil
}

// Write log data with log level was Warn.
func (l *LogWrite) WriteWarn(args ...interface{}) error {
	if l.w.level <= LevelWarn {
		return l.put(LevelWarn, args)
	}
	return nil
}

// Write log data with log level was Error.
func (l *LogWrite) WriteError(args ...interface{}) error {
	return l.put(LevelError, args)
}

// Write log data with log level was Fatal.
func (l *LogWrite) WriteFatal(args ...interface{}) error {
	return l.put(LevelDebug, args)
}

func (l *LogWrite) WritePanic(err error, args ...interface{}) {
	l.put(LevelPanic, args)
	// wirter log in file and close
	l.putPanic(nil)
	panic(err)
}

// Write log data with log level was Debug.
func (l *LogWrite) WriteDebugf(messages string, args ...interface{}) error {

	return l.putf(LevelDebug, messages, args)
}

// Write log data with log level was Info.
func (l *LogWrite) WriteInfof(messages string, args ...interface{}) error {
	if l.w.level <= LevelInfo {
		return l.putf(LevelInfo, messages, args)
	}
	return nil
}

// Write log data with log level was Warn.
func (l *LogWrite) WriteWarnf(messages string, args ...interface{}) error {
	if l.w.level <= LevelWarn {
		return l.putf(LevelWarn, messages, args)
	}
	return nil
}

// Write log data with log level was Error.
func (l *LogWrite) WriteErrorf(messages string, args ...interface{}) error {
	return l.putf(LevelError, messages, args)
}

// Write log data with log level was Fatal.
func (l *LogWrite) WriteFatalf(messages string, args ...interface{}) error {
	return l.putf(LevelFatal, messages, args)
}

func (l *LogWrite) WritePanicf(err error, messages string, args ...interface{}) {
	// wirter log in file and close
	l.putf(LevelPanic, messages, args)
	l.putPanic(nil)
	panic(err)
}

// Write byte log data, not prefix.
func (l *LogWrite) WriteBytes(bts []byte) error {
	return l.w.putByte(bts)
}

func (l *LogWrite) WriteJsonOrgin(data interface{}) error {
	return l.w.json(data)
}

func (l *LogWrite) WriteJson(level LogLevel, data interface{}) error {
	return l.w.json(&JsonStruct{
		Time:  time.Now().Format(l.w.format),
		Uid:   l.uid,
		Level: levelMsg[level],
		Data:  data,
	})
}

// stop server need write log data to file and close
func CloseLogServer() {
	for _, l := range server.logmap {
		if l.cache {
			l.file.Write(l.buf)
		}
		l.file.Close()
	}
}
