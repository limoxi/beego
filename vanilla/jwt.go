package vanilla

import (
	"crypto/hmac"
	"time"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
)

const SALT string = "030e2cf548cf9da683e340371d1a74ee"

func EncodeJWT(data Map) string {
	//header := "{'typ':'JWT','alg':'HS256'}"
	headerBase64Code := "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9"
	
	now := time.Now()
	data["exp"] = Strftime(&now, "%Y-%m-%d %H:%M")
	
	timedelta := Timedelta{Days:365}
	expData := now.Add(timedelta.Duration())
	data["iat"] = Strftime(&expData, "%Y-%m-%d %H:%M")
	
	payload := ToJsonString(data)
	payloadBase64Code := base64.StdEncoding.EncodeToString([]byte(payload))
	
	message := fmt.Sprintf("%s.%s", headerBase64Code, payloadBase64Code)
	h := hmac.New(sha256.New, []byte(SALT))
	h.Write([]byte(message))
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))
	
	return fmt.Sprintf("%s.%s.%s", headerBase64Code, payloadBase64Code, signature)
}

func init() {
}