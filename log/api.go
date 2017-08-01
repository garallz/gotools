package llz_log

func (l *LogStruct) Init() {
	l.open()
}

func (l *LogStruct) WriteInfo(messages ...interface{}) error {
	if l.Level == LevelInfo {
		return l.put(" INFO  ", messages)
	}
	return nil
}

func (l *LogStruct) WriteDebug(messages ...interface{}) error {
	if l.Level <= LevelDebug {
		return l.put(" DEBUG ", messages)
	}
	return nil
}

func (l *LogStruct) WriteWarn(messages ...interface{}) error {
	if l.Level <= LevelWarn {
		return l.put(" WARN  ", messages)
	}
	return nil
}

func (l *LogStruct) WriteError(messages ...interface{}) error {
	if l.Level <= LevelError {
		return l.put(" ERROR ", messages)
	}
	return nil
}

func (l *LogStruct) WriteFatal(err error, messages ...interface{}) {
	if l.Level <= LevelFatal {
		l.put(" FATAL ", messages)
		panic(err)
	}
}
