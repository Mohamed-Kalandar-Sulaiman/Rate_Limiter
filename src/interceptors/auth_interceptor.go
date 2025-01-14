package interceptors

import (
	"context"
	"crypto/rsa"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var rsaPublicKey *rsa.PublicKey

func init() {
	publicKeyData, err := os.ReadFile("./secure/public_key.pem") 
	if err != nil {
		log.Fatalf("Error loading public key: %v", err)
	}

	rsaPublicKey, err = jwt.ParseRSAPublicKeyFromPEM(publicKeyData)
	if err != nil {
		log.Fatalf("Error parsing RSA public key: %v", err)
	}
}

func AuthInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, fmt.Errorf("missing metadata")
		}

		tokenString := ""
		if val, ok := md["authorization"]; ok && len(val) > 0 {
			tokenString = val[0]
		}

		if tokenString == "" {
			return nil, fmt.Errorf("authorization token is required")
		}

		if len(tokenString) > 7 && strings.HasPrefix(tokenString, "Bearer ") {
			tokenString = tokenString[7:]
		}

		claims := &jwt.StandardClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return rsaPublicKey, nil
		})

		if err != nil {
			log.Println("Error parsing token:", err)
			return nil, fmt.Errorf("invalid token")
		}

		if token.Valid {
			log.Println("Token is valid")
		} else {
			return nil, fmt.Errorf("invalid token")
		}

		return handler(ctx, req)
	}
}
