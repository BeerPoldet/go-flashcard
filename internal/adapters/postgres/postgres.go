package postgres

type Postgres struct{}

func NewPostgres() *Postgres {
	postgres := &Postgres{}
	return postgres
}

func (p *Postgres) ListFlashcards() []string {
	return []string{
		"Go", "Rust", "TypeScript",
	}
}
