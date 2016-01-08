package lib

import (
	"log"
	"net/http"
)

type Httperror struct {
	errinfo string
	Err     error
}

var httperrorcode = 406

func Handleerr(extrainfo string, respwriter *http.ResponseWriter, rawerr error) {
	if extrainfo != "" {
		log.Println(extrainfo)
	}

	var errinfo string
	if rawerr != nil {
		errinfo = `{"errorMessage":"` + extrainfo + rawerr.Error() + `"}`

	} else {

		errinfo = `{"errorMessage":"` + extrainfo + `"}`
	}

	http.Error(*respwriter, errinfo, httperrorcode)
	//fmt.Fprintln(w, error)
	return
}
