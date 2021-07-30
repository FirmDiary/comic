package services

import (
	"comic/common"
	"github.com/streadway/amqp"
	"log"
)

const (
	DelImageExchange = "my-delay-exchange"
	DelImageQueue    = "test"
)

func QueueListen() {
	DelImgHandler()
}

//删除用户上传的文件
func DelImg(name string, second int64) {
	rabbitmq := common.NewRabbitMQSimple(DelImageExchange)

	// 声明一个主要使用的 exchange
	err := rabbitmq.Channel.ExchangeDeclare(
		rabbitmq.QueueName,  // name
		"x-delayed-message", // type
		true,                // durable
		false,               // auto-deleted
		false,               // internal
		false,               // no-wait
		map[string]interface{}{
			"x-delayed-type": "direct",
		}, // arguments
	)

	common.FailOnError(err, "Failed to declare an exchange")

	// 将消息发送到延时队列上
	err = rabbitmq.Channel.Publish(
		rabbitmq.QueueName, // exchange 这里为空则不选择 exchange
		"",                 // routing key
		false,              // mandatory
		false,              // immediate
		amqp.Publishing{
			Headers: amqp.Table{
				"x-delay": second * 1000,
			},
			ContentType: "text/plain",
			Body:        []byte(name),
			//Expiration: second*2000, // 设置过期时间
		})

	log.Printf(" [ %v 秒后执行文件删除,文件名:] Sent %s", second, name)
}

func DelImgHandler() {
	rabbitmq := common.NewRabbitMQSimple(DelImageQueue)

	_, err := rabbitmq.Channel.QueueDeclare(DelImageQueue, true, false, false, false, nil)
	common.FailOnError(err, "Failed to declare a queue")
	err = rabbitmq.Channel.QueueBind(DelImageQueue, "", DelImageExchange, false, nil)
	common.FailOnError(err, "Failed to bind a queue")

	msgs, err := rabbitmq.Channel.Consume(
		rabbitmq.QueueName, // queue name, 这里指的是 test_logs
		"",                 // consumer
		true,               // auto-ack
		false,              // exclusive
		false,              // no-local
		false,              // no-wait
		nil,                // args
	)
	common.FailOnError(err, "Failed to register a consumer")
	forever := make(chan bool)
	go func() {
		log.Println("图片删除任务监听中...")
		for d := range msgs {
			err = DelUploadImg(string(d.Body))
			common.FailOnError(err, "执行文件删除失败")
			log.Printf(" [执行文件删除,文件名：] %s", d.Body)
		}
	}()
	<-forever
}
