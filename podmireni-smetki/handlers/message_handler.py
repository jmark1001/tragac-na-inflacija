import json


def on_message(ch, method, properties, body):
    try:
        message = json.loads(body)
        print(f"Received message: {json.dumps(message, indent=4)}")
    except json.JSONDecodeError:
        print("Received non-JSON message:", body)
