package models

import (
	"context"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

//GetCountriesCollection ...
func GetCountriesCollection(collection *mongo.Collection) interface{} {
	var ctx context.Context
	cur, _ := collection.Find(ctx, bson.M{})
	defer cur.Close(ctx)

	resultList := make([]Country, 0)
	var result Country
	for cur.Next(ctx) {
		err := cur.Decode(&result)
		resultList = append(resultList, result)
		if err != nil {
			log.Fatal(err)
		}
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}
	return resultList
}

//CreateCountriesCollection this method is used the first time to create the countries collection
func CreateCountriesCollection(collection *mongo.Collection, countriesData string) {
	err := json.Unmarshal([]byte(countriesData), &ctry)
	if err != nil {
		log.Fatal(err)
	}

	result, err := collection.InsertMany(context.TODO(), ctry)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Total Count: ", len(result.InsertedIDs))
}

//UpdateCountriesCollection ...
func UpdateCountriesCollection(collection *mongo.Collection, countriesData string) {
	err := json.Unmarshal([]byte(countriesData), &countries)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Starting database update...")

	j := 0
	for i := range countries {
		doc := countries[i]
		filter := bson.M{"country": bson.M{"$eq": doc["country"]}}
		update := bson.M{"$set": bson.M{
			"cases":               doc["cases"],
			"todayCases":          doc["todayCases"],
			"deaths":              doc["deaths"],
			"todayDeaths":         doc["todayDeaths"],
			"recovered":           doc["recovered"],
			"active":              doc["active"],
			"critical":            doc["critical"],
			"casesPerOneMillion":  doc["casesPerOneMillion"],
			"deathsPerOneMillion": doc["deathsPerOneMillion"],
			"tests":               doc["tests"],
			"testsPerOneMillion":  doc["testsPerOneMillion"],
		}}
		result, _ := collection.UpdateOne(context.TODO(), filter, update)

		if result.MatchedCount > 0 {
			j++
		}
	}

	fmt.Printf("%v countries updated", j)
}

//UpdateCountriesPopulationDensity ...
func UpdateCountriesPopulationDensity(collection *mongo.Collection, countriesData []CountryData) {
	fmt.Println("Starting database seed...")
	j := 0
	k := len(countriesData)
	for i := range countriesData {
		doc := countriesData[i]
		filter := bson.M{"country": bson.M{"$eq": doc.Country}}
		update := bson.M{"$set": bson.M{"populationDensity": doc.PopulationDensity}}
		result, _ := collection.UpdateOne(context.TODO(), filter, update)
		if result.MatchedCount > 0 {
			j++
		}
	}
	fmt.Printf("Total successful entries were %v out of %v", j, k)
}

//UpdateWeekCases ...
func UpdateWeekCases(collection *mongo.Collection, country string, fortnightChange, weekChange []float64) {
	filter := bson.M{"country": bson.M{"$eq": country}}
	update := bson.M{"$addToSet": bson.M{
		"fortnightCases": bson.M{"$each": fortnightChange},
		"weekCases":      bson.M{"$each": weekChange},
	}}
	collection.UpdateOne(context.TODO(), filter, update)
}
