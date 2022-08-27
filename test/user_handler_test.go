package test

import (
	"cake-store/internal/cakes"
	"cake-store/internal/middlewares"
	mock_repository "cake-store/mocks/repository"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestServices(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Services Suite")
}

var errSomething = errors.New("something error")

var _ = Describe("Test Cake Service", func() {
	var (
		e                *echo.Echo
		mockCtrl         *gomock.Controller
		serviceInterface cakes.SvcInterface
		repo             *mock_repository.MockRepoInterface
		mockData         cakes.Cake
		mockDataList     []cakes.Cake
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockCtrl.Finish()
		repo = mock_repository.NewMockRepoInterface(mockCtrl)
		serviceInterface = cakes.NewHandler(repo)
		e = echo.New()
		middlewares.UseCustomValidatorHandler(e)
		image := "https://img.taste.com.au/ynYrqkOs/w720-h480-cfill-q80/taste/2016/11/sunny-lemon-cheesecake-102220-1.jpeg"
		mockData = cakes.Cake{
			ID:          1,
			Title:       "Lemon cheesecake",
			Description: "A cheesecake made of lemon",
			Rating:      7,
			Image:       &image,
			CreatedAt:   time.Now(),
			UpdatedAt:   nil,
		}
		mockDataList = []cakes.Cake{
			{
				ID:          1,
				Title:       "Lemon cheesecake",
				Description: "A cheesecake made of lemon",
				Rating:      7,
				Image:       &image,
				CreatedAt:   time.Now(),
				UpdatedAt:   nil,
			},
			{
				ID:          2,
				Title:       "Blueberry cheesecake",
				Description: "A cheesecake made of blueberry",
				Rating:      8,
				Image:       nil,
				CreatedAt:   time.Now(),
				UpdatedAt:   nil,
			},
			{
				ID:          3,
				Title:       "Apple cheesecake",
				Description: "A cheesecake made of apple",
				Rating:      8,
				Image:       nil,
				CreatedAt:   time.Now(),
				UpdatedAt:   nil,
			},
		}
	})

	Describe("Create Cake", func() {
		request := `{"title": "Cinnamon Cheesecake", "description": "A cheesecake with a hint of cinnamon", "rating": 7, "image": "https://www.elmundoeats.com/wp-content/uploads/2020/10/FP-Cinnamon-Roll-Cheesecake.jpg"}`

		It("return succeed", func() {
			repo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(request))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			err := serviceInterface.Create(c)
			Expect(err).Should(Succeed())
			Expect(rec.Code).Should(Equal(http.StatusCreated))
		})

		It("return error", func() {
			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(request))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			repo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(errSomething)
			err := serviceInterface.Create(c)
			Expect(err).Should(HaveOccurred())
		})

		It("return error on binding", func() {
			request = `{`
			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(request))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			err := serviceInterface.Create(c)
			Expect(err).Should(HaveOccurred())
		})

		It("return error on validation", func() {
			request = `{}`
			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(request))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			err := serviceInterface.Create(c)
			Expect(err).Should(HaveOccurred())
		})
	})

	Describe("Fetch Cakes", func() {
		It("return succeed", func() {
			repo.EXPECT().List(gomock.Any(), gomock.Any()).Return(mockDataList, int64(len(mockDataList)), nil)
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/cakes")
			err := serviceInterface.List(c)
			Expect(err).Should(Succeed())
			Expect(rec.Code).Should(Equal(http.StatusOK))
		})

		It("return error on validate param", func() {
			req := httptest.NewRequest(http.MethodGet, "/cakes?limit=-1", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			err := serviceInterface.List(c)
			Expect(err).Should(HaveOccurred())
		})

		It("return error", func() {
			repo.EXPECT().List(gomock.Any(), gomock.Any()).Return(nil, int64(0), errSomething)
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/cakes")
			err := serviceInterface.List(c)
			Expect(err).Should(HaveOccurred())
		})
	})

	Describe("Get Cake", func() {
		It("return succeed", func() {
			repo.EXPECT().Get(gomock.Any(), 1).Return(&mockData, nil)
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/cakes/:id")
			c.SetParamNames("id")
			c.SetParamValues("1")
			err := serviceInterface.Get(c)
			Expect(err).Should(Succeed())
			Expect(rec.Code).Should(Equal(http.StatusOK))
		})

		It("return error on invalid id", func() {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/cakes/:id")
			c.SetParamNames("id")
			c.SetParamValues("tes")
			err := serviceInterface.Get(c)
			Expect(err).Should(HaveOccurred())
		})

		It("return data not found", func() {
			repo.EXPECT().Get(gomock.Any(), 1).Return(nil, nil)
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/cakes/:id")
			c.SetParamNames("id")
			c.SetParamValues("1")
			err := serviceInterface.Get(c)
			Expect(err).Should(HaveOccurred())
		})

		It("return error", func() {
			repo.EXPECT().Get(gomock.Any(), 1).Return(&mockData, errSomething)
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/cakes/:id")
			c.SetParamNames("id")
			c.SetParamValues("1")
			err := serviceInterface.Get(c)
			Expect(err).Should(HaveOccurred())
		})
	})

	Describe("Update Cake", func() {
		request := `{
			"title": "Cinnamon Cheesecake",
			"description": "A cheesecake with a hint of cinnamon",
			"rating": 7,
			"image": "https://www.elmundoeats.com/wp-content/uploads/2020/10/FP-Cinnamon-Roll-Cheesecake.jpg"
		}`
		It("return succeed", func() {
			repo.EXPECT().Get(gomock.Any(), 1).Return(&mockData, nil)
			repo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)
			req := httptest.NewRequest(http.MethodPatch, "/", strings.NewReader(request))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/cakes/:id")
			c.SetParamNames("id")
			c.SetParamValues("1")
			err := serviceInterface.Update(c)
			Expect(err).Should(Succeed())
			Expect(rec.Code).Should(Equal(http.StatusOK))
		})

		It("return error", func() {
			repo.EXPECT().Get(gomock.Any(), 1).Return(&mockData, nil)
			repo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(errSomething)
			req := httptest.NewRequest(http.MethodPatch, "/", strings.NewReader(request))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/cakes/:id")
			c.SetParamNames("id")
			c.SetParamValues("1")
			err := serviceInterface.Update(c)
			Expect(err).Should(HaveOccurred())
		})

		It("return error on not found", func() {
			repo.EXPECT().Get(gomock.Any(), 1).Return(nil, nil)
			req := httptest.NewRequest(http.MethodPatch, "/", strings.NewReader(request))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/cakes/:id")
			c.SetParamNames("id")
			c.SetParamValues("1")
			err := serviceInterface.Update(c)
			Expect(err).Should(HaveOccurred())
		})

		It("return error on getting data", func() {
			repo.EXPECT().Get(gomock.Any(), 1).Return(&mockData, errSomething)
			req := httptest.NewRequest(http.MethodPatch, "/", strings.NewReader(request))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/cakes/:id")
			c.SetParamNames("id")
			c.SetParamValues("1")
			err := serviceInterface.Update(c)
			Expect(err).Should(HaveOccurred())
		})

		It("return error on binding body", func() {
			request = `{`
			req := httptest.NewRequest(http.MethodPatch, "/", strings.NewReader(request))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/cakes/:id")
			c.SetParamNames("id")
			c.SetParamValues("1")
			err := serviceInterface.Update(c)
			Expect(err).Should(HaveOccurred())
		})

		It("return error on validation", func() {
			request = `{"image":"plain"}`
			req := httptest.NewRequest(http.MethodPatch, "/", strings.NewReader(request))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/cakes/:id")
			c.SetParamNames("id")
			c.SetParamValues("1")
			err := serviceInterface.Update(c)
			Expect(err).Should(HaveOccurred())
		})
	})

	Describe("Delete Cake", func() {
		It("return succeed", func() {
			repo.EXPECT().Get(gomock.Any(), 1).Return(&mockData, nil)
			repo.EXPECT().Delete(gomock.Any(), 1).Return(nil)
			req := httptest.NewRequest(http.MethodDelete, "/", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/cakes/:id")
			c.SetParamNames("id")
			c.SetParamValues("1")
			err := serviceInterface.Delete(c)
			Expect(err).Should(Succeed())
			Expect(rec.Code).Should(Equal(http.StatusOK))
		})

		It("return error on invalid id", func() {
			req := httptest.NewRequest(http.MethodDelete, "/", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/cakes/:id")
			c.SetParamNames("id")
			c.SetParamValues("dor")
			err := serviceInterface.Delete(c)
			Expect(err).Should(HaveOccurred())
		})

		It("return error on not found", func() {
			repo.EXPECT().Get(gomock.Any(), gomock.Any()).Return(nil, errSomething)
			req := httptest.NewRequest(http.MethodDelete, "/", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/cakes/:id")
			c.SetParamNames("id")
			c.SetParamValues("18")
			err := serviceInterface.Delete(c)
			Expect(err).Should(HaveOccurred())
		})

		It("return internal server error on get data", func() {
			repo.EXPECT().Get(gomock.Any(), 1).Return(&mockData, errSomething)
			req := httptest.NewRequest(http.MethodDelete, "/", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/cakes/:id")
			c.SetParamNames("id")
			c.SetParamValues("1")
			err := serviceInterface.Delete(c)
			Expect(err).Should(HaveOccurred())
		})

		It("return error", func() {
			repo.EXPECT().Get(gomock.Any(), 1).Return(&mockData, nil)
			req := httptest.NewRequest(http.MethodDelete, "/", nil)
			repo.EXPECT().Delete(gomock.Any(), 1).Return(errSomething)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/cakes/:id")
			c.SetParamNames("id")
			c.SetParamValues("1")
			err := serviceInterface.Delete(c)
			Expect(err).Should(HaveOccurred())
		})
	})
})
