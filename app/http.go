package main

import (
	"bytes"
	"fmt"
	"os"
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

	r.Body = []byte(strings.Trim(requestArray[len(requestArray)-1], "\x00"))
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

		headers := map[string]string{
			"Content-Type":   "text/plain",
			"Content-Length": strconv.Itoa(len(matches[1])),
		}

		if request.Headers["Accept-Encoding"] == "gzip" {
			headers["Content-Encoding"] = "gzip"
		}

		response := HttpResponse{
			StatusLine: "HTTP/1.1 200 OK",
			Headers:    headers,
			Body:       []byte(matches[1]),
		}
		return response
	case strings.HasPrefix(request.RequestTarget, "/files"):
		r, err := regexp.Compile("/files/(.*)")
		if err != nil {
			fmt.Println("Error compiling regex")
		}
		matches := r.FindStringSubmatch(request.RequestTarget)

		filePath := FileDirectory + matches[1]

		if request.Method == "GET" {
			fileContent, err := os.ReadFile(filePath)
			if err != nil {
				fmt.Println("Error reading file:", err)
				return HttpResponse{
					StatusLine: "HTTP/1.1 404 Not Found",
				}
			}
			response := HttpResponse{
				StatusLine: "HTTP/1.1 200 OK",
				Headers: map[string]string{
					"Content-Type":   "application/octet-stream",
					"Content-Length": strconv.Itoa(len(fileContent)),
				},
				Body: []byte(fileContent),
			}
			return response
		} else if request.Method == "POST" {
			file, err := os.Create(filePath)
			if err == nil {
				_, err = file.Write(request.Body)
				file.Close()
			}
			if err != nil {
				fmt.Println("Error writing file:", err)
				return HttpResponse{
					StatusLine: "HTTP/1.1 500 Internal Server Error",
				}
			}
			response := HttpResponse{
				StatusLine: "HTTP/1.1 201 Created",
			}
			return response
		} else {
			return HttpResponse{
				StatusLine: "HTTP/1.1 405 Method Not Allowed",
			}
		}

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
