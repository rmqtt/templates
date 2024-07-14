use salvo::http::Method;
use salvo::prelude::*;
use std::collections::HashMap;

#[handler]
async fn auth(req: &mut Request, res: &mut Response) {
    let params = match parse_params(req).await {
        Ok(params) => params,
        Err(e) => {
            return res.set_status_error(StatusError::bad_request().with_detail(e.to_string()))
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
    if username == "admin" {
        let _ = res.add_header("X-Superuser", "true", true);
    }

    // allow
    res.render(Text::Plain("allow"));
}

#[handler]
async fn acl(req: &mut Request, res: &mut Response) {
    let params = match parse_params(req).await {
        Ok(params) => params,
        Err(e) => {
            return res.set_status_error(StatusError::bad_request().with_detail(e.to_string()))
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
    Server::new(TcpListener::bind("0.0.0.0:9090"))
        .serve(router)
        .await;
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
