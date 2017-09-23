# Sestro RMS (Restaurant Management)

## API Docs

https://api.stage.sestro.io/v1/visits/apidocs.json

## Config

Please change the env variables in app.yaml 
at https://github.com/SestroAI/envrionments/blob/master/stage.sestro-165123/visits.yaml

## API Specs

```
Base Path: /v1/visits/

1) POST /

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
```
