# Accounts Service Backend
Accounts service handles authentication, registration, forgotten password, and more.  
You can inherit from base-accounts-service and extend and develop the accounts system you need.  
Recommend this service is not publicly available. Only serve accounts-service-frontend.

# TODO List
- [ ] Enhanced password requirements.
- [ ] System email supports multi languages.
- [x] /auth/v1/signup/confirmation add new param "continue".
- [x] /auth/v1/signup/tokeninfo response add "continue" attribute.
- [x] /auth/v1/forgetpassword/confirmation add new param "continue".
- [x] /auth/v1/forgetpassword/tokeninfo response add "continue" attribute.

## Quick start
```bash
Customize your .env file.
$ cp ./.env.example ./.env
Launch with docker-compose on development environment.
$ docker-compose -f dev.yml up --build -d
```

## Endpoint
### HealthCheck
#### GET /healthcheck/v1/ping
- Params
  - None
- Resonse
  - 200
	```json
	{
  	"message": "OK"
	}
	```

### Auth
#### POST /auth/v1/login
- Params
  - Headers
    - Content-Type : application/json
  - Body
    - identity
      - Required : True
      - Type : String
      - Example : "xxx@mail.com"
    - password
      - Required : True
      - Type : String
      - Example : "IamPassword"
- Response
  - 200
	```json
	{
	  "id": "d655af53-e544-4ae7-a6b9-0f91d19b327a",
	  "identity": "xxx@mail.com",
	  "status": "enabled",
	  "createdAt": "2022-07-19 07:44:29",
	  "updatedAt": "2022-07-19 07:44:29"
	}
	```
  - 400 | 401 | 403 | 404 | 500
	```json
	{
	  "message": "Error Message"
	}
	```

#### POST /auth/v1/signup/confirmation
- Params
  - Headers
    - Content-Type : application/json
  - Body
    - email
      - Required : True
      - Type : String
      - Example : "xxx@mail.com"
    - verifyPageURL
      - Required : True
      - Type : String
      - Example : "https://www.example.com/"
    - continue
      - Required : False
      - Type : String
      - Example : "https://www.continue.com/"
    - subject
      - Required : False
      - Type : String
      - Example : "Signup Email Confirmation"
    - html
      - Required : False
      - Type : String
      - Example : "<!DOCTYPE html> <html> <body> <div class=\"test\"> <p>This is email confirmation, please follow below link to complete sign up flow.</p> <a href={{ .RealVerifyPageURI }}>Click to complete this flow</a> </div> </body> </html>"
			- Note : Please remember compress and encode your html for JSON format. And Keep `{{ VerifyPageURI }}` this variable in your html, so that the service will replace with the actual value.
- Response
  - 202
	```json
	{
	  "message": "Accepted"
	}
	```
  - 400 | 401 | 403 | 409 | 500
	```json
	{
	  "message": "Error Message"
	}
	```

#### GET /auth/v1/signup/tokeninfo
- Params
  - Headers
  - QueryString
    - token
      - Required : True
      - Type : String
      - Example : "JWT"
- Response
  - 200
	```json
	{
	  "email": "xxx@mail.com",
	  "continue": "https://www.continue.com/"
	}
	```
  - 400 | 401 | 403 | 500
	```json
	{
	  "message": "Error Message"
	}
	```

#### POST /auth/v1/signup
- Params
  - Headers
    - Content-Type : application/json
  - Body
    - token
      - Required : True
      - Type : String
      - Example : "JWT"
    - password
      - Required : True
      - Type : String
      - Example : "IamPassword"
- Response
  - 201
	```json
	{
	  "id": "9cfa987b-022d-4461-82c6-f7f12d706163",
	  "identity": "xxx@mail.com",
	  "status": "enabled",
	  "createdAt": "2022-11-01 07:08:34",
	  "updatedAt": "2022-11-01 07:08:34"
	}
	```
  - 400 | 401 | 403 | 409 | 500
	```json
	{
	  "message": "Error Message"
	}
	```

#### POST /auth/v1/forgetpassword/confirmation
- Params
  - Headers
    - Content-Type : application/json
  - Body
    - email
      - Required : True
      - Type : String
      - Example : "xxx@mail.com"
    - verifyPageURL
      - Required : True
      - Type : String
      - Example : "https://www.example.com/"
    - continue
      - Required : False
      - Type : String
      - Example : "https://www.continue.com/"
    - subject
      - Required : False
      - Type : String
      - Example : "Forget Password Email Confirmation"
    - html
      - Required : False
      - Type : String
      - Example : "<!DOCTYPE html> <html> <body> <div class=\"test\"> <p>This is email confirmation, please follow below link to complete forget password flow.</p> <a href={{ .RealVerifyPageURI }}>Click to complete this flow</a> </div> </body> </html>"
			- Note : Please remember compress and encode your html for JSON format. And Keep `{{ VerifyPageURI }}` this variable in your html, so that the service will replace with the actual value.
- Response
  - 202
	```json
	{
	  "message": "Accepted"
	}
	```
  - 400 | 401 | 403 | 409 | 500
	```json
	{
	  "message": "Error Message"
	}
	```

#### GET /auth/v1/forgetpassword/tokeninfo
- Params
  - Headers
  - QueryString
    - token
      - Required : True
      - Type : String
      - Example : "JWT"
- Response
  - 200
	```json
	{
	  "email": "xxx@mail.com",
	  "continue": "https://www.continue.com/"
	}
	```
  - 400 | 401 | 403 | 500
	```json
	{
	  "message": "Error Message"
	}
	```

#### PUT /auth/v1/password
- Params
  - Headers
    - Content-Type : application/json
  - Body
    - token
      - Required : True
      - Type : String
      - Example : "JWT"
    - password
      - Required : True
      - Type : String
      - Example : "IamPassword"
- Response
  - 204
  - 400 | 401 | 403 | 404 | 409 | 500
	```json
	{
	  "message": "Error Message"
	}
	```
