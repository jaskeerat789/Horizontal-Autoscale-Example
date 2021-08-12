var amqp = require('amqplib/callback_api')

amqp.connect(`amqp://${process.env.RABBIT_USERNAME}:${process.env.RABBIT_PASSWORD}@${process.env.RABBIT_HOST}:${process.env.RABBIT_PORT}`, (err, connection) => {
    if (err) {
        throw err
    }

    connection.createChannel((err, channel) => {
        if (err) {
            throw err
        }
        
        var queue = "Orders"
        channel.assertQueue(queue, {
            durable: false
        })

        console.log("Waiting for messages...")
        
        channel.prefetch(1)
        channel.consume(queue, async msg => {
            await new Promise(resolve => setTimeout(resolve, 5000));
            console.log("Received %s", msg.content.toString());
            channel.ack(msg)
        })
    })
})