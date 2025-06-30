from datetime import timedelta, datetime, timezone
from typing import Annotated

from jose import jwt, JWTError
from fastapi import APIRouter, Depends, HTTPException
from passlib.context import CryptContext
from sqlalchemy.orm import Session
from starlette import status

from database.db import SessionLocal
from features.authentication.schema import TokenRequest
from features.users.models import Users
from main import oauth2_scheme

SECRET_KEY = '197b2c37c391bed93fe80344fe73b806947a65e36206e05a1a23c2fa12702fe3'
ALGORITHM = 'HS256'

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

def create_access_token(username:str, user_id:int, role:int, expires_delta: timedelta ):
    expires = datetime.now(timezone.utc) + expires_delta
    encode = {'sub': username, 'id': user_id, 'role': role, 'exp':expires}
    return jwt.encode(encode, SECRET_KEY, ALGORITHM)

async def get_current_user(token: str = Depends(oauth2_scheme)):
    credentials_exception = HTTPException(
        status_code=status.HTTP_401_UNAUTHORIZED,
        detail="Could not validate credentials",
        headers={"WWW-Authenticate": "Bearer"},
    )
    try:
        payload = jwt.decode(token, SECRET_KEY, algorithms=[ALGORITHM])
        username: str = payload.get("sub")
        user_id: int = payload.get('id')
        role: int = payload.get('role')
        if username is None:
            raise credentials_exception
    except JWTError:
        raise credentials_exception
    return {'username': username, 'user_id': user_id, 'role': role}

@router.post('/token')
async def create_access_token(body:TokenRequest, db:db_dependency):

    authenticated_user = authenticate_user(username=body.username,password=body.password, db=db)
    if not authenticated_user:
        raise HTTPException(status_code=400, detail="Wrong username or password")
    token = create_access_token(body.username, authenticated_user.id, authenticated_user.role, timedelta(minutes=20))
    return {'access_token': token, 'token_type': 'bearer'}