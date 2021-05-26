package common

import (
    "comic/config"
    "fmt"
    "github.com/streadway/amqp"
    "log"
    "sync"
)

//rabbitMQ结构体
type RabbitMQ struct {
    Conn    *amqp.Connection
    Channel *amqp.Channel
    //队列名称
    QueueName string
    //交换机名称
    Exchange string
    //bind Key 名称
    Key string
    //连接信息
    Mqurl string
    sync.Mutex
}

//创建结构体实例
func NewRabbitMQ(queueName string, exchange string, key string) *RabbitMQ {
    cfg := config.GetConfig().RabbitMQ
    mqUrl := fmt.Sprintf("amqp://%s:%s@localhost:%s/", cfg.User, cfg.Password, cfg.Host)
    return &RabbitMQ{QueueName: queueName, Exchange: exchange, Key: key, Mqurl: mqUrl}
}

//断开channel 和 connection
func (r *RabbitMQ) Destory() {
    r.Channel.Close()
    r.Conn.Close()
}

//错误处理函数
func (r *RabbitMQ) failOnErr(err error, message string) {
    if err != nil {
        log.Fatalf("%s:%s", message, err)
        panic(fmt.Sprintf("%s:%s", message, err))
    }
}

//创建RabbitMQ实例
func NewRabbitMQSimple(queueName string) *RabbitMQ {
    //创建RabbitMQ实例
    rabbitmq := NewRabbitMQ(queueName, "", "")
    var err error
    //获取connection
    rabbitmq.Conn, err = amqp.Dial(rabbitmq.Mqurl)
    rabbitmq.failOnErr(err, "failed to connect rabbitmq!")
    //获取channel
    rabbitmq.Channel, err = rabbitmq.Conn.Channel()
    rabbitmq.failOnErr(err, "failed to open a channel")

    return rabbitmq
}
