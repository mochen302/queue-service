package queue

import (
	"fmt"
	"github.com/emirpasic/gods/lists/singlylinkedlist"
	"runtime"
	"sync"
	"time"
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
	handleChan     chan UserQueueStateInfo
	waitChan       chan UserQueueStateInfo
	waitList       *singlylinkedlist.List
	maxWaitCount   int
	maxHandleCount int
	userInfoMap    *sync.Map
	lock           *sync.RWMutex
}

func (q *QueueService) handleWaitChan() {
	for {
		select {
		case userStateInfo := <-q.waitChan:
			{
				join2TheWaitList(q, userStateInfo)
			}
		default:
			Error("q.waitChan is not all userStateInfo")
		}
	}
}

func join2TheWaitList(q *QueueService, info UserQueueStateInfo) {
	q.lock.Lock()
	defer q.lock.Unlock()

	info.stateInfo.state = ING
	info.stateInfo.extInfo = fmt.Sprint(q.waitList.Size())
	q.waitList.Append(info)
}

func (q *QueueService) handleWaitList() {
	for {
		func() {
			q.lock.Lock()
			defer q.lock.Unlock()

			count := q.maxHandleCount
			for ; count > 0; count-- {
				q.handleWaitList0()
			}
		}()

		time.Sleep(time.Duration(10) * time.Millisecond)
	}
}

func (q *QueueService) handleWaitList0() {
	if q.waitList.Size() == 0 {
		return
	}

	userStateInfo1, _ := q.waitList.Get(0)
	userStateInfo := userStateInfo1.(UserQueueStateInfo)
	userStateInfo.stateInfo.extInfo = "1"
	q.waitList.Remove(0)

	go func() {
		q.handleChan <- userStateInfo
	}()
}

func (q *QueueService) handleHandleChan() {
	for {
		select {
		case userStateInfo := <-q.handleChan:
			{
				handleToken(q, userStateInfo)
			}
		default:
			Error("q.handleChan is not all userStateInfo")
		}
	}
}

func handleToken(q *QueueService, info UserQueueStateInfo) {
	info.stateInfo.state = COMPLETE
	info.stateInfo.extInfo = "token"
}

func New(maxHandleCount int, maxWaitCount int) *QueueService {
	queueService := new(QueueService)
	queueService.maxHandleCount = maxHandleCount
	queueService.handleChan = make(chan UserQueueStateInfo, maxHandleCount)
	queueService.waitChan = make(chan UserQueueStateInfo, runtime.NumCPU()*2)
	queueService.maxWaitCount = maxWaitCount
	queueService.waitList = singlylinkedlist.New()
	queueService.lock = new(sync.RWMutex)
	queueService.userInfoMap = new(sync.Map)

	go queueService.handleWaitChan()
	go queueService.handleWaitList()
	go queueService.handleHandleChan()

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
