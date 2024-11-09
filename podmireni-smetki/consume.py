from clients.rabbitmq_client import RabbitMQClient
from handlers import message_handler


def execute():
    rabbitmq_client = RabbitMQClient()
    rabbitmq_client.consume_messages(callback=message_handler.on_message)
    rabbitmq_client.close()


if __name__ == "__main__":
    execute()
