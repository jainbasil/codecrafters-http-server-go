package main

import (
	"bytes"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type HttpRequest struct {
	Method        string
	RequestTarget string
	HttpVersion   string
	Headers       map[string]string
	Body          []byte
}

type HttpResponse struct {
	StatusLine string
	Headers    map[string]string
	Body       []byte
}

func (r *HttpResponse) String() string {
	responseStr := r.StatusLine + "\r\n" // \r\n marks end of status line
	for k, v := range r.Headers {
		responseStr += k + ": " + v + "\r\n"
	}
	responseStr += "\r\n" + string(bytes.NewBuffer(r.Body).String())
	return responseStr
}

func ParseRequest(buf []byte) HttpRequest {
	r := HttpRequest{}

	reqStr := string(bytes.NewBuffer(buf).String())
	requestArray := strings.Split(reqStr, "\r\n")

	requestLine := requestArray[0]
	fields := strings.Fields(requestLine)
	if len(fields) == 3 {
		r.Method = fields[0]
		r.RequestTarget = fields[1]
		r.HttpVersion = fields[2]
	} else {
		fmt.Println("Invalid request line:", requestLine)
	}

	r.Headers = make(map[string]string)
	for i := 1; i < len(requestArray)-2; i++ { // -2 to exclude the last item in the array
		header := strings.Split(requestArray[i], ": ")
		r.Headers[header[0]] = header[1]
	}

	// last item in the array will be the body
	r.Body = []byte(requestArray[len(requestArray)-1])

	return r
}

func HandleRequest(request HttpRequest) HttpResponse {
	switch {
	case strings.HasPrefix(request.RequestTarget, "/echo"):
		r, err := regexp.Compile("/echo/(.*)")
		if err != nil {
			fmt.Println("Error compiling regex")
		}
		matches := r.FindStringSubmatch(request.RequestTarget)

		response := HttpResponse{
			StatusLine: "HTTP/1.1 200 OK",
			Headers: map[string]string{
				"Content-Type":   "text/plain",
				"Content-Length": strconv.Itoa(len(matches[1])),
			},
			Body: []byte(matches[1]),
		}
		return response
	case request.RequestTarget == "/user-agent":
		userAgent := request.Headers["User-Agent"]
		response := HttpResponse{
			StatusLine: "HTTP/1.1 200 OK",
			Headers: map[string]string{
				"Content-Type":   "text/plain",
				"Content-Length": strconv.Itoa(len(userAgent)),
			},
			Body: []byte(userAgent),
		}
		return response
	case request.RequestTarget == "/":
		response := HttpResponse{
			StatusLine: "HTTP/1.1 200 OK",
			Headers:    make(map[string]string),
			Body:       []byte(""),
		}
		return response
	default:
		return HttpResponse{
			StatusLine: "HTTP/1.1 404 Not Found",
		}
	}
}

func SerializeReponse(response HttpResponse) []byte {
	return []byte(response.String())
}
