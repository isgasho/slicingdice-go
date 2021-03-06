package main

import (
	"../slicingdice"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"math"
	"os/signal"
	"strings"
	"time"
	"reflect"
)

type SlicingDiceTester struct {
	client           *slicingdice.SlicingDice
	columnTranslation map[string]string
	sleepTime        int
	path             string
	extension        string
	numSuccess       int
	numFails         int
	failedTests      []string
	verbose          bool
	perTestInsert	 bool
}

// Run all the tests of determined query type
func (s *SlicingDiceTester) runTests(queryType string) {
	testData := s.loadTestData(queryType, "").([]interface{})
	numTests := len(testData)

	singleInsert := testData[0].(map[string]interface{})
	s.perTestInsert = singleInsert["insert"] != nil

	for i, test := range testData {
		queryType_ := queryType
		var err error
		var result map[string]interface{}
		testConverted := test.(map[string]interface{})
		s.emptyColumnTranslation()

		fmt.Printf("(%v/%v) Executing test \"%v\"\n", i+1, numTests, testConverted["name"])

		if _, ok := testConverted["description"]; ok {
			fmt.Printf("  Description: %v\n", testConverted["description"])
		}

		fmt.Printf("  Query type: %v\n", queryType)
		if s.perTestInsert {
			err = s.createColumns(testConverted)
			if err != nil {
				s.compareResult(testConverted, nil, err, queryType)
				continue
			}
			err = s.insertData(testConverted)
			if err != nil {
				s.compareResult(testConverted, nil, err, queryType)
				continue
			}	
		}

		if queryType == "delete" || queryType == "update" {
			result_additional := s.runAdditionalOperations(queryType, testConverted)
			if !result_additional {
				continue
			}
			queryType_ = "count_entity"
		}

		result, err = s.executeQuery(queryType_, testConverted)
		if err != nil {
			if queryType == "delete" || queryType == "update" {
				s.numFails += 1
				s.failedTests = append(s.failedTests, testConverted["name"].(string))
				expectedData, _ := json.Marshal(testConverted["expected"])
				fmt.Printf("  Expected: \"%v\": %v\n", "result", string(expectedData))
				fmt.Printf("  Result: \"result\": %v\n", err)
				fmt.Println("  Status: Failed\n")
				continue
			}
			s.compareResult(testConverted, nil, err, queryType)
			continue
		}

		s.compareResult(testConverted, result, nil, queryType)
	}
}

// Method used to run delete and update operations, this operations
// are executed before the query and the result comparison
func (s *SlicingDiceTester) runAdditionalOperations(queryType string, test map[string]interface{}) bool {
	queryData := s.translateColumnNames(test["additional_operation"].(map[string]interface{}), false)
	if queryType == "delete" {
		fmt.Println("  Deleting")
	} else {
		fmt.Println("  Updating")
	}

	if s.verbose {
		fmt.Println("    - %v", queryData)
	}

	var result map[string]interface{}
	var err error
	if queryType == "delete" {
		result, err = s.client.Delete(queryData)
	} else if queryType == "update" {
		result, err = s.client.Update(queryData)
	}

	expected := s.translateColumnNames(test["result_additional"].(map[string]interface{}), false)
	if result == nil {
		s.numFails += 1
		s.failedTests = append(s.failedTests, test["name"].(string))
		fmt.Printf("  Expected: \"%v\": %v\n", "result", expected)
		fmt.Printf("  Result: \"result\": %v\n", err)
		fmt.Println("  Status: Failed\n")
		return false
	}

	for key, value := range expected {
		if value == "ignore" {
			continue
		}

		if !s.compareJSONValue(result[key], expected[key]) {
			s.numFails += 1
			s.failedTests = append(s.failedTests, test["name"].(string))

			fmt.Printf("  Expected: \"%v\": %v\n", key, expected)
			fmt.Printf("  Result: \"result\": %v\n", result)
			fmt.Println("  Status: Failed\n")
			return false
		} else {
			s.numSuccess += 1
			fmt.Println("  Status: Passed\n")
			return true
		}
	}

	return true
}

// Create columns on Slicing Dice API
func (s *SlicingDiceTester) createColumns(test map[string]interface{}) error {
	var columnOrColumns string
	columns := test["columns"].([]interface{})
	isSingular := len(columns) == 1

	if isSingular {
		columnOrColumns = "column"
	} else {
		columnOrColumns = "columns"
	}

	fmt.Printf("  Creating %v %v\n", len(columns), columnOrColumns)

	for _, column := range columns {
		newColumn := s.appendTimestampToColumnName(column.(map[string]interface{}))
		_, err := s.client.CreateColumn(newColumn)

		if err != nil {
			return err
		}

		if s.verbose {
			fmt.Printf("    - %v\n", newColumn["api-name"])
		}
	}
	return nil
}

/* Append a timestamp to column name
This technique allows the same test suite to be executed over and over
again, since each execution will use different column names.
*/
func (s *SlicingDiceTester) appendTimestampToColumnName(column map[string]interface{}) map[string]interface{} {
	oldName := fmt.Sprintf("\"%v\"", column["api-name"])

	timestamp := s.getTimestamp()
	column["name"] = column["name"].(string) + timestamp
	column["api-name"] = column["api-name"].(string) + timestamp
	newName := fmt.Sprintf("\"%v\"", column["api-name"])

	s.columnTranslation[oldName] = newName
	return column
}

// Get actual timestamp on string
func (s *SlicingDiceTester) getTimestamp() string {
	now := time.Now().UnixNano() / int64(time.Millisecond)
	return fmt.Sprintf("%v", now)
}

// Erase column translation map
func (s *SlicingDiceTester) emptyColumnTranslation() {
	s.columnTranslation = map[string]string{}
}

// Insert data on Slicing Dice API
func (s *SlicingDiceTester) insertData(test map[string]interface{}) error {
	var entityOrEntities string
	insert := test["insert"].(map[string]interface{})
	isSingular := len(insert) == 1

	if isSingular {
		entityOrEntities = "entity"
	} else {
		entityOrEntities = "entities"
	}

	fmt.Printf("  Inserting %v %v\n", len(insert), entityOrEntities)

	insertDataTranslated := s.translateColumnNames(insert, true)

	if s.verbose {
		fmt.Printf("    - %v\n", insertDataTranslated)
	}

	_, err := s.client.Insert(insertDataTranslated)
	if err != nil {
		fmt.Println(err)
		return err
	}

	time.Sleep(time.Duration(s.sleepTime) * time.Second)

	return nil
}

// Translate column names to match the name with timestamp
func (s *SlicingDiceTester) translateColumnNames(jsonData map[string]interface{}, isRequest bool) map[string]interface{} {
	dataConverted, _ := json.Marshal(jsonData)
	dataString := string(dataConverted)

	for oldName, newName := range s.columnTranslation {
		dataString = strings.Replace(dataString, oldName, newName, -1)
	}

	if isRequest {
		return s.decodeWithNumberJSON(dataString).(map[string]interface{})
	} else {
		res, _ := s.client.DecodeJSON(dataString)
		return res.(map[string]interface{})
	}
}

// Execute a query of a determined type on Slicing Dice API
func (s *SlicingDiceTester) executeQuery(queryType string, test map[string]interface{}) (map[string]interface{}, error) {
	var result interface{}
	var err error

	if queryType == "sql" {
		result, err = s.client.Sql(test["query"].(string))

		if result == nil {
			return nil, err
		}
		return result.(map[string]interface{}), err
	}

	query := test["query"].(map[string]interface{})
	queryDataTranslated := s.translateColumnNames(query, true)

	fmt.Println("  Querying")
	if s.verbose {
		fmt.Printf("    - %v\n", queryDataTranslated)
	}

	if queryType == "count_entity" {
		result, err = s.client.CountEntity(queryDataTranslated)
	} else if queryType == "count_event" {
		result, err = s.client.CountEvent(queryDataTranslated)
	} else if queryType == "top_values" {
		result, err = s.client.TopValues(queryDataTranslated)
	} else if queryType == "aggregation" {
		result, err = s.client.Aggregation(queryDataTranslated)
	} else if queryType == "result" {
		result, err = s.client.Result(queryDataTranslated)
	} else if queryType == "score" {
		result, err = s.client.Score(queryDataTranslated)
	} 
	if result == nil {
		return nil, err
	}
	return result.(map[string]interface{}), err
}

// Compare and assert result received from Slicing Dice API
func (s *SlicingDiceTester) compareResult(test map[string]interface{}, result map[string]interface{}, err error,
	queryType string) {
	expected := test["expected"].(map[string]interface{})
	if s.perTestInsert {
		expected = s.translateColumnNames(test["expected"].(map[string]interface{}), false)
	}
	if err != nil {
		s.numFails += 1
		s.failedTests = append(s.failedTests, test["name"].(string))
		expectedData, _ := json.Marshal(expected["result"])
		fmt.Printf("  Expected: \"%v\": %v\n", "result", string(expectedData))
		fmt.Printf("  Result: \"result\": %v\n", err)
		fmt.Println("  Status: Failed\n")
	} else if result != nil {
		for key, value := range expected {
			if value == "ignore" {
				continue
			}

			if !s.compareJSONValue(result[key], expected[key]) {
				if !s.testSecondTime(test, expected, key, queryType) {
					resultData, _ := json.Marshal(result[key])
					expectedData, _ := json.Marshal(expected[key])

					s.numFails += 1
					s.failedTests = append(s.failedTests, test["name"].(string))

					fmt.Printf("  Expected: \"%v\": %v\n", key, string(expectedData))
					fmt.Printf("  Result: \"result\": %v\n", string(resultData))
					fmt.Println("  Status: Failed\n")
					return
				}
			} else {
				s.numSuccess += 1
				fmt.Println("  Status: Passed\n")
			}
		}
	}
}

func (s *SlicingDiceTester) testSecondTime(test map[string]interface{}, expected map[string]interface{},
	key string, queryType string) bool {
	fmt.Printf("Trying again")
	time.Sleep(time.Duration(s.sleepTime) * time.Second)

	result, err := s.executeQuery(queryType, test)
	if err != nil {
		return false
	}

	if !s.compareJSONValue(result[key], expected[key]) {
		return false
	} else {
		fmt.Printf("  Passed at second try!")
		s.numSuccess += 1
		fmt.Println("  Status: Passed\n")
		return true
	}
}

func (s *SlicingDiceTester) compareJSON(expected map[string]interface{}, got map[string]interface{}) bool {
	if len(expected) != len(got) {
		return false
	}

	for key, value := range expected {
		valueExpected := value
		valueGot := got[key]

		if !s.compareJSONValue(valueExpected, valueGot) {
			return false
		}
	}

	return true
}

func (s *SlicingDiceTester) compareJSONArray(expected []interface{}, got []interface{}) bool {
	if len(expected) != len(got) {
		return false
	}

	for _, valueExpected := range expected {
		hasEqual := false
		for _, valueGot := range got {
			if s.compareJSONValue(valueExpected, valueGot) {
				hasEqual = true
			}
		}

		if !hasEqual {
			return false
		}
	}

	return true
}

var floatType = reflect.TypeOf(float64(0))

func (s *SlicingDiceTester) getFloat(unk interface{}) (float64, error) {
    v := reflect.ValueOf(unk)
    v = reflect.Indirect(v)
    if !v.Type().ConvertibleTo(floatType) {
        return 0, fmt.Errorf("cannot convert %v to float64", v.Type())
    }
    fv := v.Convert(floatType)
    return fv.Float(), nil
}

func (s *SlicingDiceTester) compareJSONValue(expected interface{}, got interface{}) bool {
	if reflect.ValueOf(expected).Kind() == reflect.Map {
		expectedMap := expected.(map[string]interface{})
		gotMap := got.(map[string]interface{})
		return s.compareJSON(expectedMap, gotMap)
	} else if reflect.ValueOf(expected).Kind() == reflect.Slice {
		expectedArray := expected.([]interface{})
		gotArray := got.([]interface{})
		return s.compareJSONArray(expectedArray, gotArray)
	} else {
		expected_type := reflect.TypeOf(expected)
		expected_kind := expected_type.Kind()
		if expected_kind == reflect.Int && s.isJsonNumber(got) {
			f, _ := got.(json.Number).Int64()
			return expected == f
		} else if expected_kind == reflect.Float64 && s.isJsonNumber(got) {
			f, _ := got.(json.Number).Float64()
			return s.floatIsClose(expected, f)
		} else if expected_kind == reflect.Float64 {
			return s.floatIsClose(expected, got)
		} else {
			return reflect.DeepEqual(expected, got)
		}
	}
}

func (s *SlicingDiceTester) floatIsClose(a interface{}, b interface{}) bool {
	a_float, _ := s.getFloat(a)
	b_float, _ := s.getFloat(b)
	return math.Abs(a_float - b_float) <= math.Max(1e-09 * math.Max(math.Abs(a_float), math.Abs(b_float)), 0.0)
}

func (s *SlicingDiceTester) isJsonNumber(to_test interface{}) bool {
	switch to_test.(type) {
	case json.Number:
		return true
	default:
		return false
	}
}

// Load test data from examples folder
func (s *SlicingDiceTester) loadTestData(queryType string, suffix string) interface{} {
	filename := s.path + queryType + suffix + s.extension
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	return s.decodeWithNumberJSON(string(file))
}

func (s *SlicingDiceTester) decodeWithNumberJSON(jsonData string) interface{} {
	var f interface{}
	d := json.NewDecoder(strings.NewReader(jsonData))
	d.UseNumber()
	if err := d.Decode(&f); err != nil {
		log.Fatal(err)
	}
	return f
}


func newSlicingDiceTester(apiKey string, verboseOption bool) (t *SlicingDiceTester) {
	keys := new(slicingdice.APIKey)
	keys.MasterKey = apiKey
	sdTester := new(SlicingDiceTester)
	sdTester.client = slicingdice.New(keys, 60)
	sdTester.verbose = verboseOption

	// Sleep Time in seconds
	sdTester.sleepTime = 10
	// Path for examples
	sdTester.path = "examples/"
	// Examples files extension
	sdTester.extension = ".json"

	sdTester.numSuccess = 0
	sdTester.numFails = 0
	return sdTester
}

func printResults(sdTester *SlicingDiceTester) {
	fmt.Println("\nResults:")
	fmt.Println("  Successes:", sdTester.numSuccess)
	fmt.Println("  Fails:", sdTester.numFails)

	for _, failedTest := range sdTester.failedTests {
		fmt.Println("    - ", failedTest)
	}
	fmt.Println("")

	if sdTester.numFails > 0 {
		var testOrTests string
		isSingular := sdTester.numFails == 1

		if isSingular {
			testOrTests = "test has"
		} else {
			testOrTests = "tests have"
		}
		fmt.Printf("FAIL: %v %v failed\n", sdTester.numFails, testOrTests)
		os.Exit(1)
	} else {
		fmt.Println("SUCCESS: All tests passed")
		os.Exit(0)
	}
}

func main() {
	// SlicingDice queries to be tested. Must match the JSON file name.
	var queryTypes = [9]string{
		"count_entity",
		"count_event",
		"top_values",
		"aggregation",
		"result",
		"score",
		"sql",
		"delete",
		"update",
	}

	apiKey, keySet := os.LookupEnv("SD_API_KEY")

	if !keySet {
		apiKey = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJfX3NhbHQiOiJkZW1vMTMzOG0iLCJwZXJtaXNzaW9uX2xldmVsIjozLCJwcm9qZWN0X2lkIjoyMTMzOCwiY2xpZW50X2lkIjoxMH0.bMUl-VKH8Psjnkmchu0ixOhJti24REVsOCKlnpq6Wws"
	}

	// Testing class with demo API key
	// You can get a new demo API key here: http://panel.slicingdice.com/docs/#api-details-api-connection-api-keys-demo-key
	sdTester := newSlicingDiceTester(apiKey, true,)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		printResults(sdTester)
	}()

	for _, queryType := range queryTypes {
		sdTester.runTests(queryType)
	}
	printResults(sdTester)
}
