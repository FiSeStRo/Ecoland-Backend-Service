from pydantic import BaseModel, Field
from sqlalchemy import Column, Integer, String, DECIMAL, ForeignKey, DateTime, func
from sqlalchemy.dialects.mysql import TIMESTAMP
from sqlalchemy.orm import relationship

from database.db import Base

class ConstructBuildingRequest(BaseModel):
    def_id: int = Field()
    display_name:str = Field(min_length=2)