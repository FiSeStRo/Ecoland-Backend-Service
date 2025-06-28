from sqlalchemy import create_engine, Engine
from sqlalchemy.orm import sessionmaker, declarative_base

from core.config import DATABASE_URL

engine: Engine = create_engine(DATABASE_URL)

SessionLocal = sessionmaker(autocommit=False, autoflush=False, bind=engine)

Base = declarative_base()