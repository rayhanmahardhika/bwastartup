package campaign

type Service interface {
	GetCampaigns(userID int) ([]Campaign, error)
	GetCampaignByID(input GetCampaignDetailInput) (Campaign, error)
}

type service struct {
	repository Repository // dependency terhadap kelas Repo
}

func NewService(repository Repository) *service {
	return &service{repository}
}

// tidak menggunakan kelas input karena tidak menggunakan request berupa JSON
func (s *service) GetCampaigns(userID int) ([]Campaign, error) {
	if userID != 0 {
		// find jika userID ada di param
		campaigns, err := s.repository.FindByUserID(userID)
		if err != nil {
			return campaigns, err
		}

		return campaigns, nil
	}

	// find jika userID tidak ada di param
	campaigns, err := s.repository.FindAll()
	if err != nil {
		return campaigns, err
	}

	return campaigns, nil
}

func (s *service) GetCampaignByID(input GetCampaignDetailInput) (Campaign, error) {
	campaign, err := s.repository.FindByID(input.ID)
	if err != nil {
		return campaign, err
	}

	return campaign, nil
}
