makeService:
	go build -o ./out/dvt_service .

makeTools:
	cd tools && go build -ldflags "-X main.apiUrl=http://localhost:3000" -o ../out/dvt_tools .
