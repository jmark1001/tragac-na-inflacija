from dotenv import load_dotenv
import os
import pika

load_dotenv()


class RabbitMQClient:
    def __init__(self):
        host = os.getenv("MQ_HOST", "localhost")
        username = os.getenv("MQ_USER", "guest")
        password = os.getenv("MQ_PASSWORD", "guest")
        self.queue_name = os.getenv("MQ_QUEUE_NAME", "file_queue")

        credentials = pika.PlainCredentials(username, password)
        self.connection = pika.BlockingConnection(
            pika.ConnectionParameters(host=host, credentials=credentials)
        )
        self.channel = self.connection.channel()

    # def publish_message(self, message: str):
    #     if self.channel is None:
    #         raise Exception
    #
    #     # Ensure the queue exists
    #     self.channel.queue_declare(queue=self.queue_name)
    #
    #     # Convert message to JSON and publish to the queue
    #     self.channel.basic_publish(
    #         exchange='',
    #         routing_key=queue_name,
    #         body=message
    #     )
    #     print(f"Sent message: {message}")

    def consume_messages(self, callback):
        if self.channel is None:
            raise Exception

        self.channel.queue_declare(queue=self.queue_name, durable=True)

        self.channel.basic_consume(
            queue=self.queue_name, on_message_callback=callback, auto_ack=True
        )
        print(f"Waiting for messages in {self.queue_name}. To exit press CTRL+C")
        self.channel.start_consuming()

    def close(self):
        """Close the connection and channel."""
        if self.connection:
            self.connection.close()
            print("Connection closed.")
