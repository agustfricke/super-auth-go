## Auth kit con Fiber

-   Verificacion de email
-   2FA
-   Social login

## Endpoints and how to use ittt :/ -> :D

### sign in

```bash
curl -X POST -H "Content-Type: application/json" -d '{"name": "agust", "email": "agustfricke@protonmail.com", "password": "agust"}' http://127.0.0.1:8080/signin
```

### verify email to sign in

```bash
curl http://127.0.0.1:8080/verify/eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTYxMzMyNjAsImlhdCI6MTY5NjA0Njg2MCwibmJmIjoxNjk2MDQ2ODYwLCJzdWIiOjJ9.c7DHULSYs3jKtOzFBt2CeGlDaDfVR78jlS3MPO7VKLI/
```

### sign in

```bash
curl -X POST -H "Content-Type: application/json" -d '{"email": "agustfricke@protonmail.com", "password": "agust"}' http://127.0.0.1:8080/signin
```

### generate otp

```bash
curl -X POST -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTYxMzM2NjAsImlhdCI6MTY5NjA0NzI2MCwibmJmIjoxNjk2MDQ3MjYwLCJzdWIiOjJ9.a_s5pkQNeUcWoXNQNdVxtkdzN6Wg0rDDYee10Q55NCc" -H "Content-Type: application/json" http://127.0.0.1:8080/generate
```

#### response

-   set the value of base32 in your authenticator app or generate QR code with the value of otpauth_url

```bash
{"base32":"52HZ62QGM3DUZJ4CVF4XR5SS","otpauth_url":"otpauth://totp/Tech%20con%20Agust:agustfricke@protonmail.com?algorithm=SHA1\u0026digits=6\u0026issuer=Tech%20con%20Agust\u0026period=30\u0026secret=52HZ62QGM3DUZJ4CVF4XR5SS"}%
```

### verify otp

#### put the code of you auth app in the token filed ;)

```bash
curl -X POST -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTYxMzM2NjAsImlhdCI6MTY5NjA0NzI2MCwibmJmIjoxNjk2MDQ3MjYwLCJzdWIiOjJ9.a_s5pkQNeUcWoXNQNdVxtkdzN6Wg0rDDYee10Q55NCc" -H "Content-Type: application/json" -d '{"token": "953446"}' http://127.0.0.1:8080/verify
```

### deactivate otp

```bash
curl -X POST -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTYxMzM2NjAsImlhdCI6MTY5NjA0NzI2MCwibmJmIjoxNjk2MDQ3MjYwLCJzdWIiOjJ9.a_s5pkQNeUcWoXNQNdVxtkdzN6Wg0rDDYee10Q55NCc" -H "Content-Type: application/json" -d '{"token": "953446"}' http://127.0.0.1:8080/disable
```

### Google auth

<a href="http://127.0.0.1:8080/auth/google">go to this link</a>
