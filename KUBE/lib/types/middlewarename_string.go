// Code generated by "stringer -type=MiddlewareName -linecomment"; DO NOT EDIT.

package types

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[MiddleWareNameTest-0]
	_ = x[MiddleWareNameBasicAuth-1]
}

const _MiddlewareName_name = "testbasicAuth"

var _MiddlewareName_index = [...]uint8{0, 4, 13}

func (i MiddlewareName) String() string {
	if i < 0 || i >= MiddlewareName(len(_MiddlewareName_index)-1) {
		return "MiddlewareName(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _MiddlewareName_name[_MiddlewareName_index[i]:_MiddlewareName_index[i+1]]
}