package generator

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"regexp"
	"strings"
	"time"
)

// ResponseGenerator generates mock responses
type ResponseGenerator struct {
}

// New creates a new response generator
func New() *ResponseGenerator {
	return &ResponseGenerator{}
}

// GenerateParams contains parameters for generating a response
type GenerateParams struct {
	Type  string
	Body  interface{}
	Faker map[string]interface{}
	Input map[string]interface{}
}

// Generate generates a response based on the configuration
func (g *ResponseGenerator) Generate(params GenerateParams) (interface{}, error) {
	switch params.Type {
	case "static":
		return params.Body, nil
	case "dynamic":
		return g.processDynamic(params.Body, params.Input)
	case "faker":
		return g.generateFaker(params.Faker)
	default:
		return params.Body, nil
	}
}

// generateFaker generates fake data based on faker config
func (g *ResponseGenerator) generateFaker(fakerConfig map[string]interface{}) (interface{}, error) {
	result := make(map[string]interface{})

	fields, ok := fakerConfig["fields"].([]interface{})
	if !ok {
		if fieldsMap, ok := fakerConfig["fields"].(map[string]interface{}); ok {
			for name, fieldConfig := range fieldsMap {
				fieldMap, ok := fieldConfig.(map[string]interface{})
				if !ok {
					continue
				}
				result[name] = g.generateField(fieldMap)
			}
			return result, nil
		}
		return nil, fmt.Errorf("invalid faker config: fields not found")
	}

	for _, field := range fields {
		fieldMap, ok := field.(map[string]interface{})
		if !ok {
			continue
		}

		name, ok := fieldMap["name"].(string)
		if !ok {
			continue
		}

		result[name] = g.generateField(fieldMap)
	}

	return result, nil
}

// generateField generates a single fake field using random data
func (g *ResponseGenerator) generateField(field map[string]interface{}) interface{} {
	fieldType, ok := field["type"].(string)
	if !ok {
		return nil
	}

	format, _ := field["format"].(string)

	switch fieldType {
	case "name":
		return g.randomName()
	case "first_name":
		return g.randomFirstName()
	case "last_name":
		return g.randomLastName()
	case "email":
		return g.randomEmail()
	case "phone", "phone_number":
		return g.randomPhone()
	case "uuid", "uuid_rfc4122":
		return g.randomUUID()
	case "date":
		if format != "" {
			return g.randomDate().Format(format)
		}
		return g.randomDate().Format("2006-01-02")
	case "datetime":
		if format != "" {
			return g.randomDateTime().Format(format)
		}
		return g.randomDateTime().Format(time.RFC3339)
	case "timestamp":
		return g.randomDate().Unix()
	case "word":
		return g.randomWord()
	case "sentence":
		return g.randomSentence()
	case "paragraph":
		return g.randomParagraph()
	case "address":
		return g.randomAddress()
	case "city":
		return g.randomCity()
	case "country":
		return g.randomCountry()
	case "state":
		return g.randomState()
	case "zipcode":
		return g.randomZipcode()
	case "latitude":
		return g.randomLatitude()
	case "longitude":
		return g.randomLongitude()
	case "url":
		return g.randomURL()
	case "ip", "ipv4":
		return g.randomIPv4()
	case "ipv6":
		return g.randomIPv6()
	case "credit_card", "credit_card_number":
		return g.randomCreditCardNumber()
	case "credit_card_type":
		return g.randomCreditCardType()
	default:
		return g.randomWord()
	}
}

// processDynamic processes dynamic response templates
func (g *ResponseGenerator) processDynamic(body interface{}, input map[string]interface{}) (interface{}, error) {
	bodyJSON, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	result := string(bodyJSON)

	// Replace {{input.field}} placeholders
	inputPattern := regexp.MustCompile(`\{\{input\.(\w+)\}\}`)
	result = inputPattern.ReplaceAllStringFunc(result, func(match string) string {
		key := inputPattern.FindStringSubmatch(match)[1]
		if val, ok := input[key]; ok {
			return fmt.Sprintf("%v", val)
		}
		return match
	})

	// Replace {{faker.field}} placeholders
	fakerPattern := regexp.MustCompile(`\{\{faker\.(\w+)\}\}`)
	result = fakerPattern.ReplaceAllStringFunc(result, func(match string) string {
		fieldType := fakerPattern.FindStringSubmatch(match)[1]
		fieldMap := map[string]interface{}{
			"type": fieldType,
		}
		val := g.generateField(fieldMap)
		return fmt.Sprintf("%v", val)
	})

	// Replace {{timestamp}} placeholder
	result = strings.ReplaceAll(result, "{{timestamp}}", fmt.Sprintf("%d", time.Now().Unix()))

	// Replace {{date}} placeholder
	result = strings.ReplaceAll(result, "{{date}}", time.Now().Format("2006-01-02"))

	// Replace {{datetime}} placeholder
	result = strings.ReplaceAll(result, "{{datetime}}", time.Now().Format(time.RFC3339))

	// Replace {{uuid}} placeholder
	result = strings.ReplaceAll(result, "{{uuid}}", g.randomUUID())

	// Replace {{random.int.min-max}} placeholder
	randIntPattern := regexp.MustCompile(`\{\{random\.int.(\d+)-(\d+)\}\}`)
	result = randIntPattern.ReplaceAllStringFunc(result, func(match string) string {
		groups := randIntPattern.FindStringSubmatch(match)
		min := atoi(groups[1])
		max := atoi(groups[2])
		return fmt.Sprintf("%d", rand.Intn(max-min+1)+min)
	})

	// Replace {{random.string.N}} placeholder
	randStrPattern := regexp.MustCompile(`\{\{random\.string.(\d+)\}\}`)
	result = randStrPattern.ReplaceAllStringFunc(result, func(match string) string {
		groups := randStrPattern.FindStringSubmatch(match)
		length := atoi(groups[1])
		letters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
		str := make([]byte, length)
		for i := range str {
			str[i] = letters[rand.Intn(len(letters))]
		}
		return string(str)
	})

	var parsed interface{}
	if err := json.Unmarshal([]byte(result), &parsed); err != nil {
		return nil, err
	}

	return parsed, nil
}

// Helper functions for random data generation
func (g *ResponseGenerator) randomName() string {
	names := []string{"John Doe", "Jane Smith", "Michael Johnson", "Emily Brown", "David Wilson"}
	return names[rand.Intn(len(names))]
}

func (g *ResponseGenerator) randomFirstName() string {
	names := []string{"John", "Jane", "Michael", "Emily", "David", "Sarah", "James", "Emma", "Robert", "Olivia"}
	return names[rand.Intn(len(names))]
}

func (g *ResponseGenerator) randomLastName() string {
	names := []string{"Doe", "Smith", "Johnson", "Brown", "Wilson", "Davis", "Miller", "Moore", "Taylor", "Anderson"}
	return names[rand.Intn(len(names))]
}

func (g *ResponseGenerator) randomEmail() string {
	domains := []string{"gmail.com", "yahoo.com", "outlook.com", "example.com", "test.com"}
	return fmt.Sprintf("user%d@%s", rand.Intn(1000), domains[rand.Intn(len(domains))])
}

func (g *ResponseGenerator) randomPhone() string {
	return fmt.Sprintf("+1-%d-%d-%d", rand.Intn(900)+100, rand.Intn(900)+100, rand.Intn(9000)+1000)
}

func (g *ResponseGenerator) randomUUID() string {
	return fmt.Sprintf("%x-%x-%x-%x-%x",
		rand.Intn(0xFFFFFFFF),
		rand.Intn(0xFFFF),
		rand.Intn(0xFFFF),
		rand.Intn(0xFFFF),
		rand.Intn(0xFFFFFFFFFFFF))
}

func (g *ResponseGenerator) randomDate() time.Time {
	now := time.Now()
	offset := rand.Intn(365*24*60*60) - 180*24*60*60
	return now.Add(time.Duration(offset) * time.Second)
}

func (g *ResponseGenerator) randomDateTime() time.Time {
	now := time.Now()
	offset := rand.Intn(365*24*60*60) - 180*24*60*60
	return now.Add(time.Duration(offset) * time.Second)
}

func (g *ResponseGenerator) randomWord() string {
	words := []string{"hello", "world", "test", "data", "mock", "api", "response", "request", "server", "client"}
	return words[rand.Intn(len(words))]
}

func (g *ResponseGenerator) randomSentence() string {
	return fmt.Sprintf("This is a random test sentence with %s.", g.randomWord())
}

func (g *ResponseGenerator) randomParagraph() string {
	return fmt.Sprintf("%s %s %s", g.randomSentence(), g.randomSentence(), g.randomSentence())
}

func (g *ResponseGenerator) randomAddress() string {
	return fmt.Sprintf("%d %s Street", rand.Intn(999)+1, g.randomWord())
}

func (g *ResponseGenerator) randomCity() string {
	cities := []string{"New York", "Los Angeles", "Chicago", "Houston", "Phoenix", "Philadelphia"}
	return cities[rand.Intn(len(cities))]
}

func (g *ResponseGenerator) randomCountry() string {
	countries := []string{"United States", "Canada", "United Kingdom", "Germany", "France", "Japan"}
	return countries[rand.Intn(len(countries))]
}

func (g *ResponseGenerator) randomState() string {
	states := []string{"NY", "CA", "IL", "TX", "AZ", "PA"}
	return states[rand.Intn(len(states))]
}

func (g *ResponseGenerator) randomZipcode() string {
	return fmt.Sprintf("%05d", rand.Intn(90000)+10000)
}

func (g *ResponseGenerator) randomLatitude() float64 {
	return float64(rand.Intn(180*100000)-90*100000) / 100000
}

func (g *ResponseGenerator) randomLongitude() float64 {
	return float64(rand.Intn(360*100000)-180*100000) / 100000
}

func (g *ResponseGenerator) randomURL() string {
	return fmt.Sprintf("https://%s.example.com/%s", g.randomWord(), g.randomWord())
}

func (g *ResponseGenerator) randomIPv4() string {
	return fmt.Sprintf("%d.%d.%d.%d", rand.Intn(256), rand.Intn(256), rand.Intn(256), rand.Intn(256))
}

func (g *ResponseGenerator) randomIPv6() string {
	return fmt.Sprintf("2001:0db8:%04x:%04x:%04x:%04x:%04x:%04x",
		rand.Intn(0xFFFF), rand.Intn(0xFFFF), rand.Intn(0xFFFF),
		rand.Intn(0xFFFF), rand.Intn(0xFFFF), rand.Intn(0xFFFF))
}

func (g *ResponseGenerator) randomCreditCardNumber() string {
	return fmt.Sprintf("%04d %04d %04d %04d", rand.Intn(9999)+1, rand.Intn(9999), rand.Intn(9999), rand.Intn(9999))
}

func (g *ResponseGenerator) randomCreditCardType() string {
	types := []string{"Visa", "MasterCard", "American Express", "Discover"}
	return types[rand.Intn(len(types))]
}

// Helper to convert string to int
func atoi(s string) int {
	n := 0
	for _, c := range s {
		n = n*10 + int(c-'0')
	}
	return n
}

// Generator is a global instance
var Generator = New()

// GenerateStatic generates a static response
func GenerateStatic(body interface{}) (interface{}, error) {
	return Generator.Generate(GenerateParams{
		Type: "static",
		Body: body,
	})
}

// GenerateDynamic generates a dynamic response
func GenerateDynamic(body interface{}, input map[string]interface{}) (interface{}, error) {
	return Generator.Generate(GenerateParams{
		Type:  "dynamic",
		Body:  body,
		Input: input,
	})
}

// GenerateFaker generates a faker response
func GenerateFaker(fakerConfig map[string]interface{}) (interface{}, error) {
	return Generator.Generate(GenerateParams{
		Type:  "faker",
		Faker: fakerConfig,
	})
}
