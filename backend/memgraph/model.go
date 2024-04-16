package memgraph

import (
	"fmt"
	"time"
)

var RelationshipTypes = []string{
	"Parent",
	"Child",
	"Spouse",
	"Sibling",
}

type Person struct {
	ID                  string              `json:"id"`
	Firstname           string              `json:"first_name"`
	Middlename          string              `json:"middle_name"`
	Lastname            string              `json:"last_name"`
	Titles              []string            `json:"titles"`   // e.g. Jr., Sr., III
	Suffixes            []string            `json:"suffixes"` // e.g. Ph.D., M.D.
	ExtraNames          []string            `json:"extra_names"`
	Aliases             []string            `json:"aliases"`
	MothersFirstname    string              `json:"mothers_first_name"`
	MothersLastname     string              `json:"mothers_last_name"`
	Born                time.Time           `json:"born"`
	Birthplace          string              `json:"birthplace"`
	Residence           string              `json:"residence"`
	Death               time.Time           `json:"death"`
	Deathplace          string              `json:"deathplace"`
	LifeEvents          []map[string]string `json:"life_events"`
	Occupations         []string            `json:"occupation"`
	OccupationToDisplay string              `json:"occupation_to_display"`
	OthersSaid          map[string]string   `json:"others_said"`
	Photos              map[string]string   `json:"photos"`
	ProfilePicture      string              `json:"profile_picture"`
}

func (p *Person) ToString() string {
	result := fmt.Sprintf("ID: '%s'", p.ID)
	if p.Firstname != "" {
		result = fmt.Sprintf("%s, Firstname: '%s'", result, p.Firstname)
	}
	if p.Lastname != "" {
		result = fmt.Sprintf("%s, Lastname: '%s'", result, p.Lastname)
	}
	if p.Middlename != "" {
		result = fmt.Sprintf("%s, Middlename: '%s'", result, p.Middlename)
	}
	if p.MothersFirstname != "" {
		result = fmt.Sprintf("%s, MothersFirstname: '%s'", result, p.MothersFirstname)
	}
	if p.MothersLastname != "" {
		result = fmt.Sprintf("%s, MothersLastname: '%s'", result, p.MothersLastname)
	}
	if !p.Born.IsZero() {
		result = fmt.Sprintf("%s, Born: date({year:%d, month:%d, day:%d})", result, p.Born.Year(), p.Born.Month(), p.Born.Day())
	}
	if !p.Death.IsZero() {
		result = fmt.Sprintf("%s, Death: date({year:%d, month:%d, day:%d})", result, p.Death.Year(), p.Death.Month(), p.Death.Day())
	}
	if p.Birthplace != "" {
		result = fmt.Sprintf("%s, Birthplace: '%s'", result, p.Birthplace)
	}
	if p.Residence != "" {
		result = fmt.Sprintf("%s, Residence: '%s'", result, p.Residence)
	}
	if p.Deathplace != "" {
		result = fmt.Sprintf("%s, Deathplace: '%s'", result, p.Deathplace)
	}
	if p.OccupationToDisplay != "" {
		result = fmt.Sprintf("%s, OccupationToDisplay: '%s'", result, p.OccupationToDisplay)
	}
	if p.ProfilePicture != "" {
		result = fmt.Sprintf("%s, ProfilePicture: '%s'", result, p.ProfilePicture)
	}

	if p.Titles != nil && len(p.Titles) > 0 {
		result = fmt.Sprintf("%s, Titles: [", result)
		for _, title := range p.Titles {
			result = fmt.Sprintf("%s'%s', ", result, title)
		}
		result = fmt.Sprintf("%s]", result[:len(result)-2])
	}

	if p.Suffixes != nil && len(p.Suffixes) > 0 {
		result = fmt.Sprintf("%s, Suffixes: [", result)
		for _, suffix := range p.Suffixes {
			result = fmt.Sprintf("%s'%s', ", result, suffix)
		}
		result = fmt.Sprintf("%s]", result[:len(result)-2])
	}

	if p.ExtraNames != nil && len(p.ExtraNames) > 0 {
		result = fmt.Sprintf("%s, ExtraNames: [", result)
		for _, extraName := range p.ExtraNames {
			result = fmt.Sprintf("%s'%s', ", result, extraName)
		}
		result = fmt.Sprintf("%s]", result[:len(result)-2])
	}

	if p.Aliases != nil && len(p.Aliases) > 0 {
		result = fmt.Sprintf("%s, Aliases: [", result)
		for _, alias := range p.Aliases {
			result = fmt.Sprintf("%s'%s', ", result, alias)
		}
		result = fmt.Sprintf("%s]", result[:len(result)-2])
	}

	if p.LifeEvents != nil && len(p.LifeEvents) > 0 {
		result = fmt.Sprintf("%s, LifeEvents: [", result)
		for i := 0; i < len(p.LifeEvents); i++ {
			date, dok := p.LifeEvents[i]["date"]
			event, eok := p.LifeEvents[i]["event"]
			if dok && eok {
				result = fmt.Sprintf("%s{date: '%s', event: '%s'}, ", result, date, event)
			}
		}
		result = fmt.Sprintf("%s]", result[:len(result)-2])
	}

	if p.Occupations != nil && len(p.Occupations) > 0 {
		result = fmt.Sprintf("%s, Occupations: [", result)

		for _, occupation := range p.Occupations {
			result = fmt.Sprintf("%s'%s', ", result, occupation)
		}

		result = fmt.Sprintf("%s]", result[:len(result)-2])
	}

	if p.OthersSaid != nil {
		result = fmt.Sprintf("%s, OthersSaid: {", result)
		for key, value := range p.OthersSaid {
			result = fmt.Sprintf("%s%s: '%s', ", result, key, value)
		}
		result = fmt.Sprintf("%s}", result[:len(result)-2])
	}

	if p.Photos != nil && len(p.Photos) > 0 {
		result = fmt.Sprintf("%s, Photos: {", result)
		for key, value := range p.Photos {
			result = fmt.Sprintf("%s%s: '%s', ", result, key, value)
		}
		result = fmt.Sprintf("%s}", result[:len(result)-2])
	}

	return result
}

// Verify checks if the person is valid and does not contain cypher injection it also escapes the delimiters contained in any of the strings
func (p *Person) Verify() error {
	if err := VerifyString(p.ID); err != nil {
		return fmt.Errorf("invalid ID type %s", err)
	}

	p.Firstname = EscapeString(p.Firstname)
	p.Middlename = EscapeString(p.Middlename)
	p.Lastname = EscapeString(p.Lastname)
	p.MothersFirstname = EscapeString(p.MothersFirstname)
	p.MothersLastname = EscapeString(p.MothersLastname)
	p.Birthplace = EscapeString(p.Birthplace)
	p.Residence = EscapeString(p.Residence)
	p.Deathplace = EscapeString(p.Deathplace)
	p.OccupationToDisplay = EscapeString(p.OccupationToDisplay)
	p.ProfilePicture = EscapeString(p.ProfilePicture)

	for i, title := range p.Titles {
		p.Titles[i] = EscapeString(title)
	}

	for i, suffix := range p.Suffixes {
		p.Suffixes[i] = EscapeString(suffix)
	}

	for i, extraName := range p.ExtraNames {
		p.ExtraNames[i] = EscapeString(extraName)
	}

	for i, alias := range p.Aliases {
		p.Aliases[i] = EscapeString(alias)
	}

	for i, lifeEvent := range p.LifeEvents {
		for key, value := range lifeEvent {
			if key != "date" && key != "event" {
				return fmt.Errorf("invalid key in life event")
			}
			p.LifeEvents[i][key] = EscapeString(value)
		}
	}

	for i, occupation := range p.Occupations {
		p.Occupations[i] = EscapeString(occupation)
	}

	for key, value := range p.OthersSaid {
		if err := VerifyString(key); err != nil {
			return fmt.Errorf("invalid key in others said %s", err)
		}
		p.OthersSaid[key] = EscapeString(value)
	}

	for key, value := range p.Photos {
		if err := VerifyString(key); err != nil {
			return fmt.Errorf("invalid key in photos %s", err)
		}
		p.Photos[key] = EscapeString(value)
	}

	return nil
}

type Relationship struct {
	FirstPersonID  string `json:"first_person_id"`
	SecondPersonID string `json:"second_person_id"`
	Relationship   string `json:"relationship"`
	Direction      string `json:"direction"`
}

// Verify checks if the relationship is valid and does not contain cypher injection
func (r *Relationship) Verify() error {
	if r.Direction != "->" && r.Direction != "<-" && r.Direction != "-" {
		return fmt.Errorf("invalid direction for relationship")
	}

	// Check if the relationship is in the list of valid relationships
	found := false
	for _, relationship := range RelationshipTypes {
		if r.Relationship == relationship {
			found = true

			break
		}
	}
	if !found {
		return fmt.Errorf("invalid relationship type")
	}

	if err := VerifyString(r.FirstPersonID); err != nil {
		return fmt.Errorf("invalid FirstPersonID type %s", err)
	}

	if err := VerifyString(r.SecondPersonID); err != nil {
		return fmt.Errorf("invalid SecondPersonID type %s", err)
	}

	return nil
}

type RelationshipAndPerson struct {
	Relationship Relationship `json:"relationship"`
	Person       Person       `json:"person"`
}

func (r *RelationshipAndPerson) Verify() error {
	if err := r.Relationship.Verify(); err != nil {
		return err
	}

	if err := r.Person.Verify(); err != nil {
		return err
	}

	return nil
}
