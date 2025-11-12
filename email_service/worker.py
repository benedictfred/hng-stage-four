import json
import pika
from datetime import datetime

from db import SessionLocal
from models import EmailMessage, EmailStatus
from services import send_email
from task_queue import get_connection, QUEUE_NAME

def process_email(email_id: str):
    """
    Process a single email:
    1. Get email from database
    2. Send it via SMTP
    3. Update status in database
    """
    db = SessionLocal()
    
    try:
        # Get email from database
        email = db.query(EmailMessage).filter(EmailMessage.id == email_id).first()
        
        if not email:
            print(f"[worker] Email {email_id} not found in database")
            return
        
        # Check if already sent
        if email.status == EmailStatus.sent:
            print(f"[worker] Email {email_id} already sent, skipping")
            return
        
        # Update status to processing
        email.status = EmailStatus.processing
        db.commit()
        print(f"[worker] Processing email {email_id} to {email.to_email}")
        
        # Send the email
        send_email(
            to_email=email.to_email,
            subject=email.subject,
            body=email.body
        )
        
        # Update status to sent
        email.status = EmailStatus.sent
        email.sent_at = datetime.utcnow()
        db.commit()
        
        print(f"[worker] ✓ Email {email_id} sent successfully!")
        
    except Exception as e:
        # If anything fails, mark as failed
        print(f"[worker] ✗ Failed to send email {email_id}: {str(e)}")
        email.status = EmailStatus.failed
        email.error_message = str(e)
        db.commit()
        
    finally:
        db.close()

def callback(ch, method, properties, body):
    """
    Callback function called by RabbitMQ when a message arrives.
    """
    try:
        # Parse the message
        message = json.loads(body)
        email_id = message['email_id']
        
        print(f"\n[worker] Received email job: {email_id}")
        
        # Process the email
        process_email(email_id)
        
        # Acknowledge message (remove from queue)
        ch.basic_ack(delivery_tag=method.delivery_tag)
        print(f"[worker] Message acknowledged\n")
        
    except Exception as e:
        print(f"[worker] Error processing message: {str(e)}")
        # Acknowledge anyway to remove bad message from queue
        ch.basic_ack(delivery_tag=method.delivery_tag)

def start_worker():
    """
    Start the worker - listens to the queue and processes emails.
    """
    print("=" * 50)
    print("EMAIL WORKER STARTING")
    print("=" * 50)
    
    # Connect to RabbitMQ
    connection = get_connection()
    channel = connection.channel()
    
    # Make sure queue exists
    channel.queue_declare(queue=QUEUE_NAME, durable=True)
    
    # Process one message at a time
    channel.basic_qos(prefetch_count=1)
    
    # Start consuming messages
    channel.basic_consume(
        queue=QUEUE_NAME,
        on_message_callback=callback
    )
    
    print(f"[worker] Listening to queue: {QUEUE_NAME}")
    print("[worker] Waiting for emails... Press CTRL+C to exit\n")
    
    try:
        channel.start_consuming()
    except KeyboardInterrupt:
        print("\n[worker] Shutting down...")
        channel.stop_consuming()
    finally:
        connection.close()

if __name__ == "__main__":
    start_worker()