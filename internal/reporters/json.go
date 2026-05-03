package reporters

import (
	"encoding/json"
	"io"

	"github.com/olddognewflex/graphql-painkiller/internal/models"
)

func JSON(w io.Writer, reports []models.Report) error {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(reports)
}
