from pydantic import BaseModel, EmailStr
from uuid import UUID
from datetime import datetime
from typing import Optional

class EmailCreate(BaseModel):
    """
    Request format to queue a new email.
    This is what the API Gateway will send to us.
    """
    user_id: UUID
    to_email: EmailStr
    subject: str
    body: str
    
    class Config:
        json_schema_extra = {
            "example": {
                "user_id": "123e4567-e89b-12d3-a456-426614174000",
                "to_email": "user@example.com",
                "subject": "Welcome to our app!",
                "body": "<h1>Welcome!</h1><p>Thanks for signing up.</p>"
            }
        }

class EmailResponse(BaseModel):
    """
    Response format when returning email info.
    """
    id: UUID
    user_id: UUID
    to_email: str
    subject: str
    status: str
    created_at: datetime
    sent_at: Optional[datetime] = None
    error_message: Optional[str] = None
    
    class Config:
        from_attributes = True

class StandardResponse(BaseModel):
    """
    API response format
    """
    success: bool
    data: Optional[EmailResponse] = None
    error: Optional[str] = None
    message: str