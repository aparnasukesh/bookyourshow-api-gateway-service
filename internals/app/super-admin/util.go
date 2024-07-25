package superadmin

import "strings"

func ExtractErrorMessage(err error) string {
	if err == nil {
		return ""
	}
	errMsg := err.Error()
	if index := strings.Index(errMsg, "desc = "); index != -1 {
		return errMsg[index+len("desc = "):]
	}
	return errMsg
}
