vcl 4.1;

import dynamic;

backend default none;

sub vcl_init {
    new dir = dynamic.director();
}

sub vcl_backend_fetch {
    set bereq.backend = dir.backend(bereq.http.host, 80);
}

sub vcl_backend_response {
    set beresp.uncacheable = true;
}