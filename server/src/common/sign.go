package common

type SignGeneratorInterface interface {
	GenerateSalt() string
	GenerateSign(source, salt string) (string, error)
	GenerateSource(src map[string]string, timestamp uint32) string
	GenerateMSign(src map[string]string, timestamp uint32, salt string) (string, error)
}

type SignGenerator struct {
}

func NewSignGenerator() *SignGenerator {
	return &SignGenerator{}
}

var AppSignMgr SignGeneratorInterface = NewSignGenerator()

func (m *SignGenerator) GenerateSalt() string {
	return RandomStr(5)
}

func (m *SignGenerator) GenerateSign(source, salt string) (string, error) {
	return GenerateDoubleMD5(source, salt)
}

func (m *SignGenerator) GenerateMSign(src map[string]string, timestamp uint32, salt string) (string, error) {
	source := m.GenerateSource(src, timestamp)
	return m.GenerateSign(source, salt)
}

func (m *SignGenerator) GenerateSource(src map[string]string, timestamp uint32) string {
	return GenerateSource(src, timestamp)
}
