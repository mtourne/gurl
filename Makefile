TARGET= gurl

all: $(TARGET)

gurl: gurl.go
	go build gurl.go


clean:
	rm -f $(TARGET)

distclean: clean
	rm -rf .gostuff

