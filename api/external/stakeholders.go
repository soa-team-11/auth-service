package external

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"

	"github.com/soa-team-11/auth-service/utils"
)

type StakeholderServiceExt interface {
	CreateProfile(userId uuid.UUID) (bool, error)
}

type StakeholderService struct{}

func (s *StakeholderService) CreateProfile(userId uuid.UUID) (bool, error) {
	stakeholdersUrl := utils.Getenv("STAKEHOLDERS_SERVICE_URL", "http://localhost:8081")

	data := map[string]string{
		"user_id": userId.String(),
	}

	jsonData, _ := json.Marshal(data)

	resp, err := http.Post(fmt.Sprintf("%s/profiles", stakeholdersUrl), "application/json", bytes.NewBuffer(jsonData))

	if err != nil {
		fmt.Println("Request error:", err)
		return false, err
	}
	defer resp.Body.Close()

	return true, nil
}
