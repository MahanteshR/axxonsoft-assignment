## Dev Notes

This service implements 2 endpoints:
* **POST /task**
  * The POST endpoint to make a call to an external third party service by mentioning in the request body the curl in the form of request attributes.
  * The following are the accepted attributes and its rules:
    * `method` -> Indicates the type of HTTP method being called. Allowed ones are [GET, POST, PATCH, PUT, DELETE].
    * `url` -> The host url of the third party service, for simplicity only the following schemes are allowed [http, https].
    * `headers` -> Headers can also be passed, an example of this:
    ```
    "headers": {
        "Content-Type": "application/json"
    },
    ```
    * `data` -> If the method is PATCH/PUT/POST, data attribute along with content-type header should be passed which indicates the request body and data attribute shouldn't be passes for GET/DELETE methods.

  * **Working**:
    * Whenever the server gets a new task, a taskID(uuid) is created, by default its status is `new` and the task detail is stored in redis cache.
    * If the pre-processing operations to the external service fail, the task's status is updated to `error`, since an error has occurred.
    * The moment a http call to the external third party service is made, the status is updated to `in_process`.
    * After receiving the response successfully, the status code is checked and is updated accordingly. If successful, information from the response is also captured in the cache.


* **GET /task/{{taskID}}**
  * The GET fetches task details from the cache given the taskID in path param.
  * If a taskID does not exist in the cache, it returns an empty object.

The following steps are to be followed to run/test the service locally.
- Repository Setup:
    * Get all the dependencies by using:
        ```
        go get -v -d ./...
        ```
- To setup the DB: redis
    - Initialize docker container
      ```shell
      docker run -d --name my-redis-container -p 6379:6379 redis
      ```


