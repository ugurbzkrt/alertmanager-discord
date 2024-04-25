package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/namsral/flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

//json structure from Alertmanager
type alertManOut struct {
	Alerts            []alertManAlert `json:"alerts"`
	CommonAnnotations struct {
		Summary string `json:"summary"`
	} `json:"commonAnnotations"`
	CommonLabels struct {
		Alertname string `json:"alertname"`
	} `json:"commonLabels"`
	ExternalURL string `json:"externalURL"`
	GroupKey    string `json:"groupKey"`
	GroupLabels struct {
		Alertname string `json:"alertname"`
	} `json:"groupLabels"`
	Receiver string `json:"receiver"`
	Status   string `json:"status"`
	Version  string `json:"version"`
}

type alertManAlert struct {
	Annotations struct {
		Description string `json:"description"`
		Summary     string `json:"summary"`
	} `json:"annotations"`
	EndsAt       string            `json:"endsAt"`
	GeneratorURL string            `json:"generatorURL"`
	Labels       map[string]string `json:"labels"`
	StartsAt     string            `json:"startsAt"`
	Status       string            `json:"status"`
}

//data sent to Discord. Where Name - is bot Name
type discordOut struct {
	Content string `json:"content"`
	Name    string `json:"username"`
}

var (
	config      string
	address     string
	webhookUrl  string
	discordName string
	logger      *log.Logger
	v           bool
)

func handler(w http.ResponseWriter, r *http.Request) {

	//parse json from Alertmanager
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	if v == true {
		logger.Println("body:", string(b), "\n")
	}

	amo := alertManOut{}
	err = json.Unmarshal(b, &amo)
	if err != nil {
		panic(err)
	}

	groupedAlerts := make(map[string][]alertManAlert)

	for _, alert := range amo.Alerts {
		groupedAlerts[alert.Status] = append(groupedAlerts[alert.Status], alert)
	}

	for status, alerts := range groupedAlerts {
		DO := discordOut{
			//Name: status,
			Name: discordName,
		}

		Content := "```"
		if amo.CommonAnnotations.Summary != "" {
			Content = fmt.Sprintf(" === %s === \n```", amo.CommonAnnotations.Summary)
		}

		for _, alert := range alerts {
			realname := alert.Labels["instance"]
			//if strings.Contains(realname, "localhost") && alert.Labels["exported_instance"] != "" {
			//	realname = alert.Labels["exported_instance"]
			//}
			Content += fmt.Sprintf("[%s]: %s on %s\n%s\n", strings.ToUpper(status), alert.Labels["alertname"], realname, alert.Annotations.Description)
			if alert.Labels["severity"] != "" {
				Content += fmt.Sprintf("Severity: %s\n\n", alert.Labels["severity"])
			}
		}

		DO.Content = Content + "```"

		DOD, _ := json.Marshal(DO)
		http.Post(webhookUrl, "application/json", bytes.NewReader(DOD))

	}
}

func main() {
	flag.StringVar(&config, "config", "", "Config file with variables - optional. Can parse both variables from config and CLI")
	flag.StringVar(&address, "address", ":9095", "Service listen address and port")
	flag.StringVar(&webhookUrl, "discord_webhook", "", "DISCORD_WEBHOOK to push messages")
	flag.StringVar(&discordName, "discord_name", "AlertManager", "DISCORD_NAME of bot pushing messages")
	flag.BoolVar(&v, "verbose", false, "Verbose mode")
	flag.Parse()

	if webhookUrl == "" {
		fmt.Fprintf(os.Stderr, "error: environment variable DISCORD_WEBHOOK not found\n")
		os.Exit(1)
	}

	logger = log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)
	logger.Println("Starting listening on ", address, "\n")

	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(address, nil))
}
