package service

type User struct {
	id       int64
	nickName int64
}

type QueueStateInfo struct {
	state   QueueState
	extInfo string
}

type QueueState int8

const (
	/*处理成功*/
	COMPLETE QueueState = 0
	/*正在排队*/
	ING QueueState = 1
	/*等待加入队列*/
	WAIT QueueState = 2
)

type IQueueService interface {
	tryJoin(user User) bool
	queryState(user User) QueueStateInfo
}

type QueueService struct {
	queue *chan User
	IQueueService
}

func (q *QueueService) tryJoin(user User) bool {
	var queue = *q.queue
	queue <- user

	return true
}

func (q *QueueService) queryState(user User) QueueStateInfo {

	return QueueStateInfo{COMPLETE, "token"}
}
