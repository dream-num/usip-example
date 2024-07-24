

class User():
    def __init__(self, user_id: str, name: str, avatar: str):
        self.user_id = user_id
        self.name = name
        self.avatar = avatar

    def normalize(self) -> dict[str, any]:
        return {
            "userID": self.user_id,
            "name": self.name,
            "avatar": self.avatar,
        }


class Member():
    def __init__(self, unit_id: str, user_id: str, role: str):
        self.unit_id = unit_id
        self.user_id = user_id
        self.role = role

    def normalize(self) -> dict[str, any]:
        return {
            "unitID": self.unit_id,
            "userID": self.user_id,
            "role": self.role,
        }


class Role():
    OWNER = "owner"
    EDITOR = "editor"
    READER = "reader"


x_authorization = {
    "token:1": "1",
    "token:2": "2",
    "token:3": "3",
}


users = [
    User("1", "user1", "avatar1"),
    User("2", "user2", "avatar2"),
    User("3", "user3", "avatar3"),
]


members = [
    Member("unit1", "1", Role.OWNER),
    Member("unit1", "2", Role.EDITOR),
    Member("unit2", "2", Role.OWNER),
    Member("unit2", "3", Role.READER),
    Member("unit3", "3", Role.OWNER),
    Member("unit3", "1", Role.EDITOR),
]


# return user_id if verify token success
def verify_token(token: str) -> str|None:
    return x_authorization.get(token, None)


def get_user(user_id: str) -> User|None:
    for user in users:
        if user.user_id == user_id:
            return user
    return None


def get_members(unit_id: str) -> list[Member]:
    return [m for m in members if m.unit_id == unit_id]


def get_role(unit_id: str, user_id: str) -> Member|None:
    for m in members:
        if m.unit_id == unit_id and m.user_id == user_id:
            return m
    return None