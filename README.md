# IN DEVELOPMENT

# Study case for some exceptions that can happen when handles with notification system

My favourite thing about a resilience system, is the way we handle the user while the flames its everywhere. Something was broken and the user dont need to know about this. The only thing the user needs be occupied its with experience of the product.

This test case its a representation how to a queue can be usefull to reprocess request using channels and gourutines. This gives a most chance to a trully response if something dont acts like expected.

In my case, i am implementing a API with GoLang that needs notify multiple users with Firebase and Discord. The main problem is notify the user about his progress on the application. Like a game when you do some experience progression and needs be notified when you up the level. In this case, these integrations can be resumed as simple HTTP request with header, body and especifies authentication rules.

The queue needs receive a representation of this request to retry in some moment when status is between 400 and 500.

## Example 1 - Sending message to a group channel with Firebase HTTP/1

    {
        "notify": {
            "request": {
                "url": "https://webhook.site/39c3e318-d788-49a6-8eb6-b11e9dddc4cd",
                "timeout": 30,
                "header": {
                    "content-type": "application/json",
                    "teste": "blablabla"
                }
            }
        },
        "request": {
            "url": "https://webhook.site/7b2d3a05-c05b-450c-b38c-725b3fabc1a5",
            "verb": "POST",
            "timeout": 30,
            "header": {
                "content-type": "application/json",
                "teste": "blablabla"
            },
            "body": {
                "firebaseToken": "blablabla"
            }
        }
    }
