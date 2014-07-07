package main

type ServerConfig struct {
	id         int64
	Name       string
	Logfile    string
	Configfile string
	PortRTMP   int
	PortFLV    int
}

/*
func Create_config_file(id int) string {
	pnsEnv[id].config.logfileName()
	configfile_name := pnsEnv[id].config.configfilePath()
	config_file_tmpl, err := template.New("configfile").Parse(TMPL)
	if err != nil {
		log.Printf("%v", err)
	}
	file, err := os.Create(configfile_name)
	w := bufio.NewWriter(file)
	err = config_file_tmpl.Execute(w, pnsEnv[id].config)
	if err != nil {
		log.Printf("%v", err)
	}
	w.Flush()
	file.Close()
	return configfile_name
}

func (serverConfig *ServerConfig) configfilePath() string {
	if serverConfig.Configfile == "" {
		var configfile_name bytes.Buffer
		configfile_name.WriteString(CONFIG_PATH)
		configfile_name.WriteString(serverConfig.Name)
		configfile_name.WriteString(".lua")
		serverConfig.Configfile = configfile_name.String()
	}
	return serverConfig.Configfile
}

func (serverConfig *ServerConfig) logfileName() string {
	if serverConfig.Logfile == "" {
		var logfile_name bytes.Buffer
		logfile_name.WriteString(strconv.FormatInt(time.Now().Unix(), 10))
		logfile_name.WriteString(serverConfig.Name)
		serverConfig.Logfile = logfile_name.String()
	}
	return serverConfig.Logfile
}
*/
