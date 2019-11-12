package queue

import "sync"

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

type QueueService struct {
	queue            chan User
	maxUserWaitCount int64
	userMap          map[int64]User
	lock             *sync.RWMutex
}

func New(maxUserCount int64, maxUserWaitCount int64) *QueueService {

	queueService := new(QueueService)
	queueService.queue = make(chan User, maxUserCount)
	queueService.userMap = make(map[int64]User)
	queueService.maxUserWaitCount = maxUserWaitCount
	queueService.lock = new(sync.RWMutex)

	return queueService
}

func (q *QueueService) TryJoin(user User) bool {
	var queue = q.queue
	queue <- user

	return true
}

func (q *QueueService) QueryState(user User) QueueStateInfo {

	return QueueStateInfo{COMPLETE, "token"}
}
