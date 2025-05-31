package resourcegroup

type errInternalServer struct {
}

func (e *errInternalServer) Error() string {
	return "internal server error"
}
