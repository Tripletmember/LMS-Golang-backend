package service

import (
	"context"

	"github.com/zhashkevych/creatly-backend/pkg/payment"
	"github.com/zhashkevych/creatly-backend/pkg/payment/fondy"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/zhashkevych/creatly-backend/internal/domain"
	"github.com/zhashkevych/creatly-backend/internal/repository"
	"github.com/zhashkevych/creatly-backend/pkg/cache"
)

type SchoolsService struct {
	repo  repository.Schools
	cache cache.Cache
	ttl   int64
}

func NewSchoolsService(repo repository.Schools, cache cache.Cache, ttl int64) *SchoolsService {
	return &SchoolsService{repo: repo, cache: cache, ttl: ttl}
}

func (s *SchoolsService) Create(ctx context.Context, name string) (primitive.ObjectID, error) {
	return s.repo.Create(ctx, name)
}

func (s *SchoolsService) GetByDomain(ctx context.Context, domainName string) (domain.School, error) {
	if value, err := s.cache.Get(domainName); err == nil {
		return value.(domain.School), nil
	}

	school, err := s.repo.GetByDomain(ctx, domainName)
	if err != nil {
		return domain.School{}, err
	}

	err = s.cache.Set(domainName, school, s.ttl)

	return school, err
}

func (s *SchoolsService) GetById(ctx context.Context, id primitive.ObjectID) (domain.School, error) {
	return s.repo.GetById(ctx, id)
}

func (s *SchoolsService) UpdateSettings(ctx context.Context, schoolId primitive.ObjectID, inp UpdateSchoolSettingsInput) error {
	return s.repo.UpdateSettings(ctx, repository.UpdateSchoolSettingsInput{
		SchoolID:            schoolId,
		Name:                inp.Name,
		Color:               inp.Color,
		Domains:             inp.Domains,
		Email:               inp.Email,
		ContactInfo:         toContactInfo(inp.ContactInfo),
		Pages:               toPages(inp.Pages),
		ShowPaymentImages:   inp.ShowPaymentImages,
		GoogleAnalyticsCode: inp.GoogleAnalyticsCode,
		LogoURL:             inp.LogoURL,
	})
}

func (s *SchoolsService) ConnectFondy(ctx context.Context, input ConnectFondyInput) error {
	client := fondy.NewFondyClient(input.MerchantID, input.MerchantPassword)

	id := primitive.NewObjectID()

	_, err := client.GeneratePaymentLink(payment.GeneratePaymentLinkInput{
		OrderId:   id.Hex(),
		Amount:    1000,
		Currency:  "USD",
		OrderDesc: "CREATLY - TESTING FONDY CREDENTIALS",
	})
	if err != nil {
		return err
	}

	creds := domain.Fondy{
		MerchantPassword: input.MerchantPassword,
		MerchantId:       input.MerchantID,
		Connected:        true,
	}

	return s.repo.SetFondyCredentials(ctx, input.SchoolID, creds)
}

func toContactInfo(inp *UpdateSchoolSettingsContactInfo) *repository.UpdateSchoolSettingsContactInfo {
	if inp == nil {
		return nil
	}

	contactInfo := new(repository.UpdateSchoolSettingsContactInfo)

	if inp.Address != nil {
		contactInfo.Address = inp.Address
	}

	if inp.BusinessName != nil {
		contactInfo.BusinessName = inp.BusinessName
	}

	if inp.Email != nil {
		contactInfo.Email = inp.Email
	}

	if inp.Phone != nil {
		contactInfo.Phone = inp.Phone
	}

	if inp.RegistrationNumber != nil {
		contactInfo.RegistrationNumber = inp.RegistrationNumber
	}

	return contactInfo
}

func toPages(inp *UpdateSchoolSettingsPages) *repository.UpdateSchoolSettingsPages {
	if inp == nil {
		return nil
	}

	pages := new(repository.UpdateSchoolSettingsPages)

	if inp.Confidential != nil {
		pages.Confidential = inp.Confidential
	}

	if inp.NewsletterConsent != nil {
		pages.NewsletterConsent = inp.NewsletterConsent
	}

	if inp.ServiceAgreement != nil {
		pages.ServiceAgreement = inp.ServiceAgreement
	}

	return pages
}
