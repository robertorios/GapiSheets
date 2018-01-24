package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/adam-hanna/randomstrings"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"
)

var tpl *template.Template

type Application struct {
	Priority            string
	NewLocation         string
	SolutionName        string
	CurrentSite         string
	SolutionDescription string
	Server              string
	Dependancies        string
	Department          string
	OurVersion          string
	LatestVersion       string
	CostToUpgrade       string
	StepsToSolution     string
	Progress            string
	// milestone           Milestone
	Percentage int
}

// type Milestone struct {
// 	step string
// }

func init() {
	tpl = template.Must(template.ParseGlob("templates/*.html"))
}

// getClient uses a Context and Config to retrieve a Token
// then generate a Client. It returns the generated Client.
func getClient(ctx context.Context, config *oauth2.Config) *http.Client {
	cacheFile, err := tokenCacheFile()
	if err != nil {
		log.Fatalf("Unable to get path to cached credential file. %v", err)
	}
	tok, err := tokenFromFile(cacheFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(cacheFile, tok)
	}
	return config.Client(ctx, tok)
}

// getTokenFromWeb uses Config to request a Token.
// It returns the retrieved Token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var code string
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
	}

	tok, err := config.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	return tok
}

// tokenCacheFile generates credential file path/filename.
// It returns the generated credential path/filename.
func tokenCacheFile() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	tokenCacheDir := filepath.Join(usr.HomeDir, ".credentials")
	os.MkdirAll(tokenCacheDir, 0700)
	return filepath.Join(tokenCacheDir,
		url.QueryEscape("sheets.googleapis.com-go-quickstart.json")), err
}

// tokenFromFile retrieves a Token from a given file path.
// It returns the retrieved Token and any read error encountered.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	t := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(t)
	defer f.Close()
	return t, err
}

// saveToken uses a file path to create a file and store the
// token in it.
func saveToken(file string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", file)
	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func main() {
	http.HandleFunc("/", index)
	http.ListenAndServe(":8081", nil)

}

func index(w http.ResponseWriter, req *http.Request) {
	ctx := context.Background()

	b, err := ioutil.ReadFile("client_secret.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved credentials
	// at ~/.credentials/sheets.googleapis.com-go-quickstart.json
	config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets.readonly")
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(ctx, config)

	srv, err := sheets.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets Client %v", err)
	}

	// Prints the names and majors of students in a sample spreadsheet:
	// https://docs.google.com/spreadsheets/d/1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms/edit
	spreadsheetId := "1h3SaiPhvLDkL4rPQkBE3uTq6pKXxUmHcGL2jXJXJ3MI"
	readRange := "Applications!A3:N"
	readRange1 := "Milestones!A2:G"

	resp, err := srv.Spreadsheets.Values.Get(spreadsheetId, readRange).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet. %v", err)
	}

	resp1, err1 := srv.Spreadsheets.Values.Get(spreadsheetId, readRange1).Do()
	if err1 != nil {
		log.Fatalf("Unable to retrieve data from sheet. %v", err1)
	}

	apps := make([]Application, 0)
	// Miles := make([]Milestone, 0)
	// // fmt.Println(resp1.Values[0])
	// // fmt.Println("resp1 type: ", reflect.TypeOf(resp1))

	// // numRows := 10

	// // Initialize a ten length slice of empty slices
	// grid := make([][]string, len(resp1.Values))

	// // Verify it is a slice of ten empty slices
	// fmt.Println(grid)

	// // Initialize those 10 empty slices
	// for i := 0; i < len(resp1.Values); i++ {
	// 	grid[i] = make([]string, 0)
	// }

	// // grid is a 2d slice of ints with dimensions 10x4
	// fmt.Println(grid)
	// fmt.Printf("%#v", resp1.Values[0])
	// if len(resp1.Values) > 0 {
	// 	for _, resp1 := range resp1.Values {
	// 		// grid[0] = reflect.ValueOf(resp1)
	// 		// fmt.Println(reflect.TypeOf(grid[0]))
	// 		// fmt.Println("reflect value:", reflect.ValueOf(resp[0]))
	// 		// fmt.Println("rest:", resp1)
	// 		// fmt.Println("rest type:", reflect.TypeOf(resp1))
	// 		// fmt.Println("resp1 type inside for: ", reflect.TypeOf(resp1))
	// 		// for _, resp1inner := range [0:len(resp1)] {
	// 		// 	l[0] = resp1inner
	// 		// }
	// 		// for i := 1; i <= len(resp1); i++ {
	// 		// 	fmt.Println(resp1[0])
	// 		// fmt.Println("resp1 type inside for: ", reflect.TypeOf(resp1.Values))
	// 		// l[0] = resp1[i]

	// 		// fmt.Println(len(resp1))
	// 		// fmt.Println("resp1 type inside for: ", reflect.TypeOf(resp1))
	// 		// fmt.Println("resp1 type inside for inside sliceofinterfaces: ", reflect.TypeOf(resp1[0]))
	// 		// s := resp1.([]string)
	// 		mile := Milestone{
	// 			step: resp1,
	// 		}
	// 		// fmt.Println("resp type inside for struct: ", reflect.TypeOf(resp1[0]))
	// 		// fmt.Println(resp1[0])
	// 		Miles = append(Miles, mile)
	// 		// i = i + 1
	// 	}
	// } else {
	// 	fmt.Print("No data found.")
	// }
	// fmt.Println("MILES type: ", reflect.TypeOf(&Miles))
	// fmt.Println("MILES type: ", reflect.TypeOf(Miles[0]))
	// fmt.Println("MILES: ", Miles[0])

	// fmt.Println(&Miles[0][0])
	// fmt.Println(Miles[0].(string))

	// fmt.Println("resp type: ", reflect.TypeOf(resp))
	if len(resp.Values) > 0 {
		// fmt.Println("Name, Major:")

		for _, resp := range resp.Values {
			// fmt.Println("resp type: ", reflect.TypeOf(resp[0]))
			// Print columns A and E, which correspond to indices 0 and 4.
			// var result []interface{}
			var result1 []string
			var result2 []string
			var result_string string
			var percentage int
			if len(resp1.Values) > 0 {
				for _, resp1 := range resp1.Values {
					// fmt.Println("resp Solution Name inside", reflect.ValueOf(resp[3]))
					// fmt.Println("resp1 application inside", reflect.ValueOf(resp1[0]))
					a := resp1[0]
					b := resp[3]
					// fmt.Println(reflect.TypeOf(a))
					// fmt.Println(reflect.TypeOf(b))
					if a == b {
						// result = make([]interface{}, 0)
						// result = append(result, resp1)
						result1 = make([]string, 0)
						result2 = make([]string, 0)
						// fmt.Println(len(resp1))
						sRand, err := randomstrings.GenerateRandomString(16) // generates a 16 digit random string
						if err != nil {
							// panic!
						}
						for i := 0; i < len(resp1); i++ {
							result1 = append(result1, resp1[i].(string)+"hello"+sRand)
						}

						// result2 = append(result1)

						for i := 0; i < len(resp1); i++ {
							if strings.Contains(resp1[i].(string), "done") {
								result2 = append(result2, resp1[i].(string))
							}

						}
						// fmt.Println(result2)
						fmt.Println(len(result1))
						fmt.Println(len(result2))

						// fmt.Println(strings.Join(result1[:], ","))
						result_string = strings.Join(result1[:], ",")
						percentage = (len(result2) * 100) / len(result1)

						// fmt.Println("resp type: ", reflect.TypeOf((len(result2)*100)/len(result1)))
						// fmt.TypeOf((len(result2) * 100) / len(result1))
						fmt.Println(percentage)
						// fmt.Printf("%#v", result1)

						// for i, resp1 := range len(resp1) {
						// 	result1 = append(result1, resp1[i].(string))
						// }
						// result := resp1
						// fmt.Println("resp Solution Name inside if", resp[3])
						// fmt.Println("resp1 milestone inside if", resp1[0])
					}
				}

			}
			var barradeprogreso string
			// sRand, err := randomstrings.GenerateRandomString(16) // generates a 16 digit random string
			// if err != nil {
			// 	// panic!
			// }
			// barradeprogreso = sRand

			// fmt.Println("resp Solution Name", resp)
			app := Application{
				// a =: resp[0]
				Priority:            resp[0].(string),
				NewLocation:         resp[1].(string),
				SolutionName:        resp[3].(string),
				CurrentSite:         resp[6].(string),
				SolutionDescription: resp[4].(string),
				Server:              resp[5].(string),
				Dependancies:        resp[2].(string),
				// Department:          resp[7].(string),
				// OurVersion:          resp[8].(string),
				// LatestVersion:       resp[9].(string),
				// CostToUpgrade:       resp[10].(string),
				StepsToSolution: result_string,
				Progress:        barradeprogreso,
				Percentage:      percentage,
				// milestone: Milestone{
				// 	step: result_string,
				// },
			}
			// fmt.Println("resp type: ", reflect.TypeOf(resp))
			// fmt.Println(resp[0])
			apps = append(apps, app)
			// fmt.Printf("%s, %s\n", row[0], row[4])
		}
	} else {
		fmt.Print("No data found.")
	}
	tpl.ExecuteTemplate(w, "index.html", apps)

}
