*The code can be found in "todoapp-json" folder.

**[changes]**

[1] Implemented CSP (Communicating Sequential Processes) in todo-app

[2] Added one central go routine and many channels to read and write from the json-file

[3] Added parallel test for storage and logic layer

[4] Changed/Added unit test for storage and logic layer

**[future works]**

[1] Repl (Read-eval-print loop)

[2] multiple start-ups [cli, api and repl] functionality

[3] graceful shutdown

------------------------------------------------------------------------------------------------------------------

**#How to test API (CRUD)**

[Step 1]

Start local server in any terminal powershell or bash and run the following command i.e.

\...> go run main.go

[Step 2]

After the "Server started on :8080" log to appears open another terminal i.e. bash and input the following curl commands

[1] Get - /get

\...> curl http://localhost:8080/get

[2] Create - /create

\...> curl -X POST http://localhost:8080/create \
  -H "Content-Type: application/json" \
  -d '{"description":"Go Shopping","status":"started"}'

[3] Update - /update

\...> curl -X PUT http://localhost:8080/update \
  -H "Content-Type: application/json" \
  -d '{"id":1,"description":"Take the dog for walk","status":"completed"}'

[4] Delete - /delete

\...> curl -X DELETE "http://localhost:8080/delete?id=16"

----------------------------------------------------------------------------------------------------------------------

**#How to run static and dynamic web page**

[1] static page

Open a browser and input the following url.

\...> http://localhost:8080/static/about.html

[2] dynamic web page

Open a browser and input the following url.

\...> http://localhost:8080/list

--------------------------------------------------------------------------------------------------------------------

**#How to run unit test**

[1] Logic layer (Parallel and Unit)

\todoapp-json> go test -v ./logic

[2] Storage layer (Parallel and unit)

\todoapp-json> go test -v ./storage

---------------------------------------------------------------------------------------------------------------------

**#How to use Todo App CLI**

[1] Add item

\...> go run main.go -action=add -desc=" " -status="not started|started|completed"

[2] List items

\...> go run todo-app.go -action=list

[3] Update item

\...> go run main.go -action=update -id=" " -desc="" -status-""

[4] Delete item

\...> go run main.go -action=delete -id=" "

[5] Todo app how to use instructions is displayed.
When anything other than "add|list|update|delete" is input after the -action flag e.g.

\...> go run main.go -action= or go run main.go -action=abcdefg

--------------------------------------------------------------------------------------------------------------------




