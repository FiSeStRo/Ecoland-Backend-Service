from fastapi import APIRouter
router = APIRouter(
    prefix='/user',
    tags=['user']
)

@router.get("/resources")
async def get_resources():

    return "get resources"

@router.get("/info")
async def get_info():
    return "user info"

@router.patch("/info")
async def update_user():
    return "update user"