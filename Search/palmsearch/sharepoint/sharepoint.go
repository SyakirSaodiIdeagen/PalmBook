package sharepoint

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/patrickmn/go-cache"
	"io"
	"log"
	"net/http"
	"net/url"
	"palmsearch/elasticsearch"
	"strings"
	"time"
)

type TokenStore struct {
	AccessToken string `json:"access_token,omitempty"`
}

type Site struct {
	Id   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type List struct {
	CreatedAt   string `json:"createdDateTime,omitempty"`
	Description string `json:"description,omitempty"`
	Id          string `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	DisplayName string `json:"displayName,omitempty"`
	WebUrl      string `json:"webUrl,omitempty"`
}

type ListItem struct {
	CreatedAt      string      `json:"createdDateTime,omitempty"`
	LastModifiedAt string      `json:"lastModifiedDateTime,omitempty"`
	WebUrl         string      `json:"webUrl,omitempty"`
	Fields         Fields      `json:"fields,omitempty"`
	ContentType    ContentType `json:"contentType,omitempty"`
	Id             string      `json:"id,omitempty"`
	Name           string      `json:"Name,omitempty"`
	Source         string      `json:"source,omitempty"`
}

type Fields struct {
	Name string `json:"FileLeafRef,omitempty"`
}

type ContentType struct {
	Name string `json:"name,omitempty"`
}

type GetSiteResponse struct {
	Value []Site `json:"value"`
}

type GetListResponse struct {
	SiteId   string `json:"siteId,omitempty"`
	Value    []List `json:"value"`
	SiteName string `json:"siteName,omitempty"`
}

type GetListItemResponse struct {
	Value    []ListItem `json:"value"`
	SiteId   string     `json:"siteId,omitempty"`
	SiteName string     `json:"siteName,omitempty"`
	ListId   string     `json:"listId,omitempty"`
}

var ts = TokenStore{}
var c = cache.New(5*time.Minute, 10*time.Minute)

/*
get access token
1. get sites
2. get lists
3. get list items
4. index into elasticsearch
*/
func IndexSharepoint() {
	GetTenantAccessToken("")
	GetSites(ts.AccessToken)
}

func GetSites(token string) {

	req, err := http.NewRequest("GET", "https://graph.microsoft.com/v1.0/sites/getAllSites?select=name,id", nil)
	if err != nil {
		log.Fatal("Error creating request:", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}
	var gsr = GetSiteResponse{}
	err = json.Unmarshal(body, &gsr)
	if err != nil {
		panic(err)
	}

	fmt.Println("Response Status:", resp.Status)
	fmt.Println("Get Site Response Body:", gsr)

	selectedSites := []string{}
	filteredSites := filterSites(gsr, selectedSites)
	GetLists(filteredSites)
	log.Println("sync complete")
	cleanUp()
}

func cleanUp() {
	existDocs := elasticsearch.GetAll()

	items := c.Items()
	docToBeDeleted := []string{}

	for _, doc := range existDocs {

		if strings.HasPrefix(doc, "gd") {
			continue
		}

		if item, found := items[doc]; found {
			fmt.Printf("Found doc in cache: Key: %s, Value: %v\n", doc, item.Object)
		} else {
			fmt.Printf("Doc %s not found in cache\n", doc)
			docToBeDeleted = append(docToBeDeleted, doc)
		}
	}
	elasticsearch.DeleteDocumentsBulk(docToBeDeleted)
	c.Flush()
}

func filterSites(gsr GetSiteResponse, arr2 []string) []Site {

	if len(arr2) == 0 {
		return gsr.Value
	}

	arr1 := gsr.Value
	namesMap := make(map[string]struct{})
	for _, name := range arr2 {
		namesMap[name] = struct{}{}
	}

	var result []Site
	for _, site := range arr1 {
		if _, exists := namesMap[site.Name]; exists {
			result = append(result, site)
		}
	}

	return result
}

func filterContentType(gsr GetListItemResponse, arr2 []string) GetListItemResponse {
	arr1 := gsr.Value
	namesMap := make(map[string]struct{})
	for _, name := range arr2 {
		namesMap[name] = struct{}{}
	}

	var result []ListItem
	for _, listItem := range arr1 {
		if _, exists := namesMap[listItem.ContentType.Name]; exists {
			result = append(result, listItem)
		}
	}

	gsr.Value = result

	return gsr
}

func GetLists(gsr []Site) {
	for idx, val := range gsr {
		fmt.Printf("%v\t%v\n", idx, val)
		fmt.Printf("Processing %v", val.Name)
		req, err := http.NewRequest("GET", fmt.Sprintf("https://graph.microsoft.com/v1.0/sites/%v/lists", val.Id), nil)
		if err != nil {
			log.Fatal("Error creating request:", err)
		}

		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", ts.AccessToken))

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error sending request:", err)
			return
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading response:", err)
			return
		}
		var glr = GetListResponse{}
		err = json.Unmarshal(body, &glr)
		if err != nil {
			panic(err)
		}

		glr.SiteName = val.Name
		glr.SiteId = val.Id
		GetListItems(glr)
	}

}

func GetListItems(glr GetListResponse) {

	for _, val := range glr.Value {
		fmt.Printf("Processing List Item %v", val.Name)
		req, err := http.NewRequest("GET", fmt.Sprintf("https://graph.microsoft.com/v1.0/sites/%v/lists/%v/items?$select=contentType,lastModifiedDateTime,id,webUrl,createdDateTime,lastModifiedBy&$expand=fields($select=FileLeafRef)", glr.SiteId, val.Id), nil)
		if err != nil {
			log.Fatal("Error creating request:", err)
		}

		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", ts.AccessToken))

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error sending request:", err)
			return
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading response:", err)
			return
		}
		var glir = GetListItemResponse{}
		err = json.Unmarshal(body, &glir)
		if err != nil {
			panic(err)
		}

		fmt.Println("Response Status:", resp.Status)
		fmt.Println("Get List Item Response Body:", glir)
		esclient := elasticsearch.GetEsClient()
		glir.ListId = val.Id
		glir.SiteId = glr.SiteId
		glir.SiteName = glr.SiteName

		for i := range glir.Value {
			glir.Value[i].Name = glir.Value[i].Fields.Name
			glir.Value[i].Source = "Sharepoint"
		}

		log.Printf("glir: %v", glir.Value)

		selectedContentType := []string{"Document"}

		filteredListItems := filterContentType(glir, selectedContentType)
		buf := getElasticBuf(filteredListItems)
		elasticsearch.BulkInsert(esclient, buf)
	}

}

func GetTenantAccessToken(tenantId string) {
	formData := url.Values{}
	tenantId = "3192a717-1c36-4a32-b40f-d91972b86f32"
	formData.Set("client_id", "1afd0015-0d04-4b6f-bffa-e2d55db24f0b")
	formData.Set("client_secret", "wyi8Q~grqvwTcs9RJsb~nvtF_KiH.y9wF37jUbQC")
	formData.Set("scope", "https://graph.microsoft.com/.default")
	formData.Set("grant_type", "client_credentials")

	req, err := http.NewRequest("POST", fmt.Sprintf("https://login.microsoftonline.com/%v/oauth2/v2.0/token", tenantId), strings.NewReader(formData.Encode()))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	err = json.Unmarshal(body, &ts)
	if err != nil {
		panic(err)
	}

	fmt.Println("Response Status:", resp.Status)
	fmt.Println("Response Body:", ts.AccessToken)
}

func getElasticBuf(data GetListItemResponse) bytes.Buffer {

	var buf bytes.Buffer
	for _, doc := range data.Value {
		id := data.SiteName + "_" + data.ListId + "_" + doc.Id
		c.Set(id, "golang-bulk-index3", cache.DefaultExpiration)

		meta := []byte(fmt.Sprintf(`{ "index" : { "_index" : "golang-bulk-index3", "_id" : "%v" } }%s`, id, "\n"))
		d, err := json.Marshal(doc)
		if err != nil {
			log.Fatalf("Cannot encode document %d: %s", doc.Id, err)
		}

		buf.Grow(len(meta) + len(d) + 1)
		buf.Write(meta)
		buf.Write(d)
		buf.WriteByte('\n')
	}

	return buf
}
func VerifyUserAccessToken() {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://graph.microsoft.com/v1.0/me"), nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	req.Header.Set("Authorization", "Bearer eyJ0eXAiOiJKV1QiLCJub25jZSI6IjhOR2lYOGtVSTBaYUpHaFJhT253X0pZRWxpM2Job1JRdk9UaTBvUHZ0YTQiLCJhbGciOiJSUzI1NiIsIng1dCI6IllUY2VPNUlKeXlxUjZqekRTNWlBYnBlNDJKdyIsImtpZCI6IllUY2VPNUlKeXlxUjZqekRTNWlBYnBlNDJKdyJ9.eyJhdWQiOiIwMDAwMDAwMy0wMDAwLTAwMDAtYzAwMC0wMDAwMDAwMDAwMDAiLCJpc3MiOiJodHRwczovL3N0cy53aW5kb3dzLm5ldC8zMTkyYTcxNy0xYzM2LTRhMzItYjQwZi1kOTE5NzJiODZmMzIvIiwiaWF0IjoxNzM4NDEzMTIwLCJuYmYiOjE3Mzg0MTMxMjAsImV4cCI6MTczODQ5OTgyMCwiYWNjdCI6MCwiYWNyIjoiMSIsImFpbyI6IkFWUUFxLzhaQUFBQXBhL1l2a1d3OHExcFB0NVhIb0ozbTEzNEJteGN5RXFZdkxoOWFYdmhsYkhMRlRKRjE0c3NLdUJHRVhUb0tVckltVUdoOHAxN0NxNlhLYjc4N0FCOWtMTXhqeTZPb0g3SXo5UjhPWGZWTENRPSIsImFtciI6WyJwd2QiLCJtZmEiXSwiYXBwX2Rpc3BsYXluYW1lIjoiR3JhcGggRXhwbG9yZXIiLCJhcHBpZCI6ImRlOGJjOGI1LWQ5ZjktNDhiMS1hOGFkLWI3NDhkYTcyNTA2NCIsImFwcGlkYWNyIjoiMCIsImZhbWlseV9uYW1lIjoiRmF1emkiLCJnaXZlbl9uYW1lIjoiRmFyaGFuIiwiaWR0eXAiOiJ1c2VyIiwiaXBhZGRyIjoiMjAwMTpmNDA6OTUwOjZhY2Y6ZjU1ZDo0NmVkOmQ2MGY6NjYyMiIsIm5hbWUiOiJGYXJoYW4gRmF1emkiLCJvaWQiOiIwM2IzM2FmNC1hNTcwLTQ2OTQtOTVkNy0zMDIwYTUwZWMzZmIiLCJwbGF0ZiI6IjMiLCJwdWlkIjoiMTAwMzIwMDQzQzFCNjJFOCIsInJoIjoiMS5BYjRBRjZlU01UWWNNa3EwRDlrWmNyaHZNZ01BQUFBQUFBQUF3QUFBQUFBQUFBQy1BQzYtQUEuIiwic2NwIjoiQVBJQ29ubmVjdG9ycy5SZWFkLkFsbCBBUElDb25uZWN0b3JzLlJlYWRXcml0ZS5BbGwgQXBwbGljYXRpb24uUmVhZC5BbGwgb3BlbmlkIHByb2ZpbGUgU2hhcmVQb2ludFRlbmFudFNldHRpbmdzLlJlYWQuQWxsIFNpdGVzLkZ1bGxDb250cm9sLkFsbCBTaXRlcy5SZWFkLkFsbCBTaXRlcy5SZWFkV3JpdGUuQWxsIFVzZXIuUmVhZCBlbWFpbCIsInNpZCI6IjAwMTM0MWM5LWM2MGYtNTk4Yi1kMWY1LTc3ODAyYWFiNDc5NyIsInNpZ25pbl9zdGF0ZSI6WyJrbXNpIl0sInN1YiI6IlI5ak9xcDhBbW1FMm9xdEMxMHlsYlo4OU5vT2Z1ejAxblRzS1hQYnNrRWsiLCJ0ZW5hbnRfcmVnaW9uX3Njb3BlIjoiQVMiLCJ0aWQiOiIzMTkyYTcxNy0xYzM2LTRhMzItYjQwZi1kOTE5NzJiODZmMzIiLCJ1bmlxdWVfbmFtZSI6ImZhcmhhbi5mYXV6aUBkZXN0aW55bGluay5vbm1pY3Jvc29mdC5jb20iLCJ1cG4iOiJmYXJoYW4uZmF1emlAZGVzdGlueWxpbmsub25taWNyb3NvZnQuY29tIiwidXRpIjoiUXFya2JRbVlJazIza2ljeV9ZZTVBQSIsInZlciI6IjEuMCIsIndpZHMiOlsiNjJlOTAzOTQtNjlmNS00MjM3LTkxOTAtMDEyMTc3MTQ1ZTEwIiwiYjc5ZmJmNGQtM2VmOS00Njg5LTgxNDMtNzZiMTk0ZTg1NTA5Il0sInhtc19jYyI6WyJDUDEiXSwieG1zX2Z0ZCI6IklRR09nekFybEh5SXFwUFhGUjhhOENuNVV5Qmt5eFc4QUZTMXhwblA0TFkiLCJ4bXNfaWRyZWwiOiIxIDEyIiwieG1zX3NzbSI6IjEiLCJ4bXNfc3QiOnsic3ViIjoiQmp6NWxPRDRLbVNodFJudTdQYjZoTXpKcExETnlVRVhKQmR0V1o3MktJTSJ9LCJ4bXNfdGNkdCI6MTczNzg4Nzc1MX0.iVWhHzUWxOiJkzeUI05SlLKCNvwV4Fti98W7HwDCDICMlCiQ1jhJv6OiJby4V1ymJlncQ2B57TOPfaYm7O287UjqjYNYmV-0EYNFU4MXdJcFAqxe1kt0HcRIYiJygAmysOXonNKj_TOBBKIVRoP4kBgvhABmUyFggp7YIGJ8kurP9z10wWO6cjJyGOiPxRlqRcCgbwtDsToIq3fPEsy-HtQ8qO7oqLMoygAKyx8b9jqERLOGLq8JymT900UtFhDUcfY21VZYq3yVd8TkfBFnc7HKDTMjAdPcDbNluEWuK_Bkm5uQcBHUE5Z5CB2lCrtaVjEGMPqUlbGLIizUsIoMvg")
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	fmt.Println("Response Status:", resp.Status)
	fmt.Println("Response Body:", string(body))
}
