package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"server-monitor/api"

	"github.com/gorilla/mux"
)

// Simple type structure with single data for templates
type Page struct {
	Current string `json:"current_page"`
}

func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	wd, err := os.Getwd()
	if err != nil {
		http.Error(w, "Error getting working directory: "+err.Error(), http.StatusInternalServerError)
		return
	}
	tmplPath := filepath.Join(wd, "web", "templates", tmpl)
	t, err := template.ParseFiles(tmplPath)

	if err != nil {
		http.Error(w, "couldn't load template:  "+err.Error(), http.StatusInternalServerError)
		return
	}
	if err = t.Execute(w, data); err != nil {
		http.Error(w, "Could not render template execution "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func main() {
	api.InitServices()
	api.InitTasks()
	r := mux.NewRouter()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		renderTemplate(w, "dashboard.html", Page{Current: "dashboard"})
	})

	apiRouter := r.PathPrefix("/api").Subrouter()

	apiRouter.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		metrics, err := api.GetMetrics()
		if err != nil {
			log.Printf("Problems getting metrics: %v", err)
			http.Error(w, "Error fetching metrics", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(metrics); err != nil {
			log.Printf("Error encoding metrics JSON: %v", err)
			return
		}
	}).Methods("GET")

	apiRouter.HandleFunc("/services", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			var actionType struct {
				Name   string `json:"name"`
				Action string `json:"action"`
			}

			err := json.NewDecoder(r.Body).Decode(&actionType)
			if err != nil {
				http.Error(w, "Error decoding request body: "+err.Error(), http.StatusBadRequest)
				return
			}

			serviceName := actionType.Name
			var errAction error

			switch actionType.Action {
			case "install":
				errAction = api.Install(serviceName)
			case "start":
				errAction = api.Start(serviceName)
			case "stop":
				errAction = api.Stop(serviceName)
			case "reload":
				errAction = api.Reload(serviceName)
			case "uninstall":
				errAction = api.Uninstall(serviceName)
			default:
				http.Error(w, "Invalid action specified", http.StatusNotImplemented)
				return
			}

			if errAction != nil {
				http.Error(w, fmt.Sprintf("Error performing action '%s' on service '%s': %s", actionType.Action, serviceName, errAction.Error()), http.StatusBadRequest)
				return
			}

			w.WriteHeader(http.StatusAccepted)
			return

		} else if r.Method == http.MethodGet {
			services := api.ListServices()
			w.Header().Set("Content-Type", "application/json")
			time.Sleep(5 * time.Second) // Add a 5-second delay
			w.WriteHeader(http.StatusOK)

			if err := json.NewEncoder(w).Encode(services); err != nil {
				log.Printf("Error encoding services JSON: %v", err)
				return
			}
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}).Methods("POST", "GET")

	apiRouter.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			var task api.Task
			if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
				http.Error(w, "Error decoding task data: "+err.Error(), http.StatusBadRequest)
				return
			}

			newTask, err := api.SubmitTask(task.Description, task.RunTime)
			if err != nil {
				http.Error(w, "Error submitting task: "+err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			if err := json.NewEncoder(w).Encode(newTask); err != nil {
				log.Printf("Error encoding new task JSON: %v", err)
				return
			}

		} else if r.Method == http.MethodGet {
			tasks := api.ListTasks()
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(w).Encode(tasks); err != nil {
				log.Printf("Error encoding tasks JSON: %v", err)
				return
			}

		} else if r.Method == http.MethodDelete {
			err := r.ParseForm()
			if err != nil {
				http.Error(w, "Missing or invalid task ID", http.StatusBadRequest)
				return
			}

			taskIdString := r.Form.Get("id")
			var taskId int
			_, err = fmt.Sscan(taskIdString, &taskId)
			if err != nil {
				http.Error(w, "Invalid task ID format", http.StatusBadRequest)
				return
			}

			err = api.DeleteTask(taskId)
			if err != nil {
				http.Error(w, "Error deleting task", http.StatusNotFound)
				return
			}
			w.WriteHeader(http.StatusNoContent)

		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}).Methods("POST", "GET", "DELETE")

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))

	fmt.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
