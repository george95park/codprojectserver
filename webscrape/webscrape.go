// Webscrape for COD: Modern Warfare guns and attachments
// From: https://callofduty.fandom.com/wiki/Call_of_Duty:_Modern_Warfare_(2019)
package webscrape

import (
    "os"
    "fmt"
    "log"
    "net/http"
    "database/sql"
	"github.com/lib/pq"
	"github.com/joho/godotenv"
    "github.com/PuerkitoBio/goquery"
)

type Gun struct {
	Gun_Id int
	Type string
	Name string
}

type Attachment struct {
    Attachment_Id int
	Gun_Id int
	Name string
	SubAttachments []string
}

var gunCount = 1
var attCount = 1

// given the link of the gun, grab its attachments
func getAttachments(href string) []Attachment {
    attachments := []Attachment{}
    // request the html page
    response, err := http.Get("https://callofduty.fandom.com/" + href)
    if err != nil {
        log.Fatal(err)
    }
    defer response.Body.Close()

    // load the html document
    document, err := goquery.NewDocumentFromReader(response.Body)
    if err != nil {
        log.Fatal("Error loading HTTP response body.", err)
    }
    // count the type of attachments for this gun
    countTypeAttachments := 1
    foundTitle := false
    // find the attachments for this gun
    document.Find("h2 .mw-headline").Each(func (index1 int, s1 *goquery.Selection) {
        if s1.Text() == "Call of Duty: Modern Warfare" {
            foundTitle = true
            s1.Parent().NextAll().Find("h3 .mw-headline").Each(func (index2 int, s2 *goquery.Selection) {
                if s2.Text() == "Attachments" {
                    skip := false
                    s2.Parent().NextAll().EachWithBreak(func (index3 int, s3 *goquery.Selection) bool {
                        if s3.Text() == "BlueprintsEdit" || s3.Text() == "TriviaEdit"{
                            return false
                        }
                        if !skip {
                            typeAttachment := s3.Find(".mw-headline").Text()
                            a := Attachment{ Attachment_Id: attCount, Gun_Id: gunCount, Name: typeAttachment }
                            // count the attachments for this type of attachment
                            countAttachments := 0
                            countTypeAttachments += 1
                            // find the attachments for this type of attachment
                            s3.Next().Find("li").Each(func (index4 int, s4 *goquery.Selection) {
                                a.SubAttachments = append(a.SubAttachments, s4.Text())
                                countAttachments += 1
                            })
                            attachments = append(attachments, a)
                            attCount++
                            fmt.Println("GunCount: ", gunCount, " -- Type of Attachment: ",typeAttachment, " -- Number of Attachments: ", countAttachments)
                        }
                        skip = !skip
                        return true
                    })
                }
            })
        }
    })

    // if the html page does not have the title 'Call of Duty: Modern Warfare', then
    // simply look for the tag with the id=attachments
    if !foundTitle {
        skip := false
        document.Find("h2 #Attachments").Parent().NextAll().EachWithBreak(func (index1 int, s1 *goquery.Selection) bool {
            if s1.Text() == "BlueprintsEdit" {
                return false
            }
            if !skip {
                typeAttachment := s1.Find(".mw-headline").Text()
                a := Attachment{ Attachment_Id: attCount, Gun_Id: gunCount, Name: typeAttachment }
                // count the attachments for this type of attachment
                countAttachments := 0
                countTypeAttachments += 1
                // find the attachments for this type of attachment
                s1.Next().Find("li").Each(func (index2 int, s2 *goquery.Selection) {
                    a.SubAttachments = append(a.SubAttachments, s2.Text())
                    countAttachments += 1
                })
                attachments = append(attachments, a)
                attCount++

                fmt.Println("GunCount: ", gunCount, " -- Type of Attachment: ", typeAttachment, " -- Number of Attachments: ", countAttachments)
            }
            skip = !skip
            return true
        })
    }
    fmt.Println("GunCount: ", gunCount, " -- Number of Types of Attachments: ", countTypeAttachments)
    return attachments
}

// gets the first set of guns
func getFirstGuns() ([]Gun, []Attachment) {
    guns := []Gun{}
    attachments := []Attachment{}
    // first set of guns
    categories := []string{"Assault_Rifle", "Submachine_Gun", "Handgun"}
    // these hashes pertains to the html section for Modern Warfare
    hashes := []string{"df9bf0e2c34e5c8c1f5e719f198ad38a", "27a12326b111cb4be7b2c3e31ba82c4b", "d5edf75ef1a22a12956b3eea251d8332"}
    for i := 0; i < len(categories); i++ {
        // make a request to the html page
        response, err := http.Get("https://callofduty.fandom.com/wiki/" + categories[i])
        if err != nil {
            log.Fatal(err)
        }
        // load html document
        document, err := goquery.NewDocumentFromReader(response.Body)
        if err != nil {
            log.Fatal("Error loading HTTP response body.", err)
        }

        // count the number of guns for this type of gun
        countGuns := 1

        // find category[i] title
        document.Find("i div").Each(func (index1 int, s1 *goquery.Selection) {
            h, exists := s1.Attr("hash")
            if exists && h == hashes[i] {
                // find the a tag
                s1.Find(".wikia-gallery-item .lightbox-caption a").Each(func (index2 int, s2 *goquery.Selection) {
                    // grab the href link to this gun
                    href, exists := s2.Attr("href")
                    if exists {
                        g := Gun { Gun_Id: gunCount, Type: categories[i], Name: s2.Text() }
                        // get the attachments for this gun
                        attachments = append(attachments, getAttachments(href)...)
                        guns = append(guns, g)
                        countGuns++
                        gunCount++
                    }
                })
            }
        })
        fmt.Println("Gun Category: ", categories[i], " -- Number of guns: ", countGuns)
        response.Body.Close()
    }
    return guns,attachments
}

// gets the second set of guns
func getSecondGuns() ([]Gun, []Attachment) {
    categories := []string{"Shotgun", "Marksman_Rifle", "Sniper_Rifle"}
    guns := []Gun{}
    attachments := []Attachment{}
    for i := 0; i < len(categories); i++ {
        // request the html page
        response, err := http.Get("https://callofduty.fandom.com/wiki/" + categories[i])
        if err != nil {
            log.Fatal(err)
        }
        // load the html document
        document, err := goquery.NewDocumentFromReader(response.Body)
        if err != nil {
            log.Fatal("Error loading HTTP response body.", err)
        }
        // count the number of guns for category[i]
        countGuns := 1

        // find the h3 headers
        document.Find("h3").Each(func (index1 int, s1 *goquery.Selection) {
            // if the h3 header is the text below
            if s1.Text() == "Call of Duty: Modern WarfareEdit" {
                // find the li a html tag
                s1.Next().Find("li a").Each(func (index2 int, s2 *goquery.Selection) {
                    // grab the link to that gun
                    href, exists := s2.Attr("href")
                    if exists {
                        g := Gun { Gun_Id: gunCount, Type: categories[i], Name: s2.Text() }
                        // get the attachments for this gun
                        attachments = append(attachments, getAttachments(href)...)
                        guns = append(guns, g)
                        countGuns++
                        gunCount++
                    }
                })
            }
        })
        fmt.Println("Gun Category: ", categories[i], " -- Number of guns: ", countGuns)
        response.Body.Close()
    }
    return guns,attachments
}

// gets the light machine guns
func getLightMachineGuns() ([]Gun,[]Attachment) {
    guns := []Gun{}
    attachments := []Attachment{}
    // request html page
    response, err := http.Get("https://callofduty.fandom.com/wiki/Light_Machine_Gun")
    if err != nil {
        log.Fatal(err)
    }
    defer response.Body.Close()

    // load the html document
    document, err := goquery.NewDocumentFromReader(response.Body)
    if err != nil {
        log.Fatal("Error loading HTTP response body.", err)
    }
    // count the number of guns for this category
    countGuns := 1
    // find div tags
    document.Find("div").Each(func (index1 int, s1 *goquery.Selection) {
        // check if this tag has a hash attribute
        h, exists := s1.Attr("hash")
        // check if the hash attribute is for modern warfare
        if exists && h == "6c97371c4254fbd9cc6781e6af829032" {
            // find the a tag
            s1.Find(".wikia-gallery-item .lightbox-caption a").Each(func (index2 int, s2 *goquery.Selection) {
                // grab the link to that gun
                href, exists := s2.Attr("href")
                if exists {
                    g := Gun { Gun_Id: gunCount, Type: "Light_Machine_Gun", Name: s2.Text() }
                    // get the attachments for this gun
                    attachments = append(attachments, getAttachments(href)...)
                    guns = append(guns, g)
                    gunCount++
                }
            })
        }
    })
    fmt.Println("Gun Category: Light_Machine_Gun -- Number of guns: ", countGuns)
    return guns,attachments
}

func main() {
    // get every category, guns, and attachments
    gunsRes, attRes := getFirstGuns()
    gunsRes2, attRes2 := getSecondGuns()
    gunsRes3, attRes3 := getLightMachineGuns()
    gunsRes = append(gunsRes, gunsRes2...)
    gunsRes = append(gunsRes, gunsRes3...)
    attRes = append(attRes, attRes2...)
    attRes = append(attRes, attRes3...)

    err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}
	db, err := sql.Open("postgres", os.Getenv("POSTGRES_URL"))
	if err != nil {
		panic(err)
	}
    defer db.Close()
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Connected to DB")

    for g := 0; g < len(gunsRes); g++ {
        id := gunsRes[g].Gun_Id
        t := gunsRes[g].Type
        name := gunsRes[g].Name
        if _,err := db.Exec("insert into guns (gun_id, type, name) values ($1, $2, $3)", id, t, name); err != nil {
            panic(err)
        }
    }
    for a := 0; a < len(attRes); a++ {
        id := attRes[a].Attachment_Id
        gid := attRes[a].Gun_Id
        name := attRes[a].Name
        subatts := pq.Array(attRes[a].SubAttachments)
        if _,err := db.Exec("insert into attachments (attachment_id, gun_id, name, subattachments) values ($1, $2, $3, $4)", id, gid, name, subatts); err != nil {
            panic(err)
        }
    }

}
