package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"simplebank/token"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func addAuthorization(
	t *testing.T,
	request *http.Request,
	tokenMaker token.Maker,
	authorizationType string,
	username string,
	duration time.Duration,
){
	token,err := tokenMaker.CreateToken(username,duration)
	require.NoError(t,err)

	authorizationHeader := fmt.Sprintf("%s %s",authorizationType,token)
	request.Header.Set(authorizationHeaderKey,authorizationHeader)

}

func TestAuthMiddleware(t *testing.T) {
	testCases := []struct {
		name          string
		setAuth       func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		checkresponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			setAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t,request,tokenMaker,authorizationTypeBearer,"user",time.Minute)
			},
			checkresponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t,http.StatusOK,recorder.Code)
			},
		},
		{
			name: "NoAuthorization",
			setAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				
			},
			checkresponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t,http.StatusUnauthorized,recorder.Code)
			},
		},
		{
			name: "UnsupportedAuthorization",
			setAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t,request,tokenMaker,"unsupported","user",time.Minute)
			},
			checkresponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t,http.StatusUnauthorized,recorder.Code)
			},
		},
		{
			name: "InvalidAuthorizationFormat",
			setAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t,request,tokenMaker,"","user",time.Minute)
			},
			checkresponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t,http.StatusUnauthorized,recorder.Code)
			},
		},
		{
			name: "ExpiredToken",
			setAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t,request,tokenMaker,authorizationTypeBearer,"user",-time.Minute)
			},
			checkresponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t,http.StatusUnauthorized,recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			server := newTestServer(t, nil)

			authPath := "/auth"
			server.router.GET(
				authPath,
				authMiddlware(server.tokenMaker),
				func(ctx *gin.Context) {
					ctx.JSON(http.StatusOK, gin.H{})
				},
			)

			recorder := httptest.NewRecorder()
			request, err := http.NewRequest(http.MethodGet, authPath, nil)
			require.NoError(t, err)

			tc.setAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkresponse(t, recorder)
		})

	}
}
