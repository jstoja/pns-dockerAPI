configuration=
{
	daemon=false, -- False : Mode console. True : Mode daemon. Pour du dev, laisser en false.
	pathSeparator="/", -- Sous Windows, le path doit être à "\\" (Penser à escape le caractère). Sous Unix, le path est à "/"




	--[[
		logAppenders={}

		Définit comment les logs sont traités.
		name : Requis. Le nom du logger
		type : Requis. "console", "coloredConsole" ou "file". "console" et "coloredConsole" envoie sur la console. "file" log dans un fichier.
		level : Requis. 0-6. Définit le niveau de log. 
				0	 FATAL
				1	 ERROR
				2	 WARNING
				3	 INFO
				4	 DEBUG
				5	 FINE
				6	 FINEST

		fileName : Requis. Le path complet du fichier dans le cas d'un type "file".
	]]
	logAppenders= -- Définit comment les logs sont traités. 
	{
		{ -- Ce bloc est optionnel dans le cas d'un environnement de production. Utile pour le dev pour garder du log dans la console. Si "daemon" est à true, le console appender sera ignoré dans tous les cas.
			name="console appender", 
			type="coloredConsole",
			level=6
		},
		{ -- Ce bloc est essentiel en production pour garder des logs dans un fichier. Optionnel dans un environnement de dev.
			name="file appender", -- Permet de garder une trace dans les logs du nom. 
			type="file",
			level=6,
			fileName="./logs/crtmpserver",
			fileHistorySize=10,
			fileLength=1024*1024,
			singleLine=true	
		}
	},
	

	--[[
		applications={}

		Définit la liste des "connecteurs" / "applications" que CRTMP va fournir. 
		Structure : 
		applications {
			{application1},
			{application2},
			...
		}

		Définition d'une application : 

		name : Requis. Nom de l'application
		description : Optionnel. Une description de l'application
		protocol : Requis. "dynamiclinklibrary" signifit que l'application est une librairie partagée.
		acceptors : Optionnel. Permet de définir les protocoles d'écoutes. 
			ip : Requis. Si 0.0.0.0, permet d'écouter sur toute les adresses IP et toute les interfaces.
			port : Requis. Le port sur lequel l'appplication écoutera.
			protocol : Requis. Le protocole utilisé pour l'adresse IP et le port donné.
				Liste des protocoles utiles pour P&S : "inboundRtmp", "inboundLiveFlv". Important de laisser les deux protocoles pour un client, sous WIndows, le RTMP ne fonctionne pas, ns devons streamer en TCP directement. 
			waitForMetadata : Optionnel. Laisser à true pour le inboundLiveFLV, afin de forcer d'attendre la réception des metadonnées indiquant le stream avant de lancer.
	]]
	applications=
	{
		rootDirectory="applications",
		{
			description="FLV Playback Sample",
			name="flvplayback",
			protocol="dynamiclinklibrary",
			default=true,
			acceptors = 
			{
				{
					ip="0.0.0.0",
					port=1935,
					protocol="inboundRtmp"
				},
				{
					ip="0.0.0.0",
					port=1234,
					protocol="inboundLiveFlv",
					waitForMetadata=true,
				},
			},
			validateHandshake=false, --Toute les connections sont acceptés, le HandShake n'est pas vérifié.
			keyframeSeek=false, --Devrait toujours être à false pour du LiveStreaming. Cherche les "frames" clés de la vidéo pour uniformiser en cas de mauvaise connection. Vu qu'il n'y a pas de frames suivantes à l'actuelle à analyser dans le cas de livestream, sur false.
			seekGranularity=1.5, --in seconds, between 0.1 and 600. Utile pour la qualité de la vidéo, la "granularité". Définition CRTMP : La résolution/granularité en secondes. Par exemple, si la granularité est de 10 secondes, et le seek à t = 2:34, le seek sera de t = 2:30.
			clientSideBuffer=5, --in seconds, between 5 and 30. La quantité de tampon côté client qui sera maintenue pour chaque connexion
		},
	}
}


