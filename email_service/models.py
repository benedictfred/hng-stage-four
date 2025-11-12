import enum
from uuid import uuid4
from datetime import datetime
from sqlalchemy import Column, String, Text, Enum, DateTime
from sqlalchemy.dialects.postgresql import UUID
from db import Base

class EmailStatus(str, enum.Enum):
    """Email status enum"""
    queued = "queued"
    processing = "processing"
    sent = "sent"
    failed = "failed"

class EmailMessage(Base):
    """
    Stores all email messages in the system.
    """
    __tablename__ = "email_messages"

    id = Column(UUID(as_uuid=True), primary_key=True, default=uuid4)
    
    # Email details
    user_id = Column(UUID(as_uuid=True), nullable=False)
    to_email = Column(String(255), nullable=False)
    subject = Column(String(500), nullable=False)
    body = Column(Text, nullable=False)
    
    # Status tracking
    status = Column(Enum(EmailStatus), default=EmailStatus.queued, nullable=False)
    error_message = Column(Text, nullable=True)
    
    
    created_at = Column(DateTime, default=datetime.utcnow)
    sent_at = Column(DateTime, nullable=True)
    
    def __repr__(self):
        return f"<EmailMessage {self.id} - {self.status}>"