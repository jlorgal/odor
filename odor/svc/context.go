//
// Copyright (c) Telefonica I+D. All rights reserved.
//

package svc

// LogContext represents the log context for a base service.
type LogContext struct {
	TransactionID string `json:"trans,omitempty"`
	Correlator    string `json:"corr,omitempty"`
	Operation     string `json:"op,omitempty"`
	Service       string `json:"svc,omitempty"`
	Component     string `json:"comp,omitempty"`
	User          string `json:"user,omitempty"`
	Realm         string `json:"realm,omitempty"`
	Alarm         string `json:"alarm,omitempty"`
}

// ReqLogContext is a complementary LogContext to log information about the request (e.g. path).
type ReqLogContext struct {
	Method     string `json:"method,omitempty"`
	Path       string `json:"path,omitempty"`
	RemoteAddr string `json:"remoteaddr,omitempty"`
}

// RespLogContext is a complementary LogContext to log information about the response (e.g. status code).
type RespLogContext struct {
	Status  int `json:"status,omitempty"`
	Latency int `json:"latency,omitempty"`
}
