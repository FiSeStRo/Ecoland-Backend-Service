from datetime import timedelta, datetime, timezone
from typing import Annotated

from fastapi.security import OAuth2PasswordBearer, OAuth2PasswordRequestForm
from jose import jwt, JWTError
from fastapi import APIRouter, Depends, HTTPException
from passlib.context import CryptContext
from sqlalchemy.orm import Session
from starlette import status

from database.db import SessionLocal
from features.authentication.schema import TokenRequest
from features.authentication.service import SECRET_KEY, ALGORITHM
from features.users.models import Users




router = APIRouter(prefix='/auth', tags=['auth'])

def get_db():
    db = SessionLocal()
    try:
        yield db
    finally:
        db.close()


db_dependency = Annotated[Session, Depends(get_db)]
bcrypt_context = CryptContext(schemes=['bcrypt'], deprecated='auto')

def authenticate_user(username: str, password: str, db):
    user = db.query(Users).filter(Users.username == username).first()
    if not user:
       return False
    if not bcrypt_context.verify(password, user.password):
        return False
    return user

def generate_jwt_token(username:str, user_id:int, role:int, expires_delta: timedelta ):
    expires = datetime.now(timezone.utc) + expires_delta
    encode = {'sub': username, 'id': user_id, 'role': role, 'exp':expires}
    return jwt.encode(encode, SECRET_KEY, ALGORITHM)



@router.post('/token')
async def create_access_token(form_data: Annotated[OAuth2PasswordRequestForm, Depends()], db:db_dependency):

    authenticated_user = authenticate_user(username=form_data.username,password=form_data.password, db=db)
    if not authenticated_user:
        raise HTTPException(status_code=400, detail="Wrong username or password")
    token = generate_jwt_token(form_data.username, authenticated_user.id, authenticated_user.role, timedelta(minutes=20))
    return {'access_token': token, 'token_type': 'bearer'}