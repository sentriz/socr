module go.senan.xyz/socr

go 1.16

require (
	github.com/Masterminds/squirrel v1.5.0
	github.com/araddon/dateparse v0.0.0-20210207001429-0eec95c9db7e
	github.com/buckket/go-blurhash v1.1.0
	github.com/cenkalti/dominantcolor v0.0.0-20171020061837-df772e8dd39e
	github.com/cespare/xxhash v1.1.0
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/fsnotify/fsnotify v1.4.9
	github.com/georgysavva/scany v0.2.7
	github.com/gorilla/handlers v1.5.1
	github.com/gorilla/mux v1.8.0
	github.com/gorilla/websocket v1.4.2
	github.com/jackc/pgproto3/v2 v2.0.7 // indirect
	github.com/jackc/pgx/v4 v4.11.0
	github.com/kr/text v0.2.0 // indirect
	github.com/lib/pq v1.9.0 // indirect
	github.com/nfnt/resize v0.0.0-20180221191011-83c6a9932646
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	github.com/otiai10/gosseract/v2 v2.3.1
	go.senan.xyz/socr-frontend v0.0.0
	golang.org/x/sys v0.0.0-20210228012217-479acdf4ea46 // indirect
	golang.org/x/text v0.3.5 // indirect
	gopkg.in/check.v1 v1.0.0-20200227125254-8fa46927fb4f // indirect

)

replace go.senan.xyz/socr-frontend => ../frontend
