package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	ports "github.com/awesome-demo-app/todolist-api/core/ports"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/trace"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/slok/go-http-metrics/metrics/prometheus"
	"github.com/slok/go-http-metrics/middleware"
	middlewarestd "github.com/slok/go-http-metrics/middleware/std"
)

type HTTPHandler struct {
	todoService ports.ToDoService
}

func NewHTTPHandler(service ports.ToDoService) *HTTPHandler {
	return &HTTPHandler{
		todoService: service,
	}
}

func (handler *HTTPHandler) HandleRequests() {
	myMiddleware := middleware.New(middleware.Config{
		Recorder: prometheus.NewRecorder(prometheus.Config{}),
	})

	// Base handler register without any fancy stuff:
	// http.HandleFunc("/todo", handler.GetAllToDos)

	// With base Prometheus tooling:
	// http.Handle(
	// 	"/todo",
	// 	promhttp.Instrument(
	// 		prometheus.DefaultRegisterer,
	// 		http.HandlerFunc(handler.ToDoDispatchHandler),
	// 	),
	// )

	// With go-http-metrics as wrapping handler:
	todoHandler := http.HandlerFunc(handler.ToDoDispatchHandler)
	wrappedHandler := middlewarestd.Handler("/todo", myMiddleware, todoHandler)

	// Replacing:
	// http.Handle("/metrics", promhttp.Handler())
	// by a goroutine with a dedicated http server for the Prometheus handler
	// ... weird consequence of go-http-metrics
	go http.ListenAndServe(":9090", promhttp.Handler())

	log.Fatal(http.ListenAndServe(":10000", wrappedHandler))
}

// Dispatch depending on method
func (handler *HTTPHandler) ToDoDispatchHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Handler → Dispatching ToDo endpoint based on HTTP verb")

	tracer := otel.Tracer("github.com/awesome-demo-app/todolist-api-go")
	ctx, span := tracer.Start(r.Context(), "Handling HTTP query")
	defer span.End()

	// Span metadata
	metadataSource, _ := baggage.NewMember("source", "HTTP query")
	metadataHttpMethod, _ := baggage.NewMember("http.method", r.Method)
	metadataHttpPath, _ := baggage.NewMember("http.path", r.URL.Path)
	baggageMetadata, err := baggage.New(
		metadataSource,
		metadataHttpMethod,
		metadataHttpPath,
	)
	if err != nil {
		log.Fatalf("Could not initialize baggage metadata for span")
	}
	ctx = baggage.ContextWithBaggage(ctx, baggageMetadata)
	requestWithEnrichedContext := r.WithContext(ctx)

	switch r.Method {
	case "GET":
		fmt.Println("Handlers → Got a GET request")
		handler.GetAllToDos(w, requestWithEnrichedContext)
	case "POST":
		fmt.Println("Handlers → Got a POST request")
		handler.CreateNewToDo(w, requestWithEnrichedContext)
	case "DELETE":
		fmt.Println("Handlers → Got a DELETE request")
		handler.DeleteToDoByID(w, requestWithEnrichedContext)
	}
}

func (handler *HTTPHandler) GetAllToDos(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Handlers → Endpoint Hit: GetAllToDos")
	span := trace.SpanFromContext(ctx)

	todos, err := handler.todoService.GetAll(r.Context())
	if err != nil {
		fmt.Fprintf(w, "Failed lol")
		span.RecordError(err)
	}

	fmt.Println("Handlers → Rendering ToDos")
	for _, todo := range todos {
		fmt.Fprintf(w, "- %+v\n", todo.Summary)
	}
}

func (handler *HTTPHandler) CreateNewToDo(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Handlers → Endpoint Hit: CreateNewToDo")

	queryParameters := r.URL.Query()
	fmt.Printf("Handlers → Got the following query parameters: %v\n", queryParameters)

	summaryRawData, ok := queryParameters["summary"]
	fmt.Printf("Handlers → Got the following summary data: %v\n", summaryRawData)

	if !ok || len(summaryRawData) != 1 {
		fmt.Println("Could not find 'summary' parameter in POST data")
		return
	}

	summary := summaryRawData[0]
	fmt.Printf("Handlers → Creating new ToDo from summary by calling Service")
	todo := handler.todoService.Create(summary, r.Context())
	fmt.Printf("Handlers → Serializing newly created ToDo in JSON before returning")
	serializedToDo, err := json.Marshal(&todo)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Fprintf(w, "%v\n", string(serializedToDo))
}

func (handler *HTTPHandler) DeleteToDoByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Handlers → Endpoint Hit: DeleteToDoByID")

	queryParameters := r.URL.Query()
	fmt.Printf("Handlers → Got the following query parameters: %v\n", queryParameters)

	idRawData, ok := queryParameters["id"]
	fmt.Printf("Handlers → Got the following id data: %v\n", idRawData)

	if !ok || len(idRawData) != 1 {
		fmt.Println("Could not find 'id' parameter in POST data")
		return
	}

	id, err := strconv.ParseUint(idRawData[0], 0, 64)
	if err != nil {
		fmt.Println(err)
	}
	handler.todoService.Delete(uint(id), r.Context())
}
