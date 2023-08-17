# go-bank

## main.go 
wyglada tak:
```go
package main

import "fmt"

func main() {
	server := NewApiServer(":3000")
	server.Run()
}
```

W golang jest konwencja, ze kontruktor nazywamy z "New: na poczatku

w main tworzoymy ApiServer i go uruchamiamy nada mas

## api.go
```go
type ApiServer struct {
	listenAddr string
}

func NewApiServer(listenAddr string) *ApiServer {
	return &ApiServer{
		listenAddr: listenAddr,
	}
}
```
ApiServer to nasza prosta struct ma tylko listen Addr
To nasz wrapper na to co goland daje do http


```go
func (s *ApiServer) Run() {
	router := mux.NewRouter()
	router.HandleFunc("/account", makeHTTPHandleFunc(s.handleAccount))
	router.HandleFunc("/account/{id}", makeHTTPHandleFunc(s.handleGetAccount))

	log.Println("JSON API server running on port: ", s.listenAddr)
	http.ListenAndServe(s.listenAddr, router)
}
```
tu juz uzywamy czegos spoza
https://github.com/gorilla/mux - to projekt który jest w Go od Bety, jest to router HTTP
Nazwa pochodzi od *HTTP Request Multiplexer* 

Program rozpoczyna dzialanie w linijce
```go
	http.ListenAndServe(s.listenAddr, router)
```
https://pkg.go.dev/net/http#ListenAndServe - ta funkcje po prostu rozpoczyna Listener http
pierwszy arg to port na ktorym słucha, drugi [Handler](https://pkg.go.dev/net/http#Handler)
Handler to interface który ma func ServeHTTP(ResponseWriter, *Request)
czyli przyjmuje coś co umie pisać odpowiedzi HTTP oraz sam Request.

[ResponseWriter](https://pkg.go.dev/net/http#ResponseWriter) umie tworzyć HTTP Reply (pakiet HTTP)
z określonym nagłowkiem i polem danych.

Dobra, powrot do `ApiServer.Run`
```go
func (s *ApiServer) Run() {
	router := mux.NewRouter()
	router.HandleFunc("/account", makeHTTPHandleFunc(s.handleAccount))
	router.HandleFunc("/account/{id}", makeHTTPHandleFunc(s.handleGetAccount))

	log.Println("JSON API server running on port: ", s.listenAddr)
	http.ListenAndServe(s.listenAddr, router)
}
```

my jako handler dalismy router od gorilla/mux

jego glowna metoda to HandleFunc, ktorej podajemy *route* oraz *handleFunc* na te route.
z [dokumentacji mux](https://github.com/gorilla/mux/blob/main/mux.go) // HandleFunc registers a new route with a matcher for the URL path.

my jako handler podajemy funkcje `makeHTTPHandleFunc`, ktora transformuje nasza funkcje z api na taka ktora matchuje [type handlerFunc](https://pkg.go.dev/net/http#HandlerFunc)


No i to nizej czyli:

```go
func (s *ApiServer) handleAccount(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.handleGetAccount(w, r)
	}
	if r.Method == "POST" {
		return s.handleCreateAccount(w, r)
	}
	if r.Method == "DELETE" {
		return s.handleDeleteAccount(w, r)
	}

	return fmt.Errorf("method not allowed %s", r.Method)
}

func (s *ApiServer) handleGetAccount(w http.ResponseWriter, r *http.Request) error {

	id := mux.Vars(r)["id"]
	// db.get(id)
	idi, _ := strconv.Atoi(id)

	account := NewAccount(idi, "Anthony", "Joshua")
	return WriteJSON(w, http.StatusOK, account)
}

func (s *ApiServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *ApiServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *ApiServer) handleTransfer(w http.ResponseWriter, r *http.Request) error {
	return nil
}
```
to sa juz po prostu funkcje api, ktore handluja konkretne sciezki.
Ich sygnatura matchuje tę:
```go
`type apiFunc func(http.ResponseWriter, *http.Request) error`
```

Mamy funkcje `makeHTTPHandleFunc`, ktora ma dwa cele:
- dopasowac sygnature
- obslugiwac error dane przez apiFunc
```go
// tego kawalka nie rozumiem //juz rozumiem
// return-type tej funkcji to http.HandlerFunc czyli konkretna sygnatura funkcji `type HandlerFunc func(ResponseWriter, *Request)` cos co rozumie router HTTP
// i ta funkcja wywoluje po prostu apiFunc oraz dodatkowo handluje jej ewentualny error
// dlaczego tak? no bo http.HandlerFunc nie ma w sygnaturze mozliwosci zwracania error, a my chcielibysmy je oblsugiwac
// no i to sie dzieje wlasnie w tym miejscu
func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			// handle error here
			WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}
```

Funckje handlujace sciezki. Np. taka:

```go
func (s *ApiServer) handleGetAccount(w http.ResponseWriter, r *http.Request) error {

	id := mux.Vars(r)["id"]
	// db.get(id)
	idi, _ := strconv.Atoi(id)

	account := NewAccount(idi, "Anthony", "Joshua")
	return WriteJSON(w, http.StatusOK, account)
}
```
[ResponseWriter](https://pkg.go.dev/net/http#ResponseWriter) umie tworzyć odpowiedzi
z [Request](https://pkg.go.dev/net/http#Request) mozna wyciagac rzeczy dotyczace requestu.

No i ta func powinna uzyc http.ResponseWriter, zeby stworzyc pakiet HTTP Reply,
my chcemy oddac jsona, wiec zrobilismy sobie func:
```go
func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}
```
ktorej wsadzamy co ma zjosnic i ona tworzy JSONEncoder na podstawie strumienia w czyli http.ReponseWriter i w ten strumień enkoduje bajty ktore są JSONem obiektu danym w encode.
ale to ten stream `w` actually śle odp

