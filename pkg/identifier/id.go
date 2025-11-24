package identifier

import gonanoid "github.com/matoous/go-nanoid/v2"

func New() string {
	id, _ := gonanoid.New()
	return id
}
