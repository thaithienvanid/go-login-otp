# go-login-otp

## Env

**go1.13.4 windows/amd64**

## Run

**go run main.go**

OR

**go build main.go**

**./main.exe**

## Test

**go clean -testcache**

**go test ./... -v -cover**

OR

**go clean -testcache**

**go test -v -cover ./... > result.test.txt**

## Scope

**Should cover feature from 1 to 8**

## Usage

**POST** /user/login/phone

#### Request

    Head:
      - Content-Type: application/json
    Body:
      - phone: string

#### Response

    Head:
      - Content-Type: application/json
    Body:
      - message: string (optional)
      - payload: any (optional)

#### Note

    It will be limit 30s / request / ip

**POST** /user/login/phone/callback

#### Request

    Head:
      - Content-Type: application/json
    Body:
      - phone: string
      - token: string

#### Response

    Head:
      - Content-Type: application/json
    Body:
      - message: string (optional)
      - payload: any (optional)
