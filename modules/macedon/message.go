package macedon

// MacedonRequest is Macedon request
type MacedonRequest struct {
    Name    string
    Address string
    Ttl     int
}

// MacedonResponse is Macedon response
type MacedonResponse struct {
    Result []MacedonRequest
}


// RecValue is Etcd request
type RecValue struct {
    Host string
    Ttl  int
}

// SubNode is Etcd sub node
type SubNode struct {
    Key   string
    Value string
}

// Node is Etcd node
type Node struct {
    Key   string
    Value string
    Nodes []SubNode
}

// EtcdResponse is Etcd response
type EtcdResponse struct {
    Node Node
}

// ServerRequest is Server request
type ServerRequest struct {
	Address string
}

// ServerResponse is Server response
type ServerResponse struct {
	Result []ServerRequest
}
