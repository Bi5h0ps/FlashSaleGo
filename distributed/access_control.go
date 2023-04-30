package distributed

import (
	"FlashSaleGo/grpc/order"
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
	sourceArray map[int]*CloudData
	sync.RWMutex
}

func (a *AccessControl) GetNewRecord(uid int) *CloudData {
	a.RLock()
	defer a.RUnlock()
	return a.sourceArray[uid]
}

func (a *AccessControl) SetNewRecord(uid int) *CloudData {
	a.Lock()
	defer a.Unlock()
	newData := &CloudData{LastOrderTime: time.Now().Unix()}
	a.sourceArray[uid] = newData
	return newData
}

func (a *AccessControl) GetDataFromMap(uid string) (result *CloudData, err error) {
	uidInt, err := strconv.Atoi(uid)
	data := a.GetNewRecord(uidInt)
	if data != nil {
		return data, nil
	}
	return a.SetNewRecord(uidInt), nil
}

func (a *AccessControl) GetDataFromOtherMap(uid string) (result *CloudData, err error) {
	uidInt, err := strconv.Atoi(uid)
	var conn *grpc.ClientConn
	conn, err = grpc.Dial(":9093", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()
	c := order.NewOrderServiceClient(conn)
	response, err := c.GetUserCloudData(context.Background(), &order.UserInfo{UserID: int64(uidInt)})
	return &CloudData{LastOrderTime: response.TimeStamp}, nil
}

func (a *AccessControl) GetDistributedRight(uid string, hashConsistant *Consistent, localhost string) bool {
	//find server base on user id
	hostRequest, err := hashConsistant.Get(uid)
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
		data, err = a.GetDataFromOtherMap(uid)
	}
	if err != nil {
		return false
	}
	return checkUserOrderFrequency(data.LastOrderTime, time.Minute)
}

func checkUserOrderFrequency(userTimeStamp int64, duration time.Duration) bool {
	return time.Now().Add(-1*duration).Unix() < userTimeStamp
}
func NewAccessControlUnit() *AccessControl {
	return &AccessControl{sourceArray: make(map[int]*CloudData)}
}
