package main

import (
  "log"
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
)

func handleHome(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/home.html")
	if err != nil {
		http.Error(w, "Error loading home page", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func handleMenu(w http.ResponseWriter, r *http.Request) {

	// Here fetch the menu items from the data server
  resp, err := http.Get("http://localhost:4002/data")
  if err != nil {
    log.Fatal(err)
  }
  body, err := io.ReadAll(resp.Body)
  if err != nil {
    log.Fatal(err)
  }
  resp.Body.Close()

  if err != nil {
		http.Error(w, "Error reading menu items", http.StatusInternalServerError)
		return
	}

	var menuItems []MenuItem
	err = json.Unmarshal(body, &menuItems)
	if err != nil {
		http.Error(w, "Error decoding menu items", http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("templates/menu.html")
	if err != nil {
		http.Error(w, "Error loading menu page", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, menuItems)
}

func handleReviewForm(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/review_form.html")
	if err != nil {
		http.Error(w, "Error loading review form", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func handleReviewSubmission(w http.ResponseWriter, r *http.Request) {

	// Here parse the form data recieved from the review form
  r.ParseForm()
  name := r.FormValue("name")
  dish := r.FormValue("dish")
  rating := stringToInt(r.FormValue("rating"))
  comments := r.FormValue("comments")

	// Here create a new review object from the form data
  review := Review{
    Name : name,
    Dish : dish,
    Rating : rating,
    Comments : comments}

	reviewData, err := json.Marshal(review)
	if err != nil {
		http.Error(w, "Error encoding review data", http.StatusInternalServerError)
		return
	}

	// Here post the review data to the data server
  resp, err := http.Post("http://localhost:4002/addReview","application/json", bytes.NewBuffer(reviewData))
	// If there is an error, return an internal server error
  if err != nil {
    log.Fatal(err)
  }
  // Close the response body when the function returns
  resp.Body.Close()
	http.Redirect(w, r, "/reviews", http.StatusSeeOther)
}

func handleReviews(w http.ResponseWriter, r *http.Request) {

  resp, err := http.Get("http://localhost:4002/reviews")
	// Here fetch the reviews from the data server
  if err != nil {
    log.Fatal(err)
  }
	// If there is an error, return an internal server error
  body, err := io.ReadAll(resp.Body)
  if err != nil {
    log.Fatal(err)
  }
	// Close the response body when the function returns
  resp.Body.Close()
	// Here read the response body

  if err != nil {
		http.Error(w, "Error reading menu items", http.StatusInternalServerError)
		return
	}

	var reviews []Review
	err = json.Unmarshal(body, &reviews)
	if err != nil {
		http.Error(w, "Error decoding reviews", http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("templates/reviews.html")
	if err != nil {
		http.Error(w, "Error loading reviews page", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, reviews)
}

func stringToInt(s string) int {
	var i int
	fmt.Sscanf(s, "%d", &i)
	return i
}
