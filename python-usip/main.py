from typing import Annotated

from fastapi import FastAPI, Header, Request, Depends, Query

from data import verify_token, get_user, get_members, get_role, Member

app = FastAPI()


@app.get("/credential")
def verify_credential(
    x_authorization: Annotated[str | None, Header()] = None
):
    user_id = verify_token(x_authorization)
    if user_id is None:
        return {"error": "Invalid token"}
    
    user = get_user(user_id)
    if user is None:
        return {"error": "User not found"}

    return {"user": user.normalize()}


def get_userinfo_param(request: Request) -> list[str]:
    query = request.url.query
    user_ids = []
    for u in query.split("&"):
        (key, v) = u.split("=")
        if key == "userIDs":
            user_ids.append(v)
    return user_ids


@app.get("/userinfo")
def batch_get_userinfo(
    user_ids: list[str] = Depends(get_userinfo_param)
):
    users = []
    for user_id in user_ids:
        user = get_user(user_id)
        if user is not None:
            users.append(user.normalize())
    return {"users": users}


@app.get("/role")
def get_unit_role(
    unit_id: str = Query(alias="unitID"),
    user_id: str = Query(alias="userID"),
):
    member = get_role(unit_id, user_id)
    if member is None:
        return {"error": "Role not found"}
    
    return member.normalize()


def get_members_by_unit_ids_param(request: Request) -> list[str]:
    query = request.url.query
    unit_ids = []
    for u in query.split("&"):
        (key, v) = u.split("=")
        if key == "unitIDs":
            unit_ids.append(v)
    return unit_ids


@app.get("/collaborators")
def get_members_by_unit_ids(
    unit_ids: list[str] = Depends(get_members_by_unit_ids_param)
):
    members_map: dict[str, list[Member]] = {}
    for unit_id in unit_ids:
        members_map[unit_id] = get_members(unit_id)
    
    collaborators = []
    for unit_id, members in members_map.items():
        data = {"unitID": unit_id, "subjects": []}
        for m in members:
            user = get_user(m.user_id)
            data["subjects"].append({
                "role": m.role, 
                "subject": {
                    "id": user.user_id, 
                    "name": user.name, 
                    "avatar": user.avatar,
                    "type": "user",
                }
            })
        collaborators.append(data)
    
    return {"collaborators": collaborators}