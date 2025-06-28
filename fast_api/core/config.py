import os

DB_USER= os.getenv("DB_USER")
DB_PW=os.getenv("DB_PW")
DB_HOST= os.getenv("DB_HOST")
DB_PORT= os.getenv("DB_PORT")
DB_NAME= os.getenv("DB_NAME")

DATABASE_URL = f"mariadb+mariadb://{DB_USER}:{DB_PW}@{DB_HOST}:{DB_PORT}/{DB_NAME}"