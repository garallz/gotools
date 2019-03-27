package logfile

var server *LogServer

// Initialization Func.
func LogInit(defalut string, ls []*LogStruct) *LogServer {
	server.logmap = make(map[string]*LogData)

	if len(ls) == 0 {
		d := &LogStruct{}
		server.defalut = d.checkStruct()
		// Open file to write and init cache.
		server.defalut.open()
		// sleep time to make new file open.
		go server.defalut.upFile()

	} else if len(ls) == 1 {
		server.defalut = ls[0].checkStruct()

		server.logmap["default"] = server.defalut
		// Open file to write and init cache.
		server.defalut.open()
		// sleep time to make new file open.
		go server.defalut.upFile()

	} else {
		if defalut == "" {
			defalut = ls[0].Name
		}

		for _, l := range ls {
			if l.Name == "" {
				panic("LogStruct Name can't be null")
			}

			// Check struct data.
			d := l.checkStruct()
			// Open file to write and init cache.
			d.open()

			server.logmap[l.Name] = d

			// TODO:
			// sleep time to make new file open.
			go d.upFile()
		}

		if l, ok := server.logmap[defalut]; ok {
			server.defalut = l
		} else {
			panic("LogStruct Name not have " + defalut)
		}
	}

	return server
}

type LogServer struct {
	logmap  map[string]*LogData
	defalut *LogData
}

type LogWrite struct {
	w      *LogData
	uid    string
	logmap map[string]*LogData
}

func NewLogWrite(uid string) *LogWrite {
	return &LogWrite{
		w:      server.defalut,
		uid:    uid,
		logmap: server.logmap,
	}
}

func (l *LogWrite) Name(name string) *LogWrite {
	if d, ok := l.logmap[name]; ok {
		return &LogWrite{
			w:      d,
			uid:    l.uid,
			logmap: l.logmap,
		}
	}
	return l
}
