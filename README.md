# Minerva Backend
> Use WebSockets for managing programs on robot computers

## Requirements

*   Golang (1.19.3)

## Base concepts

Minerva is designed to manage various applications running on robots with one or more on-board computers.

On each of the on-board computers of the robot, you must run its own instance of the program. Before that, you need to select `master computer` on wich there will be the main instance. Both the client and other instances will communicate with master.

## Running one instance

1.  Clone current repository
2.  Move to root of project
3.  Install dependencies

    ```bash
    go get .
    ```

4. Configure service

    Create cofiguration for master instance (with `master` setted to `true`):

    ```yml
    service:
      host: '127.0.0.1'
      port: '8080'
      master: true
      slaves:
        - http://127.0.0.1:8080
        - http://127.0.0.1:8082

    actions:
      - id: '8080_1'
        title: 'First'
        description: 'Test back'
        command: ["echo", "Hello"]
      - id: '8080_2'
        title: 'Second'
        description: 'Test back'
        command: ["sh", "./logt_task.sh"]
    ```

    Or you can configure instance as a slave:
    ```yml
    service:
      host: '127.0.0.1'
      port: '8080'

    actions:
      - id: '8080_1'
        title: 'First'
        description: 'Test back'
        command: ["echo", "Hello"]
      - id: '8080_2'
        title: 'Second'
        description: 'Test back'
        command: ["sh", "./logt_task.sh"]
    ```

5. Run service

    ```bash
    go run .
    ```
