package memgraph

import (
	"testing"
	"time"
)

func TestPerson_ToString(t *testing.T) {
	tests := []struct {
		name string
		p    *Person
		want string
	}{
		{
			name: "Test with nil values",
			p: &Person{
				ID:               "1",
				Firstname:        "John",
				Lastname:         "Doe",
				MothersFirstname: "Jane",
				MothersLastname:  "Doe",
				Born:             time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				Birthplace:       "New York",
				Residence:        "New York",
				Death:            time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				Deathplace:       "New York",
			},
			want: "ID: '1', Firstname: 'John', Lastname: 'Doe', MothersFirstname: 'Jane', MothersLastname: 'Doe', Born: date({year:2021, month:1, day:1}), Death: date({year:2021, month:1, day:1}), Birthplace: 'New York', Residence: 'New York', Deathplace: 'New York', OccupationToDisplay: '', ProfilePicture: ''",
		}, {
			name: "Test with All values",
			p: &Person{
				ID:               "1",
				Firstname:        "John",
				Lastname:         "Doe",
				MothersFirstname: "Jane",
				MothersLastname:  "Doe",
				Born:             time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				Birthplace:       "New York",
				Residence:        "New York",
				Death:            time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				Deathplace:       "New York",
				LifeEvents: []map[string]string{
					{
						"date":  "2021-01-01",
						"event": "Event 1",
					},
					{
						"date":  "2021-01-02",
						"event": "Event 2",
					},
				},
				Occupations: []string{
					"Welder",
					"Plumber",
				},
				OccupationToDisplay: "Welder",
				OthersSaid: map[string]string{
					"Beni": "He is a good person",
					"Jani": "He is a bad person",
				},
				Photos: map[string]string{
					"Profile": "profile.jpg",
					"Family":  "family.jpg",
				},
				ProfilePicture: "profile.jpg",
			},
			want: "ID: '1', Firstname: 'John', Lastname: 'Doe', MothersFirstname: 'Jane', MothersLastname: 'Doe', Born: date({year:2021, month:1, day:1}), Death: date({year:2021, month:1, day:1}), Birthplace: 'New York', Residence: 'New York', Deathplace: 'New York', OccupationToDisplay: 'Welder', ProfilePicture: 'profile.jpg', LifeEvents: [{date: '2021-01-01', event: 'Event 1'}, {date: '2021-01-02', event: 'Event 2'}], Occupations: ['Welder', 'Plumber'], OthersSaid: {Beni: 'He is a good person', Jani: 'He is a bad person'}, Photos: {Profile: 'profile.jpg', Family: 'family.jpg'}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.ToString(); got != tt.want {
				t.Errorf("Person.ToString() = %v, want %v", got, tt.want)
			}
		})
	}
}
