import web
import json
import base64
import chardet
import os
os.environ["PORT"] = "5656"

class WebhookService:
    def POST(self):
        try:
            ctype = web.ctx.env.get("CONTENT_TYPE")
            if ctype is not None and ctype.startswith("application/json"):
                data = json.loads(web.data())
                print(data)
                if "action" in data and "payload" in data and (data["action"] == "message_publish" or data["action"] == "message_delivered" or data["action"] == "message_acked" or data["action"] == "message_dropped"):
                    try:
                        payload = base64.b64decode(data["payload"])
                        encoding = chardet.detect(payload)['encoding']
                        if encoding is not None:
                            payload = payload.decode(encoding)
                        print("payload:", payload)
                    except Exception as e:
                        print("base64.decode Exception:", e)
            else:
                print("CONTENT_TYPE is not application/json")
                return web.badrequest()
        except Exception as e:
            print(e)
            return web.badrequest()

def main():
    urls = (
        '/mqtt/webhook', 'WebhookService',
    )
    app = web.application(urls, globals())
    app.run()

if __name__ == '__main__':
    main()
    
    