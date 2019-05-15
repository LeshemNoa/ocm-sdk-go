/*
Copyright (c) 2019 Red Hat, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// This file contains tests for the methods that request tokens.

package client

import (
	"net/http"
	"time"

	// nolint
	. "github.com/onsi/ginkgo"
	// nolint
	. "github.com/onsi/gomega"
	// nolint
	. "github.com/onsi/gomega/ghttp"
)

var _ = Describe("Tokens", func() {
	// Servers used during the tests:
	var oidServer *Server
	var apiServer *Server

	// Logger used during the testss:
	var logger Logger

	BeforeEach(func() {
		var err error

		// Create the servers:
		oidServer = NewServer()
		apiServer = NewServer()

		// Create the logger:
		logger, err = NewStdLoggerBuilder().
			Streams(GinkgoWriter, GinkgoWriter).
			Debug(true).
			Build()
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		// Stop the servers:
		oidServer.Close()
		apiServer.Close()
	})

	Describe("Refresh grant", func() {
		It("Returns the access token generated by the server", func() {
			// Generate the tokens:
			accessToken := Token("Bearer", 5*time.Minute)
			refreshToken := Token("Refresh", 10*time.Hour)

			// Configure the server:
			oidServer.AppendHandlers(
				CombineHandlers(
					VerifyRefreshGrant(refreshToken),
					RespondWithTokens(accessToken, refreshToken),
				),
			)

			// Create the connection:
			connection, err := NewConnectionBuilder().
				Logger(logger).
				TokenURL(oidServer.URL()).
				URL(apiServer.URL()).
				Tokens(refreshToken).
				Build()
			Expect(err).ToNot(HaveOccurred())
			defer connection.Close()

			// Get the tokens:
			returnedAccess, returnedRefresh, err := connection.Tokens()
			Expect(err).ToNot(HaveOccurred())
			Expect(returnedAccess).To(Equal(accessToken))
			Expect(returnedRefresh).To(Equal(refreshToken))
		})

		It("Sends the token request the first time only", func() {
			// Generate the tokens:
			accessToken := Token("Bearer", 5*time.Minute)
			refreshToken := Token("Refresh", 10*time.Hour)

			// Configure the server:
			oidServer.AppendHandlers(CombineHandlers(
				VerifyRefreshGrant(refreshToken),
				RespondWithTokens(accessToken, refreshToken),
			))

			// Create the connection:
			connection, err := NewConnectionBuilder().
				Logger(logger).
				TokenURL(oidServer.URL()).
				URL(apiServer.URL()).
				Tokens(refreshToken).
				Build()
			Expect(err).ToNot(HaveOccurred())
			defer connection.Close()

			// Get the tokens the first time:
			firstAccess, firstRefresh, err := connection.Tokens()
			Expect(err).ToNot(HaveOccurred())

			// Get the tones the second time:
			secondAccess, secondRefresh, err := connection.Tokens()
			Expect(err).ToNot(HaveOccurred())
			Expect(firstAccess).To(Equal(secondAccess))
			Expect(firstRefresh).To(Equal(secondRefresh))
		})

		It("Refreshes the access token request if it is expired", func() {
			// Generate the tokens:
			expiredAccess := Token("Bearer", -5*time.Minute)
			validAccess := Token("Bearer", -5*time.Minute)
			refreshToken := Token("Refresh", 10*time.Hour)

			// Configure the server:
			oidServer.AppendHandlers(
				CombineHandlers(
					VerifyRefreshGrant(refreshToken),
					RespondWithTokens(validAccess, refreshToken),
				),
			)

			// Create the connection:
			connection, err := NewConnectionBuilder().
				Logger(logger).
				TokenURL(oidServer.URL()).
				URL(apiServer.URL()).
				Tokens(expiredAccess, refreshToken).
				Build()
			Expect(err).ToNot(HaveOccurred())
			defer connection.Close()

			// Get the tokens:
			returnedAccess, _, err := connection.Tokens()
			Expect(err).ToNot(HaveOccurred())
			Expect(returnedAccess).To(Equal(validAccess))
		})

		It("Refreshes the access token if it expires in less than one minute", func() {
			// Generate the tokens:
			firstAccess := Token("Bearer", 50*time.Second)
			secondAccess := Token("Bearer", 5*time.Minute)
			refreshToken := Token("Refresh", 10*time.Hour)

			// Configure the server:
			oidServer.AppendHandlers(
				CombineHandlers(
					VerifyRefreshGrant(refreshToken),
					RespondWithTokens(secondAccess, refreshToken),
				),
			)

			// Create the connection:
			connection, err := NewConnectionBuilder().
				Logger(logger).
				TokenURL(oidServer.URL()).
				URL(apiServer.URL()).
				Tokens(firstAccess, refreshToken).
				Build()
			Expect(err).ToNot(HaveOccurred())
			defer connection.Close()

			// Get the tokens:
			returnedAccess, _, err := connection.Tokens()
			Expect(err).ToNot(HaveOccurred())
			Expect(returnedAccess).To(Equal(secondAccess))
		})

		It("Fails if the access token is expired and there is no refresh token", func() {
			// Generate the tokens:
			accessToken := Token("Bearer", -5*time.Second)

			// Create the connection:
			connection, err := NewConnectionBuilder().
				Logger(logger).
				TokenURL(oidServer.URL()).
				URL(apiServer.URL()).
				Tokens(accessToken).
				Build()
			Expect(err).ToNot(HaveOccurred())
			defer connection.Close()

			// Get the tokens:
			_, _, err = connection.Tokens()
			Expect(err).To(HaveOccurred())
		})

		It("Fails if the refresh token is expired", func() {
			// Generate the tokens:
			refreshToken := Token("Refresh", -5*time.Second)

			// Create the connection:
			connection, err := NewConnectionBuilder().
				Logger(logger).
				TokenURL(oidServer.URL()).
				URL(apiServer.URL()).
				Tokens(refreshToken).
				Build()
			Expect(err).ToNot(HaveOccurred())
			defer connection.Close()

			// Get the tokens:
			_, _, err = connection.Tokens()
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("Password grant", func() {
		It("Returns the access and refresh tokens generated by the server", func() {
			// Generate the tokens:
			accessToken := Token("Bearer", 5*time.Minute)
			refreshToken := Token("Refresh", 10*time.Hour)

			// Configure the server:
			oidServer.AppendHandlers(
				CombineHandlers(
					VerifyPasswordGrant("myuser", "mypassword"),
					RespondWithTokens(accessToken, refreshToken),
				),
			)

			// Create the connection:
			connection, err := NewConnectionBuilder().
				Logger(logger).
				TokenURL(oidServer.URL()).
				URL(apiServer.URL()).
				User("myuser", "mypassword").
				Build()
			Expect(err).ToNot(HaveOccurred())
			defer connection.Close()

			// Get the tokens:
			returnedAccess, returnedRefresh, err := connection.Tokens()
			Expect(err).ToNot(HaveOccurred())
			Expect(returnedAccess).To(Equal(accessToken))
			Expect(returnedRefresh).To(Equal(refreshToken))
		})

		It("Refreshes access token", func() {
			// Generate the tokens:
			expiredAccess := Token("Bearer", -5*time.Second)
			validAccess := Token("Bearer", 5*time.Minute)
			refreshToken := Token("Refresh", 10*time.Hour)

			// Configure the server:
			oidServer.AppendHandlers(
				CombineHandlers(
					VerifyPasswordGrant("myuser", "mypassword"),
					RespondWithTokens(expiredAccess, refreshToken),
				),
				CombineHandlers(
					VerifyRefreshGrant(refreshToken),
					RespondWithTokens(validAccess, refreshToken),
				),
			)

			// Create the connection:
			connection, err := NewConnectionBuilder().
				Logger(logger).
				TokenURL(oidServer.URL()).
				URL(apiServer.URL()).
				User("myuser", "mypassword").
				Build()
			Expect(err).ToNot(HaveOccurred())
			defer connection.Close()

			// Get the tokens the first time:
			firstAccess, _, err := connection.Tokens()
			Expect(err).ToNot(HaveOccurred())
			Expect(firstAccess).To(Equal(expiredAccess))

			// Get the tokens the second time:
			secondAccess, _, err := connection.Tokens()
			Expect(err).ToNot(HaveOccurred())
			Expect(secondAccess).To(Equal(validAccess))
		})

		It("Requests a new refresh token when it expires", func() {
			// Generate the tokens:
			expiredAccess := Token("Bearer", -5*time.Second)
			expiredRefresh := Token("Refresh", -15*time.Second)
			validAccess := Token("Bearer", 5*time.Minute)
			validRefresh := Token("Refresh", 10*time.Hour)

			// Configure the server:
			oidServer.AppendHandlers(
				CombineHandlers(
					VerifyPasswordGrant("myuser", "mypassword"),
					RespondWithTokens(expiredAccess, expiredRefresh),
				),
				CombineHandlers(
					VerifyPasswordGrant("myuser", "mypassword"),
					RespondWithTokens(validAccess, validRefresh),
				),
			)

			// Create the connection:
			connection, err := NewConnectionBuilder().
				Logger(logger).
				TokenURL(oidServer.URL()).
				URL(apiServer.URL()).
				User("myuser", "mypassword").
				Build()
			Expect(err).ToNot(HaveOccurred())
			defer connection.Close()

			// Get the tokens the first time:
			_, firstRefresh, err := connection.Tokens()
			Expect(err).ToNot(HaveOccurred())
			Expect(firstRefresh).To(Equal(expiredRefresh))

			// Get the tokens the second time:
			_, secondRefresh, err := connection.Tokens()
			Expect(err).ToNot(HaveOccurred())
			Expect(secondRefresh).To(Equal(validRefresh))
		})

		It("Requests a new refresh token when expires in less than ten seconds", func() {
			// Generate the tokens:
			expiredAccess := Token("Bearer", -5*time.Second)
			expiredRefresh := Token("Refresh", 5*time.Second)
			validAccess := Token("Bearer", 5*time.Minute)
			validRefresh := Token("Refresh", 10*time.Hour)

			// Configure the server:
			oidServer.AppendHandlers(
				CombineHandlers(
					VerifyPasswordGrant("myuser", "mypassword"),
					RespondWithTokens(expiredAccess, expiredRefresh),
				),
				CombineHandlers(
					VerifyPasswordGrant("myuser", "mypassword"),
					RespondWithTokens(validAccess, validRefresh),
				),
			)

			// Create the connection:
			connection, err := NewConnectionBuilder().
				Logger(logger).
				TokenURL(oidServer.URL()).
				URL(apiServer.URL()).
				User("myuser", "mypassword").
				Build()
			Expect(err).ToNot(HaveOccurred())
			defer connection.Close()

			// Get the tokens the first time:
			_, firstRefresh, err := connection.Tokens()
			Expect(err).ToNot(HaveOccurred())
			Expect(firstRefresh).To(Equal(expiredRefresh))

			// Get the tokens the second time:
			_, secondRefresh, err := connection.Tokens()
			Expect(err).ToNot(HaveOccurred())
			Expect(secondRefresh).To(Equal(validRefresh))
		})

		It("Fails with wrong user name", func() {
			// Configure the server:
			oidServer.AppendHandlers(
				CombineHandlers(
					VerifyPasswordGrant("baduser", "mypassword"),
					RespondWithError("bad_user", "Bad user"),
				),
			)

			// Create the connection:
			connection, err := NewConnectionBuilder().
				Logger(logger).
				TokenURL(oidServer.URL()).
				URL(apiServer.URL()).
				User("baduser", "mypassword").
				Build()
			Expect(err).ToNot(HaveOccurred())
			defer connection.Close()

			// Get the tokens:
			_, _, err = connection.Tokens()
			Expect(err).To(HaveOccurred())
		})

		It("Fails with wrong password", func() {
			// Configure the server:
			oidServer.AppendHandlers(
				CombineHandlers(
					VerifyPasswordGrant("myuser", "badpassword"),
					RespondWithError("bad_password", "Bad password"),
				),
			)

			// Create the connection:
			connection, err := NewConnectionBuilder().
				Logger(logger).
				TokenURL(oidServer.URL()).
				URL(apiServer.URL()).
				User("myuser", "badpassword").
				Build()
			Expect(err).ToNot(HaveOccurred())
			defer connection.Close()

			// Get the tokens:
			_, _, err = connection.Tokens()
			Expect(err).To(HaveOccurred())
		})
	})

	When("Only the access token is provided", func() {
		It("Returns the access token if it hasn't expired", func() {
			// Generate the token:
			accessToken := Token("Bearer", 5*time.Minute)

			// Create the connection:
			connection, err := NewConnectionBuilder().
				Logger(logger).
				TokenURL(oidServer.URL()).
				URL(apiServer.URL()).
				Tokens(accessToken).
				Build()
			Expect(err).ToNot(HaveOccurred())
			defer connection.Close()

			// Get the tokens:
			returnedAccess, returnedRefresh, err := connection.Tokens()
			Expect(err).ToNot(HaveOccurred())
			Expect(returnedAccess).To(Equal(accessToken))
			Expect(returnedRefresh).To(BeEmpty())
		})

		It("Returns an error if the access token has expired", func() {
			// Generate the token:
			accessToken := Token("Bearer", -5*time.Minute)

			// Create the connection:
			connection, err := NewConnectionBuilder().
				Logger(logger).
				TokenURL(oidServer.URL()).
				URL(apiServer.URL()).
				Tokens(accessToken).
				Build()
			Expect(err).ToNot(HaveOccurred())
			defer connection.Close()

			// Get the tokens:
			returnedAccess, returnedRefresh, err := connection.Tokens()
			Expect(err).To(HaveOccurred())
			Expect(returnedAccess).To(BeEmpty())
			Expect(returnedRefresh).To(BeEmpty())
		})
	})

})

func VerifyPasswordGrant(user, password string) http.HandlerFunc {
	return CombineHandlers(
		VerifyRequest(http.MethodPost, "/"),
		VerifyContentType("application/x-www-form-urlencoded"),
		VerifyFormKV("grant_type", "password"),
		VerifyFormKV("client_id", "uhc"),
		VerifyFormKV("username", user),
		VerifyFormKV("password", password),
	)
}

func VerifyRefreshGrant(refreshToken string) http.HandlerFunc {
	return CombineHandlers(
		VerifyRequest(http.MethodPost, "/"),
		VerifyContentType("application/x-www-form-urlencoded"),
		VerifyFormKV("grant_type", "refresh_token"),
		VerifyFormKV("client_id", "uhc"),
		VerifyFormKV("refresh_token", refreshToken),
	)
}

func RespondWithTokens(accessToken, refreshToken string) http.HandlerFunc {
	return RespondWithJSONTemplate(
		http.StatusOK,
		`{
			"access_token": "{{ .AccessToken }}",
			"refresh_token": "{{ .RefreshToken }}"
		}`,
		"AccessToken", accessToken,
		"RefreshToken", refreshToken,
	)
}

func RespondWithError(err, description string) http.HandlerFunc {
	return RespondWithJSONTemplate(
		http.StatusUnauthorized,
		`{
			"error": "{{ .Error }}",
			"error_description": "{{ .Description }}"
		}`,
		"Error", err,
		"Description", description,
	)
}
