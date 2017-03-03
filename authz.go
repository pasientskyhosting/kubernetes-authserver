package main

import (
	"encoding/json"
	"net/http"
)

func AuthzFailed(w http.ResponseWriter, r *http.Request, reason string) {
	response := &SubjectAccessReviewResponse{
		APIVersion: APIVERSION,
		Kind:       "SubjectAccessReview",
		Status: &SubjectAccessReviewResponseStatus{
			Allowed: false,
			Reason:  reason,
		},
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		panic(err)
	}
}

func AuthzOK(w http.ResponseWriter, r *http.Request) {
	response := &SubjectAccessReviewResponse{
		APIVersion: APIVERSION,
		Kind:       "SubjectAccessReview",
		Status: &SubjectAccessReviewResponseStatus{
			Allowed: true,
		},
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		panic(err)
	}
}
