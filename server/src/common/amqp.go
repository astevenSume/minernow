package common

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/mailru/easyjson"
	"github.com/streadway/amqp"
	"reflect"
	"sync"
	"time"
)

type rabbitMQConfigItem struct {
	Exchange  string `json:"exchange"`           //exchange name
	Func      string `json:"handle,omitempty"`   //alternate exchange callback function name
	IsDelay   bool   `json:"is_delay,omitempty"` //延迟队列
	AutoAck   bool   `json:"auto_ack"`           //if auto ack
	NeedAe    bool   `json:"need_ae"`            //if need alternate exchange
	NoDurable bool   `json:"no_durable"`         //if no durable
}

type rabbitMQConfig struct {
	address, appName   string
	regionId, serverId int64
	producers          map[string]rabbitMQConfigItem
	consumers          map[string]rabbitMQConfigItem
	broadcastProducers map[string]rabbitMQConfigItem
	broadcastConsumers map[string]rabbitMQConfigItem
	isRpcServer        bool   //is open rpc server
	rpcFuncName        string //rpc server func name
}

func NewRabbitMQConfig() *rabbitMQConfig {
	return &rabbitMQConfig{
		producers:          make(map[string]rabbitMQConfigItem),
		consumers:          make(map[string]rabbitMQConfigItem),
		broadcastProducers: make(map[string]rabbitMQConfigItem),
		broadcastConsumers: make(map[string]rabbitMQConfigItem),
	}
}

//
type rabbitMQManager struct {
	address       string
	producers     map[string]*rabbitMQProducer
	producerLock  sync.Mutex
	consumers     map[string]*rabbitMQConsumer
	consumerLock  sync.Mutex
	funcContainer interface{}

	rpcClients    map[string]*rabbitMQRpcClient //
	rpcClientLock sync.Mutex

	rpcServer *rabbitMQRpcServer
}

func NewRabbitMQManager() *rabbitMQManager {
	return &rabbitMQManager{
		producers:  make(map[string]*rabbitMQProducer),
		consumers:  make(map[string]*rabbitMQConsumer),
		rpcClients: make(map[string]*rabbitMQRpcClient),
	}
}

// global manager instance
var rabbitMQMgr = NewRabbitMQManager()

func (m *rabbitMQManager) Init(config *rabbitMQConfig, funcContainer interface{}) (err error) {
	m.funcContainer = funcContainer
	m.address = config.address

	// init producers
	for businessName, conf := range config.producers {
		err = m.addProducer(businessName, conf, RabbitMQExchangeTypeTopic)
		if err != nil {
			return
		}
	}

	// init consumers
	for businessName, conf := range config.consumers {
		err = m.addConsumer(businessName, conf, RabbitMQExchangeTypeTopic)
		if err != nil {
			return
		}
	}

	// init fanout producers
	for businessName, conf := range config.broadcastProducers {
		err = m.addProducer(businessName, conf, RabbitMQExchangeTypeFanout)
		if err != nil {
			return
		}
	}

	// init fanout consumers
	for businessName, conf := range config.broadcastConsumers {
		err = m.addConsumer(businessName, conf, RabbitMQExchangeTypeFanout)
		if err != nil {
			return
		}
	}

	// init rpc server
	if config.isRpcServer {
		// check func name
		f := m.getRpcFunc(config.rpcFuncName)
		if f == nil {
			err = fmt.Errorf("rpc server function %s no found", config.rpcFuncName)
			return
		}

		m.rpcServer = NewRabbitMQRpcServer()
		err = m.rpcServer.Init(config.address, config.appName, config.regionId, config.serverId, f)
		if err != nil {
			return
		}
	}

	return
}

func (m *rabbitMQManager) addProducer(businessName string, config rabbitMQConfigItem, exchangeType int) (err error) {
	m.producerLock.Lock()
	defer m.producerLock.Unlock()

	return m.addProducerUnsafe(businessName, config, exchangeType)
}

func (m *rabbitMQManager) addProducerIfNoExist(businessName string, config rabbitMQConfigItem, exchangeType int) (err error) {
	m.producerLock.Lock()
	defer m.producerLock.Unlock()

	if _, ok := m.producers[businessName]; ok {
		return
	}

	return m.addProducerUnsafe(businessName, config, exchangeType)
}

func (m *rabbitMQManager) addProducerUnsafe(businessName string, config rabbitMQConfigItem, exchangeType int) (err error) {
	producer := NewRabbitMQProducer(exchangeType)
	err = producer.Init(m.address, config.Exchange, config.AutoAck, config.NeedAe, !config.NoDurable, config.IsDelay, m.getFunc(config.Func))
	if err != nil {
		LogFuncError("rabbitMQManager.addProducer failed :%v, %v", businessName, err)
		return
	}

	m.producers[businessName] = producer
	return
}

func (m *rabbitMQManager) addConsumer(businessName string, config rabbitMQConfigItem, exchangeType int) (err error) {
	consumer := NewRabbitMQConsumer(exchangeType)
	err = consumer.Init(m.address, config.Exchange, config.AutoAck, config.NeedAe, !config.NoDurable, m.getFunc(config.Func))
	m.consumerLock.Lock()
	m.consumers[businessName] = consumer
	m.consumerLock.Unlock()
	return
}

func (m *rabbitMQManager) addRpcClientIfNoExist(appName string, regionId, serverId int64) (err error) {
	server := rpcServer(appName, regionId, serverId)

	m.rpcClientLock.Lock()
	defer m.rpcClientLock.Unlock()

	if _, ok := m.rpcClients[server]; ok {
		return
	}

	return m.addRpcClientIfNoExistUnsafe(m.address, appName, server, regionId, serverId)
}

func (m *rabbitMQManager) addRpcClientIfNoExistUnsafe(address, appName, server string, regionId, serverId int64) (err error) {
	c := NewRabbitMQRpcClient()
	err = c.Init(address, appName, regionId, serverId)
	if err != nil {
		LogFuncError("%v", err)
		return
	}

	m.rpcClients[server] = c

	return
}

func (m *rabbitMQManager) rpcClient(appName string, regionId, serverId int64) *rabbitMQRpcClient {
	m.rpcClientLock.Lock()
	defer m.rpcClientLock.Unlock()

	if c, ok := m.rpcClients[rpcServer(appName, regionId, serverId)]; ok {
		return c
	} else {
		return nil
	}
}

func (m *rabbitMQManager) getFunc(fName string) (f MessageFunc) {
	if len(fName) <= 0 {
		fName = "Default"
	}
	fc := reflect.ValueOf(m.funcContainer)
	fvZero := reflect.Value{}
	fv := fc.MethodByName(fName)
	if fv == fvZero { //while the function no found, just panic to gain attention of the administrator.
		panic(fmt.Sprintf("function cron.FunctionContainer.%s no found", fName))
	}

	f = func(delivery amqp.Delivery) (err error) {
		rets := fv.Call([]reflect.Value{reflect.ValueOf(delivery)})
		if len(rets) != 1 {
			panic(fmt.Sprintf("need return %d values but %d", 1, len(rets)))
		}

		if nil == rets[0].Interface() {
			return nil
		} else {
			return rets[0].Interface().(error)
		}
	}

	return
}

func (m *rabbitMQManager) getRpcFunc(fName string) (f RpcMessageFunc) {
	if len(fName) <= 0 {
		return nil
	}
	fc := reflect.ValueOf(m.funcContainer)
	fvZero := reflect.Value{}
	fv := fc.MethodByName(fName)
	if fv == fvZero { //while the function no found, just panic to gain attention of the administrator.
		panic(fmt.Sprintf("function cron.FunctionContainer.%s no found", fName))
	}

	f = func(delivery amqp.Delivery) (body []byte, err error) {
		rets := fv.Call([]reflect.Value{reflect.ValueOf(delivery)})
		if len(rets) != 2 {
			panic(fmt.Sprintf("need return %d values but %d", 2, len(rets)))
		}

		if nil == rets[1].Interface() {
			return rets[0].Interface().([]byte), nil
		} else {
			return rets[0].Interface().([]byte), rets[1].Interface().(error)
		}
	}

	return
}

var ErrProducerNoFound = errors.New("producer no found")

//
func (m *rabbitMQManager) publish(businessName, key string, msg []byte) (err error) {
	m.producerLock.Lock()
	p, ok := m.producers[businessName]
	m.producerLock.Unlock()

	if !ok {
		err = ErrProducerNoFound
		LogFuncWarning("producer of %s no found.", businessName)
		return
	}

	err = p.Publish(key, msg)

	return
}

func (m *rabbitMQManager) publishDelay(businessName, key string, msg []byte, t string) (err error) {
	m.producerLock.Lock()
	p, ok := m.producers[businessName]
	m.producerLock.Unlock()

	if !ok {
		err = ErrProducerNoFound
		LogFuncWarning("producer of %s no found.", businessName)
		return
	}

	err = p.PublishDelay(key, msg, t)

	return
}

// @Description init rabbit mq module.
func RabbitMQInit(aeFuncContainer interface{}) (err error) {
	config := NewRabbitMQConfig()

	{
		if regionId, err := beego.AppConfig.Int64("RegionId"); err != nil {
			panic("no specific RegionId")
		} else {
			config.regionId = regionId
		}
	}

	{
		if serverId, err := beego.AppConfig.Int64("ServerId"); err != nil {
			panic("no specific ServerId")
		} else {
			config.serverId = serverId
		}
	}

	{
		config.appName = beego.AppConfig.String("appname")
	}

	config.address = beego.AppConfig.String("rabbitmq::address")

	// get configuration of producers\consumers\broadcastProducers\broadcastConsumers
	keys := []string{"producer", "consumer", "broadcastProducer", "broadcastConsumer"}
	values := []*(map[string]rabbitMQConfigItem){&config.producers, &config.consumers, &config.broadcastProducers, &config.broadcastConsumers}

	for i := 0; i < len(keys); i++ {
		err = json.Unmarshal([]byte(beego.AppConfig.String("rabbitmq::"+keys[i])), values[i])
		if err != nil {
			panic(err.Error())
		}
	}

	// get configuration of rpc
	config.isRpcServer, _ = beego.AppConfig.Bool("rabbitmq::isRpcServer")
	config.rpcFuncName = beego.AppConfig.String("rabbitmq::rpcServerFuncName")

	return rabbitMQMgr.Init(config, aeFuncContainer)
}

// RabbitMQPublish publish specific business message.
func RabbitMQPublish(businessName, key string, msg []byte) (err error) {
	return rabbitMQMgr.publish(businessName, key, msg)
}

func RabbitMQPublishDelay(businessName, key string, msg []byte, t string) (err error) {
	return rabbitMQMgr.publishDelay(businessName, key, msg, t)
}

// RabbitMQAddProducer add producer
func RabbitMQAddProducer(businessName, exchange string) (err error) {
	return rabbitMQMgr.addProducer(businessName, rabbitMQConfigItem{
		Exchange: exchange,
		NeedAe:   true,
	}, RabbitMQExchangeTypeTopic)
}

// RabbitMQAddProducerIfNoExist add producer if no exist
func RabbitMQAddProducerIfNoExist(businessName, exchange string) (err error) {
	return rabbitMQMgr.addProducerIfNoExist(businessName, rabbitMQConfigItem{
		Exchange: exchange,
		NeedAe:   true,
	}, RabbitMQExchangeTypeTopic)
}

// RabbitMQAddRpcClient add rpc client
func RabbitMQAddRpcClient(appName string, regionId, serverId int64) (err error) {
	return rabbitMQMgr.addRpcClientIfNoExist(appName, regionId, serverId)
}

// RabbitMQAddRpcClientIfNoExist add rpc client if no exist
func RabbitMQAddRpcClientIfNoExist(appName string, regionId, serverId int64) (err error) {
	return rabbitMQMgr.addRpcClientIfNoExist(appName, regionId, serverId)
}

// RabbitMQRpcSend send rpc request to specific server
// return response, or timeout error if server no response in timeOut period
func RabbitMQRpcSend(appName string, regionId, serverId int64, msg *RabbitMQRpcMsg, timeOut int) (resp RabbitMQRpcMsg, err error) {
	server := rpcServer(appName, regionId, serverId)

	//get rpc client entity
	c := rabbitMQMgr.rpcClient(appName, regionId, serverId)
	if c == nil {
		err = errors.New(fmt.Sprintf("rpc client of %s no found", server))
		return
	}

	var buf []byte
	//encode request message
	buf, err = easyjson.Marshal(msg)
	if err != nil {
		LogFuncError("%v", err)
		return
	}

	//send
	buf, err = c.Publish(buf, time.Duration(timeOut))
	if err != nil {
		return
	}

	//decode response bytes
	err = easyjson.Unmarshal(buf, &resp)
	if err != nil {
		LogFuncError("%v", err)
		return
	}

	return
}

type rabbitMQBase struct {
	address string
	conn    *amqp.Connection
	ch      *amqp.Channel
	chErr   chan *amqp.Error //channel to listening close notify
	fReset  rabbitMQResetFunc
}

type rabbitMQResetFunc func()

func NewRabbitMQBase() *rabbitMQBase {
	return &rabbitMQBase{}
}

func processErr(r *rabbitMQBase) {
	LogFuncDebug("rabbitMQClient Listener Begin...")
	var rabbitErr *amqp.Error
	for {
		rabbitErr = <-r.chErr
		if rabbitErr != nil {
			LogFuncWarning("rabbitMQClient connect to %s", r.address)
			r.conn = reconnect(r.address) //keep trying until connected
			r.chErr = make(chan *amqp.Error)
			r.conn.NotifyClose(r.chErr)
			r.fReset()
		}
	}
}
