package queue

import (
	"fmt"
	"github.com/emirpasic/gods/lists/singlylinkedlist"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

/**
线程模型如下
                             head->... waitList(最大maxWaitCount)  ...tail
handleChan(handleChanLength) ============================================wait2JoinChan(核数*2)
userInfoMap根据玩家Id保存状态信息

1.先将新来的请求放到waitJoinChan则返回
2.waitJoinChan将请求信息放到waitList结尾
3.handleChan不断从waitList头获取请求信息处理
*/
type Queue struct {
	/*最终处理的chan*/
	handleChan chan *UserQueueStateInfo
	/*最终处理的chan缓冲区大小*/
	handleChanLength int
	/*等待加入waitList的chan*/
	wait2JoinChan chan *UserQueueStateInfo
	/*等待被处理的List*/
	waitList *singlylinkedlist.List
	/*等待被处理的List的最大大小*/
	maxWaitCount int
	/*用户信息Map*/
	userInfoMap *sync.Map
	/*读写锁*/
	lock *sync.RWMutex
	/*请求次数*/
	requestCount int32
	/*处理成功次数*/
	handleSuccessCount int32
}

type User struct {
	id       int64
	nickName string
	fmt.Stringer
}

func (u *User) String() string {
	return fmt.Sprintf("id:%v nickName:%v", u.id, u.nickName)
}

type StateInfo struct {
	state   State
	extInfo string
	fmt.Stringer
}

func (stateInfo *StateInfo) String() string {
	return fmt.Sprintf("state:%v extInfo:%v", stateInfo.state, stateInfo.extInfo)
}

type State int8

const (
	/*处理成功 返回token*/
	COMPLETE State = 0
	/*正在排队 返回ranking*/
	ING State = 1
	/*等待加入队列*/
	WAIT State = 2
)

type UserQueueStateInfo struct {
	user      *User
	stateInfo *StateInfo
	fmt.Stringer
}

func (userStateInfo *UserQueueStateInfo) String() string {
	return fmt.Sprintf("user:{%v} stateInfo:{%v}", userStateInfo.user.String(), userStateInfo.stateInfo.String())
}

type StatInfo struct {
	requestCount       int32
	handleSuccessCount int32

	waitJoinChanCount int
	waitCount         int
	handleCount       int
	fmt.Stringer
}

func (statInfo *StatInfo) String() string {
	return fmt.Sprintf("requestCount:%v handleSuccessCount:%v waitJoinChanCount:%v waitCount:%v handleCount:%v",
		statInfo.requestCount, statInfo.handleSuccessCount, statInfo.waitJoinChanCount, statInfo.waitCount, statInfo.handleCount)
}

func (q *Queue) handleWaitChan() {
	for {
		select {
		case userStateInfo := <-q.wait2JoinChan:
			join2TheWaitList(q, userStateInfo)
		}
	}
}

func join2TheWaitList(q *Queue, info *UserQueueStateInfo) {
	/*todo 此处锁的粒度太大了，时间限制以后优化*/
	q.lock.Lock()
	defer q.lock.Unlock()
	defer func() {
		err := recover()
		if err != nil {
			Error(info.String(), " join2TheWaitList error", err)
		}
	}()

	info.stateInfo.state = ING
	info.stateInfo.extInfo = fmt.Sprint(q.waitList.Size())
	q.waitList.Append(info)
	Info(info.String(), " join wait chan suc!")
	atomic.AddInt32(&q.handleSuccessCount, 1)
}

func (q *Queue) handleWaitList() {
	for {
		func() {
			/*todo 此处锁的粒度太大了，时间限制以后优化*/
			q.lock.Lock()
			defer q.lock.Unlock()
			defer func() {
				err := recover()
				if err != nil {
					Error(" handleWaitList error", err)
				}
			}()

			count := q.handleChanLength
			for ; count > 0; count-- {
				q.handleWaitList0()
			}
		}()

		time.Sleep(time.Duration(10) * time.Millisecond)
	}
}

func (q *Queue) handleWaitList0() {
	if q.waitList.Size() == 0 {
		return
	}

	userStateInfo1, _ := q.waitList.Get(0)
	userStateInfo := userStateInfo1.(*UserQueueStateInfo)
	userStateInfo.stateInfo.extInfo = "1"
	q.waitList.Remove(0)

	go func() {
		q.handleChan <- userStateInfo
	}()
}

func (q *Queue) handleHandleChan() {
	for {
		select {
		case userStateInfo := <-q.handleChan:
			handleToken(q, userStateInfo)
		}
	}
}

func handleToken(q *Queue, info *UserQueueStateInfo) {
	defer func() {
		err := recover()
		if err != nil {
			Error(info.String(), " handleToken error", err)
		}
	}()

	info.stateInfo.state = COMPLETE
	info.stateInfo.extInfo = "token"
	Info(info.String(), " handle suc!")
}

func (q *Queue) updateUserRanking(info *UserQueueStateInfo) {
	q.lock.RLock()
	defer q.lock.RUnlock()

	ranking := 1
	iterator := q.waitList.Iterator()
	for ; iterator.Next(); ranking++ {
		temp := iterator.Value().(*UserQueueStateInfo)
		if temp.user.id == info.user.id {
			break
		}
	}

	info.stateInfo.extInfo = fmt.Sprint(ranking)
}

func New(handleChanLength int, maxWaitCount int) *Queue {
	queueService := new(Queue)
	queueService.handleChanLength = handleChanLength
	queueService.handleChan = make(chan *UserQueueStateInfo, handleChanLength)
	queueService.wait2JoinChan = make(chan *UserQueueStateInfo, runtime.NumCPU()*2)
	queueService.maxWaitCount = maxWaitCount
	queueService.waitList = singlylinkedlist.New()
	queueService.lock = new(sync.RWMutex)
	queueService.userInfoMap = new(sync.Map)
	queueService.requestCount = 0
	queueService.handleSuccessCount = 0

	go queueService.handleWaitChan()
	go queueService.handleWaitList()
	go queueService.handleHandleChan()

	return queueService
}

func (q *Queue) TryJoin(p ...interface{}) interface{} {
	p0 := p[0].([]interface{})
	id := p0[0].(int64)
	nickname := p0[1].(string)

	currentUser := new(User)
	currentUser.id = id
	currentUser.nickName = nickname

	Info(currentUser.String(), " try join")

	existUserStateInfo, ok := q.userInfoMap.Load(id)
	if ok {
		/**
		 * 这个地方主要是看具体的业务
		 */
		Info(existUserStateInfo.(*UserQueueStateInfo).String(), "has join before")
		return true
	}

	atomic.AddInt32(&q.requestCount, 1)

	waitSize := q.waitList.Size()
	if waitSize >= q.maxWaitCount {
		Error(currentUser.String(), " try join fail cause the waitSize:", waitSize, ">", q.maxWaitCount)
		return false
	}

	userStateInfo := &UserQueueStateInfo{
		user: currentUser,
		stateInfo: &StateInfo{
			state:   WAIT,
			extInfo: "",
		},
	}

	q.userInfoMap.Store(id, userStateInfo)
	q.wait2JoinChan <- userStateInfo
	Info(userStateInfo.String(), " join wait chan suc!")
	return true
}

func (q *Queue) QueryState(p ...interface{}) interface{} {
	p0 := p[0].([]interface{})
	id := p0[0].(int64)
	userStateInfo1, ok := q.userInfoMap.Load(id)
	if !ok {
		Error("id:", id, " QueryState failed")
		panic(fmt.Sprint("can not find userInfo:", id))
	}

	userStateInfo := userStateInfo1.(*UserQueueStateInfo)
	if userStateInfo.stateInfo.state == ING {
		q.updateUserRanking(userStateInfo)
	}

	Info(fmt.Sprintf("QueryState id:%v result:%v", id, userStateInfo.String()))
	return userStateInfo
}

func (q *Queue) StatInfo(p ...interface{}) interface{} {
	statInfo := new(StatInfo)
	statInfo.requestCount = q.requestCount
	statInfo.handleSuccessCount = q.handleSuccessCount
	statInfo.waitJoinChanCount = len(q.wait2JoinChan)
	statInfo.waitCount = q.waitList.Size()
	statInfo.handleCount = len(q.handleChan)
	return statInfo
}

func (q *Queue) Close() {
	defer func() {
		err := recover()
		if err != nil {
			Error("close error", err)
		}
	}()

	close(q.wait2JoinChan)
	close(q.handleChan)
	q.waitList = nil
	q.userInfoMap = nil
	Info("queue.close success")
}
