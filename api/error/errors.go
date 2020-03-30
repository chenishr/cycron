package error

import "strconv"

type ServerError string

func (e ServerError) Error() string {
	return "httpserver error: " + strconv.Quote(string(e))
}
