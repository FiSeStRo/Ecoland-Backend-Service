from typing import Annotated

from fastapi import APIRouter, Depends
from passlib.context import CryptContext
from sqlalchemy.orm import Session

from starlette import status

from database.db import SessionLocal
from features.users.models import Users
from features.users.schemas import CreateUser

bcrypt_context = CryptContext(schemes=['bcrypt'], deprecated='auto')

router = APIRouter(prefix="/users", tags=["users"])

def get_db():
    db =SessionLocal()
    try:
        yield db
    finally:
        db.close()

db_dependency = Annotated[Session, Depends(get_db)]

@router.post("/create", status_code=status.HTTP_201_CREATED)
async def create_new_user(body: CreateUser, db: db_dependency):

    user_model = Users(
    username=body.username,
    email= body.email,
    password = bcrypt_context.hash(body.password),
    role = 2
    )
    db.add(user_model)
    db.commit()
    return

# TODO: This should only be reachable for admin accounts
@router.get("/all")
async def get_all_users(db: db_dependency):
    return db.query("users").all()