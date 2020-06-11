package jam

type StackCreateConfig struct {
	Repository      string
	CustomSubDomain string
	IsThirdParty    bool
}

type StackDestroyConfig struct {
	ID           string
	IsThirdParty bool
}

// RoutePrefix is the JAM endpoint
const RoutePrefix = "/jam"

// // HandleHTTP handles the JAM http endpoint
// func (j *Creator) ServeHTTP(w http.ResponseWriter, req *http.Request) {
// 	methodHandlerMap := map[string]http.HandlerFunc{
// 		http.MethodGet:    j.handleGET,
// 		http.MethodPost:   j.handlePOST,
// 		http.MethodDelete: j.handleDELETE,
// 	}

// 	if handler, ok := methodHandlerMap[req.Method]; ok {
// 		handler(w, req)
// 		return
// 	}

// 	w.WriteHeader(http.StatusMethodNotAllowed)
// }

// func (j *Creator) handlePOST(w http.ResponseWriter, req *http.Request) {
// 	var config stackConfig

// 	err := json.NewDecoder(req.Body).Decode(&config)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	updatedList, err := j.build(req.Context(), config)
// 	if err != nil {
// 		http.Error(w, err.Error(), 500)
// 		return
// 	}

// 	w.Write(updatedList)
// }

// func (j *Creator) handleGET(w http.ResponseWriter, req *http.Request) {
// 	stacks, _ := j.list(req.Context())
// 	stackJSON, err := json.Marshal(stacks)
// 	if err != nil {
// 		log.Println(err)
// 	}

// 	w.Write(stackJSON)
// }

// func (j *Creator) handleDELETE(w http.ResponseWriter, req *http.Request) {
// 	urlParams := req.URL.Query()
// 	idQueries, ok := urlParams["id"]
// 	if !ok {
// 		io.WriteString(w, "Invalid Query. Missing \"id\"")
// 		return
// 	}

// 	bucketID := idQueries[0]
// 	newlist, err := j.delete(req.Context(), bucketID)
// 	if err != nil {
// 		http.Error(w, err.Error(), 500)
// 		return
// 	}

// 	w.Write(newlist)
// }

// // const RouteCredentialsPrefix = "/creds"

// // type ThirdPartyAWSCredentials struct {
// // 	SecretKey string
// // 	AccessKey string
// // }

// // func SetThirdPartyAWSCredentials(w http.ResponseWriter, req *http.Request) {
// // 	if req.Method != http.MethodPost {
// // 		w.WriteHeader(http.StatusMethodNotAllowed)
// // 		return
// // 	}

// // 	var creds ThirdPartyAWSCredentials

// // 	err := json.NewDecoder(req.Body).Decode(&creds)
// // 	if err != nil {
// // 		http.Error(w, err.Error(), http.StatusBadRequest)
// // 		return
// // 	}

// // 	if creds.AccessKey == "" || creds.SecretKey == "" {
// // 		http.Error(w, "AccessKey or SecretKey missing in request", http.StatusBadRequest)
// // 		return
// // 	}

// // 	awsSess, err := session.AWSSessionThirdParty(creds.AccessKey, creds.SecretKey)
// // 	if err != nil {
// // 		log.Fatal("could not authenticate against AWS", err)
// // 	}

// // 	buc := s3.New(awsSess)

// // 	res, err := buc.ListBuckets(&s3.ListBucketsInput{})
// // 	if err != nil {
// // 		http.Error(w, err.Error(), http.StatusBadRequest)
// // 		return
// // 	}

// // 	log.Println(res.String())

// // 	log.Println("CREDS!", creds.AccessKey, creds.SecretKey)

// // 	w.WriteHeader(http.StatusOK)
// // }
