package main

func main() {
	cfg := Config{
		Addr: ":3000",
	}
	s := GetServer(&cfg)
	s.Serve()
}
