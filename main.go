package main

import (
	fmt "github.com/jhunt/go-ansi"
	"net/http"
	"os"
	"strconv"

	"github.com/jhunt/go-route"
	"github.com/jhunt/vcaptive"
	"github.com/rs/cors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Item struct {
	gorm.Model `json:"-"`
	ID         uint   `json:"id"`
	Position   int    `json:"position"`
	Text       string `json:"text"`
	Done       bool   `json:"done"`
}

func main() {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" && os.Getenv("VCAP_SERVICES") != "" {
		services, err := vcaptive.ParseServices(os.Getenv("VCAP_SERVICES"))
		if err != nil {
			fmt.Fprintf(os.Stderr, "💥 unable to parse VCAP_SERVICES: %s\n", err)
			os.Exit(1)
		}
		instance, found := services.Tagged("mysql")
		if !found {
			fmt.Fprintf(os.Stderr, "💥 VCAP_SERVICES: no 'mysql' service found\n")
			os.Exit(2)
		}

		hostname, ok := instance.GetString("hostname")
		if !ok {
			hostname, ok = instance.GetString("host")
			if !ok {
				fmt.Fprintf(os.Stderr, "💥 VCAP_SERVICES: '%s' service has no 'hostname' credential\n", instance.Label)
				os.Exit(3)
			}
		}
		port, ok := instance.GetString("port")
		if !ok {
			if n, found := instance.GetUint("port"); found && n > 0 {
				port = fmt.Sprintf("%d", n)
			} else {
				fmt.Fprintf(os.Stderr, "💥 VCAP_SERVICES: '%s' service has no 'port' credential\n", instance.Label)
				os.Exit(3)
			}
		}
		username, ok := instance.GetString("username")
		if !ok {
			fmt.Fprintf(os.Stderr, "💥 VCAP_SERVICES: '%s' service has no 'username' credential\n", instance.Label)
			os.Exit(3)
		}
		password, ok := instance.GetString("password")
		if !ok {
			fmt.Fprintf(os.Stderr, "💥 VCAP_SERVICES: '%s' service has no 'password' credential\n", instance.Label)
			os.Exit(3)
		}
		name, ok := instance.GetString("name")
		if !ok {
			name, ok = instance.GetString("database")
			if !ok {
				fmt.Fprintf(os.Stderr, "💥 VCAP_SERVICES: '%s' service has no 'name' (database name) credential\n", instance.Label)
				os.Exit(3)
			}
		}

		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, hostname, port, name)
	}
	if dsn == "" {
		fmt.Fprintf(os.Stderr, "💥 unable to determine database dsn from environment: No DB_DSN or VCAP_SERVICE environment variables found\n")
		os.Exit(1)
	}
	//dsn = "root:foo@tcp(127.0.0.1:3306)/todos?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		fmt.Fprintf(os.Stderr, "💥 unable to connect to backend database: %s\n", err)
		os.Exit(1)
	}

	db.AutoMigrate(&Item{})

	bind := os.Getenv("BIND")
	if bind == "" {
		port := os.Getenv("PORT")
		if port == "" {
			port = "3003"
		}
		bind = ":" + port
	}

	r := route.Router{
		Name:  "core",
		Debug: os.Getenv("DEBUG") == "yes",
	}

	r.Dispatch("GET /v1/todos", func(req *route.Request) {
		var items []Item
		db.Find(&items)
		req.OK(items)
	})

	r.Dispatch("POST /v1/todos", func(req *route.Request) {
		var in Item
		if !req.Payload(&in) {
			return
		}
		in.ID = 0

		db.Create(&in)
		req.OK(in)
	})

	r.Dispatch("PUT /v1/todos/:id", func(req *route.Request) {
		var in Item
		if !req.Payload(&in) {
			return
		}
		id64, err := strconv.ParseInt(req.Args[1], 10, 32)
		if err != nil {
			req.Fail(route.Oops(err, "unable to parse todo item id as numeric"))
			return
		}
		in.ID = uint(id64)

		db.Save(in)
		req.OK(in)
	})

	r.Dispatch("DELETE /v1/todos/:id", func(req *route.Request) {
		id64, err := strconv.ParseInt(req.Args[1], 10, 32)
		if err != nil {
			req.Fail(route.Oops(err, "unable to parse todo item id as numeric"))
			return
		}
		id := uint(id64)

		db.Delete(&Item{}, id)
		req.Success("item deleted")
	})

	public := os.Getenv("WEBROOT")
	if public == "" {
		public = "ux/dist"
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/v1/", r.ServeHTTP)
	mux.HandleFunc("/", func (w http.ResponseWriter, req *http.Request) {
		if req.Method != "GET" {
			w.WriteHeader(405)
			return
		}

		path := req.URL.Path
		if path == "/" {
			path = "/index.html"
		}
		if f, ok := files[path]; ok {
			w.Header().Set("Content-Type", f.t)
			w.Write(f.b)
			return
		}
		w.WriteHeader(404)
		fmt.Fprintf(w, "oops... 404 not found.\n")
	})

	fmt.Printf("✅ @G{cf-todo} starting @C{up} on @M{%s}\n", bind)
	err = http.ListenAndServe(bind, cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		Debug:          os.Getenv("DEBUG") == "yes",
		AllowedMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
	}).Handler(mux))
	if err != nil {
		fmt.Fprintf(os.Stderr, "💥 failed to spin up http listener: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("✅ @G{cf-todo} shutting @C{down}\n")
	os.Exit(0)
}
