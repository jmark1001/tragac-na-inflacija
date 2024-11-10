from clients.rabbitmq_client import RabbitMQClient


def execute():
    rabbitmq_client = RabbitMQClient()
    rabbitmq_client.consume_messages()
    rabbitmq_client.close()


if __name__ == "__main__":
    execute()
