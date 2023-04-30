package distributed

import (
	"net/http"
	"strconv"
	"sync"
	"time"
)

type CloudData struct {
	lastOrderTime int64
}

// AccessControl stores some user info, like timestamp
type AccessControl struct {
	sourceArray map[int]CloudData
	sync.RWMutex
}

var accessControl = AccessControl{sourceArray: make(map[int]CloudData)}

func (a *AccessControl) GetNewRecord(uid int) interface{} {
	a.RLock()
	defer a.RUnlock()
	return a.sourceArray[uid]
}

func (a *AccessControl) SetNewRecord(uid int) {
	a.Lock()
	defer a.Unlock()
	a.sourceArray[uid] = CloudData{lastOrderTime: time.Now().Unix()}
}

func (a *AccessControl) GetDataFromMap(uid string) (isOk bool) {
	uidInt, err := strconv.Atoi(uid)
	if err != nil {
		return false
	}
	data := a.GetNewRecord(uidInt)
	if data != nil {
		return true
	}
	return false
}

func (a *AccessControl) GetDataFromOtherMap(host string, request *http.Request) (isOk bool) {
	//hostUrl := fmt.Sprintf("http://%s:%s/checkRight", host, port)
	//response, body, err := GetCurl(hostUrl, request)
	//if err != nil {
	//	return false
	//}
	//if response.StatusCode == 200 {
	//	if string(body) == "true" {
	//		return true
	//	} else {
	//		return
	//	}
	//}
	return false
}
