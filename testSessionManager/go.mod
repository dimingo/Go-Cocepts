module example/testSessionManager

go 1.20

replace example.com/sessions => ../sessions

replace example/sessions => ../sessions

require example/sessions v0.0.0-00010101000000-000000000000

require github.com/gorilla/securecookie v1.1.1 // indirect
