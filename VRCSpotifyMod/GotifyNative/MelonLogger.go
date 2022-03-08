package main

type ConsoleColor uintptr

const (
	ConsoleColorBlack       ConsoleColor = 0
	ConsoleColorDarkBlue    ConsoleColor = 1
	ConsoleColorDarkGreen   ConsoleColor = 2
	ConsoleColorDarkCyan    ConsoleColor = 3
	ConsoleColorDarkRed     ConsoleColor = 4
	ConsoleColorDarkMagenta ConsoleColor = 5
	ConsoleColorDarkYellow  ConsoleColor = 6
	ConsoleColorGray        ConsoleColor = 7
	ConsoleColorDarkGray    ConsoleColor = 8
	ConsoleColorBlue        ConsoleColor = 9
	ConsoleColorGreen       ConsoleColor = 10
	ConsoleColorCyan        ConsoleColor = 11
	ConsoleColorRed         ConsoleColor = 12
	ConsoleColorMagenta     ConsoleColor = 13
	ConsoleColorYellow      ConsoleColor = 14
	ConsoleColorWhite       ConsoleColor = 15
)

var mono, funnyError = NewMono()

type LoggerInstance struct {
	pointer_      uintptr
	handle_       uintptr
	classPointer_ uintptr

	name_ string

	ctorString_ uintptr

	msgObj_        uintptr
	msgString_     uintptr
	warningString_ uintptr
	errorString_   uintptr
}

func NewLoggerInstance(name string) (instance *LoggerInstance, err error) {
	instance = &LoggerInstance{}
	instance.name_ = name

	MelonLoader, err2 := mono.GetAssemblyImage("MelonLoader")
	if err2 != nil {
		return nil, err2
	}
	if instance.classPointer_ == 0 {
		instance.classPointer_, err2 = mono.GetClass(MelonLoader, "MelonLoader", "MelonLogger/Instance")
		if err2 != nil {
			return nil, err2
		}
	}

	if instance.ctorString_ == 0 {
		instance.ctorString_, err2 = mono.FindMethod(instance.classPointer_, ":.ctor(string)", false)
		if err2 != nil {
			return nil, err2
		}
	}

	nameString, stringHandle, err := mono.NewString(name)
	if err != nil {
		return nil, err
	}

	defer mono.mono_gchandle_free_(stringHandle)

	instance.pointer_, err2 = mono.NewObject(instance.classPointer_)
	if err2 != nil {
		return nil, err2
	}

	instance.handle_, _, err = mono.mono_gchandle_new_(instance.pointer_)
	if err != nil && instance.handle_ == 0 {
		return nil, err
	}

	args := make([]uintptr, 1)
	args[0] = nameString
	mono.RuntimeInvoke(instance.ctorString_, instance.pointer_, args)

	return instance, nil
}

func (instance *LoggerInstance) Name() string {
	return instance.name_
}

func (instance *LoggerInstance) Destroy() {
	mono.mono_gchandle_free_(instance.handle_)
}

func (instance *LoggerInstance) MsgObj(obj uintptr) (err error) {
	if instance.msgObj_ == 0 {
		instance.msgObj_, err = mono.FindMethod(instance.classPointer_, ":Msg(object)", false)
		if err != nil {
			return err
		}
	}

	args := make([]uintptr, 1)
	args[0] = obj

	mono.RuntimeInvoke(instance.msgObj_, instance.pointer_, args)
	return nil
}

func (instance *LoggerInstance) MsgString(str string) (err error) {
	if instance.msgString_ == 0 {
		instance.msgString_, err = mono.FindMethod(instance.classPointer_, ":Msg(string)", false)
		if err != nil {
			return err
		}
	}

	args := make([]uintptr, 1)
	var handle uintptr
	args[0], handle, err = mono.NewString(str)
	if err != nil {
		return err
	}
	defer mono.mono_gchandle_free_(handle)

	mono.RuntimeInvoke(instance.msgString_, instance.pointer_, args)
	return nil
}

func (instance *LoggerInstance) WarningString(str string) (err error) {
	if instance.warningString_ == 0 {
		instance.warningString_, err = mono.FindMethod(instance.classPointer_, ":Warning(string)", false)
		if err != nil {
			return err
		}
	}

	args := make([]uintptr, 1)
	var handle uintptr
	args[0], handle, err = mono.NewString(str)
	if err != nil {
		return err
	}
	defer mono.mono_gchandle_free_(handle)

	mono.RuntimeInvoke(instance.warningString_, instance.pointer_, args)
	return nil
}

func (instance *LoggerInstance) ErrorString(str string) (err error) {
	if instance.errorString_ == 0 {
		instance.errorString_, err = mono.FindMethod(instance.classPointer_, ":Error(string)", false)
		if err != nil {
			return err
		}
	}

	args := make([]uintptr, 1)
	var handle uintptr
	args[0], handle, err = mono.NewString(str)
	if err != nil {
		return err
	}
	defer mono.mono_gchandle_free_(handle)
	mono.RuntimeInvoke(instance.errorString_, instance.pointer_, args)
	return nil
}
