package feed

type Session interface {
}

type Storage interface {
}

type Server struct {
	session Session
	storage Storage
}

func New(session Session, storage Storage) *Server {
	s := &Server{
		session: session,
		storage: storage,
	}
	return s
}
