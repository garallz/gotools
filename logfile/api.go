package logfile

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
func (l *LogData) WriteDebug(messages ...interface{}) error {
	if l.level == LevelDebug {
		return l.put(" DEBUG ", messages)
	}
	return nil
}

// Write log data with log level was Info.
func (l *LogData) WriteInfo(messages ...interface{}) error {
	if l.level <= LevelInfo {
		return l.put(" INFO  ", messages)
	}
	return nil
}

// Write log data with log level was Warn.
func (l *LogData) WriteWarn(messages ...interface{}) error {
	if l.level <= LevelWarn {
		return l.put(" WARN  ", messages)
	}
	return nil
}

// Write log data with log level was Error.
func (l *LogData) WriteError(messages ...interface{}) error {
	if l.level <= LevelError {
		return l.put(" ERROR ", messages)
	}
	return nil
}

// Write log data with log level was Fatal.
func (l *LogData) WriteFatal(err error, messages ...interface{}) {
	if l.level <= LevelFatal {
		l.put(" FATAL ", messages)
		panic(err)
	}
}

// Write byte log data, not prefix.
func (l *LogData) WriteByte(bts []byte) error {
	return l.putByte(append(bts, []byte("\n")...))
}

// Change log error level.
func (l *LogData) ChangeErrLevel(level LogLevel) {
	l.level = level
}
