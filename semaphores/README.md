## Swimming pool semaphores

See online at http://35.224.214.137:5050/

### Prerequisites

#### To run the presentation

- Go ≥ 1.11
- [Demoit](https://github.com/dgageot/demoit)

### Installation

```
go get -u github.com/dgageot/demoit
go get -u github.com/Deleplace/samples
cd ~/go/src/github.com/Deleplace/samples/semaphores
demoit
```

Then open your browser at localhost:8888

### To modify the simulation

Install [GopherJS](https://github.com/gopherjs/gopherjs) and build the JavaScript parth of the application.

```
go get -u github.com/gopherjs/gopherjs
gopherjs build -o js/semaphores.js
```
