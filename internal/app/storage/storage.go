package storage

type Repository interface {
	Add(ID string, URL string)
	Get(ID string) string
}

type Storage struct {
	List map[string]string
}

func (s Storage) Add(ID string, URL string) {
	s.List[ID] = URL
}

func (s Storage) Get(ID string) string {
	return s.List[ID]
}

var instance *Storage = nil

func GetInstance() *Storage {
	if instance == nil {
		instance = &Storage{make(map[string]string)}
	}

	return instance
}
