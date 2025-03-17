# URL Shortener

This is a simple URL shortening service implemented in Go. It allows users to shorten a long URL and then redirect to the original URL when the shortened URL is accessed.

---


## Usage

### 1. Shorten a URL

To shorten a URL, make a POST request to the /shorten endpoint with the original URL in the request body.

Example with curl:
```bash
curl -X POST -d "http://myurl.com" http://localhost:8080/shorten
```
Response
```bash
http://localhost:8080/12310
```

### 2. Access the Original URL
To access the original URL, make a GET request to the shortened URL.

Example with curl:
```bash
curl -v http://localhost:8080/12310
```
Response:

The server will respond with a 307 Temporary Redirect status code and a Location header pointing to the original URL.

```
< HTTP/1.1 307 Temporary Redirect
< Location: http://myurl.com
< Date: Mon, 17 Mar 2025 19:00:47 GMT
< Content-Length: 0
```

