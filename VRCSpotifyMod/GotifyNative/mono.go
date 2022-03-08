package main

//#cgo LDFLAGS:
//#include <stdio.h>
//#include <stdlib.h>
import "C"
import (
	"syscall"
	"unsafe"
)

type Mono struct {
	monoLib                            *syscall.LazyDLL
	mono_get_root_domain_              func(a ...uintptr) (r1 uintptr, r2 uintptr, err error)
	mono_thread_attach_                func(a ...uintptr) (r1 uintptr, r2 uintptr, err error)
	mono_domain_assembly_open_         func(a ...uintptr) (r1 uintptr, r2 uintptr, err error)
	mono_assembly_get_image_           func(a ...uintptr) (r1 uintptr, r2 uintptr, err error)
	mono_class_from_name_              func(a ...uintptr) (r1 uintptr, r2 uintptr, err error)
	mono_method_desc_new_              func(a ...uintptr) (r1 uintptr, r2 uintptr, err error)
	mono_method_desc_search_in_class_  func(a ...uintptr) (r1 uintptr, r2 uintptr, err error)
	mono_class_get_method_from_name_   func(a ...uintptr) (r1 uintptr, r2 uintptr, err error)
	mono_string_new_                   func(a ...uintptr) (r1 uintptr, r2 uintptr, err error)
	mono_object_new_                   func(a ...uintptr) (r1 uintptr, r2 uintptr, err error)
	mono_runtime_invoke_               func(a ...uintptr) (r1 uintptr, r2 uintptr, err error)
	mono_object_get_class_             func(a ...uintptr) (r1 uintptr, r2 uintptr, err error)
	mono_object_to_string_             func(a ...uintptr) (r1 uintptr, r2 uintptr, err error)
	mono_property_get_get_method_      func(a ...uintptr) (r1 uintptr, r2 uintptr, err error)
	mono_runtime_object_init_          func(a ...uintptr) (r1 uintptr, r2 uintptr, err error)
	mono_runtime_object_init_checked_  func(a ...uintptr) (r1 uintptr, r2 uintptr, err error)
	mono_image_get_name_               func(a ...uintptr) (r1 uintptr, r2 uintptr, err error)
	mono_class_get_name_               func(a ...uintptr) (r1 uintptr, r2 uintptr, err error)
	mono_class_get_namespace_          func(a ...uintptr) (r1 uintptr, r2 uintptr, err error)
	mono_class_get_property_from_name_ func(a ...uintptr) (r1 uintptr, r2 uintptr, err error)
	mono_assembly_get_name_            func(a ...uintptr) (r1 uintptr, r2 uintptr, err error)
	mono_gchandle_new_                 func(a ...uintptr) (r1 uintptr, r2 uintptr, err error)
	mono_gchandle_free_                func(a ...uintptr) (r1 uintptr, r2 uintptr, err error)
	Domain                             uintptr
	ImageCache                         map[string]uintptr
}

func NewMono() (mono *Mono, err error) {
	mono = &Mono{}
	mono.monoLib = syscall.NewLazyDLL("MelonLoader/Dependencies/MonoBleedingEdge.x64/mono-2.0-bdwgc.dll")
	mono.mono_get_root_domain_ = mono.monoLib.NewProc("mono_get_root_domain").Call
	mono.mono_thread_attach_ = mono.monoLib.NewProc("mono_thread_attach").Call
	mono.mono_domain_assembly_open_ = mono.monoLib.NewProc("mono_domain_assembly_open").Call
	mono.mono_assembly_get_image_ = mono.monoLib.NewProc("mono_assembly_get_image").Call
	mono.mono_class_from_name_ = mono.monoLib.NewProc("mono_class_from_name").Call
	mono.mono_method_desc_new_ = mono.monoLib.NewProc("mono_method_desc_new").Call
	mono.mono_method_desc_search_in_class_ = mono.monoLib.NewProc("mono_method_desc_search_in_class").Call
	mono.mono_class_get_method_from_name_ = mono.monoLib.NewProc("mono_class_get_method_from_name").Call
	mono.mono_string_new_ = mono.monoLib.NewProc("mono_string_new").Call
	mono.mono_object_new_ = mono.monoLib.NewProc("mono_object_new").Call
	mono.mono_runtime_invoke_ = mono.monoLib.NewProc("mono_runtime_invoke").Call
	mono.mono_object_get_class_ = mono.monoLib.NewProc("mono_object_get_class").Call
	mono.mono_object_to_string_ = mono.monoLib.NewProc("mono_object_to_string").Call
	mono.mono_property_get_get_method_ = mono.monoLib.NewProc("mono_property_get_get_method").Call
	mono.mono_runtime_object_init_ = mono.monoLib.NewProc("mono_runtime_object_init").Call
	mono.mono_runtime_object_init_checked_ = mono.monoLib.NewProc("mono_runtime_object_init_checked").Call
	mono.mono_image_get_name_ = mono.monoLib.NewProc("mono_image_get_name").Call
	mono.mono_class_get_name_ = mono.monoLib.NewProc("mono_class_get_name").Call
	mono.mono_class_get_namespace_ = mono.monoLib.NewProc("mono_class_get_namespace").Call
	mono.mono_class_get_property_from_name_ = mono.monoLib.NewProc("mono_class_get_property_from_name").Call
	mono.mono_assembly_get_name_ = mono.monoLib.NewProc("mono_assembly_get_name").Call
	mono.mono_gchandle_new_ = mono.monoLib.NewProc("mono_gchandle_new").Call
	mono.mono_gchandle_free_ = mono.monoLib.NewProc("mono_gchandle_free").Call

	mono.ImageCache = make(map[string]uintptr)

	domain_, _, err := mono.mono_get_root_domain_()
	if err != nil && domain_ == 0 {
		return nil, err
	}

	mono.Domain = domain_

	return mono, nil
}

func (mono *Mono) GetRootDomain() (ret uintptr, err error) {
	if mono.Domain != 0 {
		return mono.Domain, nil
	}

	ret, _, err = mono.mono_get_root_domain_()
	if err != nil && ret == 0 {
		return 0, err
	}

	mono.Domain = ret

	return ret, nil
}

func (mono *Mono) FindMethod(klass uintptr, name string, includeNameSpace bool) (ret uintptr, err error) {
	desc, err := mono.NewMethodDesc(name, includeNameSpace)
	if err != nil {
		return 0, err
	}

	ret, err = mono.SearchDescInClass(desc, klass)
	if err != nil {
		return 0, err
	}

	return ret, nil
}

func (mono *Mono) ThreadAttach() (err error) {
	ret, _, err := mono.mono_thread_attach_(mono.Domain)
	if err != nil && ret == 0 {
		return err
	}

	return nil
}

func (mono *Mono) GetAssemblyImage(name string) (ret uintptr, err error) {

	if image, ok := mono.ImageCache[name]; ok && image != 0 {
		return image, nil
	}

	assembly, err := mono.OpenAssembly(name)
	if err != nil {
		return 0, err
	}

	ret, err = mono.GetImage(assembly)
	if err != nil {
		return 0, err
	}

	mono.ImageCache[name] = ret

	return ret, nil
}

func (mono *Mono) OpenAssembly(name string) (assembly uintptr, err error) {
	str := C.CString(name)
	ptr := uintptr(unsafe.Pointer(str))
	defer C.free(unsafe.Pointer(str))

	assembly, _, err = mono.mono_domain_assembly_open_(mono.Domain, ptr)
	if err != nil && assembly == 0 {
		return 0, err
	}

	return assembly, nil
}

func (mono *Mono) GetImage(assembly uintptr) (image uintptr, err error) {

	image, _, err = mono.mono_assembly_get_image_(assembly)
	if err != nil && image == 0 {
		return 0, err
	}

	return image, nil
}

func (mono *Mono) GetClass(image uintptr, namespace string, name string) (klass uintptr, err error) {
	str1 := C.CString(namespace)
	str2 := C.CString(name)
	ptr1 := uintptr(unsafe.Pointer(str1))
	ptr2 := uintptr(unsafe.Pointer(str2))
	defer C.free(unsafe.Pointer(str1))
	defer C.free(unsafe.Pointer(str2))

	klass, _, err = mono.mono_class_from_name_(image, ptr1, ptr2)
	if err != nil && klass == 0 {
		return 0, err
	}

	return klass, nil
}

func (mono *Mono) NewMethodDesc(method string, includeNamespace bool) (desc uintptr, err error) {
	str := C.CString(method)
	ptr := uintptr(unsafe.Pointer(str))

	var ptr2 uintptr

	if includeNamespace {
		ptr2 = 1
	} else {
		ptr2 = 0
	}

	defer C.free(unsafe.Pointer(str))

	desc, _, err = mono.mono_method_desc_new_(ptr, ptr2)
	if err != nil && desc == 0 {
		return 0, err
	}

	return desc, nil
}

func (mono *Mono) SearchDescInClass(klass uintptr, desc uintptr) (method uintptr, err error) {
	method, _, err = mono.mono_method_desc_search_in_class_(klass, desc)
	if err != nil && method == 0 {
		return 0, err
	}

	return method, nil
}

func (mono *Mono) GetMethod(klass uintptr, name string, paramCount uintptr) (method uintptr, err error) {
	str := C.CString(name)
	ptr := uintptr(unsafe.Pointer(str))
	defer C.free(unsafe.Pointer(str))

	method, _, err = mono.mono_class_get_method_from_name_(klass, ptr, paramCount)
	if err != nil && method == 0 {
		return 0, err
	}

	return method, nil
}

func (mono *Mono) NewObject(klass uintptr) (object uintptr, err error) {
	object, _, err = mono.mono_object_new_(mono.Domain, klass)
	if err != nil && object == 0 {
		return 0, err
	}

	return object, nil
}

func (mono *Mono) NewString(str string) (object uintptr, gcHandle uintptr, err error) {
	str1 := C.CString(str + "\000")
	ptr := uintptr(unsafe.Pointer(str1))

	object, _, err = mono.mono_string_new_(mono.Domain, ptr)
	if err != nil && object == 0 {
		return 0, 0, err
	}

	gcHandle, _, err = mono.mono_gchandle_new_(object)
	if err != nil && gcHandle == 0 {
		return 0, 0, err
	}

	return object, gcHandle, nil
}

func (mono *Mono) RuntimeInvoke(method uintptr, obj uintptr, params []uintptr) (ret uintptr, err error) {
	var args uintptr
	if params == nil {
		args = 0
	} else {
		args = uintptr(unsafe.Pointer(&params[0]))
	}
	ret, _, err = mono.mono_runtime_invoke_(method, obj, args)
	if err != nil && ret == 0 {
		return 0, err
	}

	return ret, nil
}
