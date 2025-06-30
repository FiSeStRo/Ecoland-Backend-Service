from pydantic import BaseModel, Field


class TokenRequest(BaseModel):
    username:str = Field(min_length=3)
    password:str = Field(min_length=8)