# Sestro Visit (Manage user visits and sessions)

## API Docs

https://api.stage.sestro.io/v1/visits/apidocs.json

## Config

Please change the env variables in app.yaml 
at https://github.com/SestroAI/envrionments/blob/master/stage.sestro-165123/visits.yaml

## API Specs

```
Base Path: /v1/vms/

1) POST /visits/

   data: 
    {
        "merchantId" : "grren-barn",
        "tableId" : "green-barn-1",
        "geoLatitude" : "",
        "geoLongitude" : ""
    }

    Output:
    {
        "visitId" : "<uuid>",
        "sessionId" : "<uuid>"
    }

Will create a new visit if this is the first user, otherwise add to existing


2) GET /sessions/{session_id}

3) PUT /sessions/{session_id}/orders

    data:

    {
        "items" : [<item_id_1>, <item_id_2>]
    }

    Output:
    Session Object

4) PUT /sessions/{session_id}/orders/{order_id}/{status}

    data :
    {}

    Output:
    Session Object

 5) POST /visits/{visit-id}/end

    data:

    {
        "guestRating" : {
            "reviewerId" : "", //Current user ID
            "value" : 5,
            "comments" : ""
        }
    }

    Output:

    200 OK

```
