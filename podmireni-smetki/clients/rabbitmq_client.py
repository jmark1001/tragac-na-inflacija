from pathlib import Path

from dotenv import load_dotenv

import json
from typing import Dict
import os
import pika

from clients import logger
from handlers import file_handler

load_dotenv()


class RabbitMQClient:
    def __init__(self):
        host = os.getenv("MQ_HOST", "rabbitmq")
        username = os.getenv("MQ_USER", "guest")
        password = os.getenv("MQ_PASSWORD", "guest")
        self.pending_queue = os.getenv("MQ_PENDING_QUEUE", "pending_files")
        self.processed_queue = os.getenv("MQ_PROCESSED_QUEUE", "processed_files")

        credentials = pika.PlainCredentials(username, password)
        self.connection = pika.BlockingConnection(
            pika.ConnectionParameters(host=host, credentials=credentials)
        )
        self.channel = self.connection.channel()

    def publish_message(self, message: Dict):
        if self.channel is None:
            raise Exception

        self.channel.queue_declare(queue=self.processed_queue, durable=True)

        self.channel.basic_publish(
            exchange='',
            routing_key=self.processed_queue,
            body=json.dumps(message).encode('utf-8')
        )
        logger.info(f"Sent message: {message}")

    def consume_messages(self):
        if self.channel is None:
            raise Exception

        self.channel.queue_declare(queue=self.pending_queue, durable=True)

        self.channel.basic_consume(
            queue=self.pending_queue, on_message_callback=self.on_message, auto_ack=True
        )
        logger.info(f"Waiting for messages in {self.pending_queue}. To exit press CTRL+C")
        self.channel.start_consuming()

    def close(self):
        """Close the connection and channel."""
        if self.connection:
            self.connection.close()
            logger.info("Connection closed.")

    def on_message(self, ch, method, properties, body):
        input_message = expense = None
        try:
            input_message = json.loads(body)
            logger.debug(input_message)
        except json.JSONDecodeError:
            logger.error("Received non-JSON message:", body)
        if input_message:
            try:
                expense = file_handler.process_file(Path(input_message.get("shared_path")))
            except Exception as e:
                logger.error("Error while parsing file: ", e)
            output_message = {
                "receipt_id": input_message.get("receipt_id"),
                "shared_path": input_message.get("shared_path"),
                "status": "failure"
            }
            if expense:
                output_message["status"] = "success"
                output_message["data"] = expense
            logger.debug(output_message)
            self.publish_message(output_message)
