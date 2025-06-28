from pydantic import BaseModel, Field, EmailStr


class CreateUser(BaseModel):
    username: str = Field(min_length=3)
    email:EmailStr = Field(min_length=3)
    password:str = Field(min_length=8)
