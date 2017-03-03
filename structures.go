package main

type Json_decode_error struct {
	Status  string `json:"status"`
	Details string `json:"details"`
}

type Auth_token struct {
	APIVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
	Spec       *Token `json:"spec"`
}

type Token struct {
	Token string `json:"token"`
}

type Authenticated struct {
	Authenticated bool `json:"authenticated"`
}

type AStatus struct {
	Authenticated bool      `json:"authenticated"`
	Userinfo      *Userinfo `json:"user"`
}

type Userinfo struct {
	Username string   `json:"username"`
	UID      string      `json:"uid"`
	Groups   []string `json:"groups"`
}

type Auth_response_successfull struct {
	APIVersion string   `json:"apiVersion"`
	Kind       string   `json:"kind"`
	Status     *AStatus `json:"status"`
}

type Auth_unsuccessfull struct {
	APIVersion string         `json:"apiVersion"`
	Kind       string         `json:"kind"`
	Status     *Authenticated `json:"status"`
}

type SubjectAccessReview struct {
	APIVersion string                   `json:"apiVersion"`
	Kind       string                   `json:"kind"`
	Spec       *SubjectAccessReviewSpec `json:"spec"`
}

type SubjectAccessReviewSpec struct {
	ResourceAttributes    *SubjectAccessReviewSpecResourceAttributes    `json:"resourceAttributes,omitempty"`
	NonResourceAttributes *SubjectAccessReviewSpecNonResourceAttributes `json:"nonresourceAttributes,omitempty"`
	User                  string                                        `json:"user"`
	Group                 []string                                      `json:"group"`
}

type SubjectAccessReviewSpecNonResourceAttributes struct {
	Path string `json:"path,omitempty"`
	Verb string `json:"verb,omitempty"`
}

type SubjectAccessReviewSpecResourceAttributes struct {
	Namespace string `json:"namespace,omitempty"`
	Verb      string `json:"verb,omitempty"`
	Group     string `json:"group,omitempty"`
	Resource  string `json:"resource,omitempty"`
}

type SubjectAccessReviewResponse struct {
	APIVersion string                             `json:"apiVersion"`
	Kind       string                             `json:"kind"`
	Status     *SubjectAccessReviewResponseStatus `json:"status"`
}

type SubjectAccessReviewResponseStatus struct {
	Allowed bool   `json:"allowed"`
	Reason  string `json:"reason,omitempty"`
}
