package parser

const (
	ACC_PUBLIC       = 0x0001 // class, field, method
	ACC_PRIVATE      = 0x0002 // class, field, method
	ACC_PROTECTED    = 0x0004 // class, field, method
	ACC_STATIC       = 0x0008 // field, method
	ACC_FINAL        = 0x0010 // class, field, method, parameter
	ACC_SUPER        = 0x0020 // class
	ACC_SYNCHRONIZED = 0x0020 // method
	ACC_OPEN         = 0x0020 // module
	ACC_TRANSITIVE   = 0x0020 // module requires
	ACC_VOLATILE     = 0x0040 // field
	ACC_BRIDGE       = 0x0040 // method
	ACC_STATIC_PHASE = 0x0040 // module requires
	ACC_VARARGS      = 0x0080 // method
	ACC_TRANSIENT    = 0x0080 // field
	ACC_NATIVE       = 0x0100 // method
	ACC_INTERFACE    = 0x0200 // class
	ACC_ABSTRACT     = 0x0400 // class, method
	ACC_STRICT       = 0x0800 // method
	ACC_SYNTHETIC    = 0x1000 // class, field, method, parameter, module *
	ACC_ANNOTATION   = 0x2000 // class
	ACC_ENUM         = 0x4000 // class(?) field inner
	ACC_MANDATED     = 0x8000 // field, method, parameter, module, module *
	ACC_MODULE       = 0x8000 // class
)
