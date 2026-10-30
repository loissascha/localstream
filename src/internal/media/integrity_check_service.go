package media

type IntegrityCheckService struct {
}

func NewIntegrityCheckService() *IntegrityCheckService {
	return &IntegrityCheckService{}
}

func (s *IntegrityCheckService) CheckFFMpegErrors(path string) {

}

func (s *IntegrityCheckService) ReencodeFile(path string, output string) {

}
