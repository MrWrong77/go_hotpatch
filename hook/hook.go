package hook

import (
	"reflect"
	"syscall"
	"unsafe"
)

// reflect.Value未导出的替代品
type value struct {
	_   uintptr
	ptr unsafe.Pointer
}

func UintptrToFuncPtr[T any](ptr uintptr) T {
	return *(*T)(unsafe.Pointer(ptr))
}

func GetPtr(v reflect.Value) unsafe.Pointer {
	return (*value)(unsafe.Pointer(&v)).ptr
}

// from is a pointer to the actual function
// to is a pointer to a go funcvalue
func ReplaceFunction(from, to uintptr) (original []byte) {
	jumpData := jmpToFunctionValue(to)
	f := rawMemoryAccess(from, len(jumpData))
	original = make([]byte, len(f))
	copy(original, f)

	copyToLocation(from, jumpData)
	return
}

func RecoverFunction(from uintptr, original []byte) {
	copyToLocation(from, original)
}

// Assembles a jump to a function value
func jmpToFunctionValue(to uintptr) []byte {
	return []byte{
		0x48, 0xBA,
		byte(to),
		byte(to >> 8),
		byte(to >> 16),
		byte(to >> 24),
		byte(to >> 32),
		byte(to >> 40),
		byte(to >> 48),
		byte(to >> 56), // movabs rdx,to
		0xFF, 0x22,     // jmp QWORD PTR [rdx]
	}
}

func rawMemoryAccess(p uintptr, length int) []byte {
	return *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: p,
		Len:  length,
		Cap:  length,
	}))
}

func mprotectCrossPage(addr uintptr, length int, prot int) {
	pageSize := syscall.Getpagesize()
	for p := pageStart(addr); p < addr+uintptr(length); p += uintptr(pageSize) {
		page := rawMemoryAccess(p, pageSize)
		err := syscall.Mprotect(page, prot)
		if err != nil {
			panic(err)
		}
	}
}

// this function is super unsafe
// aww yeah
// It copies a slice to a raw memory location, disabling all memory protection before doing so.
func copyToLocation(location uintptr, data []byte) {
	f := rawMemoryAccess(location, len(data))

	mprotectCrossPage(location, len(data), syscall.PROT_READ|syscall.PROT_WRITE|syscall.PROT_EXEC)
	copy(f, data[:])
	mprotectCrossPage(location, len(data), syscall.PROT_READ|syscall.PROT_EXEC)
}

func pageStart(ptr uintptr) uintptr {
	return ptr & ^(uintptr(syscall.Getpagesize() - 1))
}

// 函数前面是否一致
func AreFunctionTypesEqual(fn1, fn2 interface{}) bool {
	type1 := reflect.TypeOf(fn1)
	type2 := reflect.TypeOf(fn2)

	// 检查是否都是函数
	if type1.Kind() != reflect.Func || type2.Kind() != reflect.Func {
		return false
	}

	// 检查输入参数
	if type1.NumIn() != type2.NumIn() {
		return false
	}
	for i := 0; i < type1.NumIn(); i++ {
		if type1.In(i) != type2.In(i) {
			return false
		}
	}

	// 检查返回值
	if type1.NumOut() != type2.NumOut() {
		return false
	}
	for i := 0; i < type1.NumOut(); i++ {
		if type1.Out(i) != type2.Out(i) {
			return false
		}
	}

	return true
}
