package dac

type PostgresDatabase struct {
	DataAccess
	PostgresOperator
}

type PostgresOperator struct {
	OperatorI
}

func (m PostgresOperator) BuildQuery(condition Condition, qf *QueryFilter) {
	switch condition.Operator {
	case Equal:
		m.Equal(condition, qf)
	case NotEqual:
		m.NotEqual(condition, qf)
	}
}

func (m PostgresOperator) Equal(condition Condition, qf *QueryFilter) {
	qf.And(condition.Key+" = ?", condition.Value)
}
func (m PostgresOperator) NotEqual(condition Condition, qf *QueryFilter) {
	qf.And(condition.Key+" != ?", condition.Value)
}

type PostgresProvider struct {
	Field
}

func init() {
	RegisterDatabase(Postgres, &PostgresDatabase{})
	RegisterOperator(Postgres, &PostgresOperator{})
	RegisterFunctionProvider(Postgres, &PostgresProvider{})
}
func (p *PostgresProvider) Length() string {
	return "LENGTH(%s)"
}
func (p *PostgresProvider) Max() string {
	return "MAX(%s)"
}
