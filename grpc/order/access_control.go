package order

import (
	"FlashSaleGo/distributed"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"strconv"
	"sync"
	"time"
)

type CloudData struct {
	LastOrderTime int64
}

// AccessControl stores some user info, like timestamp
type AccessControl struct {
	sourceArray map[int64]*CloudData
	sync.RWMutex
}

func (a *AccessControl) GetNewRecord(uid int64) *CloudData {
	a.RLock()
	defer a.RUnlock()
	return a.sourceArray[uid]
}

func (a *AccessControl) SetNewRecord(uid int64) *CloudData {
	a.Lock()
	defer a.Unlock()
	newData := &CloudData{LastOrderTime: time.Now().Unix()}
	a.sourceArray[uid] = newData
	return newData
}

func (a *AccessControl) GetDataFromMap(uid int64) (result *CloudData, err error) {
	data := a.GetNewRecord(uid)
	if data != nil {
		return data, nil
	}
	return nil, nil
}

func (a *AccessControl) GetDataFromOtherMap(uid int64, hostRequest string) (result *CloudData, err error) {
	var conn *grpc.ClientConn
	conn, err = grpc.Dial(hostRequest+":9093", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()
	c := NewOrderServiceClient(conn)
	response, err := c.GetUserCloudData(context.Background(), &UserInfo{UserID: uid})
	if err != nil {
		return nil, err
	}
	if response == nil {
		return nil, nil
	}
	return &CloudData{LastOrderTime: response.TimeStamp}, nil
}

func (a *AccessControl) GetDistributedRight(uid int64, hashConsistant *distributed.Consistent, localhost string) bool {
	uidString := strconv.FormatInt(uid, 10)
	//find server base on user id
	hostRequest, err := hashConsistant.Get(uidString)
	if err != nil {
		return false
	}

	var data *CloudData
	if hostRequest == localhost {
		//get data from localhost
		data, err = a.GetDataFromMap(uid)
	} else {
		//user data not on this server, go to the target server and retrieve user info
		//act like a delegate and return result
		data, err = a.GetDataFromOtherMap(uid, hostRequest)
	}
	if err != nil {
		return false
	}
	//user note allowed to make another purchase in a one-minute window
	if data != nil {
		result := checkUserOrderFrequency(data.LastOrderTime, time.Minute)
		if result {
			a.SetNewRecord(uid)
		}
		return result
	} else {
		a.SetNewRecord(uid)
		return true
	}
}

func checkUserOrderFrequency(userTimeStamp int64, duration time.Duration) bool {
	return time.Now().Add(-1*duration).Unix() > userTimeStamp
}
func NewAccessControlUnit() *AccessControl {
	return &AccessControl{sourceArray: make(map[int64]*CloudData)}
}
