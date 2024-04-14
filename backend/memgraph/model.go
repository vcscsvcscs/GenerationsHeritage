package memgraph

import (
	"fmt"
	"time"
)

type Person struct {
	ID                  string              `json:"id"`
	Firstname           string              `json:"first_name"`
	Lastname            string              `json:"last_name"`
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
	result := fmt.Sprintf("ID: '%s', Firstname: '%s', Lastname: '%s', MothersFirstname: '%s', MothersLastname: '%s'", p.ID, p.Firstname, p.Lastname, p.MothersFirstname, p.MothersLastname)
	result = fmt.Sprintf("%s, Born: date({year:%d, month:%d, day:%d}), Death: date({year:%d, month:%d, day:%d})", result, p.Born.Year(), p.Born.Month(), p.Born.Day(), p.Death.Year(), p.Death.Month(), p.Death.Day())
	result = fmt.Sprintf("%s, Birthplace: '%s', Residence: '%s', Deathplace: '%s', OccupationToDisplay: '%s', ProfilePicture: '%s'", result, p.Birthplace, p.Residence, p.Deathplace, p.OccupationToDisplay, p.ProfilePicture)

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
