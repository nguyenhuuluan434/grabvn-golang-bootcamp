package dao

type db struct {
}

func (db db) FetchMessage(lang string) (string, error) {
	if lang == "en" {
		return "hello", nil
	}
	if lang == "vi" {
		return "xin chào", nil
	}
	return "bzzzz", nil
}

func (db db) FetchDefaultMessage() (string, error) {
	return "default message", nil
}

type DBService interface {
	FetchMessage(lang string) (string, error)
	FetchDefaultMessage() (string, error)
}

func NewDB() DBService{
	return &db{}
}
