VARNISHSRC = /usr/include/varnish
VMODTOOL = /usr/share/varnish/vmodtool.py

VMOD = simple
VMOD_SO = libvmod_$(VMOD).so

CFLAGS = -Wall -Werror -I/usr/include/varnish -I. -fpic -shared
LDFLAGS = -shared

all: $(VMOD_SO)

# Generate VCC interface files
vcc_if.c vcc_if.h: $(VMOD).vcc
	python3 $(VMODTOOL) $(VMOD).vcc

# Build the Go code as an archive
go_archive.a: main.go vcc_if.h
	go build -buildmode=c-archive -o go_archive.a main.go

# Build the VMOD
$(VMOD_SO): vcc_if.c vcc_if.h go_archive.a
	$(CC) $(CFLAGS) -o $@ vcc_if.c go_archive.a $(LDFLAGS)

clean:
	rm -f $(VMOD_SO) vcc_if.c vcc_if.h go_archive.a vmod_$(VMOD).rst vmod_$(VMOD).man.rst go_archive.h

.PHONY: all clean 