import os
from smtplib import SMTP
from email.mime.text import MIMEText
from email.mime.multipart import MIMEMultipart
from dotenv import load_dotenv

load_dotenv()

def send_email(to_email: str, subject: str, body: str):
    """
    Send an email using SMTP.
    
    Args:
        to_email: Recipient email address
        subject: Email subject
        body: Email body (can be HTML)
    
    Raises:
        Exception: If email fails to send
    """
    # SMTP configuration
    smtp_host = os.getenv("SMTP_HOST", "smtp.gmail.com")
    smtp_port = int(os.getenv("SMTP_PORT", "587"))
    smtp_user = os.getenv("SMTP_USER")
    smtp_pass = os.getenv("SMTP_PASS")
    sender = os.getenv("EMAIL_SENDER", "noreply@example.com")
    
    # Create email message
    msg = MIMEMultipart('alternative')
    msg['Subject'] = subject
    msg['From'] = sender
    msg['To'] = to_email
    
    # Add HTML body
    html_part = MIMEText(body, 'html')
    msg.attach(html_part)
    
    # Send email
    print(f"[smtp] Connecting to {smtp_host}:{smtp_port}...")
    with SMTP(smtp_host, smtp_port) as server:
        server.starttls()
        
        if smtp_user and smtp_pass:
            print(f"[smtp] Logging in as {smtp_user}...")
            server.login(smtp_user, smtp_pass)
        
        print(f"[smtp] Sending email to {to_email}...")
        server.sendmail(sender, [to_email], msg.as_string())
        print(f"[smtp] Email sent successfully!")