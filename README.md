*The code can be found in "todoapp-json" folder.

*"Using text file" folder contains work from week 1 that has been brought forward and has discontinued

**#How to use Todo App CLI**

[1] Add item

go run todo-app.go -action=add -desc=" " -status="not started|started|completed"

[2] List items

go run todo-app.go -action=list

[3] Update item

go run todo-app.go -action=update -Id=" " -desc="" -status-""

[4] Delete item

go run todo-app.go -action=delete -Id=" "

[5] Start local server

go run todo-app.go -action=serve

[6] Todo app how to use instructions is displayed.
When anything other than "add|list|update|delete|serve" is input after the -action flag e.g.

go run todo-app.go -action= or go run todo-app.go -action=abcdefg

----------------------------------------------------------------------------------------------------------------------

**#How to test API (CRUD)**

[Step 1]

Start local server in any terminal powershell or bash and run the following command i.e.

go run todo-app.go -action=serve

[Step 2]

After the "Server started on :8080" log to appears open another terminal i.e. bash and input the following curl commands

[1] Get - /get

curl http://localhost:8080/get

[2] Create - /create

curl -X POST http://localhost:8080/create \
  -H "Content-Type: application/json" \
  -d '{"description":"Go Shopping","status":"started"}'

[3] Update - /update

curl -X POST http://localhost:8080/update \
  -H "Content-Type: application/json" \
  -d '{"id":1,"description":"Take the dog for walk","status":"completed"}'

[4] Delete - /delete

curl http://localhost:8080/delete?id=1

---------------------------------------------------------------------------------------------------------------------

**#How to run static and dynamic web page**

[1] static page

Open a browser and input the following url.

http://localhost:8080/static/

[2] dynamic web page

Open a browser and input the following url.

http://localhost:8080/list


