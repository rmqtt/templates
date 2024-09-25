import web
import json
import os
os.environ["PORT"] = "9090"

CONTENT_TYPE_FORM = "application/x-www-form-urlencoded"
CONTENT_TYPE_JSON = "application/json"

class BaseService:

    def GET(self):
        try:
            params = web.input()
            return self.verify(params)
        except Exception as e:
            print(e)
            return web.badrequest()

    def POST(self):
        try:
            ctype = web.ctx.env.get("CONTENT_TYPE")
            if ctype is None:
                print("CONTENT_TYPE is None:")
                return web.badrequest()

            if ctype.startswith(CONTENT_TYPE_FORM):
                params = web.input()
            elif ctype.startswith(CONTENT_TYPE_JSON):
                params = json.loads(web.data())
            return self.verify(params)

        except Exception as e:
            print(e)
            return web.badrequest()
        
    def PUT(self):
        return self.POST()
        
    def verify(self, params):
        return web.internalerror()


class AuthService(BaseService):

    def verify(self, params):
        clientid = params["clientid"] if "clientid" in params else ""
        username = params["username"] if "username" in params else ""
        password = params["password"] if "password" in params else ""
        protocol = params["protocol"] if "protocol" in params else ""
        print("AuthService clientid:", clientid)
        print("AuthService username:", username)
        print("AuthService password:", password)
        print("AuthService protocol:", protocol)
        #
        # @TODO Verify user validity, 
        #

        #is user-ignore
        if username == "user-ignore":
            return "ignore"

        #is user-deny
        if username == "user-deny":
            return "deny"

        #is admin
        if username == "user-admin":
            web.header("X-Superuser", "true")

        #acl
        if username == "user-acl":
            data = {
                    "result":"allow",
                    "superuser": False,
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
                        "retain": False,
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
                        "retain": True
                        }
                    ]
                }
            web.header('Content-Type', CONTENT_TYPE_JSON)
            return json.dumps(data)


        #other
        return "allow"


class ACLService(BaseService):

    def verify(self, params):
        #access = "%A", username = "%u", protocol = "%r", clientid = "%c", ipaddr = "%a", topic = "%t"
        access = params["access"] if "access" in params else ""
        clientid = params["clientid"] if "clientid" in params else ""
        username = params["username"] if "username" in params else ""
        protocol = params["protocol"] if "protocol" in params else ""
        ipaddr = params["ipaddr"] if "ipaddr" in params else ""
        topic = params["topic"] if "topic" in params else ""
        print("ACLService clientid:", clientid)
        print("ACLService username:", username)
        print("ACLService protocol:", protocol)
        print("ACLService access:", access)
        print("ACLService ipaddr:", ipaddr)
        print("ACLService topic:", topic)
        #
        # @TODO Verify topic validity, 
        #

        # is Subscribe
        if access == "1":
           print("is Subscribe, topic is ", topic)

        if access == "2":
           print("is Publish, topic is ", topic)

        if topic.endswith("/cache"):
            web.header("X-Cache", "-1") #Unit millisecond, if the value is -1, it will not expire before disconnecting

        #test ignore
        if topic.startswith("/test/ignore"):
            return "ignore"

        #test deny
        if topic.startswith("/test/deny"):
            return "deny"


        #other
        return "allow"


def main():
    urls = (
        '/mqtt/auth', 'AuthService',
        '/mqtt/acl', 'ACLService',
    )
    app = web.application(urls, globals())
    app.run()

if __name__ == '__main__':
    main()
    
    