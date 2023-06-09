# IN DEVELOPMENT

# Study case for some exceptions that can happen when handles with notification system

My favourite thing about a resilience system, is the way we handle the user while the flames its everywhere. Something was broken and the user dont need to know about this. The only thing the user needs be occupied its with experience of the product.

This test case its a representation how to a queue can be usefull to reprocess request using channels and gourutines. This gives a most chance to a trully response if something dont acts like expected. 

In my case, i am implementing a API with GoLang that needs notify multiple users with Firebase and Discord. The main problem is notify the user about his progress on the application. Like a game when you do some experience progression and needs be notified when you up the level. In this case, these integrations can be resumed as simple HTTP request with header, body and especifies authentication rules. 

The queue needs receive a representation of this request to retry in some moment.

## Example 1 - Sending message to a group channel with Firebase HTTP/1 

    {
        "request": {
            "url": "https://fcm.googleapis.com/v1/projects/myproject-b5ae1/messages:send",
            "verb": "POST",
            "timeout": 30,
            "header": {
                "Authorization": "Baerer blablabla"
            },
            "body": {
                "message": {
                    "topic": "news",
                    "notification": {
                        "title": "Breaking News",
                        "body": "New news story available."
                    },
                    "data": {
                        "story_id": "story_12345"
                    }
                }
            }
        },
        "response": {
            "code": 200
        }
    }
 
