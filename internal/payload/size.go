package payload

import "fmt"

type Size int64

const (
	SizeEmpty Size = 0
	Size1KB   Size = 1 << 10
	Size16KB  Size = 16 << 10
	Size64KB  Size = 64 << 10
	Size256KB Size = 256 << 10
	Size1MB   Size = 1 << 20
	Size4MB   Size = 4 << 20
)

func (s Size) Valid() bool {
	switch s {
	case SizeEmpty, Size1KB, Size16KB, Size64KB, Size256KB, Size1MB, Size4MB:
		return true
	default:
		return false
	}
}

func (s Size) String() string {
	switch s {
	case SizeEmpty:
		return "empty"
	case Size1KB:
		return "1kb"
	case Size16KB:
		return "16kb"
	case Size64KB:
		return "64kb"
	case Size256KB:
		return "256kb"
	case Size1MB:
		return "1mb"
	case Size4MB:
		return "4mb"
	default:
		return fmt.Sprintf("%dB", s)
	}
}
