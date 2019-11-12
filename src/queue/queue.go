package queue

import (
	"fmt"
	"github.com/emirpasic/gods/lists/singlylinkedlist"
	"sync"
)

type User struct {
	id       int64
	nickName int64
	fmt.Stringer
}

func (u User) String() string {
	return fmt.Sprintf("id:%v nickName:%v", u.id, u.nickName)
}

type QueueStateInfo struct {
	state   QueueState
	extInfo string
	fmt.Stringer
}

func (stateInfo QueueStateInfo) String() string {
	return fmt.Sprintf("state:%v extInfo:%v", stateInfo.state, stateInfo.extInfo)
}

type QueueState int8

const (
	/*处理成功 返回token*/
	COMPLETE QueueState = 0
	/*正在排队 返回ranking*/
	ING QueueState = 1
	/*等待加入队列*/
	WAIT QueueState = 2
)

type UserQueueStateInfo struct {
	user      User
	stateInfo QueueStateInfo
	fmt.Stringer
}

func (userStateInfo UserQueueStateInfo) String() string {
	return fmt.Sprintf("user:%v stateInfo:%v", userStateInfo.user.String(), userStateInfo.stateInfo.String())
}

type QueueService struct {
	handleChan   chan UserQueueStateInfo
	waitChan     chan UserQueueStateInfo
	waitList     *singlylinkedlist.List
	maxWaitCount int
	userInfoMap  sync.Map
	lock         *sync.RWMutex
}

func New(maxHandleCount int, maxWaitCount int) *QueueService {
	queueService := new(QueueService)
	queueService.handleChan = make(chan UserQueueStateInfo, maxHandleCount)
	queueService.waitChan = make(chan UserQueueStateInfo)
	queueService.maxWaitCount = maxWaitCount
	queueService.waitList = singlylinkedlist.New()
	queueService.lock = new(sync.RWMutex)

	return queueService
}

func (q *QueueService) TryJoin(currentUser User) bool {
	waitSize := q.waitList.Size()
	if waitSize >= q.maxWaitCount {
		Error(currentUser.String(), "try join fail cause the waitSize:", waitSize, ">", q.maxWaitCount)
		return false
	}

	id := currentUser.id
	existUserStateInfo, ok := q.userInfoMap.Load(id)
	if ok {
		/**
		 * 这个地方主要是看具体的业务
		 */
		Info(existUserStateInfo.(UserQueueStateInfo).String(), "has join before")
		return true
	}

	userStateInfo := UserQueueStateInfo{
		user: currentUser,
		stateInfo: QueueStateInfo{
			state:   WAIT,
			extInfo: "",
		},
	}

	q.userInfoMap.Store(id, userStateInfo)
	q.waitChan <- userStateInfo
	Debug(userStateInfo.String(), "join wait queue suc!")
	return true
}

func (q *QueueService) QueryState(user User) *UserQueueStateInfo {
	userStateInfo := new(UserQueueStateInfo)
	return userStateInfo
}
