use lazy_static::lazy_static;
use std::collections::HashMap;
use core::fmt::{Display, Formatter, Result};

pub struct User {
    pub user_id: String,
    pub name: String,
    pub avatar: String,
}

#[derive(Debug)]
pub enum Role {
    Owner,
    Editor,
    Reader,
}

impl Display for Role {
    fn fmt(&self, f: &mut Formatter) -> Result {  
        match self {  
            Role::Owner => write!(f, "owner"),
            Role::Editor => write!(f, "editor"),
            Role::Reader => write!(f, "reader"),
		}  
    }
}

#[derive(Debug)]
pub struct Member {
    pub user_id: String,
    pub unit_id: String,
    pub role: Role,
}

lazy_static! {
    static ref X_AUTHORIZATION: HashMap<String, String> = {
        let mut map = HashMap::new();
        map.insert("token:1".to_string(), "1".to_string());
        map.insert("token:2".to_string(), "2".to_string());
        map.insert("token:3".to_string(), "3".to_string());
        map
    };

    static ref Users: Vec<User> = {
        let mut users = Vec::new();
        users.push(User {
            user_id: "1".to_string(),
            name: "Alice".to_string(),
            avatar: "https://example.com/alice.jpg".to_string(),
        });
        users.push(User {
            user_id: "2".to_string(),
            name: "Bob".to_string(),
            avatar: "https://example.com/bob.jpg".to_string(),
        });
        users.push(User {
            user_id: "3".to_string(),
            name: "Charlie".to_string(),
            avatar: "https://example.com/charlie.jpg".to_string(),
        });
        users
    };

    static ref Members: Vec<Member> = {
        let mut members = Vec::new();
        members.push(Member {
            unit_id: "unit1".to_string(),
            user_id: "1".to_string(),
            role: Role::Owner,
        });
        members.push(Member {
            unit_id: "unit1".to_string(),
            user_id: "2".to_string(),
            role: Role::Editor,
        });
        members.push(Member {
            unit_id: "unit2".to_string(),
            user_id: "2".to_string(),
            role: Role::Owner,
        });
        members.push(Member {
            unit_id: "unit2".to_string(),
            user_id: "3".to_string(),
            role: Role::Reader,
        });
        members
    };
}

pub fn verify_token(token: &String) -> (String, bool) {
    let user_id = X_AUTHORIZATION.get(token);
    match user_id {
        Some(id) => (id.to_string(), true),
        _none => ("".to_string(), false),
    }
}

pub fn get_user(user_id: &String) -> Option<&User> {
    Users.iter().find(|u| u.user_id == *user_id)
}

pub fn get_role<'r>(unit_id: &String, user_id: &String) -> Option<&'r Member> {
    Members.iter().find(|m| m.unit_id == *unit_id && m.user_id == *user_id)
}

pub fn get_members(unit_id: &String) -> Vec<&Member> {
    Members.iter().filter(|m| m.unit_id == *unit_id).collect()
}