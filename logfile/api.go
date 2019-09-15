package logfile

import(
	"encoding/json"
)

// Initialization Func.
func LogInit(l *LogStruct) *LogData {
	// Check struct data.
	d := l.checkStruct()
	// Open file to write and init cache.
	d.open()
	// sleep time to make new file open.
	go d.upFile()

	return d
}

// Write log data with log level was Debug.
func (l *LogData) WriteDebug(args ...interface{}) error {
	if l.level == LevelDebug {
		return l.put(" DEBUG\t", args)
	}
	return nil
}

// Write log data with log level was Info.
func (l *LogData) WriteInfo(args ...interface{}) error {
	if l.level <= LevelInfo {
		return l.put(" INFO\t", args)
	}
	return nil
}

// Write log data with log level was Warn.
func (l *LogData) WriteWarn(args ...interface{}) error {
	if l.level <= LevelWarn {
		return l.put(" WARN\t", args)
	}
	return nil
}

// Write log data with log level was Error.
func (l *LogData) WriteError(args ...interface{}) error {
		return l.put(" ERROR\t", args)
}

// Write log data with log level was Fatal.
func (l *LogData) WriteFatal(args ...interface{}) error {
	return l.put(" FATAL\t", args)
}

func (l *LogData) WritePanic(err error, args ...interface{}) {
	l.put(" FATAL\t", args)
	// wirter log in file and close
	l.putPanic(nil)
	panic(err)
}

// Write log data with log level was Debug.
func (l *LogData) WriteDebugf(messages string, args ...interface{}) error {
	if l.level == LevelDebug {
		return l.putf(" DEBUG\t"+messages, args)
	}
	return nil
}

// Write log data with log level was Info.
func (l *LogData) WriteInfof(messages string, args ...interface{}) error {
	if l.level <= LevelInfo {
		return l.putf(" INFO\t"+messages, args)
	}
	return nil
}

// Write log data with log level was Warn.
func (l *LogData) WriteWarnf(messages string, args ...interface{}) error {
	if l.level <= LevelWarn {
		return l.putf(" WARN\t"+messages, args)
	}
	return nil
}

// Write log data with log level was Error.
func (l *LogData) WriteErrorf(messages string, args ...interface{}) error {
	return l.putf(" ERROR\t"+messages, args)
}

// Write log data with log level was Fatal.
func (l *LogData) WriteFatalf(messages string, args ...interface{}) error {
	return l.putf(" FATAL\t"+messages, args)
}

func (l *LogData) WritePanicf(err error, messages string, args ...interface{}) {
	// wirter log in file and close
	l.putf(" PANIC\t"+messages, args)
	l.putPanic(nil)
	panic(err)
}

// Write byte log data, not prefix.
func (l *LogData) WriteBytes(bts []byte) error {
	return l.putByte(append(bts, []byte("\n")...))
}

func (l *LogData) WriterJson(data interface{}) error {
	if bts, err := json.Marshal(data); err != nil {
		return err
	} else {
		return l.putByte(append(bts, []byte("\n")...))
	}
}

// Change log error level.
func (l *LogData) ChangeErrLevel(level LogLevel) {
	l.level = level
}
