use salvo::http::Method;
use salvo::prelude::*;
use std::collections::HashMap;

#[macro_use]
pub extern crate serde_json;

#[handler]
async fn auth(req: &mut Request, res: &mut Response) {
    let params = match parse_params(req).await {
        Ok(params) => params,
        Err(e) => {
            res.render(StatusError::bad_request().detail(e.to_string()));
            return
        }
    };

    let clientid = params
        .get("clientid")
        .map(|v| v.as_str())
        .unwrap_or_default();
    let username = params
        .get("username")
        .map(|v| v.as_str())
        .unwrap_or_default();
    let password = params
        .get("password")
        .map(|v| v.as_str())
        .unwrap_or_default();
    let protocol = params
        .get("protocol")
        .map(|v| v.as_str())
        .unwrap_or_default();
    println!("auth clientid: {}", clientid);
    println!("auth username: {}", username);
    println!("auth password: {}", password);
    println!("auth protocol: {}", protocol);

    // @TODO Verify user validity,

    // is user-ignore
    if username == "user-ignore" {
        res.render(Text::Plain("ignore"));
        return;
    }

    // is user-deny
    if username == "user-deny" {
        res.render(Text::Plain("deny"));
        return;
    }

    // is admin
    if username == "user-admin" {
        let _ = res.add_header("X-Superuser", "true", true);
    }

    //acl
    if username == "user-acl" {
        let json_acl = json!({
            "result":"allow",
            "superuser": false,
            "expire_at": 1827143027,
            "acl": [
                {
                  "permission": "allow",
                  "action": "all",
                  "topic": "foo/${clientid}"
                },
                {
                  "permission": "allow",
                  "action": "subscribe",
                  "topic": "eq foo/1/#",
                  "qos": [1,2]
                },
                {
                  "permission": "allow",
                  "action": "subscribe",
                  "topic": "foo/2/#",
                  "qos": 1
                },
                {
                  "permission": "allow",
                  "action": "publish",
                  "topic": "foo/2/1",
                  "qos": 1
                },
                {
                  "permission": "allow",
                  "action": "publish",
                  "topic": "foo/${username}",
                  "retain": false,
                  "qos": [0,1]
                },
                {
                  "permission": "deny",
                  "action": "all",
                  "topic": "foo/3"
                },
                {
                  "permission": "deny",
                  "action": "publish",
                  "topic": "foo/4",
                  "retain": true
                }
            ]
        });
        res.render(Json(json_acl));
        return;
    }

    // allow
    res.render(Text::Plain("allow"));
}

#[handler]
async fn acl(req: &mut Request, res: &mut Response) {
    let params = match parse_params(req).await {
        Ok(params) => params,
        Err(e) => {
            res.render(StatusError::bad_request().detail(e.to_string()));
            return;
        }
    };

    //access = "%A", username = "%u", clientid = "%c", ipaddr = "%a", topic = "%t"
    let access = params.get("access").map(|v| v.as_str()).unwrap_or_default();
    let clientid = params
        .get("clientid")
        .map(|v| v.as_str())
        .unwrap_or_default();
    let username = params
        .get("username")
        .map(|v| v.as_str())
        .unwrap_or_default();
    let protocol = params
        .get("protocol")
        .map(|v| v.as_str())
        .unwrap_or_default();
    let ipaddr = params.get("ipaddr").map(|v| v.as_str()).unwrap_or_default();
    let topic = params.get("topic").map(|v| v.as_str()).unwrap_or_default();

    println!("acl clientid: {}", clientid);
    println!("acl username: {}", username);
    println!("acl protocol: {}", protocol);
    println!("acl access: {}", access);
    println!("acl ipaddr: {}", ipaddr);
    println!("acl topic: {}", topic);

    // @TODO Verify topic validity,

    if access == "1" {
        println!("is Subscribe, topic is {}", topic);
    }

    if access == "2" {
        println!("is Publish, topic is {}", topic);
    }

    if topic.ends_with("/cache") {
        //Unit millisecond, if the value is -1, it will not expire before disconnecting
        let _ = res.add_header("X-Cache", "-1", true);
    }

    // test ignore
    if topic.starts_with("/test/ignore") {
        res.render(Text::Plain("ignore"));
        return;
    }

    // test deny
    if topic.starts_with("/test/deny") {
        res.render(Text::Plain("deny"));
        return;
    }

    // allow
    res.render(Text::Plain("allow"));
}

#[tokio::main]
async fn main() {
    let router = Router::new()
        .push(
            Router::with_path("/mqtt/auth")
                .get(auth)
                .post(auth)
                .put(auth),
        )
        .push(Router::with_path("/mqtt/acl").get(acl).post(acl).put(acl));

    let laddr = "0.0.0.0:9090";
    println!("Auth HTTP Server Listening on {}", laddr);
    let acceptor = TcpListener::new(laddr).bind().await;
    Server::new(acceptor).serve(router).await;
}

async fn parse_params(req: &mut Request) -> Result<HashMap<String, String>, anyhow::Error> {
    match req.method() {
        &Method::GET => Ok(req.parse_queries::<HashMap<String, String>>()?),
        &Method::POST | &Method::PUT => {
            if let Some(ctype) = req.content_type() {
                match ctype.essence_str() {
                    "application/x-www-form-urlencoded" => {
                        Ok(req.parse_form::<HashMap<String, String>>().await?)
                    }
                    "application/json" => Ok(req.parse_json::<HashMap<String, String>>().await?),
                    _ => Err(anyhow::Error::msg(format!(
                        "content type({:?}) not supported",
                        ctype
                    ))),
                }
            } else {
                Err(anyhow::Error::msg("content type is not exist"))
            }
        }
        _ => Err(anyhow::Error::msg(format!(
            "method({}) not supported",
            req.method()
        ))),
    }
}
