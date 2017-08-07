package llz_log

// Initialization Func.
func (l *LogStruct) Init() {
	// Check struct data.
	l.checkStruct()
	// Open file to write and init cache.
	l.open()
	// init cache data.
	if l.Cache {
		l.buf = l.buf[:0]
	}
	go l.upFile()
}

// Write log data with log level was Info.
func (l *LogStruct) WriteInfo(messages ...interface{}) error {
	if l.Level == LevelInfo {
		return l.put(" INFO  ", messages)
	}
	return nil
}

// Write log data with log level was Debug.
func (l *LogStruct) WriteDebug(messages ...interface{}) error {
	if l.Level <= LevelDebug {
		return l.put(" DEBUG ", messages)
	}
	return nil
}

// Write log data with log level was Warn.
func (l *LogStruct) WriteWarn(messages ...interface{}) error {
	if l.Level <= LevelWarn {
		return l.put(" WARN  ", messages)
	}
	return nil
}

// Write log data with log level was Error.
func (l *LogStruct) WriteError(messages ...interface{}) error {
	if l.Level <= LevelError {
		return l.put(" ERROR ", messages)
	}
	return nil
}

// Write log data with log level was Fatal.
func (l *LogStruct) WriteFatal(err error, messages ...interface{}) {
	if l.Level <= LevelFatal {
		l.put(" FATAL ", messages)
		panic(err)
	}
}

// Write byte log data.
func (l *LogStruct) WriteByte(bts []byte) error {
	return l.putByte(append(bts, []byte("\n")...))
}
