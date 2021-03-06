// Code generated by "stringer -type=ClusterIssuerType -linecomment"; DO NOT EDIT.

package types

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[ClusterIssuerTypeCALocal-0]
	_ = x[ClusterIssuerTypeLetsEncryptStaging-1]
	_ = x[ClusterIssuerTypeLetsEncryptProd-2]
}

const _ClusterIssuerType_name = "ca-localletsencrypt-stagingletsencrypt-prod"

var _ClusterIssuerType_index = [...]uint8{0, 8, 27, 43}

func (i ClusterIssuerType) String() string {
	if i < 0 || i >= ClusterIssuerType(len(_ClusterIssuerType_index)-1) {
		return "ClusterIssuerType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _ClusterIssuerType_name[_ClusterIssuerType_index[i]:_ClusterIssuerType_index[i+1]]
}
