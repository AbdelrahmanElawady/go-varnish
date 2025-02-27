vcl 4.1;

import etcd;

backend default none;

sub vcl_init {
    new client = etcd.client("localhost:2379");
}

sub vcl_recv {
    if (req.url == "/set") {
        client.set(req.http.key, req.http.value);
        return (synth(10200, "OK"));
    }
    
    if (req.url == "/get") {
        return (synth(10200, client.get(req.http.key)));
    }
}

sub vcl_synth {
    if (resp.status == 10200) {
        synthetic(resp.reason + {"
"} ); // new line to look okay in terminal
        return (deliver);
    }
}