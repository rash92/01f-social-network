package handlefuncs

import (
	"backend/pkg/db/dbfuncs"
	"errors"
	"net/http"
	"testing"
)

type mockResponse struct {
	//add fields that actual response may expect or we need for testing
	statusCode int
	headers	http.Header
}

func (m *mockResponse) Header() http.Header {
	if m.headers == nil {
		m.headers = http.Header{}
	}
	return m.headers
}

func (m *mockResponse) Write(input []byte) (int, error) {
	return 0, nil
}

func (m *mockResponse) WriteHeader(inputStatusCode int) {
	m.statusCode = inputStatusCode
}

// could do individual test cases like this, or one big test 
// that tests individual parts one after another in order expected to encounter errors

func TestHandleAddCommentwithEmptyRequest(t *testing.T) {
	//arrange
	repo = defaultRepository()
	repo.AddComment = func(*dbfuncs.Comment, *dbfuncs.File) (string, error) {
		return "", errors.New("default error")
	}
	//implement what the fake request and response writer needs for specific test
	mockRequest := &http.Request{}
	mockResponseWriter := &mockResponse{}
	expectedError := 405
	errorMessage := "shoud be error 405"

	//act
	HandleAddComment(mockResponseWriter, mockRequest)

	//assert
	if mockResponseWriter.statusCode != expectedError {
		t.Error(errorMessage)
		t.Fatal(errorMessage)
	}

}

// many tests in one function testing one after the other
func TestHandleAddCommentWithBadInputRequests(t *testing.T) {

	// setup mocks - general arrange
	repo = defaultRepository()
	repo.AddComment = func(*dbfuncs.Comment, *dbfuncs.File) (string, error) {
		return "", errors.New("default error")
	}
	mockResponseWriter := &mockResponse{}
	// endof mocks
	
	// scenario: empty request
	// scenario arrange
	rq := &http.Request{}
	// scenario act
	HandleAddComment(mockResponseWriter, rq)
	// scenario assert
	if mockResponseWriter.statusCode != http.StatusMethodNotAllowed {
		t.Error("expected method post")
		t.Fatal("expected method post")
	}


	// scenario empty everything except post
	rq.Method = http.MethodPost
	mockResponseWriter = &mockResponse{}
	HandleAddComment(mockResponseWriter, rq)
	
	if mockResponseWriter.statusCode != http.StatusInternalServerError {
		t.Error("expected it to panic on empty body", mockResponseWriter.statusCode)
		t.Fatal("expected it to panic on empty body", mockResponseWriter.statusCode)
	}

	// rq.Body = actual body and check next test
	// HandleAddComment(&mockResponse{}, rq)
	// if mockResponseWriter.statusCode != StatusNoContent {
	// 	t.Error("expected non empty body")
	// 	t.Fatal("expected non empty body")
	// }
}
