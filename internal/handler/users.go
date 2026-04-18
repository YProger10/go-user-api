package handler

import (
    "encoding/json"
    "net/http"
    "strconv"
    "sync"
    "time"
)

type User struct {
    ID        int64     `json:"id"`
    Name      string    `json:"name"`
    Email     string    `json:"email"`
    CreatedAt time.Time `json:"created_at"`
}

var (
    mu    sync.RWMutex
    users = map[int64]User{}
    seq   int64
)

func List(w http.ResponseWriter, r *http.Request) {
    mu.RLock(); defer mu.RUnlock()
    list := make([]User, 0, len(users))
    for _, u := range users { list = append(list, u) }
    json.NewEncoder(w).Encode(list)
}

func Create(w http.ResponseWriter, r *http.Request) {
    var req struct{ Name, Email string }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "bad request", 400); return
    }
    mu.Lock(); seq++
    u := User{ID: seq, Name: req.Name, Email: req.Email, CreatedAt: time.Now()}
    users[seq] = u; mu.Unlock()
    w.WriteHeader(201); json.NewEncoder(w).Encode(u)
}

func GetByID(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
    if err != nil { http.Error(w, "bad id", 400); return }
    mu.RLock(); u, ok := users[id]; mu.RUnlock()
    if !ok { http.Error(w, "not found", 404); return }
    json.NewEncoder(w).Encode(u)
}
