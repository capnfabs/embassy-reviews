#!/usr/bin/env python3
import csv

def main():
    capitals = set()
    with open('data/countryInfo.txt', newline='') as csvfile:
        fieldnames = [
            'iso',
            'iso3',
            'iso-numeric',
            'fips',
            'country',
            'capital',
            'area',
            'population',
            'continent',
            'tld',
            'currencycode',
            'currencyname',
            'phone',
            'postal_code_format',
            'postal_code_regex',
            'languages',
            'geonameid',
            'neighbours',
            'equivalentfipscode',
        ]
        reader = csv.DictReader(filter(lambda row: row[0]!='#', csvfile), fieldnames=fieldnames, dialect=csv.excel_tab )
        for row in reader:
            capitals.add({
                'country': row['country'],
                'capital': row['capital']
            })

    with open('data/cities15000.txt', newline='') as csvfile:
        fieldnames = [
            'geonameid',
            'name',
            'asciiname',
            'alternatenames',
            'latitude',
            'longitude',
            'feature_class',
            'feature_code',
            'country_code',
            'cc2',
            'admin1_code',
            'admin2_code',
            'admin3_code',
            'admin4_code',
            'population',
            'elevation',
            'dem',
            'timezone',
            'modification_date',
        ]
        reader = csv.DictReader(csvfile, fieldnames=fieldnames, dialect=csv.excel_tab)
        count = 0
        for row in reader:
            print(row['asciiname'], row['longitude'], row['latitude'])
            count+=1
        print('Cities: %d' % count)

if __name__ == "__main__":
    main()
