package server_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"sync"
	"testing"
	"time"
	"togo/config"
	"togo/server"
	"togo/utils"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const contentTypeJSON = "application/json"

const headerContentType = "Content-Type"

func signUp(t *testing.T, router *gin.Engine, user map[string]interface{}) {
	jsonValue, err := json.Marshal(user)
	require.NoError(t, err)

	req, err := http.NewRequest("POST", SignUpURL, bytes.NewBuffer(jsonValue))
	require.NoError(t, err)
	req.Header.Set(headerContentType, contentTypeJSON)

	signupRes := httptest.NewRecorder()
	router.ServeHTTP(signupRes, req)
}

func login(t *testing.T, router *gin.Engine, user map[string]interface{}) string {
	loginReq, err := json.Marshal(user)
	require.NoError(t, err)

	login, err := http.NewRequest("POST", LoginURL, bytes.NewBuffer(loginReq))
	require.NoError(t, err)
	login.Header.Set(headerContentType, contentTypeJSON)

	loginRes := httptest.NewRecorder()
	router.ServeHTTP(loginRes, login)

	require.Equal(t, http.StatusOK, loginRes.Code)

	loginResBody := loginRes.Body.String()
	var loginResBodyMap = map[string]interface{}{}
	err = json.Unmarshal([]byte(loginResBody), &loginResBodyMap)
	require.NoError(t, err)

	return loginResBodyMap["data"].(map[string]interface{})["token"].(string)
}

func createRandomTask(t *testing.T, router *gin.Engine, token string) {
	randomTask := map[string]interface{}{
		"name":   utils.RandomName(),
		"status": utils.RandomInt(1, 3),
	}

	taskReq, err := json.Marshal(randomTask)
	require.NoError(t, err)

	req, err := http.NewRequest("POST", TaskURL, bytes.NewBuffer(taskReq))
	require.NoError(t, err)
	req.Header.Set(headerContentType, contentTypeJSON)
	req.Header.Set("Authorization", "Bearer "+token)

	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)
	assert.True(t, res.Code == http.StatusOK || res.Code == http.StatusTooManyRequests,
		"Expected status code %d or %d, but got %d", http.StatusOK, http.StatusTooManyRequests, res.Code)
}

func autoCreateTask(t *testing.T, router *gin.Engine, token string, syn *sync.WaitGroup) {
	defer syn.Done()
	createRandomTask(t, router, token)
}

func checkTaskAmount(t *testing.T, router *gin.Engine, token string) {
	req, err := http.NewRequest("GET", TaskURL, nil)
	require.NoError(t, err)
	req.Header.Set(headerContentType, contentTypeJSON)
	req.Header.Set("Authorization", "Bearer "+token)

	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)

	assert.Equal(t, http.StatusOK, res.Code)

	ResBody := res.Body.String()
	var ResBodyMap = map[string]interface{}{}
	err = json.Unmarshal([]byte(ResBody), &ResBodyMap)
	require.NoError(t, err)

	defaultTaskLimitPerDay, err := strconv.Atoi(os.Getenv("DEFAULT_TASK_LIMIT_PER_DAY"))
	if err != nil {
		defaultTaskLimitPerDay = 5
	}
	assert.Equal(t, defaultTaskLimitPerDay,
		len(ResBodyMap["data"].(map[string]interface{})["tasks"].([]interface{})))
}

func TestCreateTask(t *testing.T) {
	cfg := config.LoadConfig("../config.yml")
	router := server.SetupRouter(cfg)

	user := map[string]interface{}{
		"email":    utils.RandomEmail(),
		"password": utils.RandomPassword(),
	}

	signUp(t, router, user)
	token := login(t, router, user)

	var wg sync.WaitGroup
	time.Sleep(time.Second / 2)
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go autoCreateTask(t, router, token, &wg)
	}
	wg.Wait()

	checkTaskAmount(t, router, token)
}
