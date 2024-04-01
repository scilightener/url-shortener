A simple go service for shortening urls

Some current features:
- postgresql as a database
- fully logged with slog
- http-framework-free
- fully dockerized
- configured github actions with deploy on a remote server (check it out: http://5.42.100.122:8083/hey)

API:
- POST /api/url
  - request:
    ```json
    {
      "alias": "hello",
      "url": "https://google.com"
    }
    ```
  - response:
    ```json
    {
      "status": "Ok",
      "alias": "hello",
      "valid_until_utc": "2024-05-01T01:33:50.182094575Z"
    }
    ```
- GET /{alias}

Base address: http://5.42.100.122:8083

Well, this is obviously not much, but it's only a start and I have some bigger plans for this project 

Planning:
- add a frontend
- add more tests
- add redis for caching
- export logs to grafana/kibana
- add authorization (authorized users will be able to choose their own alias for short url, while not authorized users will get a random one)
