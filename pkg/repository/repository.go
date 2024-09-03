package repository

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"gitlab.com/golib4/logger/logger"
	"reflect"
)

type Params struct {
	LogInfo bool
}

type BaseRepository struct {
	connection *sql.DB
	logger     logger.Logger
	params     Params
}

func NewBaseRepository(connection *sql.DB, logger logger.Logger, params Params) BaseRepository {
	return BaseRepository{
		connection: connection,
		logger:     logger,
		params:     params,
	}
}

type Model interface {
	GetTableName() string
}

type ClearableModel interface {
	GetTableName() string
	Clear() Model
}

type FindManyByOptions struct {
	Expression  exp.Expression
	Expressions []exp.Expression
	Limit       uint
	Offset      uint
	OrderBy     exp.OrderedExpression
}

func FindManyBy[T ClearableModel](r *BaseRepository, options FindManyByOptions, model T, tx driver.Tx) ([]T, error) {
	modelMetaData, err := r.getModelMetaData(model)
	if err != nil {
		return nil, fmt.Errorf("can not get model metadata: %s", err)
	}
	baseQuery := r.createBaseSelectQuery(model, *modelMetaData)

	if options.Expression != nil {
		baseQuery = baseQuery.Where(options.Expression)
	}
	if options.Expressions != nil {
		baseQuery = baseQuery.Where(options.Expressions...)
	}

	baseQuery = baseQuery.Limit(options.Limit).
		Offset(options.Offset)

	if options.OrderBy != nil {
		baseQuery = baseQuery.Order(options.OrderBy)
	}

	query, _, err := baseQuery.ToSQL()

	if err != nil {
		return nil, fmt.Errorf("can not generate query in find many by method: %s", err)
	}
	rows, err := r.QueryRows(tx, query)
	if err != nil {
		return nil, fmt.Errorf("can not execute find many by query: %s", err)
	}
	if rows.Err() != nil {
		return nil, fmt.Errorf("can not execute query result: %s", rows.Err())
	}
	models, err := scanRowsFields(r, model, *modelMetaData, rows)
	if err != nil {
		return nil, fmt.Errorf("can not scan in find many by method: %s", err)
	}

	return models, nil
}

func (r BaseRepository) FindOneBy(options goqu.Ex, model Model, tx driver.Tx) error {
	modelMetaData, err := r.getModelMetaData(model)
	if err != nil {
		return fmt.Errorf("can not get model metadata: %s", err)
	}
	query, _, err := r.createBaseSelectQuery(model, *modelMetaData).Where(options).ToSQL()
	if err != nil {
		return fmt.Errorf("can not generate query in find one by method: %s", err)
	}
	row, err := r.QueryRow(tx, query)
	if err != nil {
		return fmt.Errorf("can not execute find one by query: %s", err)
	}
	if row.Err() != nil {
		return fmt.Errorf("can not execute query result: %s", row.Err())
	}
	err = r.scanRowFields(model, *modelMetaData, row)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("can not scan query result: %s", err)
	}

	return nil
}

func (r BaseRepository) DeleteBy(options goqu.Ex, model Model, tx driver.Tx) error {
	query, _, err := r.createBaseDeleteQuery(model).Where(options).ToSQL()
	if err != nil {
		return fmt.Errorf("can not generate query in delete by by method: %s", err)
	}
	_, err = r.Exec(tx, query)
	if err != nil {
		return fmt.Errorf("can not execute delete by by query: %s", err)
	}

	return nil
}

func (r BaseRepository) FindOne(idValue any, model Model, tx driver.Tx) error {
	modelMetaData, err := r.getModelMetaData(model)
	if err != nil {
		return fmt.Errorf("can not get model metadata: %s", err)
	}
	query, _, err := r.createBaseSelectQuery(model, *modelMetaData).Where(goqu.Ex{
		modelMetaData.id.name: idValue,
	}).ToSQL()
	if err != nil {
		return fmt.Errorf("can not generate query in find one method: %s", err)
	}
	row, err := r.QueryRow(tx, query)
	if err != nil {
		return fmt.Errorf("can not execute find one query: %s", err)
	}
	if row.Err() != nil {
		return fmt.Errorf("can not execute query result: %s", row.Err())
	}
	err = r.scanRowFields(model, *modelMetaData, row)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("can not scan query result: %s", err)
	}

	return nil
}

func (r BaseRepository) Refresh(model Model, tx driver.Tx) error {
	modelMetaData, err := r.getModelMetaData(model)
	if err != nil {
		return fmt.Errorf("can not get model metadata: %s", err)
	}
	if modelMetaData.isNew {
		return fmt.Errorf("trying to refresh not persisted model %s", model.GetTableName())
	}
	query, _, err := r.createBaseSelectQuery(model, *modelMetaData).Where(goqu.Ex{
		modelMetaData.id.name: modelMetaData.id.value,
	}).ToSQL()
	if err != nil {
		return fmt.Errorf("can not generate query in refresh method: %s", err)
	}

	row, err := r.QueryRow(tx, query)
	if err != nil {
		return fmt.Errorf("can not execute refresh query: %s", err)
	}
	if row.Err() != nil {
		return fmt.Errorf("can not execute query result: %s", row.Err())
	}
	err = r.scanRowFields(model, *modelMetaData, row)
	if err != nil {
		return fmt.Errorf("can not scan query result: %s", err)
	}

	return nil
}

func (r BaseRepository) CreateOrUpdate(model Model, tx driver.Tx) error {
	modelMetaData, err := r.getModelMetaData(model)
	if err != nil {
		return fmt.Errorf("can not get model metadata: %s", err)
	}
	record := make(map[string]interface{})
	for _, fieldData := range modelMetaData.fields {
		record[fieldData.name] = fieldData.value
	}
	if modelMetaData.isNew {
		err = r.execInsertQuery(model, record, *modelMetaData, tx)
	} else {
		err = r.execUpdateQuery(model, record, *modelMetaData, tx)
	}
	if err != nil {
		return fmt.Errorf("can not execute query: %s", err)
	}

	return nil
}

func (r BaseRepository) execInsertQuery(model Model, record goqu.Record, modelMetaData metaData, tx driver.Tx) error {
	query, _, err := goqu.Insert(model.GetTableName()).Rows(record).Returning(modelMetaData.id.name).ToSQL()
	if err != nil {
		return fmt.Errorf("can not generate insert query: %s", err)
	}

	row, err := r.QueryRow(tx, query)
	if err != nil {
		return fmt.Errorf("can not execute insert query: %s", err)
	}
	if row.Err() != nil {
		return fmt.Errorf("can not execute insert query: %s", row.Err())
	}
	err = r.scanRowId(model, modelMetaData, row)
	if err != nil {
		return fmt.Errorf("can not scan insert query: %s", err)
	}

	return nil
}

func (r BaseRepository) execUpdateQuery(model Model, record goqu.Record, modelMetaData metaData, tx driver.Tx) error {
	query, _, err := goqu.Update(model.GetTableName()).Set(record).Where(goqu.Ex{modelMetaData.id.name: modelMetaData.id.value}).ToSQL()
	if err != nil {
		return fmt.Errorf("can not generate insert query: %s", err)
	}

	_, err = r.Exec(tx, query)
	if err != nil {
		return fmt.Errorf("can not execute update query: %s", err)
	}

	return nil
}

func (r BaseRepository) createBaseSelectQuery(model Model, modelMetaData metaData) *goqu.SelectDataset {
	var fields []interface{}
	fields = append(fields, modelMetaData.id.name)
	for _, fieldData := range modelMetaData.fields {
		fields = append(fields, fieldData.name)
	}

	return goqu.Select(fields...).From(model.GetTableName())
}

func (r BaseRepository) createBaseDeleteQuery(model Model) *goqu.DeleteDataset {
	return goqu.Delete(model.GetTableName())
}

func (r BaseRepository) scanRowId(model Model, modelMetaData metaData, row *sql.Row) error {
	valType := reflect.TypeOf(modelMetaData.id.value)
	newValue := reflect.New(valType).Interface()

	err := row.Scan(&newValue)
	if err != nil {
		return fmt.Errorf("can not scan id: %s", row.Err())
	}

	err = r.setFieldValue(model, modelMetaData.id.index, newValue)
	if err != nil {
		return fmt.Errorf("can not set id after insert query id: %s", err)
	}

	return nil
}

func scanRowsFields[T ClearableModel](r *BaseRepository, model T, modelMetaData metaData, rows *sql.Rows) ([]T, error) {
	defer rows.Close()

	values := make([]interface{}, len(modelMetaData.fields)+1)
	valType := reflect.TypeOf(modelMetaData.id.value)
	newValue := reflect.New(valType).Interface()
	values[0] = newValue
	for i, fieldData := range modelMetaData.fields {
		j := i + 1
		valType := reflect.TypeOf(fieldData.value)
		newValue := reflect.New(valType).Interface()
		values[j] = newValue
	}
	var models []T
	for rows.Next() {
		newModel := model.Clear().(T)
		err := rows.Scan(values...)
		if err != nil {
			return nil, err
		}

		for i, value := range values {
			err := r.setFieldValue(newModel, i, value)
			if err != nil {
				return nil, fmt.Errorf("can not set value while scanning: %s", err)
			}
		}

		models = append(models, newModel)
	}

	return models, nil
}

func (r BaseRepository) scanRowFields(model Model, modelMetaData metaData, row *sql.Row) error {
	values := make([]interface{}, len(modelMetaData.fields)+1)
	valType := reflect.TypeOf(modelMetaData.id.value)
	newValue := reflect.New(valType).Interface()
	values[0] = newValue
	for i, fieldData := range modelMetaData.fields {
		j := i + 1
		valType := reflect.TypeOf(fieldData.value)
		newValue := reflect.New(valType).Interface()
		values[j] = newValue
	}

	err := row.Scan(values...)
	if err != nil {
		return err
	}
	for i, value := range values {
		err := r.setFieldValue(model, i, value)
		if err != nil {
			return fmt.Errorf("can not set value while scanning: %s", err)
		}
	}

	return nil
}

type field struct {
	name  string
	value any
}

type idField struct {
	value any
	name  string
	index int
}

type metaData struct {
	isNew  bool
	id     idField
	fields []field
}

func (r BaseRepository) getModelMetaData(model Model) (*metaData, error) {
	originalValue, err := r.getModelValue(model)
	if err != nil {
		return nil, fmt.Errorf("can not get model value: %s", err)
	}
	modelType := originalValue.Type()

	var id idField
	var fields []field
	var isNew bool
	for i := 0; i < originalValue.NumField(); i++ {
		fieldData := modelType.Field(i)

		fieldName := fieldData.Tag.Get("db")
		if fieldName == "-" || fieldName == "" {
			continue
		}

		if !originalValue.Field(i).CanInterface() {
			return nil, fmt.Errorf("field %s of model %s has unsupportable type or field is private", fieldName, model.GetTableName())
		}

		fieldValue := originalValue.Field(i).Interface()

		primaryFieldName := fieldData.Tag.Get("primary")
		if primaryFieldName == "true" {
			fieldValueAsDigit := fmt.Sprintf("%d", fieldValue)
			isNew = fieldValue == "" || fieldValueAsDigit == "0"
			id.value = fieldValue
			id.name = fieldName
			id.index = i
		}

		mustGenerate := fieldData.Tag.Get("must_generate")
		if mustGenerate == "true" {
			continue
		}

		fields = append(fields, field{
			name:  fieldName,
			value: fieldValue,
		})
	}

	return &metaData{
		isNew:  isNew,
		id:     id,
		fields: fields,
	}, nil
}

func (BaseRepository) getModelValue(model Model) (*reflect.Value, error) {
	value := reflect.ValueOf(model)
	if value.Kind() != reflect.Ptr || value.Elem().Kind() != reflect.Struct {
		return nil, errors.New("model must be a pointer to a struct")
	}

	result := reflect.Indirect(value)
	return &result, nil
}

func (r BaseRepository) setFieldValue(model Model, index int, value any) error {
	originalValue, err := r.getModelValue(model)
	if err != nil {
		return fmt.Errorf("can not get model value: %s", err)
	}

	field := originalValue.Field(index)
	if field.IsValid() && field.CanSet() {
		val := reflect.ValueOf(value)
		if val.Kind() != reflect.Ptr {
			field.Set(reflect.ValueOf(value))
		}

		field.Set(reflect.Indirect(val))
		return nil
	}

	return fmt.Errorf("cannot set field value: unknown error")
}

func (r BaseRepository) QueryRow(tx driver.Tx, query string, args ...any) (*sql.Row, error) {
	r.logQuery(query)

	if tx == nil {
		return r.connection.QueryRow(query, args...), nil
	}

	sqlTx, isTypeOfSql := tx.(*sql.Tx)
	if !isTypeOfSql {
		return nil, fmt.Errorf("provided type (%T) of driver.Tx is not supported yet", tx)
	}

	return sqlTx.QueryRow(query, args...), nil
}

func (r BaseRepository) QueryRows(tx driver.Tx, query string, args ...any) (*sql.Rows, error) {
	r.logQuery(query)

	if tx == nil {
		return r.connection.Query(query, args...)
	}

	sqlTx, isTypeOfSql := tx.(*sql.Tx)
	if !isTypeOfSql {
		return nil, fmt.Errorf("provided type (%T) of driver.Tx is not supported yet", tx)
	}

	return sqlTx.Query(query, args...)
}

func (r BaseRepository) Exec(tx driver.Tx, query string, args ...any) (sql.Result, error) {
	r.logQuery(query)

	if tx == nil {
		return r.connection.Exec(query, args...)
	}

	sqlTx, isTypeOfSql := tx.(*sql.Tx)
	if !isTypeOfSql {
		return nil, fmt.Errorf("provided type (%T) of driver.Tx is not supported yet", tx)
	}

	return sqlTx.Exec(query, args...)
}

func (r BaseRepository) logQuery(message string) {
	if !r.params.LogInfo {
		return
	}

	r.logger.Infof("sql query was: %s", message)
}
