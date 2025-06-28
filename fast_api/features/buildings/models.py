from sqlalchemy import Column, Integer, String, Float, ForeignKey, TIMESTAMP, func
from sqlalchemy.orm import relationship

from database.db import Base


class DefBuildings(Base):

    __tablename__ = 'def_buildings'

    id = Column(Integer,primary_key=True, index=True)
    token_name = Column(String)
    base_construction_cost = Column(Float)
    base_construction_time = Column(Integer)

    buildings = relationship("Buildings", back_populates="definition")

class Buildings(Base):

    __tablename__ = 'buildings'

    id = Column(Integer, primary_key=True, index=True)
    user_id = Column(Integer,ForeignKey("users.id"), nullable=False)
    def_id = Column(Integer, ForeignKey("def_buildings.id"), nullable=False)
    name = Column(String)
    time_build = Column(TIMESTAMP(timezone=True), server_default=func.now(), nullable=False)

    owner = relationship("Users", back_populates="buildings")
    definition = relationship("DefBuildings", back_populates="buildings")