package ip_region

import (
	"fmt"
	"slices"

	ipdb "github.com/ipipdotnet/ipdb-go"
)

type ipdbSearcher struct {
	db       *ipdb.City
	language string
}

func newIPDBSearcher(dbPath string) (*ipdbSearcher, error) {
	db, err := ipdb.NewCity(dbPath)
	if err != nil {
		return nil, err
	}
	for _, field := range []string{"country_name", "region_name", "city_name"} {
		if !slices.Contains(db.Fields(), field) {
			return nil, fmt.Errorf("IPDB is missing field %q; use a standard city database", field)
		}
	}
	// qqwry standard databases use CN; also accept English city databases.
	for _, language := range []string{"CN", "EN"} {
		if slices.Contains(db.Languages(), language) {
			return &ipdbSearcher{db: db, language: language}, nil
		}
	}
	return nil, fmt.Errorf("IPDB has no supported language (CN or EN)")
}

func (s *ipdbSearcher) Search(ip string) (region, error) {
	fields, err := s.db.FindMap(ip, s.language)
	if err != nil {
		return region{}, err
	}
	return region{
		country:  fields["country_name"],
		province: fields["region_name"],
		city:     fields["city_name"],
	}, nil
}

func (s *ipdbSearcher) Close() {
	// The IPDB library loads the file into memory and keeps no open handles.
}
