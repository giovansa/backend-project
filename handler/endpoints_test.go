package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/SawitProRecruitment/UserService/internal"
	"github.com/SawitProRecruitment/UserService/model"
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http"
	"net/http/httptest"
)

var _ = Describe("Handler", func() {
	var (
		e        *echo.Echo
		server   *Server
		ctrl     *gomock.Controller
		mockRepo *repository.MockRepositoryInterface
		recorder *httptest.ResponseRecorder
	)

	BeforeEach(func() {
		e = echo.New()
		ctrl = gomock.NewController(GinkgoT())
		mockRepo = repository.NewMockRepositoryInterface(ctrl)
		server = &Server{
			Cfg: internal.Config{
				App: internal.AppConfig{
					Env: "unit_test",
				},
			},
			Repository: mockRepo,
		}
		recorder = httptest.NewRecorder()
	})

	Context("Create User", func() {
		It("return fail 400 Bad Request - wrong value request", func() {
			userReq := model.RegisterUserReq{
				Phone:    "0821",
				Name:     "John",
				Password: "Test123456!",
			}

			reqBody, _ := json.Marshal(userReq)
			req, err := http.NewRequest("POST", "/user", bytes.NewReader(reqBody))
			Expect(err).NotTo(HaveOccurred())
			req.Header.Set("Content-Type", "application/json")

			c := e.NewContext(req, recorder)
			err = server.Register(c)
			Expect(recorder.Code).Should(Equal(400))
		})

		It("return error 500 - RegisterUser error", func() {
			userReq := model.RegisterUserReq{
				Phone:    "+62821111121",
				Name:     "John",
				Password: "Test123456!",
			}

			reqBody, _ := json.Marshal(userReq)
			req, err := http.NewRequest("POST", "/user", bytes.NewReader(reqBody))
			Expect(err).NotTo(HaveOccurred())
			req.Header.Set("Content-Type", "application/json")

			c := e.NewContext(req, recorder)
			mockRepo.EXPECT().RegisterUser(gomock.Any(), gomock.Any()).
				Return("1", errors.New("err"))

			err = server.Register(c)
			Expect(recorder.Code).Should(Equal(500))
		})

		It("return success 200 Ok", func() {
			userReq := model.RegisterUserReq{
				Phone:    "+62821111121",
				Name:     "John",
				Password: "Test123456!",
			}

			reqBody, _ := json.Marshal(userReq)
			req, err := http.NewRequest("POST", "/user", bytes.NewReader(reqBody))
			Expect(err).NotTo(HaveOccurred())
			req.Header.Set("Content-Type", "application/json")

			c := e.NewContext(req, recorder)
			mockRepo.EXPECT().RegisterUser(gomock.Any(), gomock.Any()).Return("1", nil)
			err = server.Register(c)
			Expect(recorder.Code).Should(Equal(200))

			var responseBody map[string]interface{}
			err = json.Unmarshal(recorder.Body.Bytes(), &responseBody)
			Expect(err).NotTo(HaveOccurred())

			Expect(responseBody).To(HaveKey("data"))
			data := responseBody["data"].(map[string]interface{})
			Expect(data["user_id"]).Should(Equal("1"))
		})
	})

	Context("Login", func() {
		It("return fail 500 Internal Server Error - get user error", func() {
			userReq := model.LoginRequest{
				Phone:    "0821",
				Password: "Test123456!",
			}

			reqBody, _ := json.Marshal(userReq)
			req, err := http.NewRequest("POST", "/login", bytes.NewReader(reqBody))
			Expect(err).NotTo(HaveOccurred())
			req.Header.Set("Content-Type", "application/json")

			mockRepo.EXPECT().GetUserByPhone(gomock.Any(), "0821").Return(repository.User{}, errors.New("err"))

			c := e.NewContext(req, recorder)
			err = server.Login(c)
			Expect(recorder.Code).Should(Equal(500))
		})

		It("return fail 401 Bad Request - wrong value request", func() {
			userReq := model.LoginRequest{
				Phone:    "0821",
				Password: "Test123456!!",
			}

			reqBody, _ := json.Marshal(userReq)
			req, err := http.NewRequest("POST", "/login", bytes.NewReader(reqBody))
			Expect(err).NotTo(HaveOccurred())
			req.Header.Set("Content-Type", "application/json")

			mockRepo.EXPECT().GetUserByPhone(gomock.Any(), "0821").
				Return(repository.User{
					Password: "$2a$10$iaxku3cUornCGSUMi8x7tu6NLeTaaWMjcSpU0T3HFb2IUG4toz1gS",
				}, nil)

			c := e.NewContext(req, recorder)
			err = server.Login(c)
			Expect(recorder.Code).Should(Equal(401))
		})

		It("return fail 500 Internal Server Error - increment login err", func() {
			userReq := model.LoginRequest{
				Phone:    "0821",
				Password: "Test123456!",
			}

			reqBody, _ := json.Marshal(userReq)
			req, err := http.NewRequest("POST", "/login", bytes.NewReader(reqBody))
			Expect(err).NotTo(HaveOccurred())
			req.Header.Set("Content-Type", "application/json")

			mockRepo.EXPECT().GetUserByPhone(gomock.Any(), "0821").
				Return(repository.User{
					Phone:    "0821",
					Password: "$2a$10$iaxku3cUornCGSUMi8x7tu6NLeTaaWMjcSpU0T3HFb2IUG4toz1gS",
				}, nil)
			mockRepo.EXPECT().IncrSuccessLogin(gomock.Any(), "0821").Return(errors.New("err"))
			c := e.NewContext(req, recorder)
			err = server.Login(c)
			Expect(recorder.Code).Should(Equal(500))
		})

		It("return success 200 Ok", func() {
			userReq := model.LoginRequest{
				Phone:    "0821",
				Password: "Test123456!",
			}

			reqBody, _ := json.Marshal(userReq)
			req, err := http.NewRequest("POST", "/login", bytes.NewReader(reqBody))
			Expect(err).NotTo(HaveOccurred())
			req.Header.Set("Content-Type", "application/json")

			mockRepo.EXPECT().GetUserByPhone(gomock.Any(), "0821").
				Return(repository.User{
					Phone:    "0821",
					Password: "$2a$10$iaxku3cUornCGSUMi8x7tu6NLeTaaWMjcSpU0T3HFb2IUG4toz1gS",
				}, nil)
			mockRepo.EXPECT().IncrSuccessLogin(gomock.Any(), "0821").Return(nil)
			c := e.NewContext(req, recorder)
			err = server.Login(c)
			Expect(recorder.Code).Should(Equal(200))

			var responseBody map[string]interface{}
			err = json.Unmarshal(recorder.Body.Bytes(), &responseBody)
			Expect(err).NotTo(HaveOccurred())
			Expect(responseBody).To(HaveKey("token"))
		})
	})

	Context("Get Profile", func() {
		It("return success 200 Ok", func() {
			req, err := http.NewRequest("GET", "/profile", nil)
			Expect(err).NotTo(HaveOccurred())
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTM3OTg5MjAsInBob25lIjoiKzYyODIxMTExMTIxIn0.uqRn1Vzshcew6pxWoA-7uEzIVJ88cHlEOQlybohlYNzdHCH8rV9R6zpILAm3c89-ERBxMObzXiguyLxq-rZj4gg7o0BPZrIBxkUnWmS_QNX2LUcpkx6dy8sP_2PNVifeT6LeAFobxLrZnQw1Y9SkXbopmHcNZelURE-bziSLnj3a0Dw3GMR_sEy7F60xhT-vG0lxaJgvbcNpOJJBiVN3v9Tfv50jnX9zT3S5HadmcPElrZKJ4BTxjbUtmcLPlIOL-IUiVrA5VlJ7y8dTDv-Gely9vJ0gkr3SvDCorxdJr_eSLg0C_UkRhfDGeOKdajT9p3C3cOasdVhopc6_4RxBgw")

			mockRepo.EXPECT().GetUserByPhone(gomock.Any(), "+62821111121").Return(repository.User{
				Phone: "+62821111121",
				Name:  "Test",
			}, nil)

			c := e.NewContext(req, recorder)
			c.Set("claims", &model.Claims{
				Phone: "+62821111121",
			})
			err = server.GetProfile(c)
			Expect(recorder.Code).Should(Equal(200))

			var responseBody map[string]interface{}
			err = json.Unmarshal(recorder.Body.Bytes(), &responseBody)
			Expect(err).NotTo(HaveOccurred())
			Expect(responseBody).To(HaveKey("data"))
		})
	})

	Context("Update User", func() {
		It("return error 400 invalid request", func() {
			req, err := http.NewRequest("PUT", "/user", nil)
			Expect(err).NotTo(HaveOccurred())
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTM3OTg5MjAsInBob25lIjoiKzYyODIxMTExMTIxIn0.uqRn1Vzshcew6pxWoA-7uEzIVJ88cHlEOQlybohlYNzdHCH8rV9R6zpILAm3c89-ERBxMObzXiguyLxq-rZj4gg7o0BPZrIBxkUnWmS_QNX2LUcpkx6dy8sP_2PNVifeT6LeAFobxLrZnQw1Y9SkXbopmHcNZelURE-bziSLnj3a0Dw3GMR_sEy7F60xhT-vG0lxaJgvbcNpOJJBiVN3v9Tfv50jnX9zT3S5HadmcPElrZKJ4BTxjbUtmcLPlIOL-IUiVrA5VlJ7y8dTDv-Gely9vJ0gkr3SvDCorxdJr_eSLg0C_UkRhfDGeOKdajT9p3C3cOasdVhopc6_4RxBgw")

			c := e.NewContext(req, recorder)
			c.Set("claims", &model.Claims{
				Phone: "+62821111121",
			})

			err = server.UpdateUser(c)
			Expect(recorder.Code).Should(Equal(400))
		})

		It("return error 400 - validate request error", func() {
			userReq := model.UpdateUserReq{
				Name: "a",
			}

			reqBody, _ := json.Marshal(userReq)
			req, err := http.NewRequest("PUT", "/user", bytes.NewReader(reqBody))
			Expect(err).NotTo(HaveOccurred())
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTM3OTg5MjAsInBob25lIjoiKzYyODIxMTExMTIxIn0.uqRn1Vzshcew6pxWoA-7uEzIVJ88cHlEOQlybohlYNzdHCH8rV9R6zpILAm3c89-ERBxMObzXiguyLxq-rZj4gg7o0BPZrIBxkUnWmS_QNX2LUcpkx6dy8sP_2PNVifeT6LeAFobxLrZnQw1Y9SkXbopmHcNZelURE-bziSLnj3a0Dw3GMR_sEy7F60xhT-vG0lxaJgvbcNpOJJBiVN3v9Tfv50jnX9zT3S5HadmcPElrZKJ4BTxjbUtmcLPlIOL-IUiVrA5VlJ7y8dTDv-Gely9vJ0gkr3SvDCorxdJr_eSLg0C_UkRhfDGeOKdajT9p3C3cOasdVhopc6_4RxBgw")

			c := e.NewContext(req, recorder)
			c.Set("claims", &model.Claims{
				Phone: "+62821111121",
			})

			err = server.UpdateUser(c)
			Expect(recorder.Code).Should(Equal(400))
		})

		It("return error 500 - update user error", func() {
			userReq := model.UpdateUserReq{
				Phone: "+62821111121",
				Name:  "Test",
			}

			reqBody, _ := json.Marshal(userReq)
			req, err := http.NewRequest("PUT", "/user", bytes.NewReader(reqBody))
			Expect(err).NotTo(HaveOccurred())
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTM3OTg5MjAsInBob25lIjoiKzYyODIxMTExMTIxIn0.uqRn1Vzshcew6pxWoA-7uEzIVJ88cHlEOQlybohlYNzdHCH8rV9R6zpILAm3c89-ERBxMObzXiguyLxq-rZj4gg7o0BPZrIBxkUnWmS_QNX2LUcpkx6dy8sP_2PNVifeT6LeAFobxLrZnQw1Y9SkXbopmHcNZelURE-bziSLnj3a0Dw3GMR_sEy7F60xhT-vG0lxaJgvbcNpOJJBiVN3v9Tfv50jnX9zT3S5HadmcPElrZKJ4BTxjbUtmcLPlIOL-IUiVrA5VlJ7y8dTDv-Gely9vJ0gkr3SvDCorxdJr_eSLg0C_UkRhfDGeOKdajT9p3C3cOasdVhopc6_4RxBgw")

			c := e.NewContext(req, recorder)
			c.Set("claims", &model.Claims{
				Phone: "+62821111121",
			})

			mockRepo.EXPECT().UpdateUser(gomock.Any(), repository.UpdateUser{
				Phone: "+62821111121",
				Name:  "Test",
			}, "+62821111121").Return(errors.New("err"))

			err = server.UpdateUser(c)
			Expect(recorder.Code).Should(Equal(500))

			var responseBody map[string]interface{}
			err = json.Unmarshal(recorder.Body.Bytes(), &responseBody)
			Expect(err).NotTo(HaveOccurred())
			Expect(responseBody).To(HaveKey("error"))
		})

		It("return success 200 Ok", func() {
			userReq := model.UpdateUserReq{
				Phone: "+62821111121",
				Name:  "Test",
			}

			reqBody, _ := json.Marshal(userReq)
			req, err := http.NewRequest("PUT", "/user", bytes.NewReader(reqBody))
			Expect(err).NotTo(HaveOccurred())
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTM3OTg5MjAsInBob25lIjoiKzYyODIxMTExMTIxIn0.uqRn1Vzshcew6pxWoA-7uEzIVJ88cHlEOQlybohlYNzdHCH8rV9R6zpILAm3c89-ERBxMObzXiguyLxq-rZj4gg7o0BPZrIBxkUnWmS_QNX2LUcpkx6dy8sP_2PNVifeT6LeAFobxLrZnQw1Y9SkXbopmHcNZelURE-bziSLnj3a0Dw3GMR_sEy7F60xhT-vG0lxaJgvbcNpOJJBiVN3v9Tfv50jnX9zT3S5HadmcPElrZKJ4BTxjbUtmcLPlIOL-IUiVrA5VlJ7y8dTDv-Gely9vJ0gkr3SvDCorxdJr_eSLg0C_UkRhfDGeOKdajT9p3C3cOasdVhopc6_4RxBgw")

			mockRepo.EXPECT().GetUserByPhone(gomock.Any(), "+62821111121").Return(repository.User{
				Phone: "+62821111121",
				Name:  "Test",
			}, nil)

			c := e.NewContext(req, recorder)
			c.Set("claims", &model.Claims{
				Phone: "+62821111121",
			})

			mockRepo.EXPECT().UpdateUser(gomock.Any(), repository.UpdateUser{
				Phone: "+62821111121",
				Name:  "Test",
			}, "+62821111121").Return(nil)

			err = server.UpdateUser(c)
			Expect(recorder.Code).Should(Equal(200))

			var responseBody map[string]interface{}
			err = json.Unmarshal(recorder.Body.Bytes(), &responseBody)
			Expect(err).NotTo(HaveOccurred())
			Expect(responseBody).To(HaveKey("data"))
		})
	})
})
