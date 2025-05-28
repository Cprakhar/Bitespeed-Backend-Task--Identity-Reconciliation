package services

import (
	"bitespeed-identity-reconciliation/internal/database"
	"bitespeed-identity-reconciliation/internal/models"
	"fmt"
	"sort"
)

type IdentityService struct {
    contactRepo *database.ContactRepository
}

func NewIdentityService() *IdentityService {
    return &IdentityService{
        contactRepo: database.ContactRepo,
    }
}

func (s *IdentityService) IdentifyContact(req *models.IdentifyRequest) (*models.IdentifyResponse, error) {

    if (req.Email == nil || *req.Email == "") && (req.PhoneNumber == nil || *req.PhoneNumber == "") {
        return nil, fmt.Errorf("at least one of email or phoneNumber must be provided")
    }

    existingContacts, err := s.contactRepo.FindByEmailOrPhone(req.Email, req.PhoneNumber)
    if err != nil {
        return nil, fmt.Errorf("error finding existing contacts: %w", err)
    }

    if len(existingContacts) == 0 {
        return s.createNewPrimaryContact(req)
    }

    contactGroups := s.groupContactsByPrimary(existingContacts)

    if len(contactGroups) > 1 {
        return s.mergeContactGroups(contactGroups, req)
    }

    primaryID := s.getPrimaryContactID(contactGroups)
    allContacts := s.getAllContactsInGroup(primaryID)
    
    if s.contactExistsWithExactMatch(allContacts, req) {
        return s.buildResponse(allContacts), nil
    }

    return s.createSecondaryContact(primaryID, allContacts, req)
}

func (s *IdentityService) createNewPrimaryContact(req *models.IdentifyRequest) (*models.IdentifyResponse, error) {
    contact := &models.Contact{
        PhoneNumber:    req.PhoneNumber,
        Email:          req.Email,
        LinkedID:       nil,
        LinkPrecedence: "primary",
    }

    if err := s.contactRepo.Create(contact); err != nil {
        return nil, fmt.Errorf("error creating new primary contact: %w", err)
    }

    return &models.IdentifyResponse{
        Contact: models.ContactInfo{
            PrimaryContactID:     contact.ID,
            Emails:              s.getEmailsFromContact(contact),
            PhoneNumbers:        s.getPhoneNumbersFromContact(contact),
            SecondaryContactIDs: []int{},
        },
    }, nil
}

func (s *IdentityService) groupContactsByPrimary(contacts []models.Contact) map[int][]models.Contact {
    groups := make(map[int][]models.Contact)
    
    for _, contact := range contacts {
        primaryID := contact.ID
        if contact.LinkedID != nil {
            primaryID = *contact.LinkedID
        }
        groups[primaryID] = append(groups[primaryID], contact)
    }
    
    return groups
}

func (s *IdentityService) mergeContactGroups(contactGroups map[int][]models.Contact, req *models.IdentifyRequest) (*models.IdentifyResponse, error) {

    var oldestPrimary *models.Contact
    var allContacts []models.Contact
    
    for primaryID, contacts := range contactGroups {
        for _, contact := range contacts {
            allContacts = append(allContacts, contact)
            if contact.ID == primaryID && (oldestPrimary == nil || contact.CreatedAt.Before(oldestPrimary.CreatedAt)) {
                oldestPrimary = &contact
            }
        }
    }

    if oldestPrimary == nil {
        return nil, fmt.Errorf("no primary contact found")
    }

    for _, contact := range allContacts {
        if contact.LinkPrecedence == "primary" && contact.ID != oldestPrimary.ID {
            err := s.contactRepo.UpdateLinkPrecedence(contact.ID, oldestPrimary.ID, "secondary")
            if err != nil {
                return nil, fmt.Errorf("error updating contact precedence: %w", err)
            }
        }
    }

    mergedContacts := s.getAllContactsInGroup(oldestPrimary.ID)
    
    if !s.contactExistsWithExactMatch(mergedContacts, req) {
        return s.createSecondaryContact(oldestPrimary.ID, mergedContacts, req)
    }

    return s.buildResponse(mergedContacts), nil
}

func (s *IdentityService) getAllContactsInGroup(primaryID int) []models.Contact {
    var allContacts []models.Contact
    
    primary, err := s.contactRepo.FindByID(primaryID)
    if err == nil && primary != nil {
        allContacts = append(allContacts, *primary)
    }
    
    secondaries, err := s.contactRepo.FindByLinkedID(primaryID)
    if err == nil {
        allContacts = append(allContacts, secondaries...)
    }
    
    return allContacts
}

func (s *IdentityService) contactExistsWithExactMatch(contacts []models.Contact, req *models.IdentifyRequest) bool {
    for _, contact := range contacts {
        emailMatch := (req.Email == nil && contact.Email == nil) || 
                     (req.Email != nil && contact.Email != nil && *req.Email == *contact.Email)
        phoneMatch := (req.PhoneNumber == nil && contact.PhoneNumber == nil) || 
                     (req.PhoneNumber != nil && contact.PhoneNumber != nil && *req.PhoneNumber == *contact.PhoneNumber)
        
        if emailMatch && phoneMatch {
            return true
        }
    }
    return false
}

func (s *IdentityService) createSecondaryContact(primaryID int, existingContacts []models.Contact, req *models.IdentifyRequest) (*models.IdentifyResponse, error) {

    if s.hasNewInformation(existingContacts, req) {
        secondaryContact := &models.Contact{
            PhoneNumber:    req.PhoneNumber,
            Email:          req.Email,
            LinkedID:       &primaryID,
            LinkPrecedence: "secondary",
        }

        if err := s.contactRepo.Create(secondaryContact); err != nil {
            return nil, fmt.Errorf("error creating secondary contact: %w", err)
        }

        existingContacts = append(existingContacts, *secondaryContact)
    }

    return s.buildResponse(existingContacts), nil
}

func (s *IdentityService) hasNewInformation(contacts []models.Contact, req *models.IdentifyRequest) bool {
    emails := make(map[string]bool)
    phones := make(map[string]bool)
    
    for _, contact := range contacts {
        if contact.Email != nil {
            emails[*contact.Email] = true
        }
        if contact.PhoneNumber != nil {
            phones[*contact.PhoneNumber] = true
        }
    }
    
    hasNewEmail := req.Email != nil && !emails[*req.Email]
    hasNewPhone := req.PhoneNumber != nil && !phones[*req.PhoneNumber]
    
    return hasNewEmail || hasNewPhone
}

func (s *IdentityService) getPrimaryContactID(contactGroups map[int][]models.Contact) int {
    for primaryID := range contactGroups {
        return primaryID
    }
    return 0
}

func (s *IdentityService) buildResponse(contacts []models.Contact) *models.IdentifyResponse {

    var primary *models.Contact
    var secondaries []models.Contact
    
    for _, contact := range contacts {
        if contact.LinkPrecedence == "primary" {
            primary = &contact
        } else {
            secondaries = append(secondaries, contact)
        }
    }
    
    if primary == nil {
        return nil
    }
    
    emailMap := make(map[string]bool)
    phoneMap := make(map[string]bool)
    var emails []string
    var phoneNumbers []string
    var secondaryIDs []int
    
    if primary.Email != nil && *primary.Email != "" {
        emails = append(emails, *primary.Email)
        emailMap[*primary.Email] = true
    }
    if primary.PhoneNumber != nil && *primary.PhoneNumber != "" {
        phoneNumbers = append(phoneNumbers, *primary.PhoneNumber)
        phoneMap[*primary.PhoneNumber] = true
    }
    
    sort.Slice(secondaries, func(i, j int) bool {
        return secondaries[i].CreatedAt.Before(secondaries[j].CreatedAt)
    })
    
    for _, contact := range secondaries {
        secondaryIDs = append(secondaryIDs, contact.ID)
        
        if contact.Email != nil && *contact.Email != "" && !emailMap[*contact.Email] {
            emails = append(emails, *contact.Email)
            emailMap[*contact.Email] = true
        }
        if contact.PhoneNumber != nil && *contact.PhoneNumber != "" && !phoneMap[*contact.PhoneNumber] {
            phoneNumbers = append(phoneNumbers, *contact.PhoneNumber)
            phoneMap[*contact.PhoneNumber] = true
        }
    }
    
    return &models.IdentifyResponse{
        Contact: models.ContactInfo{
            PrimaryContactID:     primary.ID,
            Emails:              emails,
            PhoneNumbers:        phoneNumbers,
            SecondaryContactIDs: secondaryIDs,
        },
    }
}

func (s *IdentityService) getEmailsFromContact(contact *models.Contact) []string {
    if contact.Email != nil && *contact.Email != "" {
        return []string{*contact.Email}
    }
    return []string{}
}

func (s *IdentityService) getPhoneNumbersFromContact(contact *models.Contact) []string {
    if contact.PhoneNumber != nil && *contact.PhoneNumber != "" {
        return []string{*contact.PhoneNumber}
    }
    return []string{}
}