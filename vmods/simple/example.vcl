vcl 4.1;

import simple;

backend default {
    .host = "127.0.0.1";
    .port = "8080";
}


sub vcl_recv {
    simple.write(req.http.host);
}

sub vcl_deliver {
    set resp.http.X-Greeting = simple.hello("World");
    
    set resp.http.X-Sum = simple.add(5, 3);
}
