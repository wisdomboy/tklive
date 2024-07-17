package utils

import (
	"github.com/gorilla/websocket"
	"sync"
)

// ConnectionManager tcp conn与websocket conn通讯结构体
type ConnectionManager struct {
	chanSize       int
	reConnCount    map[string]int
	connections    map[string]chan interface{}
	lock           sync.RWMutex
	socketConnColl map[string]*websocket.Conn
}

// NewConnectionManager 创建结构体访问对象
func NewConnectionManager() *ConnectionManager {
	return &ConnectionManager{
		chanSize:       100,
		reConnCount:    make(map[string]int, 50),
		connections:    make(map[string]chan interface{}, 100),
		lock:           sync.RWMutex{},
		socketConnColl: make(map[string]*websocket.Conn),
	}
}

// SetReConnCount 设置websocket重连次数
func (cn *ConnectionManager) SetReConnCount(roomId string, count int) int {
	cn.lock.Lock()
	defer cn.lock.Unlock()

	cn.reConnCount[roomId] += count

	return cn.reConnCount[roomId]
}

// GetReConnCount 获取websocket重连次数
func (cn *ConnectionManager) GetReConnCount(roomId string) int {
	cn.lock.RLock()
	cn.lock.RUnlock()

	return cn.reConnCount[roomId]
}

func (cn *ConnectionManager) GetChanSize() int {
	return cn.chanSize
}

// AddSocketConnect 增加socket连接集合
func (cn *ConnectionManager) AddSocketConnect(roomId string, socketConn *websocket.Conn) *websocket.Conn {
	cn.lock.Lock()
	defer cn.lock.Unlock()

	if cn.socketConnColl[roomId] == nil {
		cn.socketConnColl[roomId] = socketConn
	}

	return cn.socketConnColl[roomId]
}

// GetSocketConnect 获取socket连接集合
func (cn *ConnectionManager) GetSocketConnect(roomId string) *websocket.Conn {
	cn.lock.RLock()
	defer cn.lock.RUnlock()

	return cn.socketConnColl[roomId]
}

func (cn *ConnectionManager) RemoveSocketConnect(roomId string) {
	cn.lock.Lock()
	defer cn.lock.Unlock()

	delete(cn.socketConnColl, roomId)
}

// AddConnection 每个房间创建通道
func (cn *ConnectionManager) AddConnection(roomId string) chan interface{} {
	cn.lock.Lock()
	defer cn.lock.Unlock()

	if cn.connections[roomId] == nil {
		cn.connections[roomId] = make(chan interface{}, 100)
	}
	return cn.connections[roomId]
}

// RemoveConnection 移除房间通道
func (cn *ConnectionManager) RemoveConnection(roomId string) {
	cn.lock.Lock()
	defer cn.lock.Unlock()
	delete(cn.connections, roomId)
}

func (cn *ConnectionManager) CloseConnection(roomId string) {
	cn.lock.Lock()
	defer cn.lock.Unlock()

	close(cn.connections[roomId])
}

// GetConnection 获取房间通道
func (cn *ConnectionManager) GetConnection(roomId string) chan interface{} {
	cn.lock.RLock()
	defer cn.lock.RUnlock()

	return cn.connections[roomId]
}
