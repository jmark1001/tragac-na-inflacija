# Use Python 3.12 official image
FROM python:3.12-slim as python-consumer-builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Python requirements file to install dependencies
COPY ./podmireni-smetki/requirements.txt ./

# Install Python dependencies
RUN pip install --no-cache-dir -r requirements.txt

# Copy the Python source code into the container
COPY ./podmireni-smetki /app

# Copy the wait-for-it script into the container
COPY ./wait-for-it.sh /app/wait-for-it.sh
RUN chmod +x /app/wait-for-it.sh

# Copy the .env file into the container's /app directory
COPY ./podmireni-smetki/.env /app/.env

# Expose the necessary port
EXPOSE 5000

# Command to run the Python consumer script
CMD ["sh", "-c", "/app/wait-for-it.sh ${MQ_HOST}:${MQ_PORT} -- python /app/consume.py"]