package application_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/rosenhouse/ifg/application"
	"github.com/rosenhouse/ifg/mocks"
)

var _ = Describe("Grid Handler", func() {
	var handler application.GridHandler
	var response *httptest.ResponseRecorder
	var dataStore *mocks.DataStore
	var keyGenerator *mocks.KeyGenerator

	const baseURL = "http://example.com:1234"

	BeforeEach(func() {
		dataStore = &mocks.DataStore{}
		dataStore.GetCall.Return.Value = []byte("[]")
		keyGenerator = &mocks.KeyGenerator{}

		handler = application.GridHandler{
			GridHTML:     []byte("some grid html"),
			DataStore:    dataStore,
			KeyGenerator: keyGenerator,
		}
		response = httptest.NewRecorder()
	})

	It("should redirect /grid/ root the top root", func() {
		request, _ := http.NewRequest("GET", baseURL+"/grid/", nil)
		handler.ServeHTTP(response, request)
		Expect(response.Code).To(Equal(301))
		Expect(response.HeaderMap.Get("Location")).To(Equal(baseURL + "/"))
	})

	It("should redirect /grid/new to a new random route", func() {
		keyGenerator.NewCall.Return = "some-random-path"
		request, _ := http.NewRequest("GET", baseURL+"/grid/new", nil)
		handler.ServeHTTP(response, request)
		Expect(response.Code).To(Equal(307))
		Expect(response.HeaderMap.Get("Location")).To(HaveSuffix("/grid/some-random-path"))
	})

	Context("when the grid exists", func() {
		It("should respond with the grid.html", func() {
			request, _ := http.NewRequest("GET", baseURL+"/grid/some-grid", nil)
			handler.ServeHTTP(response, request)
			Expect(response.Body.String()).To(ContainSubstring("some grid html"))
		})

		Context("when PUTing data", func() {
			It("should store the data and respond with 204 No Content", func() {
				request, _ := http.NewRequest("PUT", baseURL+"/grid/some-grid/data", strings.NewReader("[]"))
				handler.DataHandler(response, request)
				Expect(response.Code).To(Equal(204))
				Expect(dataStore.SetCall.Args.Key).To(Equal("some-grid"))
				Expect(dataStore.SetCall.Args.Value).To(Equal([]byte("[]")))
			})
		})

		Context("when GETing data", func() {
			It("should return the data for the associated grid", func() {
				dataStore.GetCall.Return.Value = []byte("some data")
				request, _ := http.NewRequest("GET", baseURL+"/grid/some-grid/data", nil)
				handler.DataHandler(response, request)
				Expect(response.Code).To(Equal(200))
				Expect(response.Body.String()).To(Equal("some data"))
				Expect(dataStore.GetCall.Args.Key).To(Equal("some-grid"))
			})
		})
	})

	Context("when the grid does not exist", func() {
		BeforeEach(func() {
			dataStore.GetCall.Return.Error = errors.New("not found")
			dataStore.GetCall.Return.Value = nil
		})

		It("should respond with a 404", func() {
			request, _ := http.NewRequest("GET", baseURL+"/grid/non-existent", nil)
			handler.ServeHTTP(response, request)
			Expect(response.Code).To(Equal(404))
		})
	})
})
