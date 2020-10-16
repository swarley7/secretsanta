package main

//assumes you have the configured aws through `aws configure`

import (
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
)

type Person struct {
	Name        string
	Gifts       map[int]string // e.g. {50: "john", 100: "jane"}
	PhoneNumber string
	Email       string
}

func ShuffleUsers(users []string) (shuffled []string) {
	shuffled = make([]string, len(users))
	copy(shuffled, users) // we need an array of users for each
	rand.Shuffle(len(shuffled), func(j, k int) {
		shuffled[j], shuffled[k] = shuffled[k], shuffled[j]
	})
	return shuffled
}

func UniqueMatches(a [][]string) bool {
	for k, _ := range a[0] {
		x := map[string]bool{}
		for _, v := range a {
			x[v[k]] = true
		}
		if len(x) != len(a) {
			// fmt.Println(x)
			return false
		}
	}
	return true
}

func GetMatches(users []string, gifts int, uniquematch bool) [][]string {
	matches := make([][]string, gifts+1)
	for {
		for i := 0; i < gifts+1; i++ {
			matches[i] = ShuffleUsers(users)
		}
		// fmt.Println(matches)
		if !uniquematch {
			break
		}
		if !UniqueMatches(matches) {
			// fmt.Println(UniqueMatches(matches))
			continue
		}
		break
	}
	// fmt.Println(matches)
	return matches
}

type Config struct {
	RepeatPartyMembers bool  `json:"repeat"`
	Amounts            []int `json:"amounts"`
	Party              []struct {
		Name        string `json:"name"`
		PhoneNumber string `json:"phone"`
	} `json:"party"`
}

func LoadConfiguration(file string) Config {
	var config Config
	configFile, err := os.Open(file)
	defer configFile.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)
	return config
}

func main() {
	// Parse config file to get user mappings
	c := flag.String("c", "config.json", "Specify the configuration file.")
	debug := flag.Bool("d", true, "Don't actually do the thing")

	flag.Parse()
	rand.Seed(1338) // Hardcoded seed means you can replicate the results lol! This is handy if your idiot family forgets their fucking match :(
	if *debug {
		rand.Seed(13337) // Debug mode; prints who got who to the terminal. use a different seed cos, like, NO PEEKING! ;-)
	}
	config := LoadConfiguration(*c)
	people := []string{}
	santas := []Person{}
	for _, person := range config.Party {
		people = append(people, person.Name)
		p := Person{Name: person.Name, PhoneNumber: person.PhoneNumber}
		p.Gifts = make(map[int]string, len(config.Amounts))
		santas = append(santas, p)
	}
	matches := GetMatches(people, len(config.Amounts), true)
	for _, santa := range santas { // Users
		var s int
		for i := range matches[0] {
			if matches[0][i] == santa.Name {
				s = i
				break
			}
		}
		for index, amount := range config.Amounts {
			santa.Gifts[amount] = matches[index+1][s]
		}
		if *debug {
			fmt.Printf("Person: %s has:", santa.Name)

			for k, v := range santa.Gifts {
				fmt.Printf(" [%d - %s] ", k, v)
			}
			fmt.Printf("\n")
		}
	}

	if *debug {
		return
	}
	fmt.Println("creating session")
	sess := session.Must(session.NewSession())
	fmt.Println("session created")

	svc := sns.New(sess)
	fmt.Println("service created")
	for _, santa := range santas {
		message := fmt.Sprintf("Hi %s!\nYour Secret Santa 🎅🏻 is:\n", santa.Name) // TODO change to template
		for i, v := range santa.Gifts {
			message += fmt.Sprintf("‣ %s ($%d gift 🎁)\n", v, i)
		}
		message += "Remember, who you get is a secret! 🤫"

		message += "Merry xmas! 🎅🏻"
		params := &sns.PublishInput{

			Message: aws.String(message),

			PhoneNumber: aws.String(santa.PhoneNumber),

			MessageAttributes: map[string]*sns.MessageAttributeValue{

				"AWS.SNS.SMS.SenderID": &sns.MessageAttributeValue{StringValue: aws.String("SecretSanta"), DataType: aws.String("String")}, "AWS.SNS.SMS.SMSType": &sns.MessageAttributeValue{StringValue: aws.String("Promotional"), DataType: aws.String("String")},
			},
		}
		time.Sleep(time.Duration(1 * time.Second))
		resp, err := svc.Publish(params)
		if err != nil {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
			return
		}
		fmt.Println(resp)

	}
}
