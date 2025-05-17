package postgresql

import (
	"strings"

	"github.com/BajoJajoOrg/Inkscryption-backend/pkg/filter"
	sq "github.com/Masterminds/squirrel"
)

func BuildQuery(filterOptions filter.Options, folder_id int, user_id int) (string, []interface{}, error) {

	// TODO: возможно стоит попробовать сделать билдер более универсальным?
	qb := sq.Select("id", "name", "url", "updated_at", "text", "folder_id", "user_id", "created_at").
		From("canvas").
		PlaceholderFormat(sq.Dollar)

	qb = qb.Where(sq.Eq{
		"folder_id": folder_id,
		"user_id":   user_id,
	})

	if nameFilter := filterOptions.GetField("name"); nameFilter != nil {
		qb = qb.Where(sq.ILike{"name": "%" + nameFilter.Value + "%"})
	}

	if dateFilter := filterOptions.GetField("created_at"); dateFilter != nil {
		switch dateFilter.Operator {
		case filter.OperatorEq:
			qb = qb.Where(sq.Eq{"created_at": dateFilter.Value})
		case filter.OperatorBetween:
			dates := strings.Split(dateFilter.Value, ":")
			qb = qb.Where(sq.Expr("created_at BETWEEN ? AND ?", dates[0], dates[1]))
		}
	}

	return qb.ToSql()
}

func BuildFolderQuery(filterOptions filter.Options, parent_folder_id int, user_id int) (string, []interface{}, error) {
	qb := sq.Select("id", "name", "parent_folder_id", "updated_at", "created_at").
		From("folder").
		PlaceholderFormat(sq.Dollar)

	qb = qb.Where(sq.Eq{
		"parent_folder_id": parent_folder_id,
		"user_id":          user_id,
	})

	if nameFilter := filterOptions.GetField("name"); nameFilter != nil {
		qb = qb.Where(sq.ILike{"name": "%" + nameFilter.Value + "%"})
	}

	if dateFilter := filterOptions.GetField("created_at"); dateFilter != nil {
		switch dateFilter.Operator {
		case filter.OperatorEq:
			qb = qb.Where(sq.Eq{"created_at": dateFilter.Value})
		case filter.OperatorBetween:
			dates := strings.Split(dateFilter.Value, ":")
			qb = qb.Where(sq.Expr("created_at BETWEEN ? AND ?", dates[0], dates[1]))
		}
	}

	return qb.ToSql()
}

// func BuildUniversalQuery(filterOptions filter.Options, identity string, parent_id int, user_id int) (string, []interface{}, error) {
// 	qb := sq.Select("id", "name", "parent_folder_id", "updated_at", "created_at").
// 		From("folder").
// 		PlaceholderFormat(sq.Dollar)

// 	qb = qb.Where(sq.Eq{"parent_folder_id": parent_folder_id}, sq.Eq{"user_id": user_id})

// 	if nameFilter := filterOptions.GetField("name"); nameFilter != nil {
// 		qb = qb.Where(sq.ILike{"name": "%" + nameFilter.Value + "%"})
// 	}

// 	if dateFilter := filterOptions.GetField("created_at"); dateFilter != nil {
// 		switch dateFilter.Operator {
// 		case filter.OperatorEq:
// 			qb = qb.Where(sq.Eq{"created_at": dateFilter.Value})
// 		case filter.OperatorBetween:
// 			dates := strings.Split(dateFilter.Value, ":")
// 			qb = qb.Where(sq.Expr("created_at BETWEEN ? AND ?", dates[0], dates[1]))
// 		}
// 	}

// 	return qb.ToSql()
// }
