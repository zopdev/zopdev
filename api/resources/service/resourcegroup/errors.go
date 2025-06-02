package resourcegroup

type errInternalServer struct {
}

func (*errInternalServer) Error() string {
	return "internal server error"
}
