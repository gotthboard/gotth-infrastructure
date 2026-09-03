package ops

import "testing"

func validService() Service {
	return Service{Name: "example-app", Image: "registry.example:5000/example/app@sha256:0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef", UID: 10001, GID: 10001, HostPort: 18080, ContainerPort: 8080, Secrets: []Secret{{Name: "database-url", SourceEnv: "DATABASE_URL_FILE", Path: "/run/secrets/database-url"}}}
}

func TestServiceValidation(t *testing.T) {
	if err := validService().Validate(); err != nil {
		t.Fatalf("Validate() returned error: %v", err)
	}
	mutations := map[string]func(*Service){
		"name":           func(s *Service) { s.Name = "Bad Name" },
		"image":          func(s *Service) { s.Image = "example:latest" },
		"uid":            func(s *Service) { s.UID = 0 },
		"gid":            func(s *Service) { s.GID = 0 },
		"host port":      func(s *Service) { s.HostPort = 0 },
		"container port": func(s *Service) { s.ContainerPort = 65536 },
		"secret name":    func(s *Service) { s.Secrets[0].Name = "Bad Secret" },
		"source env":     func(s *Service) { s.Secrets[0].SourceEnv = "database-url-file" },
		"target root":    func(s *Service) { s.Secrets[0].Path = "/etc/passwd" },
		"nested target":  func(s *Service) { s.Secrets[0].Path = "/run/secrets/nested/database-url" },
		"duplicate name": func(s *Service) {
			s.Secrets = append(s.Secrets, Secret{Name: "database-url", SourceEnv: "OTHER_FILE", Path: "/run/secrets/other"})
		},
		"duplicate source": func(s *Service) {
			s.Secrets = append(s.Secrets, Secret{Name: "other", SourceEnv: "DATABASE_URL_FILE", Path: "/run/secrets/other"})
		},
		"duplicate target": func(s *Service) {
			s.Secrets = append(s.Secrets, Secret{Name: "other", SourceEnv: "OTHER_FILE", Path: "/run/secrets/database-url"})
		},
	}
	for name, mutate := range mutations {
		service := validService()
		service.Secrets = append([]Secret(nil), service.Secrets...)
		mutate(&service)
		if err := service.Validate(); err == nil {
			t.Errorf("invalid service %q passed", name)
		}
	}
}
