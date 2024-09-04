#[macro_use] extern crate rocket;

use rocket::Request;
use rocket::request::{self, FromRequest};
use rocket::serde::json::{json, Value, Json};
use rocket::serde::{Deserialize, Serialize};

mod data;



#[derive(Debug)]
struct Headers {
    headers: Vec<(String, String)>,
}

#[rocket::async_trait]
impl<'r> FromRequest<'r> for Headers {
    type Error = ();

    async fn from_request(req: &'r Request<'_>) -> request::Outcome<Self, Self::Error> {
        let headers: Vec<(String, String)> = req.headers()
            .iter()
            .map(|h| (h.name().to_string(), h.value().to_string()))
            .collect();

        request::Outcome::Success(Headers { headers })
    }
}

fn get_token(header: &Headers) -> String {
    let token = header.headers.iter().find(|(k, _)| k == "x-authorization");
    match token {
        Some((_, v)) => v.to_string(),
        _none => "".to_string(),
    }
}

#[get("/credential")]
fn verify_credential(header: Headers) -> Value {
    let token = get_token(&header);
    let (user_id, ok) = data::verify_token(&token);
    if !ok {
        return json!({"error": "Invalid token"})
    }

    let user = data::get_user(&user_id);
    match user {
        Some(u) => json!({
            "users": {
                "userID": u.user_id,
                "name": u.name,
                "avatar": u.avatar,
            },
        }),
        _none => json!({"error": "user not found"})
    }
}


#[derive(Debug, Serialize, Deserialize)]
#[serde(crate = "rocket::serde")]
struct UserInfoParam {
    #[serde(alias = "userIDs")]
    user_ids: Vec<String>,
}

#[post("/userinfo", format = "json", data = "<p>")]
fn batch_get_userinfo(p: Json<UserInfoParam>) -> Value {
    // println!("{:?}", query.user_ids);
    if p.user_ids.is_empty() {
        return json!({"error": "userIDs is required"})
    }

    let mut users = Vec::new();
    for user_id in p.user_ids.iter() {
        let user = data::get_user(user_id);
        match user {
            Some(u) => users.push(json!({
                "userID": u.user_id,
                "name": u.name,
                "avatar": u.avatar,
            })),
            _none => users.push(json!({"error": "user not found"})),
        }
    }

    json!({"users": users})
}

#[derive(Debug)]
struct GetRoleParam {
    user_id: String,
    unit_id: String,
}

#[rocket::async_trait]
impl<'r> FromRequest<'r> for GetRoleParam {
    type Error = ();

    async fn from_request(req: &'r Request<'_>) -> request::Outcome<Self, Self::Error> {
        let unit_id = req.query_value::<String>("unitID");
        let user_id = req.query_value::<String>("userID");
        match (unit_id, user_id) {
            (Some(u), Some(v)) => request::Outcome::Success(GetRoleParam {
                user_id: v.unwrap(),
                unit_id: u.unwrap(),
            }),
            _ => request::Outcome::Success(GetRoleParam {
                user_id: "".to_string(),
                unit_id: "".to_string(),
            })
        }
    }
}

#[get("/role")]
fn get_role(query: GetRoleParam) -> Value {
    if query.user_id.is_empty() || query.unit_id.is_empty() {
        return json!({"error": "userID and unitID are required"})
    }

    let member = data::get_role(&query.unit_id, &query.user_id);
    match member {
        Some(m) => json!({
            "unitID": m.unit_id,
            "userID": m.user_id,
            "role": m.role.to_string(),
        }),
        _none => json!({"error": "role not found"})
    }
}

#[derive(Debug, Serialize, Deserialize)]
#[serde(crate = "rocket::serde")]
struct CollaboratorsParam {
    #[serde(alias = "unitIDs")]
    unit_ids: Vec<String>,
}

#[post("/collaborators", format = "json", data = "<p>")]
fn get_collaborators(p: Json<CollaboratorsParam>) -> Value {
    // println!("{:?}", query.unit_ids);
    if p.unit_ids.is_empty() {
        return json!({"error": "unitIDs is required"})
    }
    let mut collaborators = Vec::new();
    for unit_id in p.unit_ids.iter() {
        let members = data::get_members(unit_id);
        let mut unit_collaborators = Vec::new();
        for member in members.iter() {
            let user = data::get_user(&member.user_id);
            match user {
                Some(u) => unit_collaborators.push(json!({
                    "role": member.role.to_string(),
                    "subject": {
                        "id": u.user_id,
                        "name": u.name,
                        "avatar": u.avatar,
                        "type": "user",
                    }
                })),
                _none => {},
            }
        }

        collaborators.push(json!({
            "unitID": unit_id,
            "subjects": unit_collaborators,
        }));
    }
    json!({"collaborators": collaborators})
}

#[launch]
fn rocket() -> _ {
    rocket::build().mount("/", routes![
        verify_credential,
        batch_get_userinfo,
        get_role,
        get_collaborators,
    ])
}
