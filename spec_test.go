package ops

import "testing"

func validService() Service {
	return Service{Name: "example-app", Image: "registry.example/example/app@sha256:0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef", UID: 10001, GID: 10001, HostPort: 18080, ContainerPort: 8080, Secrets: []Secret{{Name: "database-url", Path: "/run/secrets/database-url"}}}
}

func TestServiceValidation(t *testing.T) {
	if err := validService().Validate(); err != nil {
		t.Fatalf("Validate() returned error: %v", err)
	}
	mutations := []func(*Service){
		func(s *Service) { s.Name = "Bad Name" }, func(s *Service) { s.Image = "example:latest" },
		func(s *Service) { s.UID = 0 }, func(s *Service) { s.HostPort = 0 },
		func(s *Service) { s.Secrets[0].Path = "/etc/passwd" },
		func(s *Service) { s.Secrets = append(s.Secrets, s.Secrets[0]) },
	}
	for index, mutate := range mutations {
		service := validService()
		service.Secrets = append([]Secret(nil), service.Secrets...)
		mutate(&service)
		if err := service.Validate(); err == nil {
			t.Errorf("invalid service %d passed", index)
		}
	}
}
