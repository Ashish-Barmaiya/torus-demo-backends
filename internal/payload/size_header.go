package payload

import "fmt"

const ResponseSizeHeader = "X-Response-Size"

func ParseSize(value string) (Size, error) {
	switch value {
	case "":
		return SizeEmpty, nil
	case "1kb":
		return Size1KB, nil
	case "16kb":
		return Size16KB, nil
	case "64kb":
		return Size64KB, nil
	case "256kb":
		return Size256KB, nil
	case "1mb":
		return Size1MB, nil
	case "4mb":
		return Size4MB, nil
	default:
		return SizeEmpty, fmt.Errorf("unsupported payload size %q", value)
	}
}
