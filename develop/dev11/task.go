package main

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

/*
=== HTTP server ===

Реализовать HTTP сервер для работы с календарем. В рамках задания необходимо работать строго со стандартной HTTP
библиотекой.
В рамках задания необходимо:
	1. Реализовать вспомогательные функции для сериализации объектов доменной области в JSON.
	2. Реализовать вспомогательные функции для парсинга и валидации параметров методов /create_event и /update_event.
	3. Реализовать HTTP обработчики для каждого из методов API, используя вспомогательные функции и объекты доменной
       области.
	4. Реализовать middleware для логирования запросов
Методы API:
POST /create_event
POST /update_event
POST /delete_event
GET /events_for_day
GET /events_for_week
GET /events_for_month

Параметры передаются в виде www-url-form-encoded (т.е. обычные user_id=3&date=2019-09-09).
В GET методах параметры передаются через queryString, в POST через тело запроса.
В результате каждого запроса должен возвращаться JSON документ содержащий либо {"result": "..."} в случае успешного
выполнения метода, либо {"error": "..."} в случае ошибки бизнес-логики.

В рамках задачи необходимо:
	1. Реализовать все методы.
	2. Бизнес логика НЕ должна зависеть от кода HTTP сервера.
	3. В случае ошибки бизнес-логики сервер должен возвращать HTTP 503. В случае ошибки входных данных
       (невалидный int например) сервер должен возвращать HTTP 400. В случае остальных ошибок сервер должен возвращать
       HTTP 500. Web-сервер должен запускаться на порту указанном в конфиге и выводить в лог каждый обработанный запрос.
	4. Код должен проходить проверки go vet и golint.
*/

// Event - модель тела запроса
type Event struct {
	Id     int      `json:"id"`
	UserId int      `json:"user_id"`
	Name   string   `json:"name"`
	Date   jsonTime `json:"date"`
}

// jsonTime - тип который реализует интерфейс для работы с json
type jsonTime time.Time

func (j *jsonTime) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}
	*j = jsonTime(t)
	return nil
}

func (j jsonTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(j))
}

// EventLocalStorage - хранилище данных о событиях, key - id, value - событие
type EventLocalStorage struct {
	sync.RWMutex

	events map[int]*Event
}

func NewStorage() *EventLocalStorage {
	return &EventLocalStorage{
		events: map[int]*Event{},
	}
}

func (s *EventLocalStorage) Create(event *Event) error {
	s.Lock()
	defer s.Unlock()

	if _, exist := s.events[event.Id]; exist {
		return errors.New("event already exists")
	}

	s.events[event.Id] = event
	return nil
}

func (s *EventLocalStorage) Update(event *Event) error {
	s.Lock()
	defer s.Unlock()

	if _, exist := s.events[event.Id]; !exist {
		return errors.New("event does not exist")
	}

	s.events[event.Id] = event
	return nil
}

func (s *EventLocalStorage) Delete(eventId int) error {
	s.Lock()
	defer s.Unlock()

	if _, exist := s.events[eventId]; !exist {
		return errors.New("event does not exist")
	}

	delete(s.events, eventId)
	return nil
}

func (s *EventLocalStorage) GetForDay(userId int, begin time.Time) (events []*Event) {
	s.RLock()
	defer s.RUnlock()

	end := begin.AddDate(0, 0, 1)

	for _, v := range s.events {
		if v.UserId == userId {
			if time.Time(v.Date).After(begin) && time.Time(v.Date).Before(end) {
				events = append(events, v)
			}
		}
	}

	return events
}

func (s *EventLocalStorage) GetForWeek(userId int, begin time.Time) (events []*Event) {
	s.RLock()
	defer s.RUnlock()

	end := begin.AddDate(0, 0, 7)

	for _, v := range s.events {
		if v.UserId == userId {
			if time.Time(v.Date).After(begin) && time.Time(v.Date).Before(end) {
				events = append(events, v)
			}
		}
	}

	return events
}

func (s *EventLocalStorage) GetForMonth(userId int, begin time.Time) (events []*Event) {
	s.RLock()
	defer s.RUnlock()

	end := begin.AddDate(0, 1, 0)

	for _, v := range s.events {
		if v.UserId == userId {
			if time.Time(v.Date).After(begin) && time.Time(v.Date).Before(end) {
				events = append(events, v)
			}
		}
	}

	return events
}

// eventServer - основная структура сервера
type eventServer struct {
	storage *EventLocalStorage
	server  *http.Server
}

func NewServer(port string) *eventServer {
	return &eventServer{
		storage: NewStorage(),
		server: &http.Server{
			Addr: net.JoinHostPort("localhost", port),
		},
	}
}

func (s *eventServer) Run() error {
	mux := http.NewServeMux()

	mux.HandleFunc("/create_event", s.CreateEventHandler)
	mux.HandleFunc("/update_event", s.UpdateEventHandler)
	mux.HandleFunc("/delete_event", s.DeleteEventHandler)
	mux.HandleFunc("/events_for_day", s.GetEventForDayHandler)
	mux.HandleFunc("/events_for_week", s.GetEventForWeekHandler)
	mux.HandleFunc("/events_for_month", s.GetEventForMonthHandler)

	s.server.Handler = LoggingMiddleware(mux)
	return s.server.ListenAndServe()
}

// вспомогательная функция для парсинга delete запроса
func parseId(body io.ReadCloser) (int, error) {
	data := &struct {
		Id int `json:"id"`
	}{}
	if err := json.NewDecoder(body).Decode(data); err != nil {
		return 0, err
	}

	return data.Id, nil
}

// вспомогательная функция для парсинга create/update запроса
func parseEvent(body io.ReadCloser) (event *Event, err error) {
	event = &Event{}
	if err = json.NewDecoder(body).Decode(event); err != nil {
		return nil, err
	}

	if event.Id < 0 {
		err = errors.New("wrong id")
		return nil, err
	}

	if event.UserId < 0 {
		err = errors.New("wrong user_id")
		return nil, err
	}

	if event.Name == "" {
		err = errors.New("name is required")
		return nil, err
	}

	return event, nil
}

// вспомогательная функция для парсинга параметров get запросов
func parseParams(r *http.Request) (userId int, date time.Time, err error) {
	err = r.ParseForm()
	if err != nil {
		return 0, time.Time{}, err
	}

	if len(r.Form) < 2 || len(r.Form) > 3 {
		err = errors.New("wrong params")
		return 0, time.Time{}, err
	}

	userId, err = strconv.Atoi(r.Form["user_id"][0])
	if err != nil {
		return 0, time.Time{}, err
	}

	loc, _ := time.LoadLocation("Europe/Moscow")
	date, err = time.ParseInLocation("2006-01-02", r.Form["date"][0], loc)
	if err != nil {
		return 0, time.Time{}, err
	}

	return userId, date, nil
}

// resultResponse и errorResponse функции для формирования ответа
func resultResponse(w http.ResponseWriter, event ...*Event) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	data := make(map[string][]*Event)
	data["result"] = event
	_ = json.NewEncoder(w).Encode(data)
}

func errorResponse(w http.ResponseWriter, err string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	data := make(map[string]string)
	data["error"] = err
	_ = json.NewEncoder(w).Encode(data)
}

func (s *eventServer) CreateEventHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		errorResponse(w, "wrong method", http.StatusBadRequest)
		return
	}

	event, err := parseEvent(r.Body)
	if err != nil {
		errorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = s.storage.Create(event)
	if err != nil {
		errorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	resultResponse(w, event)
}

func (s *eventServer) UpdateEventHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		errorResponse(w, "wrong method", http.StatusBadRequest)
		return
	}

	event, err := parseEvent(r.Body)
	if err != nil {
		errorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = s.storage.Update(event)
	if err != nil {
		errorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	resultResponse(w, event)
}

func (s *eventServer) DeleteEventHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		errorResponse(w, "wrong method", http.StatusBadRequest)
		return
	}

	id, err := parseId(r.Body)
	if err != nil {
		errorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = s.storage.Delete(id)
	if err != nil {
		errorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	resultResponse(w, &Event{Id: id})
}

func (s *eventServer) GetEventForDayHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		errorResponse(w, "wrong method", http.StatusBadRequest)
		return
	}

	userId, date, err := parseParams(r)
	if err != nil {
		errorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	events := s.storage.GetForDay(userId, date)

	resultResponse(w, events...)
}

func (s *eventServer) GetEventForWeekHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		errorResponse(w, "wrong method", http.StatusBadRequest)
		return
	}

	userId, date, err := parseParams(r)
	if err != nil {
		errorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	events := s.storage.GetForWeek(userId, date)

	resultResponse(w, events...)
}

func (s *eventServer) GetEventForMonthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		errorResponse(w, "wrong method", http.StatusBadRequest)
		return
	}

	userId, date, err := parseParams(r)
	if err != nil {
		errorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	events := s.storage.GetForMonth(userId, date)

	resultResponse(w, events...)
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, req)
		log.Printf("%s %s %s", req.Method, req.RequestURI, time.Since(start))
	})
}

func main() {
	if len(os.Args) < 2 {
		log.Println("enter port")
		return
	}

	s := NewServer(os.Args[1])
	if err := s.Run(); err != nil {
		log.Fatal(err)
	}
}
