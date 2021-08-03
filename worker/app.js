var amqp = require('amqplib/callback_api')
const cliProgress = require('cli-progress')

const bar1 = new cliProgress.SingleBar({},cliProgress.Presets.shades_classic)
const length = 10_000_000

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
            bar1.start(length,0)
            for (let index = 1; index <= length; index++) { 
                bar1.update(index)
            }
            bar1.stop()
            console.log("Received %s", msg.content.toString());
            channel.ack(msg)
        })
    })
})