package producer

import (
	"fmt"
	"log"
	"testing"

	"github.com/streadway/amqp"
)

func TestMQProducer(t *testing.T) {

	// 1. 尝试连接RabbitMQ，建立连接
	// 该连接抽象了套接字连接，并为我们处理协议版本协商和认证等
	conn, err := amqp.Dial("amqp://guest:guest@127.0.0.1:5672/")
	defer conn.Close()
	if err != nil {
		t.Errorf("%s has error[%+v]", t.Name(), err)
	}

	var ch *amqp.Channel
	ch, err = conn.Channel()
	defer ch.Close()

	// 3. 声明消息要发送到的队列
	//参数：
	//1.queue:队列名称
	//2.durable：是否持久化，当mq重启之后，还在
	//3.exclusive：参数有两个意思 a)是否独占即只能有一个消费者监听这个队列 b)当connection关闭时，是否删除队列
	//4.autoDelete：是否自动删除。当没有Consumer时，自动删除掉
	//5.argument：参数。配置如何删除
	//如果没有一个名字叫hello的队列，则会创建该队列，如果有则不会创建
	var q amqp.Queue
	if q, err = ch.QueueDeclare(
		"firstmq", // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	); err != nil {
		t.Errorf("%s has error[%+v]", t.Name(), err)
	}

	body := "Hello, RabbitMQ!"
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	fmt.Println(fmt.Sprintf("err:\t%+v", err))
	fmt.Printf("Sent %s\n", body)
}

func TestMQConsumer(t *testing.T) {
	// 1. 尝试连接RabbitMQ，建立连接
	conn, err := amqp.Dial("amqp://guest:guest@127.0.0.1:5672/")
	if err != nil {
		fmt.Println(err)
	}
	//defer conn.Close()

	// 2. 接下来，我们创建一个通道，大多数API都是用过该通道操作的。
	ch, err := conn.Channel()
	if err != nil {
		fmt.Println(err)
	}
	//defer ch.Close()
	// 3. 声明消息要发送到的队列
	//如果没有一个名字叫hello的队列，则会创建该队列，如果有则不会创建
	q, err := ch.QueueDeclare(
		"firstmq", // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		fmt.Println(err)
	}

	// 4.接收消息
	msgs, err := ch.Consume( // 注册一个消费者（接收消息）
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)

	for d := range msgs {
		log.Printf("Received a message: %s", d.Body)
	}
}
