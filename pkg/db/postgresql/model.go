package postgresql

import (
	"strings"

	"github.com/BajoJajoOrg/Inkscryption-backend/pkg/filter"
	sq "github.com/Masterminds/squirrel"
)

func BuildQuery(filterOptions filter.Options, id string) (string, []interface{}, error) {

	// TODO: возможно стоит попробовать сделать билдер более универсальным?
	qb := sq.Select("id", "name", "url", "updated_at", "text", "folder_id").
		From("canvas").
		PlaceholderFormat(sq.Dollar)

	qb = qb.Where(sq.Eq{"user_id": id})

	if nameFilter := filterOptions.GetField("name"); nameFilter != nil {
		qb = qb.Where(sq.ILike{"canvas_name": "%" + nameFilter.Value + "%"})
	}

	if dateFilter := filterOptions.GetField("created_at"); dateFilter != nil {
		switch dateFilter.Operator {
		case filter.OperatorEq:
			qb = qb.Where(sq.Eq{"update_time": dateFilter.Value})
		case filter.OperatorBetween:
			dates := strings.Split(dateFilter.Value, ":")
			qb = qb.Where(sq.Expr("update_time BETWEEN ? AND ?", dates[0], dates[1]))
		}
	}

	return qb.ToSql()
}
