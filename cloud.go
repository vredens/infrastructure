package infrastructure

import (
	"github.com/vredens/infrastructure/lib/aws"
)

func AWS() aws.AWS {
	return aws.New()
}
