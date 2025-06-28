from dataclasses import Field
from datetime import datetime

from sqlalchemy import Column, String, Integer, DateTime, ForeignKey, Float, TIMESTAMP, func
from sqlalchemy.orm import relationship

from database.db import Base


class Users(Base):

    __tablename__ = 'users'

    id = Column(primary_key=True, index=True)
    username = Column(String)
    email  = Column(String)
    password= Column(String)
    role = Column(Integer)
    time_created = Column(TIMESTAMP(timezone=True), server_default=func.now(), nullable=False)
    time_last_active = Column(TIMESTAMP(timezone=True), nullable=True)

    buildings = relationship("Buildings", back_populates="owner")
    resources = relationship("UserResources", back_populates="owner")

class UserResources(Base):

    __tablename__ = 'user_resources'

    user_id = Column(Integer, ForeignKey("users.id"), primary_key=True, index=True)
    money = Column(Float)
    prestige = Column(Integer)

    owner = relationship("Users", back_populates="resources")

